package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	apiv1 "github.com/loshz/platform/internal/api/v1"
	"github.com/loshz/platform/internal/config"
)

const (
	// Max no. of discovery register retries.
	MaxDiscoveryRetries = 3
)

// StartDiscovery attempts to establish a new connection to the discovery server.
func (s *Service) StartDiscovery(ctx context.Context) error {
	s.LoadDiscoveryConfig()

	// Return early if discovery not enabled.
	if !s.Config().Bool(config.KeyServiceDiscoveryEnabled) {
		return nil
	}

	// Start the discovery service with given credentials.
	return s.Discovery().Start(ctx, s.Config().String(config.KeyServiceDiscoveryAddr), s.Creds().GrpcClient())
}

// RegisterDiscovery attempts to periodically register a service with the discovery service.
func (s *Service) RegisterDiscovery(ctx context.Context) {
	// Return early if registration not enabled.
	if s.Config().Duration(config.KeyServiceRegisterInt) == 0 {
		return
	}

	s.Scheduler().Add(1)
	defer s.Scheduler().Done()

	// Create a timer with a small initial tick to allow service processes to start
	// before registering for discovery.
	t := time.NewTimer(5 * time.Second)
	defer t.Stop()

	// Get service details from config.
	interval := s.Config().Duration(config.KeyServiceRegisterInt)
	httpPort := s.Config().Uint(config.KeyHTTPPort)
	grpcPort := s.Config().Uint(config.KeyGRPCServerPort)

	// Keep track of failed retries.
	retries := 0
	for {
		select {
		case <-t.C:
			// Reset the timer to the larger periodic interval.
			t.Reset(interval)

			service := &apiv1.Service{
				Uuid:     s.ID(),
				Address:  s.Name(), // TODO: this won't work if we run more than 1 replica.
				HttpPort: uint32(httpPort),
				GrpcPort: uint32(grpcPort),
				LastSeen: time.Now().Unix(),
			}
			if err := s.Discovery().Register(context.TODO(), service); err != nil {
				retries++
				if retries == MaxDiscoveryRetries {
					s.Error(fmt.Errorf("failed to register for discovery: %w", err))
					return
				}

				log.Error().Err(err).Msg("error registering service for discovery, retrying")
				continue
			}
		case <-ctx.Done():
			// Attempt to deregister the service on shutdown.
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := s.Discovery().Deregister(ctx, s.ID()); err != nil {
				log.Error().Err(err).Msg("error deregistering service from discovery")
			}
			cancel()
			return
		}
	}
}
