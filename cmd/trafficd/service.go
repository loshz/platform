package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	apiv1 "github.com/loshz/platform/internal/api/v1"
	"github.com/loshz/platform/internal/discovery"
)

type Trafficd struct {
	MachineId string
}

func NewTrafficd() (*Trafficd, error) {
	// Get machine id.
	b, err := os.ReadFile("/etc/machine-id")
	if err != nil {
		return nil, err
	}

	// NOTE: machine id should be treated as sensitive data and therefore hashed.
	// https://man7.org/linux/man-pages/man5/machine-id.5.html
	h := sha256.New()
	h.Write(b)
	bs := h.Sum(nil)

	return &Trafficd{
		MachineId: hex.EncodeToString(bs),
	}, nil
}

func (trf *Trafficd) GetEventdAddr(sd discovery.ServiceDiscoverer) (*url.URL, error) {
	svcs, err := sd.Lookup(context.Background(), "eventd")
	if err != nil {
		return nil, err
	}

	if len(svcs) == 0 {
		return nil, fmt.Errorf("no eventd services found")
	}

	return url.Parse(fmt.Sprintf("%s:%d", svcs[0].Address, svcs[0].GrpcPort))
}

func (trf *Trafficd) RegisterHost(ctx context.Context, client apiv1.EventServiceClient) error {
	host, err := os.Hostname()
	if err != nil {
		return err
	}

	req := &apiv1.RegisterHostRequest{
		Host: &apiv1.Host{
			MachineId: trf.MachineId,
			Hostname:  host,
		},
		Timestamp: time.Now().Unix(),
	}

	_, err = client.RegisterHost(ctx, req)
	return err
}

func (trf *Trafficd) StreamEvents(ctx context.Context, client apiv1.EventServiceClient) error {
	// Cancelling the contxt here would result in an error returned from the server,
	// so we pass a new context and handle the stream close later on.
	stream, err := client.SendEvent(context.Background())
	if err != nil {
		return fmt.Errorf("error getting stream: %w", err)
	}

	t := time.NewTicker(5 * time.Second)

Loop:
	for {
		select {
		case <-t.C:
			req, err := GenerateRandomEvent(trf.MachineId)
			if err != nil {
				log.Error().Err(err).Msg("error generating event")
				continue
			}

			if err := stream.Send(req); err != nil {
				log.Error().Err(err).Msg("error sending event")
				continue
			}
		case <-ctx.Done():
			break Loop
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("error closing stream: %w", err)
	}

	log.Info().Msgf("total successful events sent: %d", res.EventsTotal)
	return nil
}
