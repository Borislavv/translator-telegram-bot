package model

import (
	"errors"
	"net/http"
)

const (
	UsernameFieldKey = "username"
	TokenFieldKey    = "token"
)

type AuthRequestData struct {
	username string
	token    string
	writer   http.ResponseWriter
}

// NewAuthRequestData - constructor of AuthRequestData struct
func NewAuthRequestData(w http.ResponseWriter, r *http.Request) (*AuthRequestData, error) {
	if r.FormValue(UsernameFieldKey) == "" && r.FormValue(TokenFieldKey) == "" {
		return nil, nil
	} else if r.FormValue(UsernameFieldKey) == "" {
		return nil, errors.New("Username cannot be omitted")
	} else if r.FormValue(TokenFieldKey) == "" {
		return nil, errors.New("Token cannot be omitted")
	}

	return &AuthRequestData{
		username: r.FormValue(UsernameFieldKey),
		token:    r.FormValue(TokenFieldKey),
		writer:   w,
	}, nil
}

// GetUsername - getter of username
func (authData *AuthRequestData) GetUsername() string {
	return authData.username
}

// GetToken - getter of token
func (authData *AuthRequestData) GetToken() string {
	return authData.token
}

// GetWriter - getter of writer (http.ResponseWriter)
func (authData *AuthRequestData) GetWriter() http.ResponseWriter {
	return authData.writer
}
