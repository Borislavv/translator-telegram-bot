package translator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/util"
)

// Languages
const EnLanguage = "en"
const RuLanguage = "ru"

// TranslatorGateway - representation of translator api as facade
type TranslatorGateway struct {
	endpoint string
}

// NewTranslatorGateway - creating a new instance of TranslatorGateway
func NewTranslatorGateway(manager *manager.Manager) *TranslatorGateway {
	return &TranslatorGateway{
		endpoint: manager.Config.Integration.Translator.ApiEndpoint,
	}
}

// Translate - trying to detect source and target languages and translate text
func (gateway *TranslatorGateway) Translate(text string) (string, error) {
	sourceLanguage, targetLanguage := DetectLanguages(text)

	// Getting updated messages from channels
	response, err := http.Get(fmt.Sprintf(gateway.endpoint, sourceLanguage, targetLanguage, url.QueryEscape(text)))
	if err != nil {
		log.Fatalln(util.Trace() + err.Error())
		return "Cannot translate: " + text, err
	}

	// Reading body to slice of bytes
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(util.Trace() + err.Error())
		return "Cannot translate: " + text, err
	}

	// Decoding json to TranslatedMessage struct
	translatedMessage := model.NewTranslatedMessage()
	if err := json.Unmarshal(body, translatedMessage); err != nil {
		log.Fatalln(util.Trace() + err.Error())
		return "Cannot translate: " + text, err
	}

	if len(translatedMessage.Translations) == 0 {
		log.Fatalln(util.Trace() + "no one translations have been received")
	}

	return translatedMessage.Translations[0].Text, nil
}
