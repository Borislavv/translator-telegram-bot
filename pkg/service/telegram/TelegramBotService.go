package telegram

import (
	"log"
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

func (bot *TelegramBot) ProcessMessages() {
	// Notifications provider which will pass messages into `notificationsChannel`
	go bot.telegramService.GetMessages()

	// Gorutine will started if `storeChannel` receive at least one message
	go bot.telegramService.StoreMessages()

	// Gorutine will started if `notificationsChannel` receive at least one message
	go bot.telegramService.SendMessages()
}

// ProcessErrors - (gorutine) simple output of errors with debug info
func (bot *TelegramBot) ProcessErrors() {
	for errMsg := range bot.errorsChannel {
		log.Println(errMsg)
	}
}
