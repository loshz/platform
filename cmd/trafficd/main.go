package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	apiv1 "github.com/loshz/platform/internal/api/v1"
	"github.com/loshz/platform/internal/credentials"
	"github.com/loshz/platform/internal/service"
)

func main() {
	s := service.New("trafficd")

	// Load required service credentials and dependencies before startup.
	s.LoadCredentials(credentials.GrpcClient)

	// Run the service.
	s.Run(run)
}

func run(ctx context.Context, s *service.Service) error {
	s.Scheduler().Add(1)

	go func() {
		defer s.Scheduler().Done()

		// TODO: refactor this whole function to use periodic refresh and retries.
		time.Sleep(5 * time.Second)

		trf, err := NewTrafficd()
		if err != nil {
			s.Error(fmt.Errorf("error initializing traffic service: %w", err))
			return
		}

		// Get eventd address.
		eventd, err := trf.GetEventdAddr(s.Discovery())
		if err != nil {
			s.Error(fmt.Errorf("error getting eventd service details from discovery: %w", err))
			return
		}

		conn, err := grpc.DialContext(ctx, eventd.String(), grpc.WithTransportCredentials(s.Creds().GrpcClient()))
		if err != nil {
			s.Error(fmt.Errorf("error dialing eventd: %w", err))
			return
		}
		defer conn.Close()
		client := apiv1.NewEventServiceClient(conn)

		// Register the machine details before sending events.
		if err := trf.RegisterHost(ctx, client); err != nil {
			s.Error(fmt.Errorf("error registering machine: %w", err))
			return
		}

		// Initiate stream and start sending events.
		if err := trf.StreamEvents(ctx, client); err != nil {
			s.Error(fmt.Errorf("error streaming events: %w", err))
			return
		}
	}()

	return nil
}
