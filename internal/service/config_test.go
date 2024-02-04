package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/loshz/platform/internal/config"
)

func TestLoadRequiredConfig(t *testing.T) {
	// Set gRPC env vars.
	t.Setenv("PLAT_SERVICE_LOG_LEVEL", "debug")
	t.Setenv("PLAT_SERVICE_STARTUP_TIMEOUT", "5s")
	t.Setenv("PLAT_SERVICE_SHUTDOWN_TIMEOUT", "10s")
	t.Setenv("PLAT_HTTP_PORT", "8001")

	// Create a new service and load required config.
	s := New("test")
	s.LoadRequiredConfig()

	// Assert loaded config is as expected.
	assert.Equal(t, s.Config.Get(config.KeyServiceLogLevel), "debug")
	assert.Equal(t, s.Config.Get(config.KeyServiceStartupTimeout), "5s")
	assert.Equal(t, s.Config.Get(config.KeyServiceShutdownTimeout), "10s")
	assert.Equal(t, s.Config.Get(config.KeyHTTPPort), "8001")
}

func TestLoadGrpcServerConfig(t *testing.T) {
	// Set gRPC env vars.
	t.Setenv("PLAT_GRPC_TLS_CA", "/path/to/ca")
	t.Setenv("PLAT_GRPC_SERVER_PORT", "8002")
	t.Setenv("PLAT_GRPC_SERVER_CERT", "/path/to/cert")
	t.Setenv("PLAT_GRPC_SERVER_KEY", "/path/to/key")
	t.Setenv("PLAT_GRPC_SERVER_CONN_TIMEOUT", "10s")

	// Create a new service and load grpc server config.
	s := New("test")
	s.LoadGrpcServerConfig()

	// Assert loaded config is as expected.
	assert.Equal(t, s.Config.Get(config.KeyGRPCTLSCA), "/path/to/ca")
	assert.Equal(t, s.Config.Get(config.KeyGRPCServerPort), "8002")
	assert.Equal(t, s.Config.Get(config.KeyGRPCServerCert), "/path/to/cert")
	assert.Equal(t, s.Config.Get(config.KeyGRPCServerKey), "/path/to/key")
	assert.Equal(t, s.Config.Get(config.KeyGRPCServerConnTimeout), "10s")
}

func TestLoadGrpcClientConfig(t *testing.T) {
	// Set gRPC env vars.
	t.Setenv("PLAT_GRPC_TLS_CA", "/path/to/ca")
	t.Setenv("PLAT_GRPC_CLIENT_CERT", "/path/to/cert")
	t.Setenv("PLAT_GRPC_CLIENT_KEY", "/path/to/key")

	// Create a new service and load grpc client config.
	s := New("test")
	s.LoadGrpcClientConfig()

	// Assert loaded config is as expected.
	assert.Equal(t, s.Config.Get(config.KeyGRPCTLSCA), "/path/to/ca")
	assert.Equal(t, s.Config.Get(config.KeyGRPCClientCert), "/path/to/cert")
	assert.Equal(t, s.Config.Get(config.KeyGRPCClientKey), "/path/to/key")
}
