package translator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
)

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
func (gateway *TranslatorGateway) RequestTranslate(
	sourceLanguage string,
	targetLanguage string,
	text string,
) (string, error) {
	// Getting updated messages from channels
	response, err := http.Get(fmt.Sprintf(gateway.endpoint, sourceLanguage, targetLanguage, url.QueryEscape(text)))
	if err != nil {
		return "", err
	}

	// Reading body to slice of bytes
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if !json.Valid(body) {
		return "", errors.New("API returned a response not in json format")
	}

	// Decoding json to TranslatedMessage struct
	translatedMessage := model.NewTranslatedMessage()
	if err := json.Unmarshal(body, translatedMessage); err != nil {
		return "", err
	}

	// No translation have been received
	if len(translatedMessage.Translations) == 0 {
		return "", errors.New("no one translations have been received")
	}

	return translatedMessage.Translations[0].Text, nil
}
