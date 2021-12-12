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
	notificationsChannel chan *modelDB.NotificationQueue
	messagesChannel      chan *model.UpdatedMessage
	errorsChannel        chan string
	storeChannel         chan *model.UpdatedMessage

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
	notificationsChannel chan *modelDB.NotificationQueue,
	messagesChannel chan *model.UpdatedMessage,
	errorsChannel chan string,
	storeChannel chan *model.UpdatedMessage,
) *TelegramService {
	return &TelegramService{
		manager:              manager,
		gateway:              gateway,
		userService:          userService,
		chatService:          chatService,
		translator:           translator,
		notificationsChannel: notificationsChannel,
		messagesChannel:      messagesChannel,
		errorsChannel:        errorsChannel,
		storeChannel:         storeChannel,
	}
}

// GetNotifications - (gorutine) receive notifications and deliver them to the `notificationChannel` for another goroutine
func (telegramService *TelegramService) GetNotifications() {
	dateTime := time.Now()
	// Handling Timezone colission
	if telegramService.manager.Config.Environment.Mode == config.ProdMode {
		dateTime = dateTime.Add(2 * time.Hour)
	}

	// Receiving notifications from database
	notifications, err := telegramService.manager.Repository.NotificationQueue().FindByScheduledDate(dateTime)
	if err != nil {
		telegramService.errorsChannel <- util.Trace(err)
		return
	}

	// Providing notifications to channel for  another gorutine
	for _, notification := range notifications {
		telegramService.notificationsChannel <- notification
	}
}

// SendNotifications - (gorutine) pick up notifications from the `notificationsChannel` and send them to the telegram chat
func (telegramService *TelegramService) SendNotifications() {
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

// GetMessages - (gorutine) receive messages and deliver them to the `messagesChannel` for another goroutine
func (telegramService *TelegramService) GetMessages(m *sync.Mutex) {
	for {
		if telegramService.lastReceivedOffset == 0 {
			offset, err := telegramService.manager.Repository.MessageQueue().GetOffset()
			if err != nil {
				telegramService.errorsChannel <- util.Trace(err)
				return
			}

			m.Lock()
			telegramService.lastReceivedOffset = offset
			m.Unlock()
		}

		messages, err := telegramService.gateway.GetUpdates(model.NewTelegramRequestMessage(telegramService.lastReceivedOffset))
		if err != nil {
			telegramService.errorsChannel <- util.Trace(err)
			return
		}

		m.Lock()
		for _, message := range messages.Messages {
			go telegramService.SendMessages(&message)
			go telegramService.StoreMessages(&message)

			// // send message for sending to telegram chat
			// telegramService.messagesChannel <- &message

			// send message for save into database
			// telegramService.storeChannel <- &message

		}

		telegramService.lastReceivedOffset = telegramService.lastReceivedOffset + int64(len(messages.Messages))
		m.Unlock()
	}
}

// SendMessages - (gorutine) pick up messages from the `messagesChannel` and send them to the telegram chat
func (telegramService *TelegramService) SendMessages(processingMessage *model.UpdatedMessage) {
	matchedValue := regexp.MustCompile(`\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}`).FindStringSubmatch(processingMessage.Data.Text)

	var err error
	var text string
	if len(matchedValue) <= 0 {
		// processing translation
		text, err = telegramService.translator.TranslateText(processingMessage.Data.Text)
		if err != nil {
			telegramService.errorsChannel <- util.Trace(err)
			return
		}

		text = "Tranlsation: " + text
	} else {
		// processing notification
		text = "Notification setted on " + matchedValue[0]
	}

	if err := telegramService.gateway.SendMessage(
		model.NewTelegramResponseMessage(
			fmt.Sprint(processingMessage.Data.Chat.ID),
			text,
		),
	); err != nil {
		telegramService.errorsChannel <- util.Trace(err)
		return
	}
}

// storeMessages - (gorutine) getting `UpdatedMessage`s from `processStoringChannel` and store it into database
func (telegramService *TelegramService) StoreMessages(storingMessage *model.UpdatedMessage) {
	// Print received message to CLI
	log.Printf("Message received: %+v\n", storingMessage)

	chat, err := telegramService.chatService.GetChat(storingMessage.Data.Chat.ID)
	if err != nil {
		telegramService.errorsChannel <- util.Trace(err)
		return
	}

	messageQueue := modelDB.NewMessageQueueConstructor(
		storingMessage.QueueId,
		storingMessage.Data.Text,
		chat.ID,
	)

	if _, err = telegramService.manager.Repository.MessageQueue().Create(messageQueue); err != nil {
		telegramService.errorsChannel <- util.Trace(err)
		return
	}

	if _, err = telegramService.userService.GetUser(storingMessage.Data.Chat.Username, chat.ID); err != nil {
		telegramService.errorsChannel <- util.Trace(err)
		return
	}

	// check if the message is notification
	matchedValue := regexp.MustCompile(`\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}`).FindStringSubmatch(messageQueue.Message)
	if len(matchedValue) > 0 {
		// processing notification
		scheduledFor, err := time.Parse("2006-01-02 15:04:05", matchedValue[0])
		if err != nil {
			telegramService.errorsChannel <- util.Trace(err)
			return
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
			return
		}
	}
}
