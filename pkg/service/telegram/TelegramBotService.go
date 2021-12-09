package telegram

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDB"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/translator"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/util"
)

type TelegramBot struct {
	manager     *manager.Manager
	gateway     *TelegramGateway
	userService *service.UserService
	translator  *translator.TranslatorService
}

// NewTelegramBot - creating a new instance of TelegramBot
func NewTelegramBot(
	manager *manager.Manager,
	gateway *TelegramGateway,
	userService *service.UserService,
	translator *translator.TranslatorService,
) *TelegramBot {
	return &TelegramBot{
		manager:     manager,
		gateway:     gateway,
		userService: userService,
		translator:  translator,
	}
}

// HandlingNotification - checking notifications on next
func (bot *TelegramBot) HandlingNotifications(chatsMap map[int64]*modelDB.Chat) {
	notifications, err := bot.manager.Repository.NotificationQueue().FindByScheduledDate(time.Now())
	if err != nil {
		log.Fatalln(util.Trace() + err.Error())
		return
	}

	for _, notification := range notifications {
		// Print that notification is sent to CLI
		log.Printf("Notification have been sent: %+v\n", notification)

		// TODO: refactor this code, because the notification already has prop. ExternalChatId
		chat, err := bot.getChat(notification.ExternalChatId, chatsMap)
		if err != nil {
			log.Fatalln(util.Trace() + err.Error())
			return
		}

		if err := bot.gateway.SendMessage(
			fmt.Sprint(chat.ExternalChatId),
			"Notification: "+notification.Message,
		); err != nil {
			log.Fatalln(util.Trace() + err.Error())
			return
		}

		_, err = bot.manager.Repository.NotificationQueue().MakeAsSent(notification)
		if err != nil {
			log.Fatalln(util.Trace() + err.Error())
			return
		}
	}
}

// HandlingMessages - main logic of processing received messages
func (bot *TelegramBot) HandlingMessages(usersMap map[string]*modelDB.User, chatsMap map[int64]*modelDB.Chat) {
	updatedMessages, err := bot.getUpdates()
	if err != nil {
		log.Fatalln(util.Trace() + err.Error())
		return
	}

	// Do nothing, if no new message have been received
	if len(updatedMessages) == 0 {
		return
	}

	for _, updatedMessage := range updatedMessages {
		// Print received message to CLI
		log.Printf("Message received: %+v\n", updatedMessage)

		chat, err := bot.getChat(updatedMessage.Data.Chat.ID, chatsMap)
		if err != nil {
			log.Fatalln(util.Trace() + err.Error())
			return
		}

		messageQueue := modelDB.MessageQueue{
			QueueId: updatedMessage.QueueId,
			Message: updatedMessage.Data.Text,
			ChatId:  chat.ID,
		}

		if _, err = bot.manager.Repository.MessageQueue().Create(&messageQueue); err != nil {
			log.Fatalln(util.Trace() + err.Error())
			return
		}

		// in the goroutine because the user is not currently being used
		go bot.getUser(updatedMessage.Data.Chat.Username, usersMap, chat)

		bot.handleMessage(chat, &messageQueue)
	}
}

// handleMessage - handle one message (right now: will send the same message with prefix)
func (bot *TelegramBot) handleMessage(chat *modelDB.Chat, messageQueue *modelDB.MessageQueue) {
	matchedValue := regexp.MustCompile(`\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}`).FindStringSubmatch(messageQueue.Message)

	var message string
	var err error
	if len(matchedValue) <= 0 {
		// processing translation
		message, err = bot.translator.TranslateText(messageQueue.Message)
		if err != nil {
			log.Fatalln(util.Trace() + err.Error())
			return
		}

		message = "Tranlsation: " + message
	} else {
		// processing notification
		scheduledFor, err := time.Parse("2006-01-02 15:04:05", matchedValue[0])
		if err != nil {
			log.Fatalln(util.Trace() + err.Error())
			return
		}

		notificationQueue := &modelDB.NotificationQueue{
			MessageQueueId: messageQueue.ID,
			ChatId:         messageQueue.ChatId,
			ScheduledFor:   scheduledFor,
		}

		_, err = bot.manager.Repository.NotificationQueue().Create(notificationQueue)
		if err != nil {
			log.Fatalln(util.Trace() + err.Error())
			return
		}

		message = "Notification setted on " + matchedValue[0]
	}

	if err := bot.gateway.SendMessage(
		fmt.Sprint(chat.ExternalChatId),
		message,
	); err != nil {
		log.Fatalln(util.Trace() + err.Error())
		return
	}
}

// getUpdates - will return a slice of UpdateMessage objects
func (bot *TelegramBot) getUpdates() ([]model.UpdatedMessage, error) {
	offset, err := bot.manager.Repository.MessageQueue().GetOffset()

	if err != nil {
		log.Fatalln(util.Trace() + err.Error())
		return nil, err
	}

	return bot.gateway.GetUpdates(offset).Messages, nil
}

// getUser - getting of target user by username, if not found, then will created
func (bot *TelegramBot) getUser(username string, usersMap map[string]*modelDB.User, chat *modelDB.Chat) (*modelDB.User, error) {
	// trying to find user in the cache
	if _, issetInCache := usersMap[username]; !issetInCache {
		// trying to find user into database
		dbUser, err := bot.manager.Repository.User().FindByUsername(username)
		if err != nil {
			log.Fatalln(util.Trace() + err.Error())
			return nil, err
		} else {
			// user was not fonud, then create and store it
			if dbUser.ID <= 0 {
				newUser := modelDB.NewUser()
				newUser.ChatId = chat.ID
				newUser.Username = username

				// store user
				dbUser, err = bot.manager.Repository.User().Create(newUser)
				if err != nil {
					log.Fatalln(util.Trace() + err.Error())
					return nil, err
				}
			}
		}

		// store user to cache
		usersMap[username] = dbUser
	}

	return usersMap[username], nil
}

// getChat - getting of target chat by externalChatId, if not found, then will created
func (bot *TelegramBot) getChat(externalChatId int64, chatsMap map[int64]*modelDB.Chat) (*modelDB.Chat, error) {
	// trying to find chat in the cache
	if _, issetInCache := chatsMap[externalChatId]; !issetInCache {
		// trying to find chat into database
		dbChat, err := bot.manager.Repository.Chat().FindByExternalChatId(externalChatId)
		if err != nil {
			return nil, err
		} else {
			// chat was not fonud, then create and store it
			if dbChat.ID <= 0 {
				newChat := modelDB.NewChat()
				newChat.ExternalChatId = externalChatId

				// store chat
				dbChat, err = bot.manager.Repository.Chat().Create(newChat)
				if err != nil {
					return nil, err
				}
			}
		}

		// store chat to cache
		chatsMap[dbChat.ExternalChatId] = dbChat
	}

	return chatsMap[externalChatId], nil
}
