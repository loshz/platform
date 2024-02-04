package main

import (
	"context"

	"google.golang.org/grpc"

	apiv1 "github.com/loshz/platform/internal/api/v1"
	"github.com/loshz/platform/internal/config"
	"github.com/loshz/platform/internal/credentials"
	pgrpc "github.com/loshz/platform/internal/grpc"
	"github.com/loshz/platform/internal/service"
)

func main() {
	s := service.New("discoveryd")

	// Load required service credentials before startup.
	s.LoadCredentials(credentials.GrpcServer)

	// Run the service.
	s.Run(run)
}

func run(ctx context.Context, s *service.Service) error {
	// Create gRPC server options including interceptors and timeout.
	opts := []grpc.ServerOption{
		grpc.Creds(s.Creds().GrpcServer()),
		grpc.UnaryInterceptor(pgrpc.UnaryInterceptor(s.ID())),
		grpc.StreamInterceptor(pgrpc.StreamInterceptor(s.ID())),
		grpc.ConnectionTimeout(s.Config().Duration(config.KeyGRPCServerConnTimeout)),
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
