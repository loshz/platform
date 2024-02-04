package main

import (
	"context"
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

	go func() {
		// TODO: don't hard code address.
		conn, err := grpc.Dial("eventd:8004", grpc.WithTransportCredentials(s.Creds().GrpcClient()))
		if err != nil {
			log.Error().Err(err).Msg("error dialing eventd")
			// TODO: s.Exit() or continually check for conn.
		}
		defer conn.Close()
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
