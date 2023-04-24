package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbv1 "github.com/loshz/platform/pkg/pb/v1"
)

func TestEvictExpiredServices(t *testing.T) {
	server := NewDiscoveryServer()

	// Manually register services with the server.
	server.services["expired-service-a"] = &pbv1.Service{
		Timestamp: time.Now().Add(-1 * time.Hour).Unix(),
	}
	server.services["expired-service-b"] = &pbv1.Service{
		Timestamp: time.Now().Add(-1 * time.Hour).Unix(),
	}
	server.services["service-a"] = &pbv1.Service{
		Timestamp: time.Now().Unix(),
	}

	server.EvictExpiredServices()

	// Assert expired services have been evicted.
	assert.Equal(t, 1, len(server.services))
	assert.Nil(t, server.services["expired-service-a"])
	assert.Nil(t, server.services["expired-service-b"])
	assert.NotNil(t, server.services["service-a"])
}

func TestRegisterService(t *testing.T) {
	server := NewDiscoveryServer()

	t.Run("TestNilService", func(t *testing.T) {
		// Create an empty request and attempt service registration.
		req := &pbv1.RegisterServiceRequest{}
		_, err := server.RegisterService(context.TODO(), req)

		// Get the error status.
		stat, _ := status.FromError(err)

		// Assert the returned error is due to an invalid argument.
		assert.NotNil(t, err)
		assert.Equal(t, codes.InvalidArgument, stat.Code())
		assert.Equal(t, fmt.Sprintf(MsgMissingRequiredField, "service"), stat.Message())
	})

	t.Run("TestNilUuid", func(t *testing.T) {
		// Create a request with an empty uuid and attempt service registration.
		req := &pbv1.RegisterServiceRequest{
			Service: &pbv1.Service{
				Uuid: "",
			},
		}
		_, err := server.RegisterService(context.TODO(), req)

		// Get the error status.
		stat, _ := status.FromError(err)

		// Assert the returned error is due to an invalid argument.
		assert.NotNil(t, err)
		assert.Equal(t, codes.InvalidArgument, stat.Code())
		assert.Equal(t, fmt.Sprintf(MsgMissingRequiredField, "uuid"), stat.Message())
	})

	t.Run("TestSuccess", func(t *testing.T) {
		// Create a valid request and attempt service registration.
		svc := &pbv1.Service{
			Uuid:      "test-service",
			HttpPort:  8001,
			GrpcPort:  8002,
			Timestamp: time.Now().Unix(),
		}
		req := &pbv1.RegisterServiceRequest{
			Service: svc,
		}
		res, err := server.RegisterService(context.TODO(), req)

		// Assert the returned error is nil and the status code is OK.
		assert.Nil(t, err)
		assert.Equal(t, status.Code(err), codes.OK)
		assert.Equal(t, svc, res.GetService())

		// Assert the service was written to the server.
		assert.Equal(t, server.services[svc.Uuid], svc)
	})
}

func TestDeregisterService(t *testing.T) {
	server := NewDiscoveryServer()

	t.Run("TestNilUuid", func(t *testing.T) {
		// Create a request with an empty uuid and attempt service deregistration.
		req := &pbv1.DeregisterServiceRequest{
			Uuid: "",
		}
		_, err := server.DeregisterService(context.TODO(), req)

		// Get the error status.
		stat, _ := status.FromError(err)

		// Assert the returned error is due to an invalid argument.
		assert.NotNil(t, err)
		assert.Equal(t, codes.InvalidArgument, stat.Code())
		assert.Equal(t, fmt.Sprintf(MsgMissingRequiredField, "uuid"), stat.Message())
	})

	t.Run("TestSuccess", func(t *testing.T) {
		// Create a valid service and manually register with server.
		uuid := "test-service"
		server.services[uuid] = &pbv1.Service{}
		req := &pbv1.DeregisterServiceRequest{
			Uuid: uuid,
		}
		res, err := server.DeregisterService(context.TODO(), req)

		// Assert the returned error is nil and the status code is OK.
		assert.Nil(t, err)
		assert.Equal(t, codes.OK, status.Code(err))
		assert.Equal(t, uuid, res.GetUuid())

		// Assert the service was deleted from the server.
		assert.Nil(t, server.services[uuid])
	})
}

func TestGetService(t *testing.T) {
	server := NewDiscoveryServer()

	t.Run("TestNilName", func(t *testing.T) {
		// Create a request with an empty name and attempt to get services.
		req := &pbv1.GetServiceRequest{
			Name: "",
		}
		_, err := server.GetService(context.TODO(), req)

		// Get the error status.
		stat, _ := status.FromError(err)

		// Assert the returned error is due to an invalid argument.
		assert.NotNil(t, err)
		assert.Equal(t, codes.InvalidArgument, stat.Code())
		assert.Equal(t, fmt.Sprintf(MsgMissingRequiredField, "name"), stat.Message())
	})

	t.Run("TestSuccess", func(t *testing.T) {
		// Manually register services with the server.
		server.services["test-service-a"] = &pbv1.Service{}
		server.services["test-service-b"] = &pbv1.Service{}
		server.services["service-a"] = &pbv1.Service{}

		// Create a valid service and manually register with server.
		req := &pbv1.GetServiceRequest{
			Name: "test-service",
		}
		res, err := server.GetService(context.TODO(), req)

		// Assert the returned error is nil and the status code is OK.
		assert.Nil(t, err)
		assert.Equal(t, codes.OK, status.Code(err))

		// Assert the expected service were returned.
		assert.Equal(t, 2, len(res.Services))
	})
}
