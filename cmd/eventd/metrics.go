package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// EventsTotal represents the total number of streamed events.
var EventsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: serviceName,
		Name:      "events_total",
		Help:      "Total number of events streamed.",
	},
	[]string{"type"},
)
