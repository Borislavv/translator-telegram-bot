package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/BurntSushi/toml"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/config"
	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDB"
	"github.com/Borislavv/Translator-telegram-bot/pkg/repository"
)

var (
	environmentMode   string
	configurationPath string
)

func main() {
	askFlags()

	// gateway := service.NewTelegramGateway()

	// updates, err := gateway.GetUpdates(0)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// for _, message := range updates.Messages {
	// 	fmt.Printf("%+v\n", message)
	// }

	////////here
	fmt.Println(configurationPath)

	config := loadConfig()
	manager := manager.New(config)

	// fmt.Printf("%s\n", config.Repository.DSN)

	chat := modelDB.NewChat()
	chatId, _ := strconv.Atoi("-1001728386516")
	chat.ExternalChatId = int64(chatId)

	chat, err := manager.Repository.Chat().Create(chat)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%+v\n", chat)

	foundChat, err := manager.Repository.ChatRepository.FindByExternalChatId("-1001728386516")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%+v\n", foundChat)
}

// askFlags - getting args. from cli
func askFlags() {
	flag.StringVar(&environmentMode, "env-mode", config.ProdMode, "one of env. modes: prod|dev")
	flag.StringVar(&configurationPath, "config-path", config.DefaultConfigPath, "path to config file, by default: config/.env.prod.toml")
	flag.Parse()
}

// loadConfig - loading a config file to struct
func loadConfig() *config.Config {
	repositoryConfig := repository.NewRepositoryConfig()

	_, err := toml.DecodeFile(configurationPath, repositoryConfig)
	if err != nil {
		log.Fatalln(err)
	}

	config := config.New()
	config.Repository = repositoryConfig

	return config
}
