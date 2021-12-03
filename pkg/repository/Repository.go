package repository

import "database/sql"

type RepositoryInterface interface {
	Open() error
	Close() error
}

type Repository struct {
	config                 *RepositoryConfig
	db                     *sql.DB
	ChatRepository         *ChatRepository
	MessageQueueRepository *MessageQueueRepository
	UserRepository         *UserRepository
}

type RepositoryConfig struct {
	Driver string `toml:"database_driver"`
	DSN    string `toml:"database_url"`
}

// Open - opening a new connection with database
func (repository *Repository) Open() error {
	db, err := sql.Open(repository.config.Driver, repository.config.DSN)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	repository.db = db

	return nil
}

// Close - closing connection with database
func (repository *Repository) Close() error {
	if err := repository.db.Close(); err != nil {
		return err
	}

	return nil
}

/**
 * Define a new repository further:
 */

// Chat - creating an instance of ChatRepository
func (repository *Repository) Chat() *ChatRepository {
	if repository.ChatRepository != nil {
		return repository.ChatRepository
	}

	repository.ChatRepository = &ChatRepository{
		repository: repository,
	}

	return repository.ChatRepository
}

// User -creating an instance of UserRepository
func (repository *Repository) User() *UserRepository {
	if repository.ChatRepository != nil {
		return repository.UserRepository
	}

	repository.UserRepository = &UserRepository{
		repository: repository,
	}

	return repository.UserRepository
}

// MessageQueue -creating an instance of MessageQueueRepository
func (repository *Repository) MessageQueue() *MessageQueueRepository {
	if repository.ChatRepository != nil {
		return repository.MessageQueueRepository
	}

	repository.MessageQueueRepository = &MessageQueueRepository{
		repository: repository,
	}

	return repository.MessageQueueRepository
}
