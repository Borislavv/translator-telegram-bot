package command

import (
	"log"
	"strconv"

	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/util"
)

type CommandService struct {
	factory CommandFactoryInterface
}

// NewCommandService - constructor of CommandService structure.
func NewCommandService(
	factory CommandFactoryInterface,
) *CommandService {
	return &CommandService{
		factory: factory,
	}
}

// ProcessCommand - execute the target command which received from client message.
func (service *CommandService) ProcessCommand(message *model.CommandMessage) *model.TelegramResponseMessage {
	command, err := service.factory.NewCommand(message)
	if err != nil {
		return model.NewTelegramResponseMessage(
			strconv.Itoa(int(message.OriginMessage.Data.Chat.ID)),
			err.Error(),
		)
	}

	response, err := command.Exec()
	if err != nil {
		log.Println(util.Trace(err))

		return model.NewTelegramResponseMessage(
			strconv.Itoa(int(message.OriginMessage.Data.Chat.ID)),
			"Sorry, now we can't process this command.",
		)
	}

	return response
}
