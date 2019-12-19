package types

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"

	cfg "github.com/tendermint/tendermint/config"
)

const MetricsSubsystem = "module_" + ModuleName

type Metrics struct {
	ActiveRequests metrics.Gauge
}

// TODO
// PrometheusMetrics returns Metrics built by the Prometheus client library.
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

	if err := stdprometheus.Register(activeRequestsVec); err != nil {
		return NopMetrics()
	}

	return &Metrics{
		ActiveRequests: prometheus.NewGauge(activeRequestsVec),
	}
}

func NopMetrics() *Metrics {
	return &Metrics{
		ActiveRequests: discard.NewGauge(),
	}
}
