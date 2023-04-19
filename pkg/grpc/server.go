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
	errCh chan error
	srv   *grpc.Server
}

func NewServer(opts []grpc.ServerOption) *Server {
	return &Server{
		errCh: make(chan error, 1),
		srv:   grpc.NewServer(opts...),
	}
}

// RegisterService registers a gRPC service to the underlying server.
func (s *Server) RegisterService(sd *grpc.ServiceDesc, svc interface{}) {
	s.srv.RegisterService(sd, svc)
}

// Error returns a receive only error channel so callers of the server
// can listen for errors.
func (s *Server) Error() <-chan error {
	return s.errCh
}

// Server starts the *grpc.Server on a given port in a goroutine. It waits for the
// server's context to be done before gracefully shutting down.
func (s *Server) Serve(ctx context.Context, port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		s.errCh <- err
	}

	go func() {
		log.Info().Msgf("grpc server running on %s", lis.Addr())
		if err := s.srv.Serve(lis); err != grpc.ErrServerStopped {
			s.errCh <- fmt.Errorf("grpc server error: %w", err)
		}
	}()

	// Wait for service to exit and shutdown.
	<-ctx.Done()
	log.Info().Msg("stopping grpc server")
	s.srv.GracefulStop()
}
