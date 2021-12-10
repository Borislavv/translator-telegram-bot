package service

import (
	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDB"
)

type ChatService struct {
	manager *manager.Manager
	Cache   map[int64]*modelDB.Chat
}

// NewChatService - constructor of ChatService
func NewChatService(manager *manager.Manager) *ChatService {
	return &ChatService{
		manager: manager,
		Cache:   map[int64]*modelDB.Chat{},
	}
}

// GetChat - getting of target chat by externalChatId, if not found, then will created
func (chatService *ChatService) GetChat(externalChatId int64) (*modelDB.Chat, error) {
	// trying to find chat in the cache
	if _, issetInCache := chatService.Cache[externalChatId]; !issetInCache {
		// trying to find chat into database
		dbChat, err := chatService.manager.Repository.Chat().FindByExternalChatId(externalChatId)
		if err != nil {
			return nil, err
		} else {
			// chat was not fonud, then create and store it
			if dbChat.ID <= 0 {
				newChat := modelDB.NewChat()
				newChat.ExternalChatId = externalChatId

				// store chat
				dbChat, err = chatService.manager.Repository.Chat().Create(newChat)
				if err != nil {
					return nil, err
				}
			}
		}

		// store chat to cache
		chatService.Cache[dbChat.ExternalChatId] = dbChat
	}

	return chatService.Cache[externalChatId], nil
}
