package commands

import (
	"strconv"

	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
)

type StartCommand struct {
	// vals.
	message *model.CommandMessage
}

// NewStartCommand - constructor of StartCommand struct.
func NewStartCommand(message *model.CommandMessage) *StartCommand {
	return &StartCommand{
		message: message,
	}
}

// Exec - will tell a user that he need setting up the local date and time.
func (command *StartCommand) Exec() (*model.TelegramResponseMessage, error) {
	return model.NewTelegramResponseMessage(
		strconv.Itoa(int(command.message.OriginMessage.Data.Chat.ID)),
		"In the first, you need set you local date and time for use notifications correctly.\n"+
			"Please, use `/set_tz {local date and time must be here}` for set up it (example: /set_tz 2022-03-11 14:27:38).",
	), nil
}
