package config

type ServerConfig struct {
	Host           string `toml:"server_host"`
	Port           string `toml:"server_port"`
	StaticFilesDir string `toml:"server_static_files_path"`
}

// NewServerConfig - creating a new instance of ServerConfig
func NewServerConfig() *ServerConfig {
	return &ServerConfig{}
}
