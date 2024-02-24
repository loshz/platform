package main

import (
	"io"

	"github.com/rs/zerolog/log"

	apiv1 "github.com/loshz/platform/internal/api/v1"
)

type grpcServer struct {
	apiv1.UnimplementedEventServiceServer
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
