package keeper

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	promutil "github.com/irisnet/irishub/tools/prometheus"
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
	Power        metrics.Gauge
}

// PrometheusMetrics returns Metrics build using Prometheus client library.
func PrometheusMetrics(config *cfg.InstrumentationConfig) *Metrics {
	if !config.Prometheus {
		return NopMetrics()
	}

	bondedTokenVec := stdprometheus.NewGaugeVec(stdprometheus.GaugeOpts{
		Namespace: config.Namespace,
		Subsystem: MetricsSubsystem,
		Name:      "bonded_token",
		Help:      "bonded token",
	}, []string{"validator_address"})

	loosenTokenVec := stdprometheus.NewGaugeVec(stdprometheus.GaugeOpts{
		Namespace: config.Namespace,
		Subsystem: MetricsSubsystem,
		Name:      "loosen_token",
		Help:      "loosen token",
	}, []string{})

	burnedTokenVec := stdprometheus.NewGaugeVec(stdprometheus.GaugeOpts{
		Namespace: config.Namespace,
		Subsystem: MetricsSubsystem,
		Name:      "burned_token",
		Help:      "burned token",
	}, []string{})

	slashedTokenVec := stdprometheus.NewCounterVec(stdprometheus.CounterOpts{
		Namespace: config.Namespace,
		Subsystem: MetricsSubsystem,
		Name:      "slashed_token",
		Help:      "slashed token",
	}, []string{"validator_address"})

	jailedVec := stdprometheus.NewGaugeVec(stdprometheus.GaugeOpts{
		Namespace: config.Namespace,
		Subsystem: MetricsSubsystem,
		Name:      "jailed",
		Help:      "jailed",
	}, []string{"validator_address"})

	powerVec := stdprometheus.NewGaugeVec(stdprometheus.GaugeOpts{
		Namespace: config.Namespace,
		Subsystem: MetricsSubsystem,
		Name:      "power",
		Help:      "power",
	}, []string{"validator_address"})

	promutil.RegisterMetrics(bondedTokenVec, loosenTokenVec, burnedTokenVec, slashedTokenVec, jailedVec, powerVec)

	return &Metrics{
		BondedToken:  prometheus.NewGauge(bondedTokenVec),
		LoosenToken:  prometheus.NewGauge(loosenTokenVec),
		BurnedToken:  prometheus.NewGauge(burnedTokenVec),
		SlashedToken: prometheus.NewCounter(slashedTokenVec),
		Jailed:       prometheus.NewGauge(jailedVec),
		Power:        prometheus.NewGauge(powerVec),
	}
}

func NopMetrics() *Metrics {
	return &Metrics{
		BondedToken:  discard.NewGauge(),
		LoosenToken:  discard.NewGauge(),
		BurnedToken:  discard.NewGauge(),
		SlashedToken: discard.NewCounter(),
		Jailed:       discard.NewGauge(),
		Power:        discard.NewGauge(),
	}
}
