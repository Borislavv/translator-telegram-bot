package repository

import (
	"database/sql"

	"github.com/Borislavv/Translator-telegram-bot/pkg/model/modelDB"
)

type UserRepository struct {
	connection *Repository
}

// Create - saving instance of User into database
func (repository *UserRepository) Create(user *modelDB.User) (*modelDB.User, error) {
	result, err := repository.connection.db.Exec(
		"INSERT INTO user (chat_id, username) VALUES (?,?)",
		user.ChatId,
		user.Username,
	)
	if err != nil {
		return nil, err
	}

	user.ID, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindByUsername - trying to find `user` into db by `username`
func (repository *UserRepository) FindByUsername(username string) (*modelDB.User, error) {
	user := modelDB.NewUser()

	if err := repository.connection.db.QueryRow(
		"SELECT id, chat_id, username, created_at FROM user WHERE username = ?",
		user.Username,
	).Scan(
		&user.ID,
		&user.ChatId,
		&user.Username,
		&user.CreatedAt,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	return user, nil
}
