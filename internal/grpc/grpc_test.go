package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/loshz/platform/internal/uuid"
)

func TestStreamInterceptor(t *testing.T) {
	expected := errors.New("handler error")

	// Mock server info.
	info := &grpc.StreamServerInfo{
		FullMethod: "test",
	}

	// Mock a handler to return an error.
	handler := func(srv interface{}, stream grpc.ServerStream) error {
		return expected
	}

	// Create a new interceptor.
	interceptor := StreamInterceptor(uuid.New("stream_service"))
	err := interceptor(nil, nil, info, handler)

	// Assert that the error from the handler is returned from the
	// interceptor.
	assert.ErrorIs(t, err, expected)
}

func TestUnaryInterceptor(t *testing.T) {
	expected := errors.New("handler error")

	// Mock server info.
	info := &grpc.UnaryServerInfo{
		FullMethod: "test",
	}

	// Mock a handler to return an error.
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return 1, expected
	}

	// Create a new interceptor.
	interceptor := UnaryInterceptor(uuid.New("stream_service"))
	res, err := interceptor(context.Background(), nil, info, handler)

	// Assert that the response and error from the handler are returned from the
	// interceptor.
	assert.Equal(t, res, 1)
	assert.ErrorIs(t, err, expected)
}
