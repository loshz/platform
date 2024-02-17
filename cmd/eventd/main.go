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
	s := service.New("eventd")

	// Load required service credentials and dependencies before startup.
	s.LoadCredentials(credentials.GrpcClient, credentials.GrpcServer)

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

	// Create a gRPC server and register the service.
	grpcSrv := pgrpc.NewServer(opts)
	grpcSrv.RegisterService(&apiv1.EventService_ServiceDesc, &grpcServer{})

	// Start the gRPC server in the background.
	go s.ServeGRPC(ctx, grpcSrv)

	return nil
}
