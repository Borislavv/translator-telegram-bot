package handler

import (
	"fmt"
	"net/http"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/handler/dashboardHandler"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/dashboardService"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/translator"
)

var (
	defaultDir  = "./web/"
	defaultPort = ":8060"
)

type Handler struct {
	manager             *manager.Manager
	dashboard           *dashboardHandler.Dashboard
	authService         *dashboardService.AuthService
	notificationService *service.NotificationService
	translatorService   *translator.TranslatorService
}

// NewHandler - constructor of Handler struct
func NewHandler(
	manager *manager.Manager,
	authService *dashboardService.AuthService,
	notificationService *service.NotificationService,
	translatorService *translator.TranslatorService,
) *Handler {
	return &Handler{
		manager: manager,
		dashboard: dashboardHandler.NewDashboard(
			manager,
			authService,
			notificationService,
			translatorService,
		),
	}
}

// HandleDashboard - handle all pages of Dashboard
func (handler *Handler) HandleDashboard() {
	// pages
	handler.dashboard.HandleTheIndexPage()
	handler.dashboard.HandleTheNotificationsPage()
	handler.dashboard.HandleTheTranslationPage()
	handler.dashboard.HandleTheDocsPage()
	handler.dashboard.HandleTheAboutPage()
	handler.dashboard.HandleTheLoginPage()
	handler.dashboard.HandleTheLogoutPage()
	// api
	handler.dashboard.HandleTheTranslationAPI()
}

// HandleStaticFiles - will serve static files in the passed dir.
func (handler *Handler) HandleStaticFiles() {
	dir := handler.manager.Config.Server.StaticFilesDir
	if dir == "" {
		dir = defaultDir
	}

	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir(dir))),
	)
}

// ListenAndServePages - starting the HTTP server.
func (handler *Handler) ListenAndServe() {
	pattern := "%v:%v"

	host := handler.manager.Config.Server.Host
	port := handler.manager.Config.Server.Port

	if port == "" {
		port = defaultPort
	}

	http.ListenAndServe(fmt.Sprintf(pattern, host, port), nil)
}
