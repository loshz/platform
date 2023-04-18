package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// Server is a wrapper around a *grpc.Server. It provides helper functions
// for starting the server and registering services.
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

// Server starts the *grpc.Server on a given port in a goroutine. It waits for the
// server's context to be done before gracefully shutting down.
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
