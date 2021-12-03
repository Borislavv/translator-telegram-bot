package types

import "time"

type User struct {
	ID        int32
	ChatId    int32
	Username  string
	CreatedAt time.Time
}

// NewUser - creating a new User instance
func NewUser() *User {
	return &User{}
}
