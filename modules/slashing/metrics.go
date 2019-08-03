package slashing

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	promutil "github.com/irisnet/irishub/tools/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	cfg "github.com/tendermint/tendermint/config"
)

const MetricsSubsystem = "module_slashing"

type Metrics struct {
	MissedBlocks metrics.Gauge
}

// PrometheusMetrics returns Metrics build using Prometheus client library.
func PrometheusMetrics(config *cfg.InstrumentationConfig) *Metrics {
	if !config.Prometheus {
		return NopMetrics()
	}
	missedBlocksVec := stdprometheus.NewGaugeVec(stdprometheus.GaugeOpts{
		Namespace: config.Namespace,
		Subsystem: MetricsSubsystem,
		Name:      "missed_blocks",
		Help:      "missed blocks",
	}, []string{"validator_address"})

	promutil.RegisterMetrics(missedBlocksVec)

	return &Metrics{
		MissedBlocks: prometheus.NewGauge(missedBlocksVec),
	}
}

func NopMetrics() *Metrics {
	return &Metrics{
		MissedBlocks: discard.NewGauge(),
	}
}
