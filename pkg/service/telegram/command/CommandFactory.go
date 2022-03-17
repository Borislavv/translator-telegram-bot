package command

import (
	"errors"
	"regexp"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/telegram/command/commands"
)

type CommandInterface interface {
	Exec() (*model.TelegramResponseMessage, error)
}

type CommandFactoryInterface interface {
	NewCommand(command *model.CommandMessage) (CommandInterface, error)
}

type CommandFactory struct {
	// deps.
	manager     *manager.Manager
	userService *service.UserService
}

// NewCommandFactory - constructor of CommandFactory struct.
func NewCommandFactory(
	manager *manager.Manager,
	userService *service.UserService,
) CommandFactoryInterface {
	return &CommandFactory{
		manager:     manager,
		userService: userService,
	}
}

// NewCommand - proxy factory method, which will create a new instance with a CommandInterface compotible type.
func (factory *CommandFactory) NewCommand(command *model.CommandMessage) (CommandInterface, error) {
	return factory.buildCommand(command)
}

// buildCommand - factory method, which will create a new instance with a CommandInterface compotible type.
func (factory *CommandFactory) buildCommand(message *model.CommandMessage) (CommandInterface, error) {
	if len(regexp.MustCompile(`\/\w+`).FindStringSubmatch(message.OriginMessage.Data.Text)) > 0 {
		// handle the /nots_{id}_d|/nots_{id}_a command (enable and disable notifications)
		if len(regexp.MustCompile(`\/nots_\d+_`).FindStringSubmatch(message.OriginMessage.Data.Text)) > 0 {
			return commands.NewNotsActivityCommand(factory.manager, message), nil

			// handle the /nots command (get list of notifications)
		} else if len(regexp.MustCompile(`\/nots`).FindStringSubmatch(message.OriginMessage.Data.Text)) > 0 {
			return commands.NewNotsCommand(factory.manager, message), nil

			// handle the /help command (get list of commands)
		} else if len(regexp.MustCompile(`\/help`).FindStringSubmatch(message.OriginMessage.Data.Text)) > 0 ||
			len(regexp.MustCompile(`\/h`).FindStringSubmatch(message.OriginMessage.Data.Text)) > 0 {
			return commands.NewHelpCommand(message), nil
		} else if len(regexp.MustCompile(`\/start`).FindStringSubmatch(message.OriginMessage.Data.Text)) > 0 {
			return commands.NewStartCommand(message), nil
		} else if len(regexp.MustCompile(`\/set_tz`).FindStringSubmatch(message.OriginMessage.Data.Text)) > 0 {
			return commands.NewTimeZoneCommand(factory.manager, message, factory.userService), nil
		}
	}

	return nil, errors.New("Unknown command. Please use /help for get more information.")
}
