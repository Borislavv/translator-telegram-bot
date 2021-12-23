package dashboard

import (
	"net/http"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
)

// Facade of Dashboard pages
type Dashboard struct {
	manager *manager.Manager
}

// NewDashboard - constructor of Dashboard struct
func NewDashboard(manager *manager.Manager) *Dashboard {
	return &Dashboard{
		manager: manager,
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
