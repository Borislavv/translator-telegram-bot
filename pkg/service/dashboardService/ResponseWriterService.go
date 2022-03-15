package dashboardService

import (
	"encoding/json"
	"net/http"

	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelAPI/dataAPI"
	response "github.com/Borislavv/Translator-telegram-bot/pkg/model/modelAPI/responseAPI"
)

type ResponseWriter struct {
}

// NewResponseWriter - constructor of ResponseWriter structure.
func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{}
}

// WriteDataIntoResponse - writing data to ResponseWriter of http package.
func (writer *ResponseWriter) WriteDataIntoResponse(w http.ResponseWriter, data dataAPI.DataInterface, code int) {
	resp := response.NewResponse(data, code)

	jsonData, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
