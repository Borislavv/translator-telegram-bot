package repository

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

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

// New - creating a new instance of origin Repository
func New(repositoryConfig *RepositoryConfig) *Repository {
	return &Repository{
		config: repositoryConfig,
	}
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

// Repositories:
//
// Chat - creating an instance of ChatRepository
func (repository *Repository) Chat() *ChatRepository {
	if repository.ChatRepository != nil {
		return repository.ChatRepository
	}

	repository.ChatRepository = &ChatRepository{
		connection: repository,
	}

	return repository.ChatRepository
}

// User -creating an instance of UserRepository
func (repository *Repository) User() *UserRepository {
	if repository.ChatRepository != nil {
		return repository.UserRepository
	}

	repository.UserRepository = &UserRepository{
		connection: repository,
	}

	return repository.UserRepository
}

// MessageQueue -creating an instance of MessageQueueRepository
func (repository *Repository) MessageQueue() *MessageQueueRepository {
	if repository.ChatRepository != nil {
		return repository.MessageQueueRepository
	}

	repository.MessageQueueRepository = &MessageQueueRepository{
		connection: repository,
	}

	return repository.MessageQueueRepository
}
