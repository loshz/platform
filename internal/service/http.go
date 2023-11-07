package service

import (
	"context"
	"encoding/json"
	"fmt"
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
func (s *Service) serveHTTP() {
	port := fmt.Sprintf(":%d", s.Config.Int(config.KeyHTTPPort))
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

	srv := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Info().Msgf("local http server running on %s", port)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			s.errCh <- fmt.Errorf("local http server error: %w", err)
			return
		}
	}()

	// Wait for service to exit and shutdown.
	<-s.ctx.Done()
	log.Info().Msg("stopping http server")
	_ = srv.Shutdown(context.Background())
}
