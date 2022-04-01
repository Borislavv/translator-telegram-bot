package main

import (
	"fmt"
	"sync"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/config"
	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/handler"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDB"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/dashboardService"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/loggerService"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/telegram"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/telegram/command"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/translator"
)

func main() {
	fmt.Println("Initialization")

	// Init. requirements.
	var mx sync.Mutex

	// Init. channels for communication with gorutines.
	messagesChannel := make(chan *model.UpdatedMessage, 128)
	notificationsChannel := make(chan *modelDB.NotificationQueue, 128)
	storeChannel := make(chan *model.UpdatedMessage, 128)
	tokensChannel := make(chan *model.TokenMessage, 128)
	commandsChannel := make(chan *model.CommandMessage, 128)
	errorsChannel := make(chan string, 756)

	// Creating an instance of Config at first and load it.
	config := config.New().Load()

	// Creating an instance of Manager (contains repos and config).
	manager := manager.New(config)

	// Creating an instance of LoggerService.
	loggerService := loggerService.NewLoggerService(manager)

	// Creating an instance of TelegramGatewayService.
	telegramGateway := telegram.NewTelegramGateway(manager)

	// Creating an instance of TimeZoneFetcherService.
	timeZonesFetcher := service.NewTimeZoneFetcherService(&mx)

	// Creating an instance of UserService.
	userService := service.NewUserService(manager, timeZonesFetcher)

	// Creating an instance of ChatService.
	chatService := service.NewChatService(manager)

	// Creating an instance of NotificationService.
	notificationService := service.NewNotificationService(manager)

	// Creating an instance of TranslatorGateway.
	translatorGateway := translator.NewTranslatorGateway(manager)

	// Creating an instance of TranslatorService.
	translator := translator.NewTranslatorService(translatorGateway)

	// Creating an instance of AuthService.
	auth := dashboardService.NewAuthService(manager, userService)

	// Creating an instance of TokenGenerator.
	tokenGenerator := dashboardService.NewTokenGenerator()

	// Creating an instance of CommandFactory.
	commandFactory := command.NewCommandFactory(manager, userService)

	// Creating an instance of CommandService.
	commandsService := command.NewCommandService(commandFactory)

	// Creating an instace of TelegramService.
	telegramService := telegram.NewTelegramService(
		manager,
		telegramGateway,
		userService,
		chatService,
		translator,
		tokenGenerator,
		commandsService,
		loggerService,
		messagesChannel,
		notificationsChannel,
		storeChannel,
		tokensChannel,
		commandsChannel,
		errorsChannel,
	)

	// Creating an instance of TelegramBotService.
	bot := telegram.NewTelegramBot(telegramService, errorsChannel)

	// Close connection with database in defer.
	defer manager.Repository.Close()

	loggerService.Info("Success initialization. Handling messages ...")

	// gRPC server.
	go runGRPCServer(manager, translator, loggerService, userService)

	// HTTP server which handle Dashboard.
	go runHTTPServer(manager, auth, notificationService, translator, loggerService)

	bot.ProcessMessages()
	bot.ProcessNotifications()
	bot.ProcessErrors()
}

// runGRPCServer - handling 8017 port (tcp) by default (! infinite loop: must be running in separate thread).
func runGRPCServer(
	manager *manager.Manager,
	translator *translator.TranslatorService,
	logger *loggerService.LoggerService,
	userService *service.UserService,
) {
	handler.
		NewHandlerGRPC(manager, translator, logger, userService).
		ListenAndServe()
}

// runHTTPServer - handling 8000 port by default (! infinite loop: must be running in separate thread).
func runHTTPServer(
	manager *manager.Manager,
	auth *dashboardService.AuthService,
	notificationService *service.NotificationService,
	translator *translator.TranslatorService,
	logger *loggerService.LoggerService,
) {
	handler.
		NewHandler(manager, auth, notificationService, translator, logger).
		HandleDashboard().
		HandleStaticFiles().
		ListenAndServe()
}
