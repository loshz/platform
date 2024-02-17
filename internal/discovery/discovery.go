package discovery

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	apiv1 "github.com/loshz/platform/internal/api/v1"
)

type Service struct {
	client apiv1.DiscoveryServiceClient
}

func New(ctx context.Context, addr string, creds credentials.TransportCredentials) (*Service, error) {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("error dialing discovery service: %w", err)
	}
	client := apiv1.NewDiscoveryServiceClient(conn)

	go func() {
		<-ctx.Done()
		// Small sleep so services can attempt to deregister.
		time.Sleep(500 * time.Millisecond)
		_ = conn.Close()
	}()

	return &Service{
		client,
	}, nil
}

func (s *Service) Register(ctx context.Context, service *apiv1.Service) error {
	req := &apiv1.RegisterServiceRequest{
		Service: service,
	}
	if _, err := s.client.RegisterService(ctx, req); err != nil {
		stat, _ := status.FromError(err)
		return errors.New(stat.Message())
	}

	return nil
}

func (s *Service) Deregister(ctx context.Context, service_id string) error {
	req := &apiv1.DeregisterServiceRequest{
		Uuid: service_id,
	}
	if _, err := s.client.DeregisterService(ctx, req); err != nil {
		stat, _ := status.FromError(err)
		return errors.New(stat.Message())
	}

	return nil
}

func (s *Service) Lookup(ctx context.Context, service string) ([]*apiv1.Service, error) {
	req := &apiv1.GetServicesRequest{
		Name: service,
	}
	res, err := s.client.GetServices(ctx, req)
	if err != nil {
		stat, _ := status.FromError(err)
		return nil, errors.New(stat.Message())
	}

	return res.Services, nil
}
