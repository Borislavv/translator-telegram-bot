package handler

import (
	"fmt"
	"net/http"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/handler/dashboardHandler"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/dashboardService"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/loggerService"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/translator"
)

var (
	defaultDir  = "./web/"
	defaultPort = ":8000"
)

type Handler struct {
	// deps
	manager   *manager.Manager
	dashboard *dashboardHandler.Dashboard
	logger    *loggerService.LoggerService
}

// NewHandler - constructor of Handler struct
func NewHandler(
	manager *manager.Manager,
	authService *dashboardService.AuthService,
	notificationService *service.NotificationService,
	translatorService *translator.TranslatorService,
	loggerService *loggerService.LoggerService,
) *Handler {
	return &Handler{
		manager: manager,
		dashboard: dashboardHandler.NewDashboard(
			manager,
			authService,
			notificationService,
			translatorService,
		),
		logger: loggerService,
	}
}

// HandleDashboard - handle all pages of Dashboard
func (handler *Handler) HandleDashboard() *Handler {
	// pages
	handler.dashboard.HandleTheIndexPage()
	handler.dashboard.HandleTheNotificationsPage()
	handler.dashboard.HandleTheTranslationPage()
	handler.dashboard.HandleTheDocsPage()
	handler.dashboard.HandleTheAboutPage()
	handler.dashboard.HandleTheLoginPage()
	handler.dashboard.HandleTheLogoutPage()

	// api
	handler.dashboard.HandleTheTranslateAPIMethod()
	handler.dashboard.HandleTheEnableNotificationAPIMethod()
	handler.dashboard.HandleTheDisableNotificationAPIMethod()

	return handler
}

// HandleStaticFiles - will serve static files in the passed dir.
func (handler *Handler) HandleStaticFiles() *Handler {
	dir := handler.manager.Config.Server.StaticFilesDir
	if dir == "" {
		dir = defaultDir
	}

	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir(dir))),
	)

	return handler
}

// ListenAndServePages - starting the HTTP server.
func (handler *Handler) ListenAndServe() {
	pattern := "%v:%v"

	host := handler.manager.Config.Server.Host
	port := handler.manager.Config.Server.Port

	if port == "" {
		port = defaultPort
	}

	if err := http.ListenAndServe(fmt.Sprintf(pattern, host, port), nil); err != nil {
		handler.logger.Critical(err.Error())
		return
	}
}
