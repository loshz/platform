package service

import "github.com/loshz/platform/internal/config"

// LoadRequiredConfig is a helper function for loading config required by
// a service.
func (s *Service) LoadRequiredConfig() {
	s.Config.MustLoad(config.KeyServiceLogLevel, "info", config.ParseLogLevel)
	s.Config.MustLoad(config.KeyServiceStartupTimeout, "10s", config.ParseDuration)
	s.Config.MustLoad(config.KeyServiceShutdownTimeout, "5s", config.ParseDuration)
	s.Config.MustLoad(config.KeyHTTPPort, 8001, config.ParseInt)
}

// LoadDiscoveryConfig is a helper function for loading service discovery config.
func (s *Service) LoadDiscoveryConfig() {
	s.Config.MustLoad(config.KeyServiceDiscoveryEnabled, true, config.ParseBool)
	s.Config.MustLoad(config.KeyServiceDiscoveryAddr, "discoveryd:8000", config.ParseString)
	s.Config.MustLoad(config.KeyServiceRegisterInt, "30s", config.ParseDuration)
}

// LoadGRPCServerConfig is a helper function for loading required gRPC
// server config.
func (s *Service) LoadGRPCServerConfig() {
	s.Config.MustLoad(config.KeyGRPCTLSCA, "", config.ParseString)
	s.Config.MustLoad(config.KeyGRPCServerPort, 8002, config.ParseInt)
	s.Config.MustLoad(config.KeyGRPCServerCert, "", config.ParseString)
	s.Config.MustLoad(config.KeyGRPCServerKey, "", config.ParseString)
	s.Config.MustLoad(config.KeyGRPCServerConnTimeout, "10s", config.ParseDuration)
}

// LoadGRPCClientConfig is a helper function for loading required gRPC
// server config.
func (s *Service) LoadGRPCClientConfig() {
	s.Config.MustLoad(config.KeyGRPCTLSCA, "", config.ParseString)
	s.Config.MustLoad(config.KeyGRPCClientCert, "", config.ParseString)
	s.Config.MustLoad(config.KeyGRPCClientKey, "", config.ParseString)
}
