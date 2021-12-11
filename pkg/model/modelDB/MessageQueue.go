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

// NewMessageQueueConstructor - constructor of MessageQueue
func NewMessageQueueConstructor(queueId int64, message string, chatId int64) *MessageQueue {
	return &MessageQueue{
		QueueId: queueId,
		Message: message,
		ChatId:  chatId,
	}
}
