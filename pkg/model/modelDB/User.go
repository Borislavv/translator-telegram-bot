package modelDB

import (
	"time"
)

type User struct {
	ID        int64  `json:"id"`
	ChatId    int64  `json:"chat_id"`
	Username  string `json:"username"`
	Token     string
	TZ        string `json:"timezone"`
	CreatedAt time.Time
}

// NewUser - creating a new User instance
func NewUser() *User {
	return &User{}
}
