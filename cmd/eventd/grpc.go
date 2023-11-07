package main

import (
	"context"

	"github.com/rs/zerolog/log"

	pbv1 "github.com/loshz/platform/internal/api/v1"
)

type grpcServer struct {
	pbv1.UnimplementedEventServiceServer
}

func (s *grpcServer) Event(ctx context.Context, req *pbv1.EventRequest) (*pbv1.EventResponse, error) {
	log.Info().Str("hostname", req.Hostname).Msg("request received")

	return &pbv1.EventResponse{
		Uuid: "test",
	}, nil
}
