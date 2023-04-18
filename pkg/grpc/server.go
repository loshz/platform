package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Server struct {
	ctx context.Context
	srv *grpc.Server
}

func NewServer(ctx context.Context, opts []grpc.ServerOption) *Server {
	return &Server{
		ctx: ctx,
		srv: grpc.NewServer(opts...),
	}
}

func (s *Server) RegisterService(desc *grpc.ServiceDesc, svc interface{}) {
	s.srv.RegisterService(desc, svc)
}

// Serve configures, registers services and starts a gRPC server on a given port.
// It is intentially not called in Start() as not every service requires a gRPC server,
// therefore it should be called directly by the service itself.
func (s *Server) Serve(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Error().Err(err).Msg("error creating tcp listener for grpc server")
	}

	go func() {
		log.Info().Msgf("grpc server running on :%d", port)
		if err := s.srv.Serve(lis); err != grpc.ErrServerStopped {
			log.Error().Err(err).Msg("grpc server error")
		}
	}()

	// Wait for service to exit and shutdown.
	<-s.ctx.Done()
	log.Info().Msg("stopping grpc server")
	s.srv.GracefulStop()
}
