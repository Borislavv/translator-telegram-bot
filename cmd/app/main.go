package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/config"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service"
)

var (
	environmentMode string
)

func main() {
	initCore()

	gateway := service.NewTelegramGateway()

	updates, err := gateway.GetUpdates(0)
	if err != nil {
		log.Fatalln(err)
	}

	for _, message := range updates.Messages {
		fmt.Printf("%+v\n", message)
	}
}

// initCore - initialization dependencies of app
func initCore() {
	askFlags()
}

// askFlags - getting args. from cli
func askFlags() {
	flag.StringVar(&environmentMode, "environment-mode", config.ProdMode, "one of env. modes: prod|dev")
}
