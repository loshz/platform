package service

import "github.com/loshz/platform/pkg/config"

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
func (s *Service) LoadGRPCClentConfig() {
	s.Config.MustLoad(config.KeyGRPCTLSCA, "", config.ParseString)
	s.Config.MustLoad(config.KeyGRPCClientCert, "", config.ParseString)
	s.Config.MustLoad(config.KeyGRPCClientKey, "", config.ParseString)
}
