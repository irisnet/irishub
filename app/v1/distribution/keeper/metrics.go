package keeper

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	promutil "github.com/irisnet/irishub/tools/prometheus"
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

	communityTaxVec := stdprometheus.NewGaugeVec(stdprometheus.GaugeOpts{
		Namespace: config.Namespace,
		Subsystem: MetricsSubsystem,
		Name:      "community_tax",
		Help:      "community tax",
	}, []string{})

	promutil.RegisterMetrics(communityTaxVec)

	return &Metrics{
		CommunityTax: prometheus.NewGauge(communityTaxVec),
	}
}

func NopMetrics() *Metrics {
	return &Metrics{
		CommunityTax: discard.NewGauge(),
	}
}
