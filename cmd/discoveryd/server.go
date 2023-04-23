package main

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbv1 "github.com/loshz/platform/pkg/pb/v1"
)

// MsgMissingRequiredField represents an error message format for
// missing request fields.
var MsgMissingRequiredField = "error: missing required '%s' field"

type Services map[string]*pbv1.Service

type DiscoveryServer struct {
	pbv1.UnimplementedDiscoveryServiceServer

	mtx      sync.RWMutex
	services Services
}

func NewDiscoveryServer() *DiscoveryServer {
	return &DiscoveryServer{
		services: make(Services),
	}
}

func (ds *DiscoveryServer) EvictExpiredServices() {
	ds.mtx.Lock()
	defer ds.mtx.Unlock()

	// Loop through all all services and check if the current timestamp is
	// after the expiry threshold.
	for uuid, svc := range ds.services {
		expired := time.Now().Add(-5 * time.Minute)
		if time.Unix(svc.Timestamp, 0).Before(expired) {
			log.Info().Str("uuid", uuid).Msg("evicting expired service")
			delete(ds.services, uuid)
		}
	}
}

func (ds *DiscoveryServer) StartEvictionProcess(ctx context.Context) {
	log.Info().Msg("starting service eviction process every 30s")

	t := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-t.C:
			ds.EvictExpiredServices()
		case <-ctx.Done():
			return
		}
	}
}

// RegisterService validates service data and stores it in the DiscoveryServer.
func (ds *DiscoveryServer) RegisterService(_ context.Context, req *pbv1.RegisterServiceRequest) (*pbv1.RegisterServiceResponse, error) {
	svc := req.GetService()
	if svc == nil {
		return nil, status.Errorf(codes.InvalidArgument, MsgMissingRequiredField, "service")
	}

	if svc.GetUuid() == "" {
		return nil, status.Errorf(codes.InvalidArgument, MsgMissingRequiredField, "uuid")
	}

	ds.mtx.Lock()
	ds.services[svc.GetUuid()] = svc
	ds.mtx.Unlock()

	return &pbv1.RegisterServiceResponse{
		Service: svc,
	}, nil
}

// RegisterService deletes a service from the DiscoveryServer.
func (ds *DiscoveryServer) DeregisterService(_ context.Context, req *pbv1.DeregisterServiceRequest) (*pbv1.DeregisterServiceResponse, error) {
	ds.mtx.Lock()
	delete(ds.services, req.GetUuid())
	ds.mtx.Unlock()

	return &pbv1.DeregisterServiceResponse{
		Uuid: req.GetUuid(),
	}, nil
}

// GetService returns all currently registered services with a given prefix.
func (ds *DiscoveryServer) GetService(_ context.Context, req *pbv1.GetServiceRequest) (*pbv1.GetServiceResponse, error) {
	if req.GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, MsgMissingRequiredField, "name")
	}

	ds.mtx.RLock()
	defer ds.mtx.RUnlock()

	var services []*pbv1.Service

	for uuid, svc := range ds.services {
		if strings.HasPrefix(uuid, req.GetName()) {
			services = append(services, svc)
		}
	}

	return &pbv1.GetServiceResponse{
		Services: services,
	}, nil
}