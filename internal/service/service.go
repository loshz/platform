// Package service implements the base structure of a platform application.
package service

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/loshz/platform/internal/config"
	"github.com/loshz/platform/internal/credentials"
	plog "github.com/loshz/platform/internal/log"
	"github.com/loshz/platform/internal/metrics"
	"github.com/loshz/platform/internal/uuid"
	"github.com/loshz/platform/internal/version"
)

const (
	ExitOK = iota
	ExitError
	ExitStartup
)

// RunFunc is a function that will be called by Run to initialize a service.
// If this function returns an error then the server will immediately shut down.
type RunFunc func(context.Context, *Service) error

// Service represents a platform application.
type Service struct {
	// Service configuration.
	Config *config.Config

	// UUID of the individual service including name prefix.
	// E.g., service-xxxx-xxxx
	id uuid.UUID

	// Channel for sending/receiving internal service errors.
	errCh chan error

	// Store the current leadership status.
	leader atomic.Bool

	// Service for storing credentials.
	creds *credentials.Store
}

// New creates a named Service with configurable dependencies.
func New(name string) *Service {
	return &Service{
		Config: config.New(),
		id:     uuid.New(name),
		errCh:  make(chan error),
		creds:  credentials.New(),
	}
}

// Service getter methods.
func (s *Service) Creds() *credentials.Store { return s.creds }
func (s *Service) ID() uuid.UUID             { return s.id }
func (s *Service) IsLeader() bool            { return s.leader.Load() }
func (s *Service) Name() string              { return s.id.Name() }

// Run starts the Service and ensures all dependencies are initialised.
//
// By default, it will start the local web server and wait for a stop
// signal to be received before attempting to gracefully shutdown.
func (s *Service) Run(run RunFunc) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	// Initialize required service config.
	s.LoadRequiredConfig()

	// Configure global logger.
	plog.ConfigureGlobalLogging(s.Config.String(config.KeyServiceLogLevel), s.ID().String(), version.Build)

	// Attempt to start the service.
	if err := s.start(ctx, run); err != nil {
		log.Error().Err(err).Msg("service startup error")
		cancel()
		s.Exit(ExitStartup)
	}

	// Register service for discovery.
	go s.RegisterDiscovery(ctx)

	// Start the local http server.
	go s.serveHTTP(ctx)

	// Wait for an exit signal or service error.
	status := s.waitSignal(ctx)
	cancel()

	// Attempt to gracefully shutdown.
	s.Exit(status)
}

// Error sends a given error to the Service's error channel.
// Services should prefer calling this method instead of manually calling s.Exit()
// so shutdown can be handled gracefully.
func (s *Service) Error(err error) { s.errCh <- err }

// Exit cancels the service's context in order to signal a shutdown to child processes.
// It sleeps for a configurable time before signalling the process to exit.
func (s *Service) Exit(status int) {
	// Exit early if startup error.
	if status == ExitStartup {
		os.Exit(status)
	}

	// Force exit after deadline.
	time.AfterFunc(s.Config.Duration(config.KeyServiceShutdownTimeout), func() {
		log.Error().Msg("service shutdown timeout expired")
		os.Exit(status)
	})

	if err := s.DeregisterDiscovery(); err != nil {
		log.Error().Err(err).Msg("error deregistering service from discovery")
	}

	os.Exit(status)
}

// waitSignal blocks waiting for operating system signals or an internal
// service error.
func (s *Service) waitSignal(ctx context.Context) int {
	// Wait for signal to be received.
	select {
	case <-ctx.Done():
		log.Info().Msg("stop signal received, starting shut down")
	case err := <-s.errCh:
		// Always return an error as nothing should be sending nil
		// to this channel.
		log.Error().Err(err).Msg("internal service error")
		return ExitError
	}

	return ExitOK
}

// start attempts to run the service with an initial timeout.
// If the deadline exceeds the time taken to run the service, it is treated
// as a failed start.
func (s *Service) start(ctx context.Context, run RunFunc) error {
	// Configure startup timeout.
	timeout := s.Config.Duration(config.KeyServiceStartupTimeout)

	// Attempt to run the main service func and record error.
	errCh := make(chan error, 1)
	go func() {
		errCh <- run(ctx, s)
	}()

	select {
	case <-time.After(timeout):
		return fmt.Errorf("startup deadline (%s) exceeded", timeout)
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("error running service: %w", err)
		}
	}

	metrics.ServiceInfo.WithLabelValues(s.ID().String(), version.Build).Inc()
	log.Info().Msg("service started")

	return nil
}
