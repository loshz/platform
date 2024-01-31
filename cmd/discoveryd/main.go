package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	apiv1 "github.com/loshz/platform/internal/api/v1"
	"github.com/loshz/platform/internal/config"
	pgrpc "github.com/loshz/platform/internal/grpc"
	"github.com/loshz/platform/internal/service"
)

func main() {
	s := service.New("discoveryd")

	// Load required service config.
	s.LoadGRPCServerConfig()

	s.Run(run)
}

func run(ctx context.Context, s *service.Service) error {
	// Load TLS credentials.
	ca := s.Config.String(config.KeyGRPCTLSCA)
	cert := s.Config.String(config.KeyGRPCServerCert)
	key := s.Config.String(config.KeyGRPCServerKey)

	// Create new gRPC server TLS credentials.
	creds, err := pgrpc.NewServerTransportCreds(ca, cert, key)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	// Create gRPC server options including interceptors and timeout.
	opts := []grpc.ServerOption{
		grpc.Creds(creds),
		grpc.UnaryInterceptor(pgrpc.UnaryInterceptor(s.ID())),
		grpc.StreamInterceptor(pgrpc.StreamInterceptor(s.ID())),
		grpc.ConnectionTimeout(s.Config.Duration(config.KeyGRPCServerConnTimeout)),
	}

	// Create a discovery server and start the service eviction process in the background.
	ds := NewDiscoveryServer()
	go ds.StartEvictionProcess(ctx)

	// Create a gRPC server and register the service.
	grpcSrv := pgrpc.NewServer(opts)
	grpcSrv.RegisterService(&apiv1.DiscoveryService_ServiceDesc, ds)

	// Start the gRPC server in the background.
	go s.ServeGRPC(ctx, grpcSrv)

	return nil
}
