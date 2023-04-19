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
	"github.com/loshz/platform/pkg/leader"
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
	// UUID of the individual service including name prefix.
	// E.g., service-xxxx-xxxx
	ID string

	// Service configuration.
	Config *config.Config

	// Global context used to signal service shutdown to spawned
	// goroutines.
	ctx       context.Context
	ctxCancel context.CancelFunc

	// Store the current leadership status.
	leader atomic.Bool
}

// New creates a named Service with configurable dependencies.
func New(name string) *Service {
	// Configure context for service shutdown signals.
	ctx, cancel := context.WithCancel(context.Background())

	return &Service{
		ID:        fmt.Sprintf("%s-%s", name, uuid.New()),
		Config:    config.New(),
		ctx:       ctx,
		ctxCancel: cancel,
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
	s.LoadRequiredConfig()

	// Configure global logger.
	plog.ConfigureGlobalLogging(s.Config.String(config.KeyServiceLogLevel), s.ID, version.Build)

	// Attempt to start the service.
	s.start(run)

	// Start the local http server.
	go s.serveHTTP()

	// Attempt to acquire leader election.
	go s.registerLeader()

	// Wait for an exit signal.
	s.waitSignal()

	// Attempt to gracefully shutdown.
	s.Exit(ExitOk)
}

// waitSignal blocks waiting for operating system signals.
//
// By default, it will handle calls to SIGINT and SIGTERM.
func (s *Service) waitSignal() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signal to be received.
	<-stop
	log.Info().Msg("stop signal received, starting shut down")
}

// Return the service context so other individual goroutines can
// listen for a stop signal.
// NOTE: this context is only cancelled upon receiving a call to s.Exit()
func (s *Service) Ctx() context.Context {
	return s.ctx
}

// Exit cancels the service's context in order to signal a shutdown to child processes.
// It sleeps for a configurable time before signalling the process to exit.
func (s *Service) Exit(status int) {
	s.ctxCancel()
	time.Sleep(s.Config.Duration(config.KeyServiceShutdownTimeout))
	os.Exit(status)
}

// registerLeader attempts to acquire election status. If successful,
// it will register itself as the leader and release leadership status
// upon service exit.
func (s *Service) registerLeader() {
	fd, err := leader.Acquire(s.Name())
	if err != nil {
		log.Error().Err(err).Msg("error atempting leader election")
		s.Exit(ExitError)
	}
	defer leader.Release(fd)

	log.Info().Msg("leadership status acquired")
	s.leader.Store(true)

	<-s.ctx.Done()
}

// IsLeader returns the status of the current service's leadership.
func (s *Service) IsLeader() bool {
	return s.leader.Load()
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

	metrics.ServiceInfo.WithLabelValues(s.ID, version.Build).Inc()
	log.Info().Msg("service started")
}

func (s *Service) Name() string {
	return strings.SplitN(s.ID, "-", 2)[0]
}
