package dashboard

import (
	"net/http"

	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDashboard"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/util"
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
		LoginByTokenLink: "/login/token",
	})

	util.RenderFromFiles(w, templates, page)
}
