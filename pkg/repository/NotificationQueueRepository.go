package repository

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDB"
	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelInterface"
)

type NotificationQueueRepository struct {
	connection *Repository
}

// FindById - trying to find the target notification by id
func (repository *NotificationQueueRepository) FindById(id int64) (*modelDB.NotificationQueue, error) {
	notification := modelDB.NewNotificationQueue()

	if err := repository.connection.db.QueryRow(
		"SELECT id, message_queue_id, chat_id, is_sent, is_active, created_at, scheduled_for FROM notification_queue WHERE id = ?",
		id,
	).Scan(
		&notification.ID,
		&notification.MessageQueueId,
		&notification.ChatId,
		&notification.IsSent,
		&notification.IsActive,
		&notification.CreatedAt,
		&notification.ScheduledFor,
	); err != nil {
		return nil, err
	}

	return notification, nil
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
func (repository *NotificationQueueRepository) FindByScheduledDate() ([]*modelDB.NotificationQueue, error) {
	var responseStack []*modelDB.NotificationQueue

	rows, err := repository.connection.db.Query(
		"SELECT nq.id, msg.message, c.external_chat_id FROM notification_queue `nq`"+
			" LEFT JOIN chat `c` ON nq.chat_id = c.id"+
			" LEFT JOIN user `u` ON u.chat_id = c.id"+
			" LEFT JOIN message_queue `msg` ON nq.message_queue_id = msg.id"+
			" WHERE nq.is_sent = 0 AND nq.is_active = 1 AND nq.scheduled_for <= ?",
		time.Now().Add(1*time.Minute).Format("2006-01-02 15:04:05"),
	)

	// Close rows anyway
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

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

// FindNotSentByUsername - trying to find actual (scheduledAt > currentTime) notification by username
func (repository *NotificationQueueRepository) FindNotSentByUsername(
	username string,
	pagination modelInterface.PaginationInterface,
) ([]*modelDB.NotificationQueue, error) {
	var responseStack []*modelDB.NotificationQueue
	var err error

	queryStr := "SELECT nq.id, msg.message, nq.scheduled_for, nq.created_at, nq.is_active FROM notification_queue `nq`" +
		" LEFT JOIN chat `c` ON nq.chat_id = c.id" +
		" LEFT JOIN user `u` ON c.id = u.chat_id" +
		" LEFT JOIN message_queue `msg` ON nq.message_queue_id = msg.id" +
		" WHERE nq.is_sent = 0 AND u.username = ? AND nq.scheduled_for > ?" +
		" ORDER BY nq.scheduled_for ASC"

	var rows *sql.Rows
	if pagination.NeedPaginate() {
		queryStr += " LIMIT ? OFFSET ?"

		rows, err = repository.connection.db.Query(
			queryStr,
			username,
			time.Now().Format("2006-01-02 15:04:05"),
			strconv.Itoa(pagination.GetLimit()),
			strconv.Itoa((pagination.GetPage()-1)*pagination.GetLimit()),
		)
	} else {
		rows, err = repository.connection.db.Query(
			queryStr,
			username,
			time.Now().Format("2006-01-02 15:04:05"),
		)
	}

	// Close rows anyway
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		notification := modelDB.NewNotificationQueue()

		if err = rows.Scan(
			&notification.ID,
			&notification.Message,
			&notification.ScheduledFor,
			&notification.CreatedAt,
			&notification.IsActive,
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

// MakeAsDisabled - set a value of `is_active` column to 0(false)
func (repository *NotificationQueueRepository) MakeAsDisabled(notification *modelDB.NotificationQueue) (*modelDB.NotificationQueue, error) {
	_, err := repository.connection.db.Exec(
		"UPDATE notification_queue SET is_active = 0 WHERE id = ?",
		notification.ID,
	)
	if err != nil {
		return nil, err
	}

	notification.IsActive = false

	return notification, nil
}

// MakeAsEnabled - set a value of `is_active` column to 1(true)
func (repository *NotificationQueueRepository) MakeAsEnabled(notification *modelDB.NotificationQueue) (*modelDB.NotificationQueue, error) {
	_, err := repository.connection.db.Exec(
		"UPDATE notification_queue SET is_active = 1 WHERE id = ?",
		notification.ID,
	)
	if err != nil {
		return nil, err
	}

	notification.IsActive = true

	return notification, nil
}
