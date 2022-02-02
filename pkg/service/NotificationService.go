package service

import (
	"log"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDB"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelInterface"
)

type NotificationService struct {
	// Dependencies
	manager *manager.Manager
}

// NewNotificationService - constructor of NotificationService struct
func NewNotificationService(manager *manager.Manager) *NotificationService {
	return &NotificationService{
		manager: manager,
	}
}

// GetListForUser - method for receive a list of notification by username
func (notificationService *NotificationService) GetListForUser(
	username string,
	pagination modelInterface.PaginationInterface,
) []*modelDB.NotificationQueue {
	list, err := notificationService.manager.Repository.NotificationQueue().FindNotSentByUsername(username, pagination)
	if err != nil {
		log.Println(err.Error())
		return []*modelDB.NotificationQueue{}
	}

	return list
}
