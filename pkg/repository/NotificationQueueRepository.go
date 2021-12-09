package repository

import (
	"time"

	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDB"
)

type NotificationQueueRepository struct {
	connection *Repository
}

// Create - adding a new `notification_queue` row into db
func (repository *NotificationQueueRepository) Create(ntfQueue *modelDB.NotificationQueue) (*modelDB.NotificationQueue, error) {
	result, err := repository.connection.db.Exec(
		"INSERT INTO notification_queue (message_queue_id, chat_id, scheduled_for) VALUES (?,?,?)",
		ntfQueue.MessageQueueId,
		ntfQueue.ChatId,
		ntfQueue.ScheduledFor,
	)
	if err != nil {
		return nil, err
	}

	ntfQueue.ID, _ = result.LastInsertId()

	return ntfQueue, nil
}

// FindByScheduledDate - trying to find notification by `scheduled_for`
func (repository *NotificationQueueRepository) FindByScheduledDate(dateTime time.Time) ([]*modelDB.NotificationQueue, error) {
	var responseStack []*modelDB.NotificationQueue

	rows, err := repository.connection.db.Query(
		"SELECT nq.id, msg.message, c.external_chat_id FROM notification_queue `nq`"+
			" LEFT JOIN chat `c` ON nq.chat_id = c.id "+
			" LEFT JOIN message_queue `msg` ON nq.message_queue_id = msg.id"+
			" WHERE is_sent = 0 AND scheduled_for BETWEEN ? AND ?",
		dateTime.Format("2006-01-02 15:04:05"),
		dateTime.Add(1*time.Minute).Format("2006-01-02 15:04:05"),
	)

	// Close rows anyway
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		notification := modelDB.NewNotificationQueue()

		if err = rows.Scan(
			&notification.ID,
			&notification.Message,
			&notification.ExternalChatId,
		); err != nil {
			return nil, err
		}

		// adding row to response slice
		responseStack = append(responseStack, notification)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return responseStack, nil
}

// MakeAsSent - set a value of `is_sent` column to 1(true)
func (repository *NotificationQueueRepository) MakeAsSent(notification *modelDB.NotificationQueue) (*modelDB.NotificationQueue, error) {
	_, err := repository.connection.db.Exec(
		"UPDATE notification_queue SET is_sent = 1 WHERE id = ?",
		notification.ID,
	)
	if err != nil {
		return nil, err
	}

	notification.IsSent = true

	return notification, nil
}
