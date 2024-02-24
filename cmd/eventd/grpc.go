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

	return &apiv1.RegisterHostResponse{}, nil
}

func (s *grpcServer) Send(stream apiv1.EventService_SendEventServer) error {
	for {
		event, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		log.Info().Msgf("event received, type: %s", event.Type)
	}

	return stream.SendAndClose(&apiv1.SendEventResponse{})
}
