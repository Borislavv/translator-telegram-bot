package commands

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDB"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service"
)

type TimeZoneCommand struct {
	// deps.
	manager     *manager.Manager
	userService *service.UserService
	// vals.
	message *model.CommandMessage
}

// NewTimeZoneCommand - constructor of TimeZoneCommand struct.
func NewTimeZoneCommand(
	manager *manager.Manager,
	message *model.CommandMessage,
	userService *service.UserService,
) *TimeZoneCommand {
	return &TimeZoneCommand{
		manager:     manager,
		message:     message,
		userService: userService,
	}
}

// Exec - will tell a user that he need setting up the local date and time.
func (command *TimeZoneCommand) Exec() (*model.TelegramResponseMessage, error) {
	matches := regexp.MustCompile(`\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}`).FindStringSubmatch(command.message.OriginMessage.Data.Text)
	if len(matches) == 0 || matches[0] == "" {
		return nil, errors.New(
			"Date and time not found in the message. " +
				"Please, check your format and try again (example: '/set_tz 2000-01-29 15:17:00').",
		)
	}

	tz, err := command.userService.GetUserTimeZone(matches[0])
	if err != nil {
		return nil, err
	}

	user := modelDB.NewUser()
	user.TZ = tz
	user.ChatId = command.message.OriginMessage.Data.Chat.ID

	_, err = command.manager.Repository.User().SetTimeZoneByExtChatId(user)
	if err != nil {
		return nil, err
	}

	return model.NewTelegramResponseMessage(
		strconv.Itoa(int(command.message.OriginMessage.Data.Chat.ID)),
		"Your timezone succefully configured. Thanks for spent time!",
	), nil
}
