package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	apiv1 "github.com/loshz/platform/internal/api/v1"
	"github.com/loshz/platform/internal/config"
	"github.com/loshz/platform/internal/discovery"
)

const (
	// Max no. of discovery register retries.
	MaxDiscoveryRetries = 3
)

func (s *Service) EnableDiscovery() {
	s.LoadDiscoveryConfig()
	s.ds = discovery.New(s.Config().String(config.KeyServiceDiscoveryAddr), s.Creds().GrpcClient())
}

// RegisterDiscovery attempts to periodically register a service with the discovery service.
func (s *Service) RegisterDiscovery(ctx context.Context) {
	// Return early if discovery not enabled.
	registerInterval := s.Config().Duration(config.KeyServiceRegisterInt)
	if registerInterval == 0 {
		return
	}

	// Create a timer with a small initial tick to allow service processes to start
	// before registering for discovery.
	t := time.NewTimer(5 * time.Second)
	defer t.Stop()

	// Keep track of failed retries.
	retries := 0
	for {
		select {
		case <-t.C:
			// Reset the timer to the larger periodic interval.
			t.Reset(registerInterval)

			service := &apiv1.Service{
				Uuid:     s.ID(),
				Address:  s.Name(), // TODO: this won't work if we run more than 1 replica.
				HttpPort: uint32(s.Config().Uint(config.KeyHTTPPort)),
				GrpcPort: uint32(s.Config().Uint(config.KeyGRPCServerPort)),
				LastSeen: time.Now().Unix(),
			}
			if err := s.ds.Register(context.TODO(), service); err != nil {
				retries++
				if retries == MaxDiscoveryRetries {
					s.Error(fmt.Errorf("failed to register for discovery: %w", err))
					return
				}

				log.Error().Err(err).Msg("error registering service for discovery, retrying")
				continue
			}
		case <-ctx.Done():
			return
		}
	}
}

// DeregisterDiscovery attempts to deregister a service with the discovery service.
func (s *Service) DeregisterDiscovery() error {
	if s.ds == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.ds.Deregister(ctx, s.ID())
}
