package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Namespace is the global qualifier of all metrics.
// It should be passed to all Prometheus metric opts.
const Namespace = "platform"

// ServiceInfo represents an indivdual service.
var ServiceInfo = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "service_info",
		Help:      "Service specific information.",
	},
	[]string{"service_id", "version"},
)
