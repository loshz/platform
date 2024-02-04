package discovery

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	apiv1 "github.com/loshz/platform/internal/api/v1"
)

type Service struct {
	addr  string
	creds credentials.TransportCredentials
}

func New(addr string, creds credentials.TransportCredentials) *Service {
	return &Service{
		addr, creds,
	}
}

func (s *Service) Register(ctx context.Context, service *apiv1.Service) error {
	conn, err := grpc.Dial(s.addr, grpc.WithTransportCredentials(s.creds))
	if err != nil {
		return fmt.Errorf("error dialing discovery service: %w", err)
	}
	defer conn.Close()
	client := apiv1.NewDiscoveryServiceClient(conn)

	req := &apiv1.RegisterServiceRequest{
		Service: service,
	}
	if _, err := client.RegisterService(ctx, req); err != nil {
		stat, _ := status.FromError(err)
		return errors.New(stat.Message())
	}

	return nil
}

func (s *Service) Deregister(ctx context.Context, service_id string) error {
	conn, err := grpc.Dial(s.addr, grpc.WithTransportCredentials(s.creds))
	if err != nil {
		return fmt.Errorf("error dialing discovery service: %w", err)
	}
	defer conn.Close()
	client := apiv1.NewDiscoveryServiceClient(conn)

	req := &apiv1.DeregisterServiceRequest{
		Uuid: service_id,
	}
	if _, err := client.DeregisterService(ctx, req); err != nil {
		stat, _ := status.FromError(err)
		return errors.New(stat.Message())
	}

	return nil
}

func (s *Service) Lookup(ctx context.Context, service string) ([]*apiv1.Service, error) {
	conn, err := grpc.Dial(s.addr, grpc.WithTransportCredentials(s.creds))
	if err != nil {
		return nil, fmt.Errorf("error dialing discovery service: %w", err)
	}
	defer conn.Close()
	client := apiv1.NewDiscoveryServiceClient(conn)

	req := &apiv1.GetServicesRequest{
		Name: service,
	}
	res, err := client.GetServices(ctx, req)
	if err != nil {
		stat, _ := status.FromError(err)
		return nil, errors.New(stat.Message())
	}

	return res.Services, nil
}
