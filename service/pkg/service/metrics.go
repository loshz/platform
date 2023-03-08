package service

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/loshz/platform/pkg/metrics"
	"github.com/loshz/platform/pkg/version"
)

// ServiceInfo represents an indivdual service.
var ServiceInfo = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: metrics.Namespace,
		Name:      "service_info",
		Help:      "Service specific info.",
	},
	[]string{"service_name", "service_id", "version"},
)

// registerDefaultMetrics registers service level metrics that are
// default on all platform services.
func (s *Service) registerDefaultMetrics() {
	prometheus.Register(ServiceInfo)

	ServiceInfo.WithLabelValues(s.Name, s.ID, version.Version).Inc()
}
