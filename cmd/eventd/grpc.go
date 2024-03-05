package main

import (
	"context"
	"io"

	"github.com/rs/zerolog/log"

	apiv1 "github.com/loshz/platform/internal/api/v1"
	pgrpc "github.com/loshz/platform/internal/grpc"
)

type grpcServer struct {
	apiv1.UnimplementedEventServiceServer
}

func (s *grpcServer) RegisterHost(_ context.Context, req *apiv1.RegisterHostRequest) (*apiv1.RegisterHostResponse, error) {
	host := req.GetHost()
	if host == nil {
		return nil, pgrpc.ErrMissingRequiredField("host")
	}
	machineId := host.GetMachineId()
	if machineId == "" {
		return nil, pgrpc.ErrMissingRequiredField("host.machine_id")
	}
	hostname := host.GetHostname()
	if hostname == "" {
		return nil, pgrpc.ErrMissingRequiredField("host.hostname")
	}

	log.Info().Str("machine_id", machineId).Msg("machine registered")

	return &apiv1.RegisterHostResponse{
		MachineId: machineId,
	}, nil
}

func (s *grpcServer) SendEvent(stream apiv1.EventService_SendEventServer) error {
	res := &apiv1.SendEventResponse{}

	for {
		event, err := stream.Recv()
		if err == io.EOF || stream.Context().Err() == context.Canceled {
			break
		}
		if err != nil {
			log.Error().Err(err).Msg("stream error")
			return err
		}

		log.Info().Msgf("event received, type: %s", event.Type)
		EventsTotal.WithLabelValues(event.Type.String()).Inc()
		res.MachineId = event.MachineId // TODO: this is inefficient.
		res.EventsTotal++
	}

	// TODO: update res with machine id.
	log.Info().Str("machine_id", res.MachineId).Msg("stream closed")
	return stream.SendAndClose(res)
}
