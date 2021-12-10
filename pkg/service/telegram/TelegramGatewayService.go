package telegram

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Borislavv/Translator-telegram-bot/pkg/app/manager"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelInterface"
)

// Method keys
const GetUpdatesMethod = "getUpdates"
const SendMessageMethod = "sendMessage"

// TelegramGateway - representation of telegram api as facade
type TelegramGateway struct {
	Endpoint string
	ApiToken string
	Methods  map[string]string
}

// NewTelegramGateway - Telegram gateway conctructor
func NewTelegramGateway(manager *manager.Manager) *TelegramGateway {
	methods := make(map[string]string)

	methods[GetUpdatesMethod] = fmt.Sprintf("%s%s", GetUpdatesMethod, "?offset=%s")
	methods[SendMessageMethod] = fmt.Sprintf("%s%s", SendMessageMethod, "?chat_id=%s&text=%s")

	return &TelegramGateway{
		Endpoint: manager.Config.Integration.Telegram.ApiEndpoint,
		ApiToken: manager.Config.Integration.Telegram.ApiToken,
		Methods:  methods,
	}
}

// GetUpdates - getting messages from telegram channels with offset
func (gateway *TelegramGateway) GetUpdates(message modelInterface.RequestMessageInterface) *model.UpdatedMessages {
	// Getting updated messages from channels
	response, err := http.Get(
		fmt.Sprintf(
			fmt.Sprintf(gateway.Endpoint, gateway.ApiToken, gateway.Methods[GetUpdatesMethod]),
			fmt.Sprint(message.GetOffset()),
		),
	)
	if err != nil {
		return nil
	}

	// Reading body to slice of bytes
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil
	}

	// Decoding json to UpdatedMessages struct
	updatedMessages := model.NewUpdatedMessages()
	if err := json.Unmarshal(body, updatedMessages); err != nil {
		return nil
	}

	return updatedMessages
}

// SendMessage - sending message to target chat
func (gateway *TelegramGateway) SendMessage(message modelInterface.ResponseMessageInterface) error {
	reqResponse, err := http.Post(
		fmt.Sprintf(
			fmt.Sprintf(
				gateway.Endpoint, gateway.ApiToken, gateway.Methods[SendMessageMethod],
			),
			message.GetChatId(),
			url.QueryEscape(message.GetMessage()),
		),
		"application/json",
		strings.NewReader(url.Values{}.Encode()),
	)
	if err != nil {
		return err
	}

	reqBody, err := ioutil.ReadAll(reqResponse.Body)
	if err != nil {
		return err
	}

	type RequestResponse struct {
		Status bool `json:"ok"`
	}

	var sendMessageResponse RequestResponse
	if err := json.Unmarshal(reqBody, &sendMessageResponse); err != nil {
		return err
	}

	if sendMessageResponse.Status {
		return nil
	}

	return errors.New("not `ok` respnse status received: " + fmt.Sprintf("%+v\n", string(reqBody)))
}
