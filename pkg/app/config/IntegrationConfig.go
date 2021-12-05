package config

type IntegrationConfig struct {
	Telegram struct {
		ApiEndpoint string `toml:"api_endpoint"`
		ApiToken    string `toml:"api_token"`
	}
	Translator struct {
		ApiEndpoint string `toml:"translator_api_endpoint"`
	}
}

// NewIntegrationConfig - creating a new instance of IntegrationConfig
func NewIntegrationConfig() *IntegrationConfig {
	return &IntegrationConfig{}
}
