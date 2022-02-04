package dashboardHandler

import (
	"net/http"
	"strconv"

	"github.com/Borislavv/Translator-telegram-bot/pkg/helper"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelAPI/dataAPI"
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
		Title                         string
		Notifications                 []*modelDB.NotificationQueue
		Funcs                         *helper.TemplateFuncsHelper
		ErrorMessage                  string
		LoginByTokenLink              string
		DisableNotificationUrlPattern string
		EnableNotificationUrlPattern  string
	}

	content := Content{
		Title:                         "Authentication",
		LoginByTokenLink:              "/login",
		Funcs:                         helper.NewTemplateFuncsHelper(),
		DisableNotificationUrlPattern: "/api/v1/notifications/disable/{id}",
		EnableNotificationUrlPattern:  "/api/v1/notifications/enable/{id}",
	}

	// checking the user have access to this section
	authData, err := model.NewAuthCookieData(w, r)
	if err != nil || !dashboard.authService.IsAuthorized(w, authData) {
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

// EnableNotificationAPIMethod - api method which serve `/api/v1/notifications/enable/{id}` (set `is_active` to 1(true|enabled)).
func (dashboard *Dashboard) EnableNotificationAPIMethod(w http.ResponseWriter, r *http.Request) {
	// will be denied to unauthorized users
	err := dashboard.authService.IsGrantedAccessByCookies(w, r)
	if err != nil {
		util.WriteResponse(w, dataAPI.NewErrorData(err.Error()), http.StatusBadRequest)
		return
	}

	// receiving id from r.URL.Path
	id, err := util.ExctractId(r.URL.Path)
	if err != nil {
		util.WriteResponse(w, dataAPI.NewErrorData(err.Error()), http.StatusBadRequest)
		return
	}

	// receiving target notification from database
	notification, err := dashboard.manager.Repository.NotificationQueue().FindById(id)
	if err != nil {
		util.WriteResponse(w, dataAPI.NewErrorData(err.Error()), http.StatusBadRequest)
		return
	}

	// updating target notification
	notification, err = dashboard.manager.Repository.NotificationQueue().MakeAsEnabled(notification)
	if err != nil {
		util.WriteResponse(w, dataAPI.NewErrorData(err.Error()), http.StatusInternalServerError)
		return
	}

	// send response to client
	util.WriteResponse(w, dataAPI.NewStatusData(notification.IsActive == true), http.StatusOK)
}

// DisableNotificationAPIMethod - api method which serve `/api/v1/notifications/disable/{id}` (set `is_active` to 0(false|disabled)).
func (dashboard *Dashboard) DisableNotificationAPIMethod(w http.ResponseWriter, r *http.Request) {
	// will be denied to unauthorized users
	err := dashboard.authService.IsGrantedAccessByCookies(w, r)
	if err != nil {
		util.WriteResponse(w, dataAPI.NewErrorData(err.Error()), http.StatusBadRequest)
		return
	}

	// receiving id from r.URL.Path
	id, err := util.ExctractId(r.URL.Path)
	if err != nil {
		util.WriteResponse(w, dataAPI.NewErrorData(err.Error()), http.StatusBadRequest)
		return
	}

	// receiving target notification from database
	notification, err := dashboard.manager.Repository.NotificationQueue().FindById(id)
	if err != nil {
		util.WriteResponse(w, dataAPI.NewErrorData(err.Error()), http.StatusBadRequest)
		return
	}

	// updating target notification
	notification, err = dashboard.manager.Repository.NotificationQueue().MakeAsDisabled(notification)
	if err != nil {
		util.WriteResponse(w, dataAPI.NewErrorData(err.Error()), http.StatusInternalServerError)
		return
	}

	// send response to client
	util.WriteResponse(w, dataAPI.NewStatusData(notification.IsActive == false), http.StatusOK)
}
