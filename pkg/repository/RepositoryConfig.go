package repository

type RepositoryConfig struct {
	Driver string `toml:"database_driver"`
	DSN    string `toml:"database_url"`
}

// NewRepositoryConfig - creating a new instance of RepositoryConfig
func NewRepositoryConfig() *RepositoryConfig {
	return &RepositoryConfig{}
}
