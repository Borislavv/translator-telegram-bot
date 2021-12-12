package telegram

import (
	"log"
	"sync"
)

type TelegramBot struct {
	// Services
	telegramService *TelegramService

	// Channels
	errorsChannel chan string
}

// NewTelegramBot - creating a new instance of TelegramBot
func NewTelegramBot(
	telegramService *TelegramService,
	errorsChannel chan string,
) *TelegramBot {
	return &TelegramBot{
		telegramService: telegramService,
		errorsChannel:   errorsChannel,
	}
}

// ProcessNotifications - trying to find nots. on next minute and send it into telegram channel
func (bot *TelegramBot) ProcessNotifications() {
	// Notifications provider which will pass messages into `notificationsChannel`
	go bot.telegramService.GetNotifications()

	// Gorutine will started if `notificationsChannel` receive at least one message
	go bot.telegramService.SendNotifications()
}

// ProcessMessages - get new messages from TelegramAPI, store it and answer
func (bot *TelegramBot) ProcessMessages(m *sync.Mutex) {
	// Notifications provider which will pass messages into `notificationsChannel`
	go bot.telegramService.GetMessages(m)

	// Gorutine will started if `storeChannel` receive at least one message
	go bot.telegramService.StoreMessages(m)

	// Gorutine will started if `notificationsChannel` receive at least one message
	go bot.telegramService.SendMessages(m)
}

// ProcessErrors - (gorutine) simple output of errors with debug info
func (bot *TelegramBot) ProcessErrors() {
	for errMsg := range bot.errorsChannel {
		log.Println(errMsg)
	}
}
