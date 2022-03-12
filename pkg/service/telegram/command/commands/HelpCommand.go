package commands

import (
	"strconv"
	"strings"

	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
)

type HelpCommand struct {
	// vals.
	message *model.CommandMessage
}

// NewHelpCommand - constructor of HelpCommand structure.
func NewHelpCommand(message *model.CommandMessage) *HelpCommand {
	return &HelpCommand{
		message: message,
	}
}

// Exec - method will return a structure with discribed commands as one string into it.
func (command *HelpCommand) Exec() (*model.TelegramResponseMessage, error) {
	commands := []string{
		"/start - info about how to configure timezone",
		"/token - get a token for authorization in the dashboard",
		"/nots - get a list of your notifications",
	}

	return model.NewTelegramResponseMessage(
		strconv.Itoa(int(command.message.OriginMessage.Data.Chat.ID)),
		strings.Join(commands, "\n"),
	), nil
}
