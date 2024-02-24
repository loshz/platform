package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
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
		time.Sleep(10 * time.Second)

		// Get eventd address.
		svcs, err := s.Discovery().Lookup(context.Background(), "eventd")
		if err != nil {
			s.Error(fmt.Errorf("error getting eventd service details from discovery: %w", err))
			return
		}

		// TODO: perform sanity check on returned eventd services.
		eventd := fmt.Sprintf("%s:%d", svcs[0].Address, svcs[0].GrpcPort)
		conn, err := grpc.DialContext(ctx, eventd, grpc.WithTransportCredentials(s.Creds().GrpcClient()))
		if err != nil {
			s.Error(fmt.Errorf("error dialing eventd: %w", err))
			return
		}
		defer conn.Close()

		client := apiv1.NewEventServiceClient(conn)

		// Register the machine details before sending events.
		req := &apiv1.RegisterHostRequest{
			Host: &apiv1.Host{
				MachineId: "blah",
				Hostname:  "p14s",
			},
			Timestamp: time.Now().Unix(),
		}
		if _, err := client.RegisterHost(ctx, req); err != nil {
			s.Error(fmt.Errorf("error registering machine: %w", err))
			return
		}

		// Initiate stream and start sending events.
		stream, err := client.SendEvent(ctx)
		if err != nil {
			s.Error(fmt.Errorf("error getting stream: %w", err))
			return
		}

		t := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-t.C:
				var buf bytes.Buffer
				enc := gob.NewEncoder(&buf)
				if err := enc.Encode(&apiv1.NetworkEvent{}); err != nil {
					log.Error().Err(err).Msg("error serializing event data")
					continue
				}
				req := &apiv1.SendEventRequest{
					Type:      apiv1.EventType_EVENT_TYPE_NETWORK,
					MachineId: "p14s",
					Data:      buf.Bytes(),
				}
				if err := stream.Send(req); err != nil {
					log.Error().Err(err).Msg("error sending event")
					continue
				}
			case <-ctx.Done():
				_, _ = stream.CloseAndRecv()
				return
			}
		}
	}()

	return nil
}
