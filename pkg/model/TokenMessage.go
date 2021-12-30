package model

type TokenMessage struct {
	UpdatedMessage *UpdatedMessage
	Token          string
}

// NewTokenMessage - constructor of TokenMessage struct
func NewTokenMessage(message *UpdatedMessage, token string) *TokenMessage {
	return &TokenMessage{
		UpdatedMessage: message,
		Token:          token,
	}
}
