package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// GRPCRequestsTotal represents the total number of gRPC requests.
var GRPCRequestsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "grpc_requests_total",
		Help:      "Total number of gRPC requests.",
	},
	[]string{"service_id", "code", "method", "type"},
)

// GRPCRequestDuration represents the duration of gRPC requests in seconds.
var GRPCRequestDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: Namespace,
		Name:      "grpc_request_duration_seconds",
		Help:      "Duration of gRPC requests in seconds.",
		Buckets:   prometheus.LinearBuckets(0.01, 0.01, 10),
	},
	[]string{"service_id", "code", "method", "type"},
)
