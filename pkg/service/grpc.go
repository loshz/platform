package service

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/loshz/platform/pkg/metrics"
)

// streamInterceptor instruments and logs information about gRPC stream calls.
func streamInterceptor(service string) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Time the underlying request.
		start := time.Now()
		err := handler(srv, ss)
		latency := time.Since(start)

		// Get the request status code.
		code := status.Code(err)

		// Record request metrics.
		labels := []string{service, code.String(), info.FullMethod, "stream"}
		metrics.GRPCRequestDuration.WithLabelValues(labels...).Observe(latency.Seconds())
		metrics.GRPCRequestsTotal.WithLabelValues(labels...).Inc()

		return err
	}
}

// unaryInterceptor instruments and logs information about gRPC unary calls.
func unaryInterceptor(service string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		res, err := handler(ctx, req)
		latency := time.Since(start)

		// Get the request status code.
		code := status.Code(err)

		// Record request metrics.
		labels := []string{"TODO", code.String(), info.FullMethod, "unary"}
		metrics.GRPCRequestDuration.WithLabelValues(labels...).Observe(latency.Seconds())
		metrics.GRPCRequestsTotal.WithLabelValues(labels...).Inc()

		return res, err
	}
}
