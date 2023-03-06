// Package service implements the base structure of a platform application.
package service

import (
	"context"
	"fmt"
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

	"github.com/loshz/platform/pkg/config"
	plog "github.com/loshz/platform/pkg/log"
	"github.com/loshz/platform/pkg/version"
)

const (
	ExitOk = iota
	ExitError
)

// Service represents a platform application.
type Service struct {
	// uuid of service including prefix.
	id string

	// server configuration information
	Config *config.Config

	// local HTTP server
	httpServer *http.Server

	// service specific exit handlers
	exitHandlers   map[string]ExitHandler
	exitHandlersMu sync.RWMutex

	// TODO: add default service deps
	//  - Queue consumers
	//  - Metric gauges
	//  - TCP listeners
	//  - etc.
}

// New creates a named Service with configurable dependencies.
func New(name string) *Service {
	// TODO: init deps here- config, DBs, etc.

	return &Service{
		id:           fmt.Sprintf("%s-%s", name, uuid.New()),
		Config:       config.New(),
		exitHandlers: make(map[string]ExitHandler),
	}
}

// String gets the Service's uuid.
func (s *Service) String() string {
	return s.id
}

// RunFunc is a function that will be called by Run to initialize a service.
// If this function returns an error then the server will immediately shut down.
// Note: the contents of this function should run in a goroutine in order for it
// not to block.
type RunFunc func(*Service) error

// Run starts the Service and ensures all dependencies are initialised.
//
// By default, it will start the local web server and wait for a stop
// signal to be received before attempting to gracefully shutdown.
func (s *Service) Run(run RunFunc) {
	// TODO: start default tasks
	//  - connect queues
	//  - etc.

	// initialize required service config
	s.initDefaultConfig()

	// configure global logger
	plog.ConfigureGlobalLogging(s.Config.String(config.KeyLogLevel), s.String(), version.Version)

	// Start the internal http server
	s.StartLocalHTTP(s.Config.Int(config.KeyHTTPPort))

	// Attempt to run the main service func
	if err := run(s); err != nil {
		log.Error().Err(err).Msg("error running service")
		s.Exit(ExitError)
	}

	log.Info().Msg("service started")

	// wait for exit signal
	s.WaitSignal()

	// exit service
	s.Exit(ExitOk)
}

// WaitSignal blocks waiting for operating system signals.
//
// By default, it will handle calls to SIGINT and SIGTERM.
func (s *Service) WaitSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	// wait for signal to be received
	<-ch
	log.Info().Msg("stop signal received, starting shut down...")
}

// Exit the current process while attempting to run any of the Service's exit handlers.
//
// Each exit handler will run in its own goroutine and will force exit after 30s.
func (s *Service) Exit(status int) {
	s.exitHandlersMu.RLock()
	defer s.exitHandlersMu.RUnlock()

	exitWait := 30 * time.Second
	deadline := time.Now().Add(exitWait)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// force exit after deadline
	time.AfterFunc(exitWait, func() { os.Exit(status) })

	// stop the local http server before calling exit handlers in case
	// we need to process any remaining data
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("error shutting down http server")
	}

	var wg sync.WaitGroup
	for name, fn := range s.exitHandlers {
		wg.Add(1)
		go func(name string, fn ExitHandler) {
			defer wg.Done()
			fn(deadline)
			log.Info().Msgf("exit handler finished: %s", name)
		}(name, fn)
	}

	wg.Wait()
	os.Exit(status)
}

// StartLocalHTTP configures and starts the local webserver.
//
// By default, it will register pprof, metrics and health endpoints.
func (s *Service) StartLocalHTTP(port int) {
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
	router.HandleFunc("/health", healthHandler(s.String(), nil))

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Info().Msgf("local server running on :%d", port)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			s.Exit(ExitError)
		}
	}()
}

// ExitHandler is a timed function ran only on service shutdown.
type ExitHandler func(deadline time.Time)

// AddExitHandler registers a function that should be called when the server is shutting down.
func (s *Service) AddExitHandler(name string, fn ExitHandler) {
	s.exitHandlersMu.Lock()
	s.exitHandlers[name] = fn
	s.exitHandlersMu.Unlock()
}

// RemoveExitHandler removes an exit function that has previously been registered.
func (s *Service) RemoveExitHandler(name string) {
	s.exitHandlersMu.Lock()
	delete(s.exitHandlers, name)
	s.exitHandlersMu.Unlock()
}

func (s *Service) initDefaultConfig() {
	s.Config.Load(config.KeyLogLevel, "info", true, config.ParseLogLevel)
	s.Config.Load(config.KeyHTTPPort, 8001, true, config.ParseInt)
}
