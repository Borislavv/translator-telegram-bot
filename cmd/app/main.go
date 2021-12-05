package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/config"
	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
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

	// Creating an instance of Config at first
	config := loadConfig()

	// Creating an instance of Manager (contains repos and config)
	manager := manager.New(config)

	// Creating an instance of TelegramGatewayService
	telegramGateway := telegram.NewTelegramGateway(manager)

	// Creating an instance of UserService
	userService := service.NewUserService(manager)

	// Creating an instance of TranslatorGateway
	translatorGateway := translator.NewTranslatorGateway(manager)

	// Creating an instance of TranslatorService
	translator := translator.NewTranslatorService(translatorGateway)

	// Creating an instance of TelegramBotService
	bot := telegram.NewTelegramBot(manager, telegramGateway, userService, translator)

	fmt.Println("Handling messages ...")

	// Close connection with database in defer
	defer manager.Repository.Close()

	usersCacheMap := make(map[string]*modelDB.User)
	chatsCacheMap := make(map[int64]*modelDB.Chat)

	for {
		// Handle batch of UpdatedMessages
		bot.HandlingMessages(usersCacheMap, chatsCacheMap)

		// Timeout before new request
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
		log.Fatalln(util.Trace() + err.Error())
	}

	// Telegram api config loading
	_, err = toml.DecodeFile(configurationPath, &config.Integration.Telegram)
	if err != nil {
		log.Fatalln(util.Trace() + err.Error())
	}

	// Translator config loading
	_, err = toml.DecodeFile(configurationPath, &config.Integration.Translator)
	if err != nil {
		log.Fatalln(util.Trace() + err.Error())
	}

	return config
}
