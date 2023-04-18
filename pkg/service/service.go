// Package service implements the base structure of a platform application.
package service

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/loshz/platform/pkg/config"
	"github.com/loshz/platform/pkg/leader"
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
	// UUID of the individual service including name prefix.
	ID   string
	name string

	// Service configuration.
	Config *config.Config

	// Global context used to signal service shutdown to spawned
	// goroutines.
	ctx    context.Context
	cancel context.CancelFunc

	// Store the current leadership status.
	leader atomic.Bool

	// Service specific exit handlers.
	exitHandlers    []ExitHandler
	exitHandlersMtx sync.RWMutex

	pbv1.UnimplementedPlatformServiceServer
}

// New creates a named Service with configurable dependencies.
func New(name string) *Service {
	// TODO: init deps here- config, DBs, etc.

	// Configure context for service shutdown signals.
	ctx, cancel := context.WithCancel(context.Background())

	return &Service{
		ID:     fmt.Sprintf("%s-%s", name, uuid.New()),
		Config: config.New(),
		name:   name,
		ctx:    ctx,
		cancel: cancel,
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
	s.loadRequiredConfig()

	// Configure global logger.
	plog.ConfigureGlobalLogging(s.Config.String(config.KeyLogLevel), s.ID, version.Build)

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

// Exit the current process while attempting to run any of the Service's exit handlers.
//
// Each exit handler will run in its own goroutine and will force exit after 30s
// regardless of completion.
func (s *Service) Exit(status int) {
	// Cancel the service context.
	s.cancel()

	exitWait := 30 * time.Second
	deadline := time.Now().Add(exitWait)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)

	// Force exit after deadline.
	time.AfterFunc(exitWait, func() {
		log.Error().Msgf("service shutdown took longer than %s, aborting", exitWait)
		cancel()
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
	cancel()
	os.Exit(status)
}

// registerLeader attempts to acquire election status. If successful,
// it will register itself as the leader and release leadership status
// upon service exit.
func (s *Service) registerLeader() {
	fd, err := leader.Acquire(s.name)
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

// ExitHandler is a timed function ran only on service shutdown.
type ExitHandler func(context.Context)

// AddExitHandler registers a function that should be called when the server is shutting down.
func (s *Service) AddExitHandler(fn ExitHandler) {
	s.exitHandlersMtx.Lock()
	s.exitHandlers = append(s.exitHandlers, fn)
	s.exitHandlersMtx.Unlock()
}

func (s *Service) loadRequiredConfig() {
	s.Config.MustLoad(config.KeyLogLevel, "info", config.ParseLogLevel)
	s.Config.MustLoad(config.KeyServiceStartupTimeout, "5s", config.ParseDuration)
	s.Config.MustLoad(config.KeyHTTPPort, 8001, config.ParseInt)
}
