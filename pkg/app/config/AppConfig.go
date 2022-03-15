package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Borislavv/Translator-telegram-bot/pkg/repository"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/util"
)

// Environment mode
const ProdMode = "prod"
const DevMode = "dev"

// Config vars. which must be defined anyway.
var ConfigRequirementsMap = []string{
	"ENV_MODE",
	"API_TOKEN",
	"API_ENDPOINT",
	"TRANSLATOR_API_ENDPOINT",
	"DATABASE_DSN",
	"DATABASE_DRIVER",
	"SERVER_HOST",
	"SERVER_PORT",
	"SERVER_STATIC_FILES_PATH",
}

type Environment struct {
	Mode string
}

// Config - main config structure, which will incapsulate all other config structs.
type Config struct {
	Repository  *repository.RepositoryConfig
	Environment *Environment
	Integration *IntegrationConfig
	Server      *ServerConfig
}

// NewConfig - creating a new instance of Config.
func New() *Config {
	return &Config{
		Repository:  repository.NewRepositoryConfig(),
		Environment: &Environment{ProdMode},
		Integration: NewIntegrationConfig(),
		Server:      NewServerConfig(),
	}
}

// Load - loading a config file to struct by received path.
func (config *Config) Load() *Config {
	if err := config.isValid(); err != nil {
		log.Fatalln(util.Trace(err))
		return nil
	}

	// app config init.
	config.Environment.Mode = os.Getenv("ENV_MODE")

	// translator config init.
	config.Integration.Telegram.ApiToken = os.Getenv("API_TOKEN")
	config.Integration.Telegram.ApiEndpoint = os.Getenv("API_ENDPOINT")
	config.Integration.Translator.ApiEndpoint = os.Getenv("TRANSLATOR_API_ENDPOINT")

	// database config init.
	config.Repository.DSN = os.Getenv("DATABASE_DSN")
	config.Repository.Driver = os.Getenv("DATABASE_DRIVER")

	// server config init.
	config.Server.Host = os.Getenv("SERVER_HOST")
	config.Server.Port = os.Getenv("SERVER_PORT")
	config.Server.StaticFilesDir = os.Getenv("SERVER_STATIC_FILES_PATH")

	return config
}

// isValid - method for validate OS-environment variables.
func (config *Config) isValid() error {
	for _, varKey := range ConfigRequirementsMap {
		if os.Getenv(varKey) == "" {
			return errors.New(fmt.Sprintf("Variable `%v` must be defined. Check your docker-compose file.", varKey))
		}
	}

	return nil
}
