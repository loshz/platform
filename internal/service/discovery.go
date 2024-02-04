package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	apiv1 "github.com/loshz/platform/internal/api/v1"
	"github.com/loshz/platform/internal/config"
)

// RegisterDiscovery attempts to periodically register a service with the discovery service.
func (s *Service) RegisterDiscovery(ctx context.Context) {
	s.LoadDiscoveryConfig()

	// Return early if discovery not enabled.
	if !s.Config().Bool(config.KeyServiceDiscoveryEnabled) {
		return
	}

	conn, err := grpc.Dial(s.Config().String(config.KeyServiceDiscoveryAddr), grpc.WithTransportCredentials(s.Creds().GrpcClient()))
	if err != nil {
		s.Error(fmt.Errorf("error dialing discovery service: %w", err))
		return
	}
	defer conn.Close()
	client := apiv1.NewDiscoveryServiceClient(conn)

	t := time.NewTicker(s.Config().Duration(config.KeyServiceRegisterInt))
	for {
		select {
		case <-t.C:
			req := &apiv1.RegisterServiceRequest{
				Service: &apiv1.Service{
					Uuid:     s.ID(),
					HttpPort: uint32(s.Config().Uint(config.KeyHTTPPort)),
					GrpcPort: uint32(s.Config().Uint(config.KeyGRPCServerPort)),
					LastSeen: time.Now().Unix(),
				},
			}
			if _, err := client.RegisterService(context.Background(), req); err != nil {
				stat, _ := status.FromError(err)
				log.Error().Err(errors.New(stat.Message())).Str("code", stat.Code().String()).Msg("error registering service for discovery")
				continue
			}
		case <-ctx.Done():
			return
		}
	}
}

// DeregisterDiscovery attempts to deregister a service with the discovery service.
func (s *Service) DeregisterDiscovery() error {
	// Return early if discovery not enabled.
	if !s.Config().Bool(config.KeyServiceDiscoveryEnabled) {
		return nil
	}

	conn, err := grpc.Dial(s.Config().String(config.KeyServiceDiscoveryAddr), grpc.WithTransportCredentials(s.Creds().GrpcClient()))
	if err != nil {
		return fmt.Errorf("error dialing discovery service: %w", err)
	}
	defer conn.Close()
	client := apiv1.NewDiscoveryServiceClient(conn)

	req := &apiv1.DeregisterServiceRequest{
		Uuid: s.ID(),
	}
	if _, err := client.DeregisterService(context.Background(), req); err != nil {
		stat, _ := status.FromError(err)
		return errors.New(stat.Message())
	}

	return nil
}
