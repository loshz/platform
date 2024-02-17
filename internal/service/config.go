package service

import "github.com/loshz/platform/internal/config"

// LoadRequiredConfig is a helper function for loading config required by
// a service.
func (s *Service) LoadRequiredConfig() {
	s.Config().MustLoad(config.KeyServiceLogLevel, "info", config.ParseLogLevel)
	s.Config().MustLoad(config.KeyServiceShutdownTimeout, "10s", config.ParseDuration)
	s.Config().MustLoad(config.KeyHTTPPort, 8001, config.ParseInt)
}

// LoadDiscoveryConfig is a helper function for loading service discovery config.
func (s *Service) LoadDiscoveryConfig() {
	s.Config().MustLoad(config.KeyServiceDiscoveryEnabled, true, config.ParseBool)
	s.Config().MustLoad(config.KeyServiceDiscoveryAddr, "discoveryd:8000", config.ParseString)
	s.Config().MustLoad(config.KeyServiceRegisterInt, "300s", config.ParseDuration)
}

// LoadGrpcServerConfig is a helper function for loading required gRPC
// server config.
func (s *Service) LoadGrpcServerConfig() {
	s.Config().MustLoad(config.KeyGRPCTLSCA, "/usr/local/share/ca-certificates/ca.crt.pem", config.ParseString)
	s.Config().MustLoad(config.KeyGRPCServerPort, 8002, config.ParseInt)
	s.Config().MustLoad(config.KeyGRPCServerCert, "/usr/local/share/ca-certificates/server.crt.pem", config.ParseString)
	s.Config().MustLoad(config.KeyGRPCServerKey, "/usr/local/share/ca-certificates/server.key.pem", config.ParseString)
	s.Config().MustLoad(config.KeyGRPCServerConnTimeout, "10s", config.ParseDuration)
}

// LoadGrpcClientConfig is a helper function for loading required gRPC
// server config.
func (s *Service) LoadGrpcClientConfig() {
	s.Config().MustLoad(config.KeyGRPCTLSCA, "/usr/local/share/ca-certificates/ca.crt.pem", config.ParseString)
	s.Config().MustLoad(config.KeyGRPCClientCert, "/usr/local/share/ca-certificates/client.crt.pem", config.ParseString)
	s.Config().MustLoad(config.KeyGRPCClientKey, "/usr/local/share/ca-certificates/client.key.pem", config.ParseString)
}
