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
	pgrpc "github.com/loshz/platform/internal/grpc"
)

// RegisterDiscovery attempts to periodically register a service with the discovery service.
func (s *Service) RegisterDiscovery() {
	s.LoadDiscoveryConfig()

	// Load TLS credentials.
	ca := s.Config.String(config.KeyGRPCTLSCA)
	cert := s.Config.String(config.KeyGRPCClientCert)
	key := s.Config.String(config.KeyGRPCClientKey)
	creds, err := pgrpc.NewClientTransportCreds(ca, cert, key)
	if err != nil {
		s.errCh <- fmt.Errorf("error loading grpc tls credentials: %w", err)
		return
	}

	conn, err := grpc.Dial(s.Config.String(config.KeyServiceDiscoveryAddr), grpc.WithTransportCredentials(creds))
	if err != nil {
		s.errCh <- fmt.Errorf("error dialing discovery service: %w", err)
		return
	}
	defer conn.Close()
	client := apiv1.NewDiscoveryServiceClient(conn)

	t := time.NewTicker(s.Config.Duration(config.KeyServiceRegisterInt))
	for {
		select {
		case <-t.C:
			req := &apiv1.RegisterServiceRequest{
				Service: &apiv1.Service{
					Uuid:      s.ID(),
					HttpPort:  uint32(s.Config.Uint(config.KeyHTTPPort)),
					GrpcPort:  uint32(s.Config.Uint(config.KeyGRPCServerPort)),
					Timestamp: time.Now().Unix(),
				},
			}
			if _, err := client.RegisterService(context.Background(), req); err != nil {
				stat, _ := status.FromError(err)
				log.Error().Err(errors.New(stat.Message())).Str("code", stat.Code().String()).Msg("error registering service for discovery")
				continue
			}
		case <-s.Ctx().Done():
			return
		}
	}
}

// DeregisterDiscovery attempts to deregister a service with the discovery service.
func (s *Service) DeregisterDiscovery() error {
	// Load TLS credentials.
	ca := s.Config.String(config.KeyGRPCTLSCA)
	cert := s.Config.String(config.KeyGRPCClientCert)
	key := s.Config.String(config.KeyGRPCClientKey)
	creds, err := pgrpc.NewClientTransportCreds(ca, cert, key)
	if err != nil {
		return fmt.Errorf("error loading grpc tls credentials: %w", err)
	}

	conn, err := grpc.Dial(s.Config.String(config.KeyServiceDiscoveryAddr), grpc.WithTransportCredentials(creds))
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
