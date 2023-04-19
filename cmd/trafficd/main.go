package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/loshz/platform/pkg/config"
	pgrpc "github.com/loshz/platform/pkg/grpc"
	pbv1 "github.com/loshz/platform/pkg/pb/v1"
	"github.com/loshz/platform/pkg/service"
)

func main() {
	s := service.New("trafficd")

	// Load required service config.
	s.LoadGRPCClientConfig()

	s.Run(run)
}

func run(s *service.Service) error {
	// Load TLS credentials.
	ca := s.Config.String(config.KeyGRPCTLSCA)
	cert := s.Config.String(config.KeyGRPCClientCert)
	key := s.Config.String(config.KeyGRPCClientKey)
	creds, err := pgrpc.NewClientTransportCreds(ca, cert, key)
	if err != nil {
		return fmt.Errorf("error loading grpc tls credentials: %w", err)
	}

	go func() {
		// TODO: don't hard code address.
		conn, err := grpc.Dial("eventd:8003", grpc.WithTransportCredentials(creds))
		if err != nil {
			log.Error().Err(err).Msg("error dialing eventd")
			// TODO: s.Exit() or continually check for conn.
		}
		defer conn.Close()
		client := pbv1.NewEventdClient(conn)

		t := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-t.C:
				res, err := client.Event(context.Background(), &pbv1.EventRequest{Hostname: "blah"})
				if err != nil {
					log.Error().Err(err).Msg("error making request to eventd")
					continue
				}

				log.Info().Msgf("eventd response: %s", res.Uuid)
			case <-s.Ctx().Done():
				conn.Close()
				return
			}
		}
	}()

	return nil
}
