package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelRepository"
)

type NotsCommand struct {
	// deps.
	manager *manager.Manager
	// vals.
	message *model.CommandMessage
}

// NewNotsCommand - constructor of NotsCommand structure.
func NewNotsCommand(manager *manager.Manager, message *model.CommandMessage) *NotsCommand {
	return &NotsCommand{
		manager: manager,
		message: message,
	}
}

// Exec - proxy method which will execute this command and call private method which will do all work.
// 	Private method receive notifications for user and convert it to string representation.
func (command *NotsCommand) Exec() (*model.TelegramResponseMessage, error) {
	message, err := command.exec()
	if err != nil {
		return nil, err
	}
	if message == "" {
		message = "List is empty."
	}

	return model.NewTelegramResponseMessage(
		strconv.Itoa(int(command.message.OriginMessage.Data.Chat.ID)),
		message,
	), nil
}

// exec - method receive notifications for user and convert it to string representation.
func (command *NotsCommand) exec() (string, error) {
	list, err := command.manager.Repository.NotificationQueue().FindNotSentByUsername(
		command.message.OriginMessage.Data.Chat.Username,
		modelRepository.NewWithoutPaginationParameters(),
	)
	if err != nil {
		return "", err
	}

	messages := []string{}
	for _, v := range list {

		var statusStr string
		if v.IsActive {
			statusStr = fmt.Sprintf("active, use /nots_%v_d for disable", v.ID)
		} else {
			statusStr = fmt.Sprintf("disabled, use /nots_%v_a for activate", v.ID)
		}

		messages = append(
			messages,
			fmt.Sprintf(
				"Text: %v\n\t\tStatus: %v\n\t\tSheduled on: %v",
				v.Message,
				statusStr,
				v.ScheduledFor,
			),
		)
	}

	return strings.Join(messages, "\n\n"), nil
}
