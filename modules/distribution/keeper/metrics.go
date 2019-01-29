package keeper

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	cfg "github.com/tendermint/tendermint/config"
)

const MetricsSubsystem = "module_distribution"

type Metrics struct {
	CommunityTax metrics.Gauge
}

// PrometheusMetrics returns Metrics build using Prometheus client library.
func PrometheusMetrics(config *cfg.InstrumentationConfig) *Metrics {
	if !config.Prometheus {
		return NopMetrics()
	}
	return &Metrics{
		CommunityTax: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: config.Namespace,
			Subsystem: MetricsSubsystem,
			Name:      "community_tax",
			Help:      "community tax",
		}, []string{"height"}),
	}
}

func NopMetrics() *Metrics {
	return &Metrics{
		CommunityTax: discard.NewGauge(),
	}
}
