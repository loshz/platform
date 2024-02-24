package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	"github.com/loshz/platform/internal/config"
)

// serveHTTP configures and starts the local webserver.
//
// By default, it will register pprof, metrics and health endpoints.
func (s *Service) serveHTTP(ctx context.Context) {
	s.Scheduler().Add(1)
	defer s.Scheduler().Done()

	router := http.NewServeMux()

	// Configure debug endpoints.
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// Expose the registered metrics via HTTP.
	router.Handle("/metrics", promhttp.Handler())

	// Expose basic health check.
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		res := struct {
			Service string `json:"service"`
			Status  string `json:"status"`
			Leader  bool   `json:"leader"`
		}{s.ID(), "OK", s.IsLeader()}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Error().Err(err).Msg("error encoding health check response data")
		}
	})

	// Configure HTTP server with sane defaults.
	timeout := 10 * time.Second
	srv := &http.Server{
		Handler:           router,
		ReadTimeout:       timeout,
		ReadHeaderTimeout: timeout,
		WriteTimeout:      timeout,
		IdleTimeout:       timeout,
	}

	go func() {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Config().Int(config.KeyHttpServerPort)))
		if err != nil {
			s.Error(fmt.Errorf("http server tcp error: %w", err))
			return
		}

		// Update config with the actual tcp listener port.
		s.Config().Set(config.KeyHttpServerPort, ln.Addr().(*net.TCPAddr).Port)

		log.Info().Msgf("http server running on %s", ln.Addr())
		if err := srv.Serve(ln); err != http.ErrServerClosed {
			s.Error(fmt.Errorf("local http server error: %w", err))
			return
		}
	}()

	// Wait for service to exit and shutdown.
	<-ctx.Done()
	log.Info().Msg("stopping http server")
	_ = srv.Shutdown(context.Background())
}
