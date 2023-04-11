package service

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/rs/zerolog/log"
)

const (
	// Indicates the service status is OK
	HealthStatusOK = "OK"

	// Indicates the service status is Failing
	HealthStatusFailing = "Failing"
)

// healthResponse is returned to the client when performing a
// full readiness check
type healthResponse struct {
	mu sync.Mutex

	Service      string             `json:"service"`
	Status       string             `json:"status"`
	Dependencies []healthDependency `json:"dependencies,omitempty"`
}

// healthDependency represents an external service dependency
type healthDependency struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

// HealthCheckedDependency defines methods for performing health checks
// on external dependencies, such as MySQL or Kafka.
type HealthCheckedDependency interface {
	Name() string
	Check() error
}

// healthHandler creates a http handler to be used as a service health check.
//
// There are 2 potential reponses based off the underlying request:
// 1. if ?status=ready - we return early with a simple 200 OK response.
// 2. Full end-to-end process of checking service dependencies based off
// the individual service configuration.
func healthHandler(service string, services []HealthCheckedDependency) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := &healthResponse{
			Service: service,
			Status:  HealthStatusOK,
		}

		// Check for a simple liveness status and return early.
		if status := r.URL.Query().Get("status"); status != "ready" {
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(res); err != nil {
				log.Error().Err(err).Msg("error encoding health check response data")
			}
			return
		}

		// Perform full readiness checks including dependency health checks.
		var wg sync.WaitGroup
		for _, service := range services {
			wg.Add(1)

			go func(service HealthCheckedDependency) {
				defer wg.Done()

				dep := healthDependency{
					Service: service.Name(),
					Status:  HealthStatusOK,
				}

				if err := service.Check(); err != nil {
					log.Error().Err(err).Msgf("error performing health check for service '%s'", service.Name())
					dep.Status = HealthStatusFailing
				}

				res.mu.Lock()
				defer res.mu.Unlock()

				// As we assume the initial status is OK, we should only update the overall status
				// if a dependency is failing, therefore not potentially overwriting a Failing staus with
				// an OK status.
				if dep.Status == HealthStatusFailing {
					res.Status = HealthStatusFailing
				}
				res.Dependencies = append(res.Dependencies, dep)
			}(service)
		}

		// Wait for all checks to finish.
		wg.Wait()

		w.Header().Set("Content-Type", "application/json")

		// Return 503 if any health checks failed.
		if res.Status == HealthStatusFailing {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Error().Err(err).Msg("error encoding health check response data")
		}
	}
}
