package main

import (
	"context"

	"github.com/rs/zerolog/log"

	apiv1 "github.com/loshz/platform/internal/api/v1"
)

type grpcServer struct {
	apiv1.UnimplementedEventServiceServer
}

func (s *grpcServer) Event(ctx context.Context, req *apiv1.EventRequest) (*apiv1.EventResponse, error) {
	log.Info().Str("hostname", req.Hostname).Msg("request received")

	return &apiv1.EventResponse{
		Uuid: "test",
	}, nil
}
