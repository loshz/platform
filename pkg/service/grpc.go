package service

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/loshz/platform/pkg/config"
	"github.com/loshz/platform/pkg/metrics"
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

// ServeGRPC configures, registers services and starts a gRPC server on a given port.
// It is intentially not called in Start() as not every service requires a gRPC server,
// therefore it should be called directly by the service itself.
func (s *Service) ServeGRPC(desc *grpc.ServiceDesc, svc interface{}) {
	port := fmt.Sprintf(":%d", s.Config.Int(config.KeyGRPCServerPort))

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Error().Err(err).Msg("error creating tcp listener for grpc server")
		s.Exit(ExitError)
	}

	// Load server TLS credentials.
	creds, err := credentials.NewServerTLSFromFile(s.Config.String(config.KeyGRPCServerCert), s.Config.String(config.KeyGRPCServerKey))
	if err != nil {
		log.Error().Err(err).Msg("error loading grpc server tls credentials")
		s.Exit(ExitError)
	}

	opts := []grpc.ServerOption{
		grpc.Creds(creds),
		grpc.UnaryInterceptor(unaryInterceptor(s.ID)),
		grpc.StreamInterceptor(streamInterceptor(s.ID)),
		grpc.ConnectionTimeout(s.Config.Duration(config.KeyGRPCServerConnTimeout)),
	}

	// Configure server and register services.
	srv := grpc.NewServer(opts...)
	srv.RegisterService(&pbv1.PlatformService_ServiceDesc, &grpcServer{})
	srv.RegisterService(desc, svc)

	go func() {
		log.Info().Msgf("grpc server running on %s", port)
		if err := srv.Serve(lis); err != grpc.ErrServerStopped {
			log.Error().Err(err).Msg("grpc server error")
			s.Exit(ExitError)
		}
	}()

	// Wait for service to exit and shutdown.
	<-s.ctx.Done()
	log.Info().Msg("stopping grpc server")
	srv.GracefulStop()
}
