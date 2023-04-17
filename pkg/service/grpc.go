package service

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/loshz/platform/pkg/config"
	pgrpc "github.com/loshz/platform/pkg/grpc"
	pbv1 "github.com/loshz/platform/pkg/pb/v1"
)

type grpcServer struct {
	pbv1.UnimplementedPlatformServiceServer

	service string
}

func (s *grpcServer) Status(context.Context, *emptypb.Empty) (*pbv1.StatusResponse, error) {
	return &pbv1.StatusResponse{
		Service: s.service,
		Status:  pbv1.Status_STATUS_OK,
	}, nil
}

// ServeGRPC configures, registers services and starts a gRPC server on a given port.
// It is intentially not called in Start() as not every service requires a gRPC server,
// therefore it should be called directly by the service itself.
func (s *Service) ServeGRPC(desc *grpc.ServiceDesc, svc interface{}) {
	port := fmt.Sprintf(":%d", s.Config.Int(config.KeyGRPCServerPort))

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Error().Err(err).Msg("error creating tcp listener for grpc server")
		s.Exit(ExitError)
	}

	// Load TLS credentials.
	cert := s.Config.String(config.KeyGRPCServerCert)
	key := s.Config.String(config.KeyGRPCServerKey)
	ca := s.Config.String(config.KeyGRPCClientCA)
	creds, err := pgrpc.NewServerTransportCreds(cert, key, ca)
	if err != nil {
		log.Error().Err(err).Msg("error loading grpc tls credentials")
		s.Exit(ExitError)
	}

	opts := []grpc.ServerOption{
		grpc.Creds(creds),
		grpc.UnaryInterceptor(pgrpc.UnaryInterceptor(s.ID)),
		grpc.StreamInterceptor(pgrpc.StreamInterceptor(s.ID)),
		grpc.ConnectionTimeout(s.Config.Duration(config.KeyGRPCServerConnTimeout)),
	}

	// Configure server and register services.
	srv := grpc.NewServer(opts...)
	srv.RegisterService(&pbv1.PlatformService_ServiceDesc, &grpcServer{})
	srv.RegisterService(desc, svc)

	go func() {
		log.Info().Msgf("grpc server running on %s", port)
		if err := srv.Serve(lis); err != grpc.ErrServerStopped {
			log.Error().Err(err).Msg("grpc server error")
			s.Exit(ExitError)
		}
	}()

	// Wait for service to exit and shutdown.
	<-s.ctx.Done()
	log.Info().Msg("stopping grpc server")
	srv.GracefulStop()
}
