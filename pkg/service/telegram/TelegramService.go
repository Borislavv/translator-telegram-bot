package telegram

import (
	"fmt"
	"log"
	"regexp"
	"sync"
	"time"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/config"
	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDB"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/translator"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/util"
)

type TelegramService struct {
	// Dependencies
	manager *manager.Manager
	gateway *TelegramGateway

	// Services
	userService *service.UserService
	chatService *service.ChatService
	translator  *translator.TranslatorService

	// Channels
	messagesChannel      chan *model.UpdatedMessage
	notificationsChannel chan *modelDB.NotificationQueue
	storechannel         chan *model.UpdatedMessage
	errorsChannel        chan string

	// Cached values
	lastReceivedOffset int64
}

// NewTelegramService - constructor of TelegramService
func NewTelegramService(
	manager *manager.Manager,
	gateway *TelegramGateway,
	userService *service.UserService,
	chatService *service.ChatService,
	translator *translator.TranslatorService,
	messagesChannel chan *model.UpdatedMessage,
	notificationsChannel chan *modelDB.NotificationQueue,
	storeChannel chan *model.UpdatedMessage,
	errorsChannel chan string,
) *TelegramService {
	return &TelegramService{
		manager:              manager,
		gateway:              gateway,
		userService:          userService,
		chatService:          chatService,
		translator:           translator,
		messagesChannel:      messagesChannel,
		notificationsChannel: notificationsChannel,
		storechannel:         storeChannel,
		errorsChannel:        errorsChannel,
	}
}

// GetNotifications - receiving notification from the database
func (telegramService *TelegramService) GetNotifications() []*modelDB.NotificationQueue {
	dateTime := time.Now()
	// Handling Timezone colission
	if telegramService.manager.Config.Environment.Mode == config.ProdMode {
		dateTime = dateTime.Add(2 * time.Hour)
	}

	// Receiving notifications from database
	notifications, err := telegramService.manager.Repository.NotificationQueue().FindByScheduledDate(dateTime)
	if err != nil {
		telegramService.errorsChannel <- util.Trace(err)
		return nil
	}

	var notificationStack []*modelDB.NotificationQueue
	for _, notification := range notifications {
		notificationStack = append(notificationStack, notification)
	}

	return notificationStack
}

// SendNotifications - sending notifications to telegram chat
func (telegramService *TelegramService) SendNotifications(notifications []*modelDB.NotificationQueue, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, notification := range notifications {
		// Print that notification is sent to CLI
		log.Printf("Notification have been sent: %+v\n", notification)

		// TODO: refactor this code, because the notification already has prop. ExternalChatId
		chat, err := telegramService.chatService.GetChat(notification.ExternalChatId)
		if err != nil {
			telegramService.errorsChannel <- util.Trace(err)
			continue
		}

		if err := telegramService.gateway.SendMessage(
			model.NewTelegramResponseMessage(
				fmt.Sprint(chat.ExternalChatId),
				"Notification: "+notification.Message,
			),
		); err != nil {
			telegramService.errorsChannel <- util.Trace(err)
			continue
		}

		_, err = telegramService.manager.Repository.NotificationQueue().MakeAsSent(notification)
		if err != nil {
			telegramService.errorsChannel <- util.Trace(err)
			continue
		}
	}
}

// GetMessages - receive messages from telegram chat
func (telegramService *TelegramService) GetMessages() {
	for {
		if telegramService.lastReceivedOffset == 0 {
			offset, err := telegramService.manager.Repository.MessageQueue().GetOffset()
			if err != nil {
				telegramService.errorsChannel <- util.Trace(err)
				return
			}

			telegramService.lastReceivedOffset = offset
		}

		messages, err := telegramService.gateway.GetUpdates(model.NewTelegramRequestMessage(telegramService.lastReceivedOffset))
		if err != nil {
			telegramService.errorsChannel <- util.Trace(err)
			return
		}

		for _, message := range messages.Messages {
			// Run `SendMessages` gorutine
			telegramService.messagesChannel <- &message

			// Run `StoreMessages` gorutine
			telegramService.storechannel <- &message

			// Log info of message received
			log.Printf("Message received: %+v\n", message)
		}

		telegramService.lastReceivedOffset = telegramService.lastReceivedOffset + int64(len(messages.Messages))

		time.Sleep(15 * time.Millisecond)
	}
}

// SendMessages - sending message to telegram chat
func (telegramService *TelegramService) SendMessages() {
	for {
		select {
		case message := <-telegramService.messagesChannel:
			matchedValue := regexp.MustCompile(`\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}`).FindStringSubmatch(message.Data.Text)

			var err error
			var text string
			if len(matchedValue) <= 0 {
				// processing translation
				text, err = telegramService.translator.TranslateText(message.Data.Text)
				if err != nil {
					telegramService.errorsChannel <- util.Trace(err)
					continue
				}

				text = "Tranlsation: " + text
			} else {
				// processing notification
				text = "Notification setted on " + matchedValue[0]
			}

			if err := telegramService.gateway.SendMessage(
				model.NewTelegramResponseMessage(
					fmt.Sprint(message.Data.Chat.ID),
					text,
				),
			); err != nil {
				telegramService.errorsChannel <- util.Trace(err)
				continue
			}

			log.Printf("Message was sent: %+v\n", message)
		}
	}
}

// storeMessages - storing messages and notification into database
func (telegramService *TelegramService) StoreMessages() {
	for {
		select {
		case message := <-telegramService.storechannel:
			chat, err := telegramService.chatService.GetChat(message.Data.Chat.ID)
			if err != nil {
				telegramService.errorsChannel <- util.Trace(err)
				continue
			}

			messageQueue := modelDB.NewMessageQueueConstructor(
				message.QueueId,
				message.Data.Text,
				chat.ID,
			)

			if _, err = telegramService.manager.Repository.MessageQueue().Create(messageQueue); err != nil {
				telegramService.errorsChannel <- util.Trace(err)
				continue
			}

			if _, err = telegramService.userService.GetUser(message.Data.Chat.Username, chat.ID); err != nil {
				telegramService.errorsChannel <- util.Trace(err)
				continue
			}

			// check if the message is notification
			matchedValue := regexp.MustCompile(`\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}`).FindStringSubmatch(messageQueue.Message)
			if len(matchedValue) > 0 {
				// processing notification
				scheduledFor, err := time.Parse("2006-01-02 15:04:05", matchedValue[0])
				if err != nil {
					telegramService.errorsChannel <- util.Trace(err)
					continue
				}

				_, err = telegramService.manager.Repository.NotificationQueue().Create(
					modelDB.NewNotificationQueueConstructor(
						messageQueue.ID,
						messageQueue.ChatId,
						scheduledFor,
					),
				)
				if err != nil {
					telegramService.errorsChannel <- util.Trace(err)
					continue
				}
			}

			log.Printf("Message was stored: %+v\n", message)
		}
	}
}
