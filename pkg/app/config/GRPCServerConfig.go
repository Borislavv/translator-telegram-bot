package config

type GRPCServerConfig struct {
	Port     string `toml:"grpc_server_port"`
	Protocol string `toml:"grpc_net_protocol"`
}

// NewServerConfig - creating a new instance of GRPCServerConfig
func NewGRPCServerConfig() *GRPCServerConfig {
	return &GRPCServerConfig{}
}
