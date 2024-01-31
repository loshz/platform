package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	apiv1 "github.com/loshz/platform/internal/api/v1"
	"github.com/loshz/platform/internal/config"
	"github.com/loshz/platform/internal/discovery"
	pgrpc "github.com/loshz/platform/internal/grpc"
)

// RegisterDiscovery attempts to periodically register a service with the discovery service.
func (s *Service) RegisterDiscovery(ctx context.Context) {
	s.LoadDiscoveryConfig()

	// Return early if discovery not enabled.
	if !s.Config.Bool(config.KeyServiceDiscoveryEnabled) {
		return
	}

	// Load TLS credentials.
	ca := s.Config.String(config.KeyGRPCTLSCA)
	cert := s.Config.String(config.KeyGRPCClientCert)
	key := s.Config.String(config.KeyGRPCClientKey)
	creds, err := pgrpc.NewClientTransportCreds(ca, cert, key)
	if err != nil {
		s.Error(fmt.Errorf("error loading grpc tls credentials: %w", err))
		return
	}

	s.ds = discovery.New(s.Config.String(config.KeyServiceDiscoveryAddr), creds)

	// Create a timer with a small initial tick to allow service processes to start
	// before registering for discovery.
	// TODO: we should attempt to register for discovery immediately, blocking further execution
	// until successful.
	t := time.NewTimer(5 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			// Reset the timer to the larger periodic interval.
			t.Reset(s.Config.Duration(config.KeyServiceRegisterInt))

			service := &apiv1.Service{
				Uuid:     s.ID().String(),
				HttpPort: uint32(s.Config.Uint(config.KeyHTTPPort)),
				GrpcPort: uint32(s.Config.Uint(config.KeyGRPCServerPort)),
				LastSeen: time.Now().Unix(),
			}
			if err := s.ds.Register(ctx, service); err != nil {
				// TODO: we should exit if this has failed multiple times.
				log.Error().Err(err).Msg("error registering service for discovery")
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
