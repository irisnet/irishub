package types

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"

	cfg "github.com/tendermint/tendermint/config"

	promutil "github.com/irisnet/irishub/tools/prometheus"
)

// MetricsSubsystem defines the metrics name for the service module
const MetricsSubsystem = "module_service"

// Metrics defines a metrics struct
type Metrics struct {
	ActiveRequests metrics.Gauge
}

// PrometheusMetrics returns Metrics built using Prometheus client library.
func PrometheusMetrics(config *cfg.InstrumentationConfig) *Metrics {
	if !config.Prometheus {
		return NopMetrics()
	}

	activeRequestsVec := stdprometheus.NewGaugeVec(stdprometheus.GaugeOpts{
		Namespace: config.Namespace,
		Subsystem: MetricsSubsystem,
		Name:      "active_requests",
		Help:      "active requests count",
	}, []string{})

	promutil.RegisterMetrics(activeRequestsVec)

	return &Metrics{
		ActiveRequests: prometheus.NewGauge(activeRequestsVec),
	}
}

// NopMetrics returns a no-op metrics
func NopMetrics() *Metrics {
	return &Metrics{
		ActiveRequests: discard.NewGauge(),
	}
}
