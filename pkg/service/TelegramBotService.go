package service

import (
	"fmt"
	"log"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDB"
)

type TelegramBot struct {
	manager     *manager.Manager
	gateway     *TelegramGateway
	userService *UserService
}

// NewTelegramBot - creating a new instance of TelegramBot
func NewTelegramBot(
	manager *manager.Manager,
	gateway *TelegramGateway,
	userService *UserService,
) *TelegramBot {
	return &TelegramBot{
		manager:     manager,
		gateway:     gateway,
		userService: userService,
	}
}

// HandlingMessages - main logic of processing received messages
func (bot *TelegramBot) HandlingMessages() {
	usersMap := make(map[string]*modelDB.User)
	chatsMap := make(map[int64]*modelDB.Chat)

	updatedMessages, err := bot.getUpdates()
	if err != nil {
		log.Fatalln(err)
		return
	}

	for _, updatedMessage := range updatedMessages {
		fmt.Printf("%+v\n", updatedMessage)

		chat, err := bot.getChat(updatedMessage.Data.Chat.ID, chatsMap)
		if err != nil {
			log.Fatalln(err)
			return
		}

		messageQueue := modelDB.MessageQueue{
			QueueId: updatedMessage.QueueId,
			Message: updatedMessage.Data.Text,
			ChatId:  chat.ID,
		}

		if _, err = bot.manager.Repository.MessageQueue().Create(&messageQueue); err != nil {
			log.Fatalln(err)
			return
		}

		// in gorutine, because right now, this method does not need an instance of User, further will
		bot.getUser(updatedMessage.Data.Chat.Username, usersMap, chat)

		bot.handleMessage(chat, &messageQueue)
	}
}

// handleMessage - handle one message (right now: will send the same message with prefix)
func (bot *TelegramBot) handleMessage(chat *modelDB.Chat, messageQueue *modelDB.MessageQueue) {
	if err := bot.gateway.SendMessage(
		fmt.Sprint(chat.ExternalChatId), "Bot answered you: "+messageQueue.Message,
	); err != nil {
		log.Fatalln(err)
		return
	}
}

// getUpdates - will return a slice of UpdateMessage objects
func (bot *TelegramBot) getUpdates() ([]model.UpdatedMessage, error) {
	offset, err := bot.manager.Repository.MessageQueue().GetOffset()
	if err != nil {
		log.Println(err.Error() + "HEREERERE")
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
			log.Fatalln(err)
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
					log.Fatalln(err)
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
