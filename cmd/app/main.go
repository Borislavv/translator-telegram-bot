package main

import (
	"flag"
	"log"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/config"
	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service"
)

var (
	environmentMode   string
	configurationPath string
)

func main() {
	askFlags()

	config := loadConfig()
	manager := manager.New(config)
	gateway := service.NewTelegramGateway(manager)
	userService := service.NewUserService(manager)
	bot := service.NewTelegramBot(manager, gateway, userService)

	defer manager.Repository.Close()

	for {
		bot.HandlingMessages()

		time.Sleep(1 * time.Second)
	}

	// //!!! Common
	// config := loadConfig()
	// manager := manager.New(config)

	// //!!! GetUpdates

	// gateway := service.NewTelegramGateway(manager)

	// updates := gateway.GetUpdates(0)

	// for _, message := range updates.Messages {
	// 	fmt.Printf("%+v\n", message)
	// }

	// //!!! SendMessage

	// gateway.SendMessage("-1001728386516", "Hello from go code!!!!")

	// //!!! Save to database

	// chat := modelDB.NewChat()
	// chatId, _ := strconv.Atoi("-1001728386516")
	// chat.ExternalChatId = int64(chatId)

	// chat, errN := manager.Repository.Chat().Create(chat)
	// if errN != nil {
	// 	log.Fatalln(errN)
	// }

	// fmt.Printf("%+v\n", chat)

	// //!!! Select one row from database

	// foundChat, err := manager.Repository.Chat().FindByExternalChatId("-1001728386516")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Printf("%+v\n", foundChat)
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

	_, err := toml.DecodeFile(configurationPath, config.Repository)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = toml.DecodeFile(configurationPath, config.Integration)
	if err != nil {
		log.Fatalln(err)
	}

	return config
}
