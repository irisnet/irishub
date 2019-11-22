package keeper

import (
	promutil "github.com/irisnet/irishub/tools/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	cfg "github.com/tendermint/tendermint/config"
)

const (
	MetricsSubsystem = "module_stake"
	ValidatorLabel   = "validator_address"
)

type Metrics struct {
	BondedToken  *stdprometheus.GaugeVec
	LoosenToken  *stdprometheus.GaugeVec
	BurnedToken  *stdprometheus.GaugeVec
	SlashedToken *stdprometheus.CounterVec
	Jailed       *stdprometheus.GaugeVec
	Power        *stdprometheus.GaugeVec
	enabled      bool
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
	}, []string{ValidatorLabel})

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
	}, []string{ValidatorLabel})

	jailedVec := stdprometheus.NewGaugeVec(stdprometheus.GaugeOpts{
		Namespace: config.Namespace,
		Subsystem: MetricsSubsystem,
		Name:      "jailed",
		Help:      "jailed",
	}, []string{ValidatorLabel})

	powerVec := stdprometheus.NewGaugeVec(stdprometheus.GaugeOpts{
		Namespace: config.Namespace,
		Subsystem: MetricsSubsystem,
		Name:      "power",
		Help:      "power",
	}, []string{ValidatorLabel})

	promutil.RegisterMetrics(bondedTokenVec, loosenTokenVec, burnedTokenVec, slashedTokenVec, jailedVec, powerVec)

	return &Metrics{
		BondedToken:  bondedTokenVec,
		LoosenToken:  loosenTokenVec,
		BurnedToken:  burnedTokenVec,
		SlashedToken: slashedTokenVec,
		Jailed:       jailedVec,
		Power:        powerVec,
		enabled:      config.Prometheus,
	}
}

func (m *Metrics) SetBondedToken(valAddr string, bondedToken float64) {
	promutil.SafeExec(func() {
		m.BondedToken.With(stdprometheus.Labels{
			ValidatorLabel: valAddr,
		}).Set(bondedToken)
	}, m.enabled)
}

func (m *Metrics) DeleteBondedToken(valAddr string) {
	promutil.SafeExec(func() {
		m.BondedToken.Delete(stdprometheus.Labels{
			ValidatorLabel: valAddr,
		})
	}, m.enabled)
}

func (m *Metrics) SetLoosenToken(loosenToken float64) {
	promutil.SafeExec(func() {
		m.LoosenToken.WithLabelValues().Set(loosenToken)
	}, m.enabled)
}

func (m *Metrics) SetBurnedToken(burnedToken float64) {
	promutil.SafeExec(func() {
		m.BurnedToken.WithLabelValues().Set(burnedToken)
	}, m.enabled)
}

func (m *Metrics) SetSlashedToken(valAddr string, slashedToken float64) {
	promutil.SafeExec(func() {
		m.SlashedToken.With(stdprometheus.Labels{
			ValidatorLabel: valAddr,
		}).Add(slashedToken)
	}, m.enabled)
}

func (m *Metrics) Jail(valAddr string) {
	promutil.SafeExec(func() {
		m.Jailed.With(stdprometheus.Labels{
			ValidatorLabel: valAddr,
		}).Set(1)
	}, m.enabled)
}

func (m *Metrics) Unjail(valAddr string) {
	promutil.SafeExec(func() {
		m.Jailed.Delete(stdprometheus.Labels{
			ValidatorLabel: valAddr,
		})
	}, m.enabled)
}

func (m *Metrics) SetVotingPower(valAddr string, power float64) {
	promutil.SafeExec(func() {
		m.BondedToken.With(stdprometheus.Labels{
			ValidatorLabel: valAddr,
		}).Set(power)
	}, m.enabled)
}

func NopMetrics() *Metrics {
	return &Metrics{
		BondedToken:  promutil.EmptyGaugeVec(),
		LoosenToken:  promutil.EmptyGaugeVec(),
		BurnedToken:  promutil.EmptyGaugeVec(),
		SlashedToken: promutil.EmptyCounterVec(),
		Jailed:       promutil.EmptyGaugeVec(),
		Power:        promutil.EmptyGaugeVec(),
	}
}
