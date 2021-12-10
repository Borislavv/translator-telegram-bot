package telegram

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/config"
	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDB"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/translator"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/util"
)

type TelegramBot struct {
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
}

// NewTelegramBot - creating a new instance of TelegramBot
func NewTelegramBot(
	manager *manager.Manager,
	gateway *TelegramGateway,
	userService *service.UserService,
	chatService *service.ChatService,
	translator *translator.TranslatorService,
	notificationsChannel chan *modelDB.NotificationQueue,
	messagesChannel chan *model.UpdatedMessage,
	errorsChannel chan string,
) *TelegramBot {
	return &TelegramBot{
		manager:              manager,
		gateway:              gateway,
		userService:          userService,
		chatService:          chatService,
		translator:           translator,
		notificationsChannel: notificationsChannel,
		messagesChannel:      messagesChannel,
		errorsChannel:        errorsChannel,
	}
}

// HandlingMessages - main logic of processing received messages
func (bot *TelegramBot) ProcessMessages() {
	updatedMessages, err := bot.getUpdates()
	if err != nil {
		bot.errorsChannel <- util.Trace(err)
		return
	}

	// Do nothing, if no new message have been received
	if len(updatedMessages) > 0 {
		go bot.sendResponseMessages()

		for _, updatedMessage := range updatedMessages {
			bot.messagesChannel <- &updatedMessage
		}
	}
}

// ProcessNotifications - checking notifications on next
func (bot *TelegramBot) ProcessNotifications() {
	dateTime := time.Now()
	if bot.manager.Config.Environment.Mode == config.ProdMode {
		dateTime = dateTime.Add(2 * time.Hour)
	}

	notifications, err := bot.manager.Repository.NotificationQueue().FindByScheduledDate(dateTime)
	if err != nil {
		bot.errorsChannel <- util.Trace(err)
		return
	}

	// Do nothing, if no new notification have been received
	if len(notifications) > 0 {
		go bot.sendNotifications()

		for _, notification := range notifications {
			bot.notificationsChannel <- notification
		}
	}
}

// ProcessErrors - simple output of errors with debug info
func (bot *TelegramBot) ProcessErrors() {
	for errMsg := range bot.errorsChannel {
		log.Println(errMsg)
	}
}

// sendNotifications - sending message to telegram channel in gorutine
func (bot *TelegramBot) sendNotifications() {
	for notification := range bot.notificationsChannel {
		// Print that notification is sent to CLI
		log.Printf("Notification have been sent: %+v\n", notification)

		// TODO: refactor this code, because the notification already has prop. ExternalChatId
		chat, err := bot.chatService.GetChat(notification.ExternalChatId)
		if err != nil {
			bot.errorsChannel <- util.Trace(err)
			continue
		}

		if err := bot.gateway.SendMessage(
			model.NewTelegramResponseMessage(
				fmt.Sprint(chat.ExternalChatId),
				"Notification: "+notification.Message,
			),
		); err != nil {
			bot.errorsChannel <- util.Trace(err)
			continue
		}

		_, err = bot.manager.Repository.NotificationQueue().MakeAsSent(notification)
		if err != nil {
			bot.errorsChannel <- util.Trace(err)
			continue
		}
	}
}

// sendResponseMessages - sending message to telegram chat from TelegramBot.messagesChannel
func (bot *TelegramBot) sendResponseMessages() {
	for message := range bot.messagesChannel {
		// Print received message to CLI
		log.Printf("Message received: %+v\n", message)

		chat, err := bot.chatService.GetChat(message.Data.Chat.ID)
		if err != nil {
			bot.errorsChannel <- util.Trace(err)
			continue
		}

		messageQueue := modelDB.MessageQueue{
			QueueId: message.QueueId,
			Message: message.Data.Text,
			ChatId:  chat.ID,
		}

		if _, err = bot.manager.Repository.MessageQueue().Create(&messageQueue); err != nil {
			bot.errorsChannel <- util.Trace(err)
			continue
		}

		if _, err = bot.userService.GetUser(message.Data.Chat.Username, chat.ID); err != nil {
			bot.errorsChannel <- util.Trace(err)
			continue
		}

		if err = bot.handleMessage(chat, &messageQueue); err != nil {
			bot.errorsChannel <- util.Trace(err)
			continue
		}
	}
}

// handleMessage - handle one message (right now: will send the same message with prefix)
func (bot *TelegramBot) handleMessage(chat *modelDB.Chat, messageQueue *modelDB.MessageQueue) error {
	matchedValue := regexp.MustCompile(`\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}`).FindStringSubmatch(messageQueue.Message)

	var message string
	var err error
	if len(matchedValue) <= 0 {
		// processing translation
		message, err = bot.translator.TranslateText(messageQueue.Message)
		if err != nil {
			return err
		}

		message = "Tranlsation: " + message
	} else {
		// processing notification
		scheduledFor, err := time.Parse("2006-01-02 15:04:05", matchedValue[0])
		if err != nil {
			return err
		}

		notificationQueue := &modelDB.NotificationQueue{
			MessageQueueId: messageQueue.ID,
			ChatId:         messageQueue.ChatId,
			ScheduledFor:   scheduledFor,
		}

		_, err = bot.manager.Repository.NotificationQueue().Create(notificationQueue)
		if err != nil {
			return err
		}

		message = "Notification setted on " + matchedValue[0]
	}

	if err := bot.gateway.SendMessage(
		model.NewTelegramResponseMessage(
			fmt.Sprint(chat.ExternalChatId),
			message,
		),
	); err != nil {
		return err
	}

	return nil
}

// getUpdates - will return a slice of UpdateMessage objects
func (bot *TelegramBot) getUpdates() ([]model.UpdatedMessage, error) {
	offset, err := bot.manager.Repository.MessageQueue().GetOffset()
	if err != nil {
		return nil, err
	}

	return bot.gateway.GetUpdates(model.NewTelegramRequestMessage(offset)).Messages, nil
}
