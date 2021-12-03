package modelDB

import "time"

type User struct {
	ID        int64
	ChatId    int64
	Username  string
	CreatedAt time.Time
}

// NewUser - creating a new User instance
func NewUser() *User {
	return &User{}
}
