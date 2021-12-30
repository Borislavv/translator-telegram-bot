package modelInterface

import "net/http"

type UserAuthInterface interface {
	GetUsername() string
	GetToken() string
	GetWriter() http.ResponseWriter
}
