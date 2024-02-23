package credentials

import (
	"fmt"

	grpc "google.golang.org/grpc/credentials"

	"github.com/loshz/platform/internal/config"
	pgrpc "github.com/loshz/platform/internal/grpc"
)

type Credential uint

const (
	GrpcClient Credential = iota
	GrpcServer
)

type Store struct {
	grpc struct {
		client, server grpc.TransportCredentials
	}
}

func (s *Store) GrpcClient() grpc.TransportCredentials { return s.grpc.client }
func (s *Store) GrpcServer() grpc.TransportCredentials { return s.grpc.server }

func (s *Store) LoadGrpcClientCreds(c *config.Config) error {
	// Load TLS credentials.
	ca := c.String(config.KeyGrpcTlsCA)
	cert := c.String(config.KeyGrpcClientCert)
	key := c.String(config.KeyGrpcClientKey)

	// Create new gRPC server TLS credentials.
	creds, err := pgrpc.NewClientTransportCreds(ca, cert, key)
	if err != nil {
		return fmt.Errorf("error loading grpc client tls credentials: %w", err)
	}

	s.grpc.client = creds
	return nil
}

func (s *Store) LoadGrpcServerCreds(c *config.Config) error {
	// Load TLS credentials.
	ca := c.String(config.KeyGrpcTlsCA)
	cert := c.String(config.KeyGrpcServerCert)
	key := c.String(config.KeyGrpcServerKey)

	// Create new gRPC server TLS credentials.
	creds, err := pgrpc.NewServerTransportCreds(ca, cert, key)
	if err != nil {
		return fmt.Errorf("error loading grpc server tls credentials: %w", err)
	}

	s.grpc.server = creds
	return nil
}
