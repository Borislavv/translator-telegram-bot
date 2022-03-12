package telegram

import (
	"fmt"
	"log"
	"regexp"
	"sync"
	"time"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDB"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/dashboardService"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/telegram/command"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/translator"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/util"
)

type TelegramService struct {
	// Dependencies
	manager *manager.Manager
	gateway *TelegramGateway

	// Services
	userService    *service.UserService
	chatService    *service.ChatService
	translator     *translator.TranslatorService
	tokenGenerator *dashboardService.TokenGenerator
	commandService *command.CommandService

	// Channels
	messagesChannel      chan *model.UpdatedMessage
	notificationsChannel chan *modelDB.NotificationQueue
	storeChannel         chan *model.UpdatedMessage
	tokensChannel        chan *model.TokenMessage
	commandsChannel      chan *model.CommandMessage
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
	tokenGenerator *dashboardService.TokenGenerator,
	commandService *command.CommandService,
	messagesChannel chan *model.UpdatedMessage,
	notificationsChannel chan *modelDB.NotificationQueue,
	storeChannel chan *model.UpdatedMessage,
	tokensChannel chan *model.TokenMessage,
	commandsChannel chan *model.CommandMessage,
	errorsChannel chan string,
) *TelegramService {
	return &TelegramService{
		manager:              manager,
		gateway:              gateway,
		userService:          userService,
		chatService:          chatService,
		translator:           translator,
		tokenGenerator:       tokenGenerator,
		commandService:       commandService,
		messagesChannel:      messagesChannel,
		notificationsChannel: notificationsChannel,
		storeChannel:         storeChannel,
		tokensChannel:        tokensChannel,
		commandsChannel:      commandsChannel,
		errorsChannel:        errorsChannel,
	}
}

// GetNotifications - receiving notification from the database
func (telegramService *TelegramService) GetNotifications(wg *sync.WaitGroup) {
	defer wg.Done()

	// Receiving notifications from database
	notifications, err := telegramService.manager.Repository.NotificationQueue().FindByScheduledDate()
	if err != nil {
		telegramService.errorsChannel <- util.Trace(err)
		return
	}

	for _, notification := range notifications {
		telegramService.notificationsChannel <- notification
	}
}

// SendNotifications - sending notifications to telegram chat
func (telegramService *TelegramService) SendNotifications(wg *sync.WaitGroup) {
	defer wg.Done()

	for notification := range telegramService.notificationsChannel {
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
				continue
			}

			telegramService.lastReceivedOffset = offset
		}

		messages, err := telegramService.gateway.GetUpdates(model.NewTelegramRequestMessage(telegramService.lastReceivedOffset))
		if err != nil {
			telegramService.errorsChannel <- util.Trace(err)
			continue
		}

		var receivedOffset int64
		for _, message := range messages.Messages {

			msg := message

			// Run `SendMessages` gorutine
			telegramService.messagesChannel <- &msg

			// Run `StoreMessages` gorutine
			telegramService.storeChannel <- &msg

			// Log info of message received
			log.Printf("Message received: %+v\n", msg)

			receivedOffset = message.QueueId
		}

		if receivedOffset != 0 {
			telegramService.lastReceivedOffset = receivedOffset + 1
		}

		time.Sleep(30 * time.Millisecond)
	}
}

// SendMessages - sending message to telegram chat
func (telegramService *TelegramService) SendMessages() {
	for {
		select {
		case message := <-telegramService.messagesChannel:
			notificationMatchedValue := regexp.MustCompile(`\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}`).FindStringSubmatch(message.Data.Text)

			var err error
			var text string

			if len(regexp.MustCompile(`\/token`).FindStringSubmatch(message.Data.Text)) > 0 {
				token, err := telegramService.tokenGenerator.Generate()
				if err != nil {
					telegramService.errorsChannel <- util.Trace(err)
					continue
				}

				// processing token request
				text = "Token: " + token

				telegramService.tokensChannel <- model.NewTokenMessage(message, token)
			} else if len(regexp.MustCompile(`\/\w+`).FindStringSubmatch(message.Data.Text)) > 0 {
				telegramService.commandsChannel <- model.NewCommandMessage(message)
			} else if len(notificationMatchedValue) > 0 {
				// processing notification
				text = "Notification setted on " + notificationMatchedValue[0]
			} else {
				// processing translation
				text, err = telegramService.translator.TranslateText(message.Data.Text)
				if err != nil {
					telegramService.errorsChannel <- util.Trace(err)
					continue
				}

				text = "Tranlsation: " + text
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

// ProcessCommands - processing commands which is available from a client.
func (telegramService *TelegramService) ProcessCommands() {
	for {
		select {
		case command := <-telegramService.commandsChannel:
			telegramService.gateway.SendMessage(
				telegramService.commandService.ProcessCommand(command),
			)
		}
	}
}

// StoreMessages - storing messages and notification into database
func (telegramService *TelegramService) StoreMessages() {
	for {
		select {
		case message := <-telegramService.storeChannel:
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

			user, err := telegramService.userService.GetUser(message.Data.Chat.Username, chat.ID)
			if err != nil {
				telegramService.errorsChannel <- util.Trace(err)
				continue
			}

			// check if the message is notification
			notificationMatchedValue := regexp.MustCompile(`\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}`).FindStringSubmatch(messageQueue.Message)
			if len(notificationMatchedValue) > 0 &&
				len(regexp.MustCompile(`\/\w+`).FindStringSubmatch(messageQueue.Message)) == 0 {

				utcLocation, err := time.LoadLocation("UTC")
				if err != nil {
					telegramService.errorsChannel <- util.Trace(err)
					continue
				}
				userLocation, err := time.LoadLocation(user.TZ)
				if err != nil {
					telegramService.errorsChannel <- util.Trace(err)
					continue
				}

				// processing notification
				scheduledFor, err := time.ParseInLocation("2006-01-02 15:04:05", notificationMatchedValue[0], userLocation)
				if err != nil {
					telegramService.errorsChannel <- util.Trace(err)
					continue
				}

				_, err = telegramService.manager.Repository.NotificationQueue().Create(
					modelDB.NewNotificationQueueConstructor(
						messageQueue.ID,
						messageQueue.ChatId,
						scheduledFor.In(utcLocation),
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

// StoreTokens - saving requested tokens into database to users
func (telegramService *TelegramService) StoreTokens() {
	for {
		select {
		case data := <-telegramService.tokensChannel:
			chat, err := telegramService.chatService.GetChat(data.UpdatedMessage.Data.Chat.ID)
			if err != nil {
				telegramService.errorsChannel <- util.Trace(err)
				continue
			}

			user, err := telegramService.userService.GetUser(data.UpdatedMessage.Data.Chat.Username, chat.ID)
			if err != nil {
				telegramService.errorsChannel <- util.Trace(err)
				continue
			}

			user, err = telegramService.userService.SetToken(user, data.Token)
			if err != nil {
				telegramService.errorsChannel <- util.Trace(err)
				continue
			}

			log.Printf("Token was stored: %+v\n", user)
		}
	}
}
