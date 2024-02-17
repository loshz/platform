package discovery

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	apiv1 "github.com/loshz/platform/internal/api/v1"
)

type MockDiscoveryServiceClient struct {
	RegisterFunc    func() (*apiv1.RegisterServiceResponse, error)
	DeregisterFunc  func() (*apiv1.DeregisterServiceResponse, error)
	GetServicesFunc func() (*apiv1.GetServicesResponse, error)
}

func (m *MockDiscoveryServiceClient) RegisterService(context.Context, *apiv1.RegisterServiceRequest, ...grpc.CallOption) (*apiv1.RegisterServiceResponse, error) {
	return m.RegisterFunc()
}

func (m *MockDiscoveryServiceClient) DeregisterService(context.Context, *apiv1.DeregisterServiceRequest, ...grpc.CallOption) (*apiv1.DeregisterServiceResponse, error) {
	return m.DeregisterFunc()
}

func (m *MockDiscoveryServiceClient) GetServices(context.Context, *apiv1.GetServicesRequest, ...grpc.CallOption) (*apiv1.GetServicesResponse, error) {
	return m.GetServicesFunc()
}

func TestRegister(t *testing.T) {
	t.Run("TestError", func(t *testing.T) {
		expected := errors.New("register error")
		svc := new(Service)
		svc.client = &MockDiscoveryServiceClient{
			RegisterFunc: func() (*apiv1.RegisterServiceResponse, error) { return nil, expected },
		}

		err := svc.Register(context.Background(), nil)
		require.ErrorContains(t, err, expected.Error())
	})

	t.Run("TestSuccess", func(t *testing.T) {
		svc := new(Service)
		svc.client = &MockDiscoveryServiceClient{
			RegisterFunc: func() (*apiv1.RegisterServiceResponse, error) { return nil, nil },
		}

		err := svc.Register(context.Background(), nil)
		require.NoError(t, err)
	})
}

func TestDeregister(t *testing.T) {
	t.Run("TestError", func(t *testing.T) {
		expected := errors.New("deregister error")
		svc := new(Service)
		svc.client = &MockDiscoveryServiceClient{
			DeregisterFunc: func() (*apiv1.DeregisterServiceResponse, error) { return nil, expected },
		}

		err := svc.Deregister(context.Background(), "service_id")
		require.ErrorContains(t, err, expected.Error())
	})

	t.Run("TestSuccess", func(t *testing.T) {
		svc := new(Service)
		svc.client = &MockDiscoveryServiceClient{
			DeregisterFunc: func() (*apiv1.DeregisterServiceResponse, error) { return nil, nil },
		}

		err := svc.Deregister(context.Background(), "service_id")
		require.NoError(t, err)
	})
}

func TestLookup(t *testing.T) {
	t.Run("TestError", func(t *testing.T) {
		expected := errors.New("lookup error")
		svc := new(Service)
		svc.client = &MockDiscoveryServiceClient{
			GetServicesFunc: func() (*apiv1.GetServicesResponse, error) { return nil, expected },
		}

		svcs, err := svc.Lookup(context.Background(), "service_id")
		require.ErrorContains(t, err, expected.Error())
		require.Nil(t, svcs)
	})

	t.Run("TestSuccess", func(t *testing.T) {
		svc := new(Service)
		svc.client = &MockDiscoveryServiceClient{
			GetServicesFunc: func() (*apiv1.GetServicesResponse, error) {
				return &apiv1.GetServicesResponse{Services: []*apiv1.Service{{Uuid: "service_id"}}}, nil
			},
		}

		svcs, err := svc.Lookup(context.Background(), "service_id")
		require.NoError(t, err)
		require.NotNil(t, svcs)
	})
}
