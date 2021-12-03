package main

import (
	"fmt"
	"log"

	"github.com/Borislavv/Translator-telegram-bot/pkg/service"
)

func main() {
	gateway := service.NewTelegramGateway()

	updates, err := gateway.GetUpdates(0)
	if err != nil {
		log.Fatalln(err)
	}

	for _, message := range updates.Messages {
		fmt.Printf("%+v\n", message)
	}
}
