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
	chatRepository         *ChatRepository
	messageQueueRepository *MessageQueueRepository
	userRepository         *UserRepository
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
	if repository.chatRepository != nil {
		return repository.chatRepository
	}

	repository.chatRepository = &ChatRepository{
		connection: repository,
	}

	return repository.chatRepository
}

// User -creating an instance of UserRepository
func (repository *Repository) User() *UserRepository {
	if repository.userRepository != nil {
		return repository.userRepository
	}

	repository.userRepository = &UserRepository{
		connection: repository,
	}

	return repository.userRepository
}

// MessageQueue -creating an instance of MessageQueueRepository
func (repository *Repository) MessageQueue() *MessageQueueRepository {
	if repository.messageQueueRepository != nil {
		return repository.messageQueueRepository
	}

	repository.messageQueueRepository = &MessageQueueRepository{
		connection: repository,
	}

	return repository.messageQueueRepository
}
