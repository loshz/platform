package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	pbv1 "github.com/loshz/platform/pkg/pb/v1"
)

type StatusServer struct {
	pbv1.UnimplementedPlatformServiceServer

	service string
}

func (s *StatusServer) Status(context.Context, *emptypb.Empty) (*pbv1.PlatformServiceStatusResponse, error) {
	return &pbv1.PlatformServiceStatusResponse{
		Service: s.service,
		Status:  pbv1.PlatformServiceStatus_STATUS_OK,
	}, nil
}
