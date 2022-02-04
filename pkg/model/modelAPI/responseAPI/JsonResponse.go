package response

import (
	"net/http"

	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelAPI/dataAPI"
)

type Response struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
}

// NewResponse - constructor of Response struct
func NewResponse(data dataAPI.DataInterface, code int) *Response {
	if code == 0 {
		code = http.StatusOK
	}

	return &Response{
		Data: data,
		Code: code,
	}
}
