package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type ServiceServer interface {
	RegisterService(sd *grpc.ServiceDesc, svc interface{})
	Serve(ctx context.Context, port int) error
	Shutdown()
}

// Server is a wrapper around a *grpc.Server. It provides helper functions
// for starting the server and registering services.
type Server struct {
	srv *grpc.Server
}

func NewServer(opts []grpc.ServerOption) *Server {
	return &Server{
		srv: grpc.NewServer(opts...),
	}
}

// RegisterService registers a gRPC service to the underlying server.
func (s *Server) RegisterService(sd *grpc.ServiceDesc, svc interface{}) {
	s.srv.RegisterService(sd, svc)
}

// Server starts the *grpc.Server on a given port in a goroutine. It waits for the
// server's context to be done before gracefully shutting down.
func (s *Server) Serve(ctx context.Context, port int) error {
	lst, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	log.Info().Msgf("grpc server running on %s", lst.Addr())
	if err := s.srv.Serve(lst); err != grpc.ErrServerStopped {
		return err
	}

	return nil
}

func (s *Server) Shutdown() { s.srv.GracefulStop() }
