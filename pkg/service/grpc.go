package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	pbv1 "github.com/loshz/platform/pkg/pb/v1"
)

func (s *Service) Status(context.Context, *emptypb.Empty) (*pbv1.PlatformServiceStatusResponse, error) {
	return &pbv1.PlatformServiceStatusResponse{
		Service: s.ID,
		Status:  pbv1.PlatformServiceStatus_STATUS_OK,
		Leader:  s.IsLeader(),
	}, nil
}
