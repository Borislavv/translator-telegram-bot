package modelDB

import "time"

type MessageQueue struct {
	ID        int64
	QueueId   int64
	Message   string
	ChatId    int64
	CreatedAt time.Time
}

// NewMessageQueue - creating a new instance of MessageQueue
func NewMessageQueue() *MessageQueue {
	return &MessageQueue{}
}
