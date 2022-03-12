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
		// Getting notification in a new thread
		go bot.telegramService.GetNotifications(&wg)

		wg.Add(1)
		// Sending notification in a new thread
		go bot.telegramService.SendNotifications(&wg)
		wg.Wait()

		time.Sleep(5 * time.Second)
	}
}

// ProcessMessages - get new messages from TelegramAPI, store it and answer
func (bot *TelegramBot) ProcessMessages() {
	// Getting messages in a new thread
	go bot.telegramService.GetMessages()

	// Sending messages in a new thread
	go bot.telegramService.SendMessages()

	// Saving messages in a new thread
	go bot.telegramService.StoreMessages()

	// Saving tokens in a new thread
	go bot.telegramService.StoreTokens()

	// Client side commands processing in a new thread
	go bot.telegramService.ProcessCommands()
}

// ProcessErrors - (gorutine) simple output of errors with debug info
func (bot *TelegramBot) ProcessErrors() {
	go func() {
		for {
			select {
			case message := <-bot.telegramService.errorsChannel:
				log.Println(message)
			default:
				time.Sleep(15 * time.Millisecond)
			}
		}
	}()
}
