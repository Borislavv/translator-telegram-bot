package config

import (
	"log"

	"github.com/BurntSushi/toml"

	"github.com/Borislavv/Translator-telegram-bot/pkg/repository"
	"github.com/Borislavv/Translator-telegram-bot/pkg/service/util"
)

// Environment mode
const ProdMode = "prod"
const DevMode = "dev"

// Environment default values
const DefaultConfigPath = "config/.env.prod.toml"

type Environment struct {
	Mode string
}

// Config - main config structure, which will incapsulate all other config structs
type Config struct {
	Repository  *repository.RepositoryConfig
	Environment *Environment
	Integration *IntegrationConfig
}

// NewConfig - creating a new instance of Config
func New() *Config {
	return &Config{
		Repository:  repository.NewRepositoryConfig(),
		Environment: &Environment{ProdMode},
		Integration: NewIntegrationConfig(),
	}
}

// Load - loading a config file to struct by received path
func (config *Config) Load(configurationPath string, environmentMode string) *Config {
	// Set environment mode
	config.Environment.Mode = environmentMode

	// Database config loading
	_, err := toml.DecodeFile(configurationPath, config.Repository)
	if err != nil {
		log.Fatalln(util.Trace(err))
	}

	// Telegram api config loading
	_, err = toml.DecodeFile(configurationPath, &config.Integration.Telegram)
	if err != nil {
		log.Fatalln(util.Trace(err))
	}

	// Translator config loading
	_, err = toml.DecodeFile(configurationPath, &config.Integration.Translator)
	if err != nil {
		log.Fatalln(util.Trace(err))
	}

	return config
}
