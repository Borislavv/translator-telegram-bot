package dashboardHandler

import (
	"net/http"

	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDashboard"
)

// Index - handler of "/" endpoint
func (dashboard *Dashboard) IndexPage(w http.ResponseWriter, r *http.Request) {

	// Порядок файлов важен, сначала дочки, далее родители
	templates := []string{
		"./pkg/templates/dashboard/pages/index/content.html.tmpl",
		"./pkg/templates/dashboard/menu/menu.html.tmpl",
		"./pkg/templates/contacts.subContent.html.tmpl",
		"./pkg/templates/base.html.tmpl",
	}

	type Content struct {
		Title            string
		TitleLink        string
		LoginByTokenLink string
	}

	page := modelDashboard.NewPage()
	page.AddMenu(modelDashboard.NewMenu())
	page.AddConent(Content{
		Title:            "Translator-telegram-bot",
		TitleLink:        "https://github.com/Borislavv/Translator-telegram-bot",
		LoginByTokenLink: "/login",
	})
	authData, err := model.NewAuthCookieData(w, r)
	if err == nil {
		cachedUser, err := dashboard.authService.GetUserFromCache(authData)
		if err == nil {
			page.AddUser(cachedUser)
		}
	}

	dashboard.renderingService.RenderFromFiles(w, templates, page)
}
