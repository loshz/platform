package service

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/loshz/platform/internal/config"
	"github.com/loshz/platform/internal/grpc"
)

func (s *Service) ServeGRPC(ctx context.Context, srv grpc.ServiceServer) {
	s.Scheduler().Add(1)
	defer s.Scheduler().Done()

	// Start the gRPC server in the background.
	go func() {
		if err := srv.Serve(ctx, s.Config().Int(config.KeyGrpcServerPort)); err != nil {
			s.Error(fmt.Errorf("grpc server error: %w", err))
			return
		}
	}()

	// Stop the gRPC server on exit.
	<-ctx.Done()
	log.Info().Msg("stopping grpc server")
	srv.Shutdown()
}
