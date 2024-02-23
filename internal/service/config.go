package service

import "github.com/loshz/platform/internal/config"

// LoadRequiredConfig is a helper function for loading config required by
// a service.
func (s *Service) LoadRequiredConfig() {
	s.Config().MustLoad(config.KeyServiceLogLevel, "info", config.ParseLogLevel)
	s.Config().MustLoad(config.KeyServiceShutdownTimeout, "10s", config.ParseDuration)
	s.Config().MustLoad(config.KeyHttpPort, 8001, config.ParseInt)
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
	s.Config().MustLoad(config.KeyGrpcTlsCA, "/usr/local/share/ca-certificates/ca.crt.pem", config.ParseString)
	s.Config().MustLoad(config.KeyGrpcServerPort, 8002, config.ParseInt)
	s.Config().MustLoad(config.KeyGrpcServerCert, "/usr/local/share/ca-certificates/server.crt.pem", config.ParseString)
	s.Config().MustLoad(config.KeyGrpcServerKey, "/usr/local/share/ca-certificates/server.key.pem", config.ParseString)
	s.Config().MustLoad(config.KeyGrpcServerConnTimeout, "10s", config.ParseDuration)
}

// LoadGrpcClientConfig is a helper function for loading required gRPC
// server config.
func (s *Service) LoadGrpcClientConfig() {
	s.Config().MustLoad(config.KeyGrpcTlsCA, "/usr/local/share/ca-certificates/ca.crt.pem", config.ParseString)
	s.Config().MustLoad(config.KeyGrpcClientCert, "/usr/local/share/ca-certificates/client.crt.pem", config.ParseString)
	s.Config().MustLoad(config.KeyGrpcClientKey, "/usr/local/share/ca-certificates/client.key.pem", config.ParseString)
	s.Config().MustLoad(config.KeyGrpcClientTimeout, "10s", config.ParseDuration)
}
