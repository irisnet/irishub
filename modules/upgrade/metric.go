package upgrade

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	cfg "github.com/tendermint/tendermint/config"
)

const (
	MetricsSubsystem = "module_upgrade"
	ValidatorLabel   = "validator_address"
	VersionLabel     = "version"
)

type Metrics struct {
	Signal  metrics.Gauge
	Upgrade metrics.Gauge
}

// PrometheusMetrics returns Metrics build using Prometheus client library.
func PrometheusMetrics(config *cfg.InstrumentationConfig) *Metrics {
	if !config.Prometheus {
		return NopMetrics()
	}
	return &Metrics{
		Signal: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: config.Namespace,
			Subsystem: MetricsSubsystem,
			Name:      "signal",
			Help:      "validator runs the new software",
		}, []string{ValidatorLabel, VersionLabel}),
		Upgrade: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: config.Namespace,
			Subsystem: MetricsSubsystem,
			Name:      "upgrade",
			Help:      "alert the validators to install new version",
		}, []string{}),
	}
}

func NopMetrics() *Metrics {
	return &Metrics{
		Signal:  discard.NewGauge(),
		Upgrade: discard.NewGauge(),
	}
}
