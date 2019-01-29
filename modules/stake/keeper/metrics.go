package keeper

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	cfg "github.com/tendermint/tendermint/config"
)

const MetricsSubsystem = "module_stake"

type Metrics struct {
	BondedToken  metrics.Gauge
	LoosenToken  metrics.Gauge
	BurnedToken  metrics.Gauge
	SlashedToken metrics.Counter
	Jailed       metrics.Gauge
}

// PrometheusMetrics returns Metrics build using Prometheus client library.
func PrometheusMetrics(config *cfg.InstrumentationConfig) *Metrics {
	if !config.Prometheus {
		return NopMetrics()
	}
	return &Metrics{
		BondedToken: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: config.Namespace,
			Subsystem: MetricsSubsystem,
			Name:      "bonded_token",
			Help:      "bonded token",
		}, []string{"validator_address"}),
		LoosenToken: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: config.Namespace,
			Subsystem: MetricsSubsystem,
			Name:      "loosen_token",
			Help:      "loosen token",
		}, []string{}),
		BurnedToken: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: config.Namespace,
			Subsystem: MetricsSubsystem,
			Name:      "burned_token",
			Help:      "burned token",
		}, []string{}),
		SlashedToken: prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: config.Namespace,
			Subsystem: MetricsSubsystem,
			Name:      "slashed_token",
			Help:      "slashed token",
		}, []string{"validator_address"}),
		Jailed: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: config.Namespace,
			Subsystem: MetricsSubsystem,
			Name:      "jailed",
			Help:      "jailed",
		}, []string{"validator_address"}),
	}
}

func NopMetrics() *Metrics {
	return &Metrics{
		BondedToken:  discard.NewGauge(),
		LoosenToken:  discard.NewGauge(),
		BurnedToken:  discard.NewGauge(),
		SlashedToken: discard.NewCounter(),
		Jailed:       discard.NewGauge(),
	}
}
