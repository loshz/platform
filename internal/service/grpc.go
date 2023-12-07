package service

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/loshz/platform/internal/config"
	"github.com/loshz/platform/internal/grpc"
)

func (s *Service) ServeGRPC(ctx context.Context, srv grpc.ServiceServer) {
	// Start the gRPC server in the background.
	go func() {
		if err := srv.Serve(ctx, s.Config.Int(config.KeyGRPCServerPort)); err != nil {
			s.SignalError(fmt.Errorf("grpc server error: %w", err))
		}
	}()

	// Stop the gRPC server on exit.
	<-ctx.Done()
	log.Info().Msg("stopping grpc server")
	srv.Shutdown()
}
