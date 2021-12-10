package manager

import (
	"log"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/config"
	"github.com/Borislavv/Translator-telegram-bot/pkg/repository"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/util"
)

type Manager struct {
	Config     *config.Config
	Repository *repository.Repository
}

// New - creating a new instance of Manager
func New(config *config.Config) *Manager {
	manager := &Manager{
		Config: config,
	}

	if err := manager.configureRepository(); err != nil {
		log.Fatalln(util.Trace(err))
		return nil
	}

	return manager
}

// configureRepository - open connection with database
func (manager *Manager) configureRepository() error {
	repository := repository.New(manager.Config.Repository)
	if err := repository.Open(); err != nil {
		return err
	}

	manager.Repository = repository

	return nil
}
