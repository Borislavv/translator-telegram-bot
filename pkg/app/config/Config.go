package config

import "github.com/Borislavv/Translator-telegram-bot/pkg/repository"

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
}

// NewConfig - creating a new instance of Config
func New() *Config {
	return &Config{
		Repository:  repository.NewRepositoryConfig(),
		Environment: &Environment{ProdMode},
	}
}
