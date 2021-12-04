package service

import (
	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
)

type UserService struct {
	manager *manager.Manager
}

// NewUserService - constructor of UserService
func NewUserService(manager *manager.Manager) *UserService {
	return &UserService{
		manager: manager,
	}
}
