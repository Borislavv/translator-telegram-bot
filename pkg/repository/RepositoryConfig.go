package repository

type RepositoryConfig struct {
	Driver string `toml:"database_driver"`
	DSN    string `toml:"database_url"`
}

// NewRepositoryConfig - creating a new instance of RepositoryConfig
func NewRepositoryConfig() *RepositoryConfig {
	return &RepositoryConfig{
		Driver: "mysql",
		DSN:    "user=username password=userpassword host=localhost dbname=translator_telegram_bot sslmode=disable",
	}
}
