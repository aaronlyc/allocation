package controllers

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	allocationMetrics = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "allocation_total_free_bad",
			Help: "The Gauge for allocation what total, free numbers and bad numbers",
		},
		[]string{"ms_name", "metrics"},
	)
)

func init() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(allocationMetrics)
}
