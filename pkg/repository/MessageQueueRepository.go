package repository

import (
	"database/sql"

	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDB"
)

type MessageQueueRepository struct {
	connection *Repository
}

// Create - adding a new `message_queue` row into db
func (repository *MessageQueueRepository) Create(msgQueue *modelDB.MessageQueue) (*modelDB.MessageQueue, error) {
	result, err := repository.connection.db.Exec(
		"INSERT INTO message_queue (queue_id, message, chat_id) VALUES (?,?,?)",
		msgQueue.QueueId,
		msgQueue.Message,
		msgQueue.ChatId,
	)
	if err != nil {
		return nil, err
	}

	msgQueue.ID, _ = result.LastInsertId()

	return msgQueue, nil
}

// GetOffset - will return an (max:`queue_id` + 1)
func (repository *MessageQueueRepository) GetOffset() (int64, error) {
	msgQueue := modelDB.NewMessageQueue()

	if err := repository.connection.db.QueryRow(
		"SELECT queue_id FROM message_queue ORDER BY queue_id DESC LIMIT 1",
	).Scan(
		&msgQueue.QueueId,
	); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		return 0, err
	}

	return msgQueue.QueueId + 1, nil
}
