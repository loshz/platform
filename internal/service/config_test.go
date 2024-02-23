package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/loshz/platform/internal/config"
)

func TestLoadRequiredConfig(t *testing.T) {
	// Set required env vars.
	t.Setenv("PLAT_SERVICE_LOG_LEVEL", "debug")
	t.Setenv("PLAT_SERVICE_SHUTDOWN_TIMEOUT", "20s")
	t.Setenv("PLAT_HTTP_PORT", "8888")

	// Create a new service and load required config.
	s := New("required")
	s.LoadRequiredConfig()

	// Assert loaded config is as expected.
	assert.Equal(t, s.Config().Get(config.KeyServiceLogLevel), "debug")
	assert.Equal(t, s.Config().Get(config.KeyServiceShutdownTimeout), "20s")
	assert.Equal(t, s.Config().Get(config.KeyHttpPort), "8888")
}

func TestLoadDiscoveryConfig(t *testing.T) {
	// Set discovery env vars.
	t.Setenv("PLAT_SERVICE_DISCOVERY_ENABLED", "false")
	t.Setenv("PLAT_SERVICE_DISCOVERY_ADDR", "discoveryd:8888")
	t.Setenv("PLAT_SERVICE_REGISTER_INTERVAL", "30s")

	// Create a new service and load required config.
	s := New("discovery")
	s.LoadDiscoveryConfig()

	// Assert loaded config is as expected.
	assert.Equal(t, s.Config().Get(config.KeyServiceDiscoveryEnabled), "false")
	assert.Equal(t, s.Config().Get(config.KeyServiceDiscoveryAddr), "discoveryd:8888")
	assert.Equal(t, s.Config().Get(config.KeyServiceRegisterInt), "30s")
}

func TestLoadGrpcServerConfig(t *testing.T) {
	// Set gRPC server env vars.
	t.Setenv("PLAT_GRPC_TLS_CA", "/path/to/ca")
	t.Setenv("PLAT_GRPC_SERVER_PORT", "8002")
	t.Setenv("PLAT_GRPC_SERVER_CERT", "/path/to/cert")
	t.Setenv("PLAT_GRPC_SERVER_KEY", "/path/to/key")
	t.Setenv("PLAT_GRPC_SERVER_CONN_TIMEOUT", "10s")

	// Create a new service and load grpc server config.
	s := New("grpc-server")
	s.LoadGrpcServerConfig()

	// Assert loaded config is as expected.
	assert.Equal(t, s.Config().Get(config.KeyGrpcTlsCA), "/path/to/ca")
	assert.Equal(t, s.Config().Get(config.KeyGrpcServerPort), "8002")
	assert.Equal(t, s.Config().Get(config.KeyGrpcServerCert), "/path/to/cert")
	assert.Equal(t, s.Config().Get(config.KeyGrpcServerKey), "/path/to/key")
	assert.Equal(t, s.Config().Get(config.KeyGrpcServerConnTimeout), "10s")
}

func TestLoadGrpcClientConfig(t *testing.T) {
	// Set gRPC client env vars.
	t.Setenv("PLAT_GRPC_TLS_CA", "/path/to/ca")
	t.Setenv("PLAT_GRPC_CLIENT_CERT", "/path/to/cert")
	t.Setenv("PLAT_GRPC_CLIENT_KEY", "/path/to/key")
	t.Setenv("PLAT_GRPC_CLIENT_TIMEOUT", "30s")

	// Create a new service and load grpc client config.
	s := New("grpc-client")
	s.LoadGrpcClientConfig()

	// Assert loaded config is as expected.
	assert.Equal(t, s.Config().Get(config.KeyGrpcTlsCA), "/path/to/ca")
	assert.Equal(t, s.Config().Get(config.KeyGrpcClientCert), "/path/to/cert")
	assert.Equal(t, s.Config().Get(config.KeyGrpcClientKey), "/path/to/key")
	assert.Equal(t, s.Config().Get(config.KeyGrpcClientTimeout), "30s")
}
