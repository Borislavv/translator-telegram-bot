package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
)

// Method keys
const GetUpdatesMethod = "getUpdates"
const SendMessageMethod = "sendMessage"

// Gateway requirements
const UrlPattern string = "https://api.telegram.org/bot%s/%s"
const ApiToken string = "5022497048:AAGQcUiyExpJXr3pjjv_cgVody3rv_MvjZ4"

// TelegramGateway - representation of telegram api as facade
type TelegramGateway struct {
	Endpoint string
	ApiToken string
	Methods  map[string]string
}

// NewTelegramGateway - Telegram gateway conctructor
func NewTelegramGateway() *TelegramGateway {
	methods := make(map[string]string)

	methods[GetUpdatesMethod] = GetUpdatesMethod
	methods[SendMessageMethod] = fmt.Sprintf("%s%s", SendMessageMethod, "?chat_id=%s&text=%s")

	return &TelegramGateway{
		Endpoint: UrlPattern,
		ApiToken: ApiToken,
		Methods:  methods,
	}
}

// GetUpdates - getting messages from telegram channels with offset
func (gateway *TelegramGateway) GetUpdates(offset int64) (*model.UpdatedMessages, error) {
	// Getting updated messages from channels
	response, err := http.Get(fmt.Sprintf(gateway.Endpoint, gateway.ApiToken, gateway.Methods[GetUpdatesMethod]))
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	// Reading body to slide of bytes
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	// Decoding json to UpdatedMessages struct
	updatedMessages := model.NewUpdatedMessages()
	if err := json.Unmarshal(body, updatedMessages); err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return updatedMessages, nil
}

// SendMessage - sending message to target chat
func (gateway *TelegramGateway) SendMessage(chatId string, message string) error {
	reqResponse, err := http.Post(
		fmt.Sprintf(
			fmt.Sprintf(
				gateway.Endpoint, gateway.ApiToken, gateway.Methods[SendMessageMethod],
			),
			chatId,
			"Hello world",
		),
		"application/json",
		strings.NewReader(url.Values{}.Encode()),
	)
	if err != nil {
		log.Panicln(err)
		return err
	}

	reqBody, err := ioutil.ReadAll(reqResponse.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	type RequestResponse struct {
		Status bool `json:"ok"`
	}

	var sendMessageResponse RequestResponse
	if err := json.Unmarshal(reqBody, &sendMessageResponse); err != nil {
		log.Fatalln(err)
		return err
	}

	if sendMessageResponse.Status {
		return nil
	}

	return errors.New("not `ok` respnse status received")
}
