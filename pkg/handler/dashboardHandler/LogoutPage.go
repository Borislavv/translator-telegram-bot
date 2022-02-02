package dashboardHandler

import (
	"net/http"

	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDashboard"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/util"
)

// LogoutPage - handler of "/Logout" endpoint
func (dashboard *Dashboard) LogoutPage(w http.ResponseWriter, r *http.Request) {

	type Content struct {
		Title        string
		ErrorMessage string
	}

	content := Content{
		Title: "Authentication",
	}

	authData, err := model.NewAuthCookieData(w, r)
	if err == nil {
		dashboard.authService.Logout(authData)

		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	}

	// Порядок файлов важен, сначала дочки, далее родители
	templates := []string{
		"./pkg/templates/dashboard/pages/auth/content.html.tmpl",
		"./pkg/templates/dashboard/menu/menu.html.tmpl",
		"./pkg/templates/empty.subContent.html.tmpl",
		"./pkg/templates/base.html.tmpl",
	}

	page := modelDashboard.NewPage()
	page.AddMenu(modelDashboard.NewMenu())
	page.AddConent(content)

	util.RenderFromFiles(w, templates, page)
}
