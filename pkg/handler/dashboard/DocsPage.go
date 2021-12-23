package dashboard

import (
	"net/http"

	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDashboard"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/util"
)

// Index - handler of "/docs" endpoint
func (dashboard *Dashboard) DocsPage(w http.ResponseWriter, r *http.Request) {

	// Порядок файлов важен, сначала дочки, далее родители
	templates := []string{
		"./pkg/templates/dashboard/pages/docs/content.html.tmpl",
		"./pkg/templates/dashboard/menu/menu.html.tmpl",
		"./pkg/templates/base.html.tmpl",
	}

	page := modelDashboard.NewPage()
	page.AddMenu(modelDashboard.NewMenu())

	util.RenderFromFiles(w, templates, page)
}
