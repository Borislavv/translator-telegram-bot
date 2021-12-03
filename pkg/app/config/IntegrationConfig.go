package config

type IntegrationConfig struct {
	ApiEndpoint string `toml:"api_endpoint"`
	ApiToken    string `toml:"api_token"`
}

// NewIntegrationConfig - creating a new instance of IntegrationConfig
func NewIntegrationConfig() *IntegrationConfig {
	return &IntegrationConfig{}
}
