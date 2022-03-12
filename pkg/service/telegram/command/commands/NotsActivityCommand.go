package commands

import (
	"errors"
	"regexp"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDB"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/util"
)

type NotsActivityCommand struct {
	//deps.
	manager *manager.Manager
	// vals.
	message *model.CommandMessage
}

// NewNotsActivityCommand - constructor of NotsActivityCommand structure.
func NewNotsActivityCommand(manager *manager.Manager, message *model.CommandMessage) *NotsActivityCommand {
	return &NotsActivityCommand{
		manager: manager,
		message: message,
	}
}

// Exec - enable and disable notifications.
func (command *NotsActivityCommand) Exec() (*model.TelegramResponseMessage, error) {
	var matches []string
	var id int64
	var err error

	matches = regexp.MustCompile(`\/nots_\d+`).FindStringSubmatch(command.message.OriginMessage.Data.Text)
	if len(matches) > 0 {
		id, err = util.ExtractInt64(matches[0])
		if err != nil {
			return nil, err
		}
	}

	notDto := modelDB.NewNotificationQueue()
	notDto.ID = id

	if len(regexp.MustCompile(`\/nots_\d+_a`).FindStringSubmatch(command.message.OriginMessage.Data.Text)) > 0 {
		notDto, err = command.manager.Repository.NotificationQueue().MakeAsEnabled(notDto)
		if err != nil {
			return nil, err
		}

		if notDto.IsActive == false {
			return nil, errors.New("Unable activate the target notification.")
		} else {
			return NewNotsCommand(command.manager, command.message).Exec()
		}
	} else if len(regexp.MustCompile(`\/nots_\d+_d`).FindStringSubmatch(command.message.OriginMessage.Data.Text)) > 0 {
		notDto, err = command.manager.Repository.NotificationQueue().MakeAsDisabled(notDto)
		if err != nil {
			return nil, err
		}

		if notDto.IsActive == true {
			return nil, errors.New("Unable disable the target notification.")
		} else {
			return NewNotsCommand(command.manager, command.message).Exec()
		}
	}

	return nil, errors.New("Undefined behavior, check regex in the previous method.")
}
