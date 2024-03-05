package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/gob"
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
	stream, err := client.SendEvent(context.Background())
	if err != nil {
		return fmt.Errorf("error getting stream: %w", err)
	}

	t := time.NewTicker(5 * time.Second)

Loop:
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
				MachineId: trf.MachineId,
				Data:      buf.Bytes(),
			}
			if err := stream.Send(req); err != nil {
				log.Error().Err(err).Msg("error sending event")
				continue
			}
		case <-ctx.Done():
			break Loop
		}
	}

	res, _ := stream.CloseAndRecv()
	log.Info().Msgf("total successful events sent: %d", res.EventsTotal)
	return nil
}
