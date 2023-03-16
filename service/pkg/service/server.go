package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	pbv1 "github.com/loshz/platform/pkg/pb/v1"
)

type grpcServer struct {
	pbv1.UnimplementedPlatformServiceServer

	service string
}

func (s *grpcServer) Status(context.Context, *emptypb.Empty) (*pbv1.StatusResponse, error) {
	return &pbv1.StatusResponse{
		Service: s.service,
		Status:  pbv1.Status_STATUS_OK,
	}, nil
}
