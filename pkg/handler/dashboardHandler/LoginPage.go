package dashboardHandler

import (
	"net/http"

	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDashboard"
)

// Index - handler of "/login" endpoint
func (dashboard *Dashboard) LoginPage(w http.ResponseWriter, r *http.Request) {

	// Main structure for passing to templates
	page := modelDashboard.NewPage()
	page.AddMenu(modelDashboard.NewMenu())

	type Content struct {
		Title        string
		ErrorMessage string
	}

	content := Content{
		Title: "Authentication",
	}

	authData, err := model.NewAuthRequestData(w, r)
	if err == nil && authData != nil {
		user, err := dashboard.authService.Login(authData)

		if err != nil {
			content.ErrorMessage = err.Error()
		} else {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}

		page.AddUser(user)
	}

	// Порядок файлов важен, сначала дочки, далее родители
	templates := []string{
		"./pkg/templates/dashboard/pages/auth/content.html.tmpl",
		"./pkg/templates/dashboard/menu/menu.html.tmpl",
		"./pkg/templates/empty.subContent.html.tmpl",
		"./pkg/templates/base.html.tmpl",
	}

	page.AddConent(content)

	dashboard.renderingService.RenderFromFiles(w, templates, page)
}
