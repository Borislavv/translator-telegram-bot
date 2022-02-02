package model

import (
	"net/http"
)

type AuthCookieData struct {
	username string
	token    string
	writer   http.ResponseWriter
}

// NewAuthCookieData - constructor of AuthData struct
func NewAuthCookieData(w http.ResponseWriter, r *http.Request) (*AuthCookieData, error) {
	usernameCookie, err := r.Cookie("username")
	if err != nil {
		return nil, err
	}

	tokenCookie, err := r.Cookie("token")
	if err != nil {
		return nil, err
	}

	return &AuthCookieData{
		username: usernameCookie.Value,
		token:    tokenCookie.Value,
		writer:   w,
	}, nil
}

// GetUsername - getter of username
func (authData *AuthCookieData) GetUsername() string {
	return authData.username
}

// GetToken - getter of token
func (authData *AuthCookieData) GetToken() string {
	return authData.token
}

// GetWriter - getter of writer (http.ResponseWriter)
func (authData *AuthCookieData) GetWriter() http.ResponseWriter {
	return authData.writer
}
