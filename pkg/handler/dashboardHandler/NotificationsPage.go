package dashboardHandler

import (
	"net/http"
	"strconv"

	"github.com/Borislavv/Translator-telegram-bot/pkg/helper"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDB"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDashboard"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelRepository"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/util"
)

// Index - handler of "/notifications" endpoint
func (dashboard *Dashboard) NotificationsPage(w http.ResponseWriter, r *http.Request) {

	page := modelDashboard.NewPage()
	page.AddMenu(modelDashboard.NewMenu())

	type Content struct {
		Title            string
		Notifications    []*modelDB.NotificationQueue
		Funcs            *helper.TemplateFuncsHelper
		ErrorMessage     string
		LoginByTokenLink string
	}

	content := Content{
		Title:            "Authentication",
		LoginByTokenLink: "/login",
		Funcs:            helper.NewTemplateFuncsHelper(),
	}

	// checking the user have access to this section
	authData, err := model.NewAuthCookieData(w, r)
	if err != nil || !dashboard.authService.IsAuthorized(authData) {
		content.ErrorMessage = "Please, visit this page for auth"
	} else {
		cachedUser, err := dashboard.authService.GetUserFromCache(authData)
		if err == nil {
			page.AddUser(cachedUser)
		}

		pageNum, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			pageNum = 0
		}
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			limit = 0
		}

		// receiving list of notifications
		content.Notifications = dashboard.notificationService.GetListForUser(
			authData.GetUsername(),
			modelRepository.NewPaginationParameters(
				pageNum,
				limit,
			),
		)
	}

	// Порядок файлов важен, сначала дочки, далее родители
	templates := []string{
		"./pkg/templates/dashboard/pages/notifications/content.html.tmpl",
		"./pkg/templates/dashboard/menu/menu.html.tmpl",
		"./pkg/templates/empty.subContent.html.tmpl",
		"./pkg/templates/base.html.tmpl",
	}

	page.AddConent(content)
	page.AddTemplateFuncs(helper.NewTemplateFuncsHelper())

	util.RenderFromFiles(w, templates, page)
}
