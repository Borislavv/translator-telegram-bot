package dashboardHandler

import (
	"net/http"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/dashboardService"
)

// Facade of Dashboard pages
type Dashboard struct {
	manager             *manager.Manager
	authService         *dashboardService.AuthService
	notificationService *service.NotificationService
}

// NewDashboard - constructor of Dashboard struct
func NewDashboard(
	manager *manager.Manager,
	authService *dashboardService.AuthService,
	notificationService *service.NotificationService,
) *Dashboard {
	return &Dashboard{
		manager:             manager,
		authService:         authService,
		notificationService: notificationService,
	}
}

// HandleTheIndexPage - handle the `index` page as controller method
func (dashboard *Dashboard) HandleTheIndexPage() {
	http.HandleFunc("/", dashboard.IndexPage)
}

// HandleTheNotificationsPage - handle the `notifications` page as controller method
func (dashboard *Dashboard) HandleTheNotificationsPage() {
	http.HandleFunc("/notifications", dashboard.NotificationsPage)
}

// HandleTheTranslationPage - handle the `translation` page as controller method
func (dashboard *Dashboard) HandleTheTranslationPage() {
	http.HandleFunc("/translation", dashboard.TranslationPage)
}

// HandleTheDocsPage - handle the `docs` page as controller method
func (dashboard *Dashboard) HandleTheDocsPage() {
	http.HandleFunc("/docs", dashboard.DocsPage)
}

// HandleTheAboutPage - handle the `about` page as controller method
func (dashboard *Dashboard) HandleTheAboutPage() {
	http.HandleFunc("/about", dashboard.AboutPage)
}

// HandleTheLoginPage - handle the `login` page as controller method
func (dashboard *Dashboard) HandleTheLoginPage() {
	http.HandleFunc("/login", dashboard.LoginPage)
}

// HandleTheLogoutPage - handle the `logout` page as controller method
func (dashboard *Dashboard) HandleTheLogoutPage() {
	http.HandleFunc("/logout", dashboard.LogoutPage)
}
