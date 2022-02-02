package util

import (
	"encoding/json"
	"net/http"
)

// WriteResponse - writing data to ResponseWriter of http package.
func WriteResponse(w http.ResponseWriter, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
