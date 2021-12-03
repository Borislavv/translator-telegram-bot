package modelDB

import "time"

type Chat struct {
	ID             int64
	ExternalChatId int64
	CreatedAt      time.Time
}

// NewChat - creating a new instance of chat
func NewChat() *Chat {
	return &Chat{}
}
