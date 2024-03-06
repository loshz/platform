package main

import (
	"context"
	"io"

	"github.com/rs/zerolog/log"

	apiv1 "github.com/loshz/platform/internal/api/v1"
	pgrpc "github.com/loshz/platform/internal/grpc"
)

type GrpcServer struct {
	apiv1.UnimplementedEventServiceServer
}

func (s *GrpcServer) RegisterHost(_ context.Context, req *apiv1.RegisterHostRequest) (*apiv1.RegisterHostResponse, error) {
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

	// TODO: store machine details in db.

	return &apiv1.RegisterHostResponse{
		MachineId: machineId,
	}, nil
}

func (s *GrpcServer) SendEvent(stream apiv1.EventService_SendEventServer) error {
	res := new(apiv1.SendEventResponse)

	for {
		event, err := stream.Recv()
		if err == io.EOF || stream.Context().Err() == context.Canceled {
			break
		}
		if err != nil {
			log.Error().Err(err).Msg("stream receive error")
			return err
		}

		// TODO: validate event.

		log.Info().Msgf("event received, type: %s", event.Type)
		EventsTotal.WithLabelValues(event.Type.String()).Inc()
		res.MachineId = event.MachineId // TODO: doing this per event is inefficient.
		res.EventsTotal++

		// TODO: process event onto queue.
	}

	log.Info().Str("machine_id", res.MachineId).Msg("stream closed")
	return stream.SendAndClose(res)
}
