package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/config"
	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDB"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/telegram"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/translator"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/util"
)

// Config vars.
var (
	environmentMode   string
	configurationPath string
)

func main() {
	askFlags()

	fmt.Println("Initialization")

	// Init. channels for communication with gorutines
	messagesChannel := make(chan *model.UpdatedMessage, 128)
	notificationsChannel := make(chan *modelDB.NotificationQueue, 128)
	errorsChannel := make(chan string, 256)

	// Creating an instance of Config at first
	config := loadConfig()

	// Creating an instance of Manager (contains repos and config)
	manager := manager.New(config)

	// Creating an instance of TelegramGatewayService
	telegramGateway := telegram.NewTelegramGateway(manager)

	// Creating an instance of UserService
	userService := service.NewUserService(manager)

	// Creating an instance of ChatService
	chatService := service.NewChatService(manager)

	// Creating an instance of TranslatorGateway
	translatorGateway := translator.NewTranslatorGateway(manager)

	// Creating an instance of TranslatorService
	translator := translator.NewTranslatorService(translatorGateway)

	// Creating an instance of TelegramBotService
	bot := telegram.NewTelegramBot(
		manager,
		telegramGateway,
		userService,
		chatService,
		translator,
		notificationsChannel,
		messagesChannel,
		errorsChannel,
	)

	// Close connection with database in defer
	defer manager.Repository.Close()

	fmt.Println("Handling messages ...")

	go func() {
		for {
			bot.ProcessErrors()
		}
	}()

	for {
		bot.ProcessMessages()
		bot.ProcessNotifications()

		// Timeout
		time.Sleep(1 * time.Second)
	}
}

// askFlags - getting args. from cli
func askFlags() {
	flag.StringVar(&environmentMode, "env-mode", config.ProdMode, "one of env. modes: prod|dev")
	flag.StringVar(&configurationPath, "config-path", config.DefaultConfigPath, "path to config file")
	flag.Parse()
}

// loadConfig - loading a config file to struct
func loadConfig() *config.Config {
	config := config.New()

	// Database config loading
	_, err := toml.DecodeFile(configurationPath, config.Repository)
	if err != nil {
		log.Fatalln(util.Trace(err))
	}

	// Telegram api config loading
	_, err = toml.DecodeFile(configurationPath, &config.Integration.Telegram)
	if err != nil {
		log.Fatalln(util.Trace(err))
	}

	// Translator config loading
	_, err = toml.DecodeFile(configurationPath, &config.Integration.Translator)
	if err != nil {
		log.Fatalln(util.Trace(err))
	}

	return config
}
