package main

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/loshz/platform/internal/config"
	pgrpc "github.com/loshz/platform/internal/grpc"
	pbv1 "github.com/loshz/platform/internal/api/v1"
	"github.com/loshz/platform/internal/service"
)

func main() {
	s := service.New("discoveryd")

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
		grpc.UnaryInterceptor(pgrpc.UnaryInterceptor(s.ID())),
		grpc.StreamInterceptor(pgrpc.StreamInterceptor(s.ID())),
		grpc.ConnectionTimeout(s.Config.Duration(config.KeyGRPCServerConnTimeout)),
	}

	// Create a discovery server and start the service eviction process in the background.
	ds := NewDiscoveryServer()
	go ds.StartEvictionProcess(s.Ctx())

	// Create a gRPC server and register the service.
	srv := pgrpc.NewServer(opts)
	srv.RegisterService(&pbv1.DiscoveryService_ServiceDesc, ds)

	// Start the gRPC server in the background.
	go srv.Serve(s.Ctx(), s.Config.Int(config.KeyGRPCServerPort))

	// Listen for gRPC server errors and exit if received.
	go func() {
		if err := <-srv.Error(); err != nil {
			log.Error().Err(err).Msg("grpc server error")
			s.Exit(service.ExitError)
		}
	}()

	return nil
}
