package main

import (
	"fmt"

	"google.golang.org/grpc"

	"github.com/loshz/platform/pkg/config"
	pgrpc "github.com/loshz/platform/pkg/grpc"
	pbv1 "github.com/loshz/platform/pkg/pb/v1"
	"github.com/loshz/platform/pkg/service"
)

func main() {
	s := service.New("eventd")

	// Load required service config.
	s.LoadGRPCServerConfig()

	s.Run(run)
}

func run(s *service.Service) error {
	// Load TLS credentials.
	ca := s.Config.String(config.KeyGRPCTLSCA)
	cert := s.Config.String(config.KeyGRPCServerCert)
	key := s.Config.String(config.KeyGRPCServerKey)

	// Create new gRPC server TLS credentials.
	creds, err := pgrpc.NewServerTransportCreds(ca, cert, key)
	if err != nil {
		return fmt.Errorf("error :%w", err)
	}

	// Create gRPC server options including interceptors and timeout.
	opts := []grpc.ServerOption{
		grpc.Creds(creds),
		grpc.UnaryInterceptor(pgrpc.UnaryInterceptor(s.ID)),
		grpc.StreamInterceptor(pgrpc.StreamInterceptor(s.ID)),
		grpc.ConnectionTimeout(s.Config.Duration(config.KeyGRPCServerConnTimeout)),
	}

	srv := pgrpc.NewServer(s.Ctx(), opts)
	srv.RegisterService(&pbv1.Eventd_ServiceDesc, &grpcServer{})
	go srv.Serve(s.Config.Int(config.KeyGRPCServerPort))

	return nil
}
