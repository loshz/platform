// Package service implements the base structure of a platform application.
package service

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/loshz/platform/pkg/config"
	plog "github.com/loshz/platform/pkg/log"
	"github.com/loshz/platform/pkg/metrics"
	pbv1 "github.com/loshz/platform/pkg/pb/v1"
	"github.com/loshz/platform/pkg/version"
)

const (
	ExitOk = iota
	ExitError
)

// Service represents a platform application.
type Service struct {
	// Explicit service name.
	Name string

	// UUID of an individual service including Name prefix.
	ID string

	// Service configuration.
	Config *config.Config

	// Service specific exit handlers.
	exitHandlers    []ExitHandler
	exitHandlersMtx sync.RWMutex
}

// New creates a named Service with configurable dependencies.
func New(name string) *Service {
	// TODO: init deps here- config, DBs, etc.

	return &Service{
		Name:   name,
		ID:     fmt.Sprintf("%s-%s", name, uuid.New()),
		Config: config.New(name),
	}
}

// RunFunc is a function that will be called by Run to initialize a service.
// If this function returns an error then the server will immediately shut down.
//
// NOTE: the contents of this function should run in a goroutine in order for it
// not to block.
type RunFunc func(*Service) error

// Run starts the Service and ensures all dependencies are initialised.
//
// By default, it will start the local web server and wait for a stop
// signal to be received before attempting to gracefully shutdown.
func (s *Service) Run(run RunFunc) {
	// Initialize required service config.
	s.initDefaultConfig()

	// Configure global logger.
	plog.ConfigureGlobalLogging(s.Config.String(config.KeyLogLevel), s.ID, version.Build)

	// Attempt to start the service.
	s.start(run)

	// Start the local http server.
	s.serveHTTP(s.Config.Int(config.KeyHTTPPort))

	// Wait for an exit signal.
	s.waitSignal()

	// Attempt to gracefully shutdown.
	s.Exit(ExitOk)
}

// waitSignal blocks waiting for operating system signals.
//
// By default, it will handle calls to SIGINT and SIGTERM.
func (s *Service) waitSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signal to be received.
	<-ch
	log.Info().Msg("stop signal received, starting shut down")
}

// Exit the current process while attempting to run any of the Service's exit handlers.
//
// Each exit handler will run in its own goroutine and will force exit after 30s
// regardless of completion.
func (s *Service) Exit(status int) {
	exitWait := 30 * time.Second
	deadline := time.Now().Add(exitWait)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// Force exit after deadline.
	time.AfterFunc(exitWait, func() {
		log.Error().Msgf("service shutdown took longer than %s, aborting", exitWait)
		os.Exit(status)
	})

	s.exitHandlersMtx.RLock()
	defer s.exitHandlersMtx.RUnlock()

	// Run each exit handler in its own goroutine.
	var wg sync.WaitGroup
	for _, fn := range s.exitHandlers {
		wg.Add(1)
		go func(fn ExitHandler) {
			defer wg.Done()
			fn(ctx)
		}(fn)
	}

	wg.Wait()
	os.Exit(status)
}

// serveHTTP configures and starts the local webserver.
//
// By default, it will register pprof, metrics and health endpoints.
func (s *Service) serveHTTP(port int) {
	router := http.NewServeMux()

	// Configure debug endpoints.
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// Expose the registered metrics via HTTP.
	router.Handle("/metrics", promhttp.Handler())

	// Expose health check
	// TODO: define service dependencies
	router.HandleFunc("/health", healthHandler(s.ID, nil))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Error().Err(err).Msg("local http server error")
			s.Exit(ExitError)
		}
	}()

	s.AddExitHandler(func(ctx context.Context) {
		log.Info().Msg("stopping http server")
		if err := srv.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("error shutting down http server")
		}
	})

	log.Info().Msgf("local http server running on :%d", port)
}

// ServeGRPC configures, registers services and starts a gRPC server on a given port.
// It is intentially not called automatically in Start() as not every service requires a gRPC server,
// therefore it should be called directly by the service itself.
func (s *Service) ServeGRPC(port int, desc *grpc.ServiceDesc, svc interface{}) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Error().Err(err).Msg("error creating tcp listener for grpc server")
		s.Exit(ExitError)
	}

	// Configure server and register services.
	// TODO: pass server opts/timeouts/etc.
	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(streamInterceptor(s.ID)),
		grpc.UnaryInterceptor(unaryInterceptor(s.ID)),
	}
	srv := grpc.NewServer(opts...)
	srv.RegisterService(&pbv1.PlatformService_ServiceDesc, &grpcServer{})
	srv.RegisterService(desc, svc)

	go func() {
		if err := srv.Serve(lis); err != grpc.ErrServerStopped {
			log.Error().Err(err).Msg("grpc server error")
			s.Exit(ExitError)
		}
	}()

	s.AddExitHandler(func(ctx context.Context) {
		log.Info().Msg("stopping grpc server")
		srv.GracefulStop()
	})

	log.Info().Msgf("grpc server running on :%d", port)
}

// start attempts to run the service with an initial timeout.
// If the deadline exceeds the time taken to run the service, it is treated
// as a failed start.
func (s *Service) start(run RunFunc) {
	// Configure startup timeout.
	timeout := s.Config.Duration(config.KeyServiceStartupTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Wait for the context to exceed the given deadline, and exit if so.
	go func() {
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			log.Error().Msgf("service startup took longer than the configured %s timeout, aborting", timeout)
			s.Exit(ExitError)
		}
	}()

	// Attempt to run the main service func.
	if err := run(s); err != nil {
		log.Error().Err(err).Msg("error running service")
		s.Exit(ExitError)
	}

	metrics.ServiceInfo.WithLabelValues(s.Name, s.ID, version.Build).Inc()
	log.Info().Msg("service started")
}

// ExitHandler is a timed function ran only on service shutdown.
type ExitHandler func(context.Context)

// AddExitHandler registers a function that should be called when the server is shutting down.
func (s *Service) AddExitHandler(fn ExitHandler) {
	s.exitHandlersMtx.Lock()
	s.exitHandlers = append(s.exitHandlers, fn)
	s.exitHandlersMtx.Unlock()
}

func (s *Service) initDefaultConfig() {
	s.Config.MustLoad(config.KeyLogLevel, "info", true, config.ParseLogLevel)
	s.Config.MustLoad(config.KeyServiceStartupTimeout, "5s", true, config.ParseDuration)
	s.Config.MustLoad(config.KeyHTTPPort, 8001, true, config.ParseInt)
}
