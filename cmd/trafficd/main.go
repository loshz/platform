package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	apiv1 "github.com/loshz/platform/internal/api/v1"
	"github.com/loshz/platform/internal/credentials"
	"github.com/loshz/platform/internal/service"
)

func main() {
	service.New("trafficd").Run(run)
}

func run(ctx context.Context, s *service.Service) error {
	// Load required service credentials before startup.
	if err := s.LoadCredentials(credentials.GrpcClient); err != nil {
		return err
	}

	// Enable the discovery service.
	s.EnableDiscovery()

	go func() {
		// TODO: add retries instead of sleeping.
		time.Sleep(30 * time.Second)

		// Get eventd address.
		svcs, err := s.Discovery().Lookup(context.Background(), "eventd")
		if err != nil {
			s.Error(fmt.Errorf("error getting eventd service details from discovery: %w", err))
			return
		}

		eventd := fmt.Sprintf("%s:%d", svcs[0].Address, svcs[0].GrpcPort)
		conn, err := grpc.Dial(eventd, grpc.WithTransportCredentials(s.Creds().GrpcClient()))
		if err != nil {
			s.Error(fmt.Errorf("error dialing eventd: %w", err))
			return
		}
		client := apiv1.NewEventServiceClient(conn)

		t := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-t.C:
				res, err := client.Event(context.Background(), &apiv1.EventRequest{Hostname: "blah"})
				if err != nil {
					log.Error().Err(err).Msg("error making request to eventd")
					continue
				}

				log.Info().Msgf("eventd response: %s", res.Uuid)
			case <-ctx.Done():
				conn.Close()
				return
			}
		}
	}()

	return nil
}
