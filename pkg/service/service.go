// Package service implements the base structure of a platform application.
package service

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/loshz/platform/pkg/config"
	plog "github.com/loshz/platform/pkg/log"
	"github.com/loshz/platform/pkg/metrics"
	"github.com/loshz/platform/pkg/version"
)

const (
	ExitOk = iota
	ExitError
)

// Service represents a platform application.
type Service struct {
	// Service configuration.
	Config *config.Config

	// UUID of the individual service including name prefix.
	// E.g., service-xxxx-xxxx
	id string

	// Global context used to signal service shutdown to spawned
	// goroutines.
	ctx       context.Context
	ctxCancel context.CancelFunc

	// Channel for sending/receiving internal service errors.
	errCh chan error

	// Store the current leadership status.
	leader atomic.Bool
}

// New creates a named Service with configurable dependencies.
func New(name string) *Service {
	// Configure context for service shutdown signals.
	ctx, cancel := context.WithCancel(context.Background())

	return &Service{
		Config:    config.New(),
		id:        fmt.Sprintf("%s-%s", name, uuid.New()),
		ctx:       ctx,
		ctxCancel: cancel,
		errCh:     make(chan error),
	}
}

// Service getter methods.
func (s *Service) Ctx() context.Context { return s.ctx }
func (s *Service) ID() string           { return s.id }
func (s *Service) IsLeader() bool       { return s.leader.Load() }
func (s *Service) Name() string         { return strings.SplitN(s.ID(), "-", 2)[0] }

// RunFunc is a function that will be called by Run to initialize a service.
// If this function returns an error then the server will immediately shut down.
type RunFunc func(*Service) error

// Run starts the Service and ensures all dependencies are initialised.
//
// By default, it will start the local web server and wait for a stop
// signal to be received before attempting to gracefully shutdown.
func (s *Service) Run(run RunFunc) {
	// Initialize required service config.
	s.LoadRequiredConfig()

	// Configure global logger.
	plog.ConfigureGlobalLogging(s.Config.String(config.KeyServiceLogLevel), s.ID(), version.Build)

	// Attempt to start the service.
	if err := s.start(run); err != nil {
		log.Error().Err(err).Msg("service startup error")
		s.Exit(ExitError)
	}

	// Start the local http server.
	go s.serveHTTP()

	// Wait for an exit signal or service error.
	status := s.waitSignal()

	// Attempt to gracefully shutdown.
	s.Exit(status)
}

// waitSignal blocks waiting for operating system signals or an internal
// service error.
func (s *Service) waitSignal() int {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signal to be received.
	select {
	case <-stop:
		log.Info().Msg("stop signal received, starting shut down")
	case err := <-s.errCh:
		// Always return an error as nothing should be sending nil
		// to this channel.
		log.Error().Err(err).Msg("internal service error")
		return ExitError
	}

	return ExitOk
}

// Exit cancels the service's context in order to signal a shutdown to child processes.
// It sleeps for a configurable time before signalling the process to exit.
func (s *Service) Exit(status int) {
	s.ctxCancel()
	time.Sleep(s.Config.Duration(config.KeyServiceShutdownTimeout))
	os.Exit(status)
}

// start attempts to run the service with an initial timeout.
// If the deadline exceeds the time taken to run the service, it is treated
// as a failed start.
func (s *Service) start(run RunFunc) error {
	// Configure startup timeout.
	timeout := s.Config.Duration(config.KeyServiceStartupTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Attempt to run the main service func and record error.
	errCh := make(chan error, 1)
	go func() {
		errCh <- run(s)
	}()

	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("startup deadline (%s) exceeded", timeout)
		}
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("error running service: %w", err)
		}
	}

	metrics.ServiceInfo.WithLabelValues(s.ID(), version.Build).Inc()
	log.Info().Msg("service started")

	return nil
}
