package types

import "time"

type MessageQueue struct {
	ID        int32
	QueueId   int64
	Message   string
	ChatId    int32
	CreatedAt time.Time
}

// NewMessageQueue - creating a new instance of MessageQueue
func NewMessageQueue() *MessageQueue {
	return &MessageQueue{}
}
