package main

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apiv1 "github.com/loshz/platform/internal/api/v1"
)

// MsgMissingRequiredField represents an error message format for
// missing request fields.
var MsgMissingRequiredField = "error: missing required '%s' field"

// Services represents a map of individually registered services keyed by the
// service uuid.
type Services map[string]*apiv1.Service

type DiscoveryServer struct {
	apiv1.UnimplementedDiscoveryServiceServer

	mtx      sync.RWMutex
	services Services
}

func NewDiscoveryServer() *DiscoveryServer {
	return &DiscoveryServer{
		services: make(Services),
	}
}

// EvictExpiredServices removes services that have a registration timestamp greater
// then 5 minutes.
func (ds *DiscoveryServer) EvictExpiredServices() {
	ds.mtx.Lock()
	defer ds.mtx.Unlock()

	// Loop through all services and check if the current timestamp is
	// after the expiry threshold.
	for uuid, svc := range ds.services {
		expired := time.Now().Add(-5 * time.Minute)
		if time.Unix(svc.LastSeen, 0).Before(expired) {
			log.Info().Msgf("expired service evicted: %s", uuid)
			delete(ds.services, uuid)
		}
	}
}

func (ds *DiscoveryServer) StartEvictionProcess(ctx context.Context) {
	log.Info().Msg("polling for expired services every 60s")
	t := time.NewTicker(60 * time.Second)
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
func (ds *DiscoveryServer) RegisterService(_ context.Context, req *apiv1.RegisterServiceRequest) (*apiv1.RegisterServiceResponse, error) {
	svc := req.GetService()
	if svc == nil {
		return nil, status.Errorf(codes.InvalidArgument, MsgMissingRequiredField, "service")
	}

	uuid := svc.GetUuid()
	if uuid == "" {
		return nil, status.Errorf(codes.InvalidArgument, MsgMissingRequiredField, "uuid")
	}

	ds.mtx.Lock()
	ds.services[uuid] = svc
	ds.mtx.Unlock()

	log.Info().Msgf("service registered: %s", uuid)

	return &apiv1.RegisterServiceResponse{
		Service: svc,
	}, nil
}

// RegisterService deletes a service from the DiscoveryServer.
func (ds *DiscoveryServer) DeregisterService(_ context.Context, req *apiv1.DeregisterServiceRequest) (*apiv1.DeregisterServiceResponse, error) {
	uuid := req.GetUuid()
	if uuid == "" {
		return nil, status.Errorf(codes.InvalidArgument, MsgMissingRequiredField, "uuid")
	}

	ds.mtx.Lock()
	delete(ds.services, uuid)
	ds.mtx.Unlock()

	log.Info().Msgf("service deregistered: %s", uuid)

	return &apiv1.DeregisterServiceResponse{
		Uuid: uuid,
	}, nil
}

// GetServices returns all currently registered services with a given prefix.
func (ds *DiscoveryServer) GetServices(_ context.Context, req *apiv1.GetServicesRequest) (*apiv1.GetServicesResponse, error) {
	services := []*apiv1.Service{}
	name := req.GetName()

	ds.mtx.RLock()
	defer ds.mtx.RUnlock()

	for uuid, svc := range ds.services {
		if name == "" || strings.HasPrefix(uuid, name) {
			services = append(services, svc)
		}
	}

	return &apiv1.GetServicesResponse{
		Services: services,
	}, nil
}
