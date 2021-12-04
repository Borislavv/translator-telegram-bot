package repository

import (
	"database/sql"

	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDB"
)

type ChatRepository struct {
	connection *Repository
}

// Create - adding a new `chat` row into db
func (repository *ChatRepository) Create(chat *modelDB.Chat) (*modelDB.Chat, error) {
	result, err := repository.connection.db.Exec("INSERT INTO chat (external_chat_id) VALUES (?)", chat.ExternalChatId)
	if err != nil {
		return nil, err
	}

	chat.ID, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return chat, nil
}

// FindByChatId - trying to find row into `chat` by id
func (repository *ChatRepository) FindByChatId(chatId string) (*modelDB.Chat, error) {
	chat := modelDB.NewChat()

	if err := repository.connection.db.QueryRow(
		"SELECT id, external_chat_id, created_at FROM chat WHERE id = ?",
		chatId,
	).Scan(
		&chat.ID,
		&chat.ExternalChatId,
		&chat.CreatedAt,
	); err != nil {
		return nil, err
	}

	return chat, nil
}

// FindByExternalChatId - trying to find row into `chat` by external_chat_id
func (repository *ChatRepository) FindByExternalChatId(externalChatId int64) (*modelDB.Chat, error) {
	chat := modelDB.NewChat()

	if err := repository.connection.db.QueryRow(
		"SELECT id, external_chat_id, created_at FROM chat WHERE external_chat_id = ?",
		externalChatId,
	).Scan(
		&chat.ID,
		&chat.ExternalChatId,
		&chat.CreatedAt,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	return chat, nil
}
