package types

import "time"

type Chat struct {
	ID        int32
	ChatId    int32
	CreatedAt time.Time
}

// NewChat - creating a new instance of chat
func NewChat() *Chat {
	return &Chat{}
}
