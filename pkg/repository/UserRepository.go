package repository

import (
	"database/sql"
	"errors"

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

	tokenField := sql.NullString{}

	if err := repository.connection.db.QueryRow(
		"SELECT id, chat_id, username, token, created_at FROM user WHERE username = ?",
		username,
	).Scan(
		&user.ID,
		&user.ChatId,
		&user.Username,
		&tokenField,
		&user.CreatedAt,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		} else {
			return nil, nil
		}
	}

	// process nullable columns
	if tokenField.Valid {
		user.Token = tokenField.String
	} else {
		user.Token = ""
	}

	return user, nil
}

// FindByUsername - trying to find `user` into db by `token`
func (repository *UserRepository) FindByToken(token string) (*modelDB.User, error) {
	user := modelDB.NewUser()

	tokenField := sql.NullString{}

	if err := repository.connection.db.QueryRow(
		"SELECT id, chat_id, username, token, created_at FROM user WHERE token = ?",
		token,
	).Scan(
		&user.ID,
		&user.ChatId,
		&user.Username,
		&tokenField,
		&user.CreatedAt,
	); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		} else {
			return nil, nil
		}
	}

	// process nullable columns
	if tokenField.Valid {
		user.Token = tokenField.String
	} else {
		user.Token = ""
	}

	return user, nil
}

// SetToken - saving user token
func (repository *UserRepository) SetToken(user *modelDB.User, token string) (*modelDB.User, error) {
	result, err := repository.connection.db.Exec(
		"UPDATE user SET token = ? WHERE id = ?",
		token,
		user.ID,
	)
	if err != nil {
		return nil, err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affectedRows <= 0 {
		return nil, errors.New("no one user was updated while saving token")
	}

	user.Token = token

	return user, nil
}
