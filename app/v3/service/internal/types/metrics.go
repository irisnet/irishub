package types

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"

	cfg "github.com/tendermint/tendermint/config"

	promutil "github.com/irisnet/irishub/tools/prometheus"
)

const MetricsSubsystem = "module_service"

type Metrics struct {
	ActiveRequests metrics.Gauge
}

// PrometheusMetrics returns Metrics build using Prometheus client library.
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

func NopMetrics() *Metrics {
	return &Metrics{
		ActiveRequests: discard.NewGauge(),
	}
}
