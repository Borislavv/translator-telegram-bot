package dashboardHandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Borislavv/Translator-telegram-bot/pkg/helper"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDashboard"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/util"
)

// TranslationPage - handler of "/translation" endpoint
func (dashboard *Dashboard) TranslationPage(w http.ResponseWriter, r *http.Request) {

	page := modelDashboard.NewPage()
	page.AddMenu(modelDashboard.NewMenu())

	type Content struct {
		Title            string
		Funcs            *helper.TemplateFuncsHelper
		ErrorMessage     string
		LoginByTokenLink string
		TranslationUrl   string
	}

	content := Content{
		Title:            "translation",
		LoginByTokenLink: "/login",
		Funcs:            helper.NewTemplateFuncsHelper(),
		TranslationUrl:   "/api/v1/translate",
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
	}

	// Порядок файлов важен, сначала дочки, далее родители
	templates := []string{
		"./pkg/templates/dashboard/pages/translation/content.html.tmpl",
		"./pkg/templates/dashboard/menu/menu.html.tmpl",
		"./pkg/templates/empty.subContent.html.tmpl",
		"./pkg/templates/base.html.tmpl",
	}

	page.AddConent(content)

	util.RenderFromFiles(w, templates, page)
}

// TranslationAPI - method for translate simple text by API
func (dashboard *Dashboard) TranslationAPI(w http.ResponseWriter, r *http.Request) {
	type RequestData struct {
		Text string `json:"text"`
	}

	type ResponseData struct {
		Text string `json:"text"`
	}

	// read body
	bytesSlice, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// parse body
	var requestData RequestData
	err = json.Unmarshal(bytesSlice, &requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// translate target text
	translatedText, err := dashboard.translatorService.TranslateText(requestData.Text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseData := ResponseData{
		Text: translatedText,
	}

	util.WriteResponse(w, responseData)
}
