package v2

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	cfg "github.com/tendermint/tendermint/config"
)

const MetricsSubsystem = "v2"

type Metrics struct {
	InvariantFailure metrics.Counter
}

// PrometheusMetrics returns Metrics build using Prometheus client library.
func PrometheusMetrics(config *cfg.InstrumentationConfig) *Metrics {
	if !config.Prometheus {
		return NopMetrics()
	}
	return &Metrics{
		InvariantFailure: prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: config.Namespace,
			Subsystem: MetricsSubsystem,
			Name:      "invariant_failure",
			Help:      "invariant failure",
		}, []string{"error"}),
	}
}

func NopMetrics() *Metrics {
	return &Metrics{
		InvariantFailure: discard.NewCounter(),
	}
}
