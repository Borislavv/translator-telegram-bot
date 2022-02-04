package model

import (
	"errors"
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
		return NewAuthCookieDataConstructor(w), errors.New("Username not found, please visit the login page.")
	}

	tokenCookie, err := r.Cookie("token")
	if err != nil {
		return NewAuthCookieDataConstructor(w), errors.New("Token not found, please visit the login page.")
	}

	return &AuthCookieData{
		username: usernameCookie.Value,
		token:    tokenCookie.Value,
		writer:   w,
	}, nil
}

// NewAuthCookieDataConstructor - empty struct constructor of AuthCookieData
func NewAuthCookieDataConstructor(w http.ResponseWriter) *AuthCookieData {
	return &AuthCookieData{
		writer: w,
	}
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

// SetWriter - setter of writer (http.ResponseWriter)
func (authData *AuthCookieData) SetWriter(w http.ResponseWriter) {
	authData.writer = w
}
