package telegram

import (
	"log"
	"sync"
	"time"
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
	var wg sync.WaitGroup

	for {
		if notifications := bot.telegramService.GetNotifications(); notifications != nil {
			wg.Add(1)

			go bot.telegramService.SendNotifications(notifications, &wg)
		}

		wg.Wait()
	}
}

// ProcessMessages - get new messages from TelegramAPI, store it and answer
func (bot *TelegramBot) ProcessMessages() {
	var wg sync.WaitGroup

	for {
		if messages := bot.telegramService.GetMessages(); messages != nil {
			wg.Add(2)

			go bot.telegramService.SendMessages(messages, &wg)
			go bot.telegramService.StoreMessages(messages, &wg)
		}

		wg.Wait()
	}
}

// ProcessErrors - (gorutine) simple output of errors with debug info
func (bot *TelegramBot) ProcessErrors() {
	for {
		select {
		case message := <-bot.telegramService.errorsChannel:
			log.Println(message)
		default:
			time.Sleep(15 * time.Millisecond)
		}
	}
}
