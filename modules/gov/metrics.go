package gov

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	cfg "github.com/tendermint/tendermint/config"
)

const (
	MetricsSubsystem = "module_gov"
	ValidatorLabel   = "validator_address"
	ProposalIDLabel  = "proposalID"
)

type Metrics struct {
	ProposalStatus metrics.Gauge // 0:DepositPeriod 1:VotingPeriod 2:Pass 3:Reject 4:Other
	Vote           metrics.Gauge // 0:Yes 1:No 2:NoWithVeto 3:Abstain
}

// PrometheusMetrics returns Metrics build using Prometheus client library.
func PrometheusMetrics(config *cfg.InstrumentationConfig) *Metrics {
	if !config.Prometheus {
		return NopMetrics()
	}
	return &Metrics{
		ProposalStatus: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: config.Namespace,
			Subsystem: MetricsSubsystem,
			Name:      "proposal_status",
			Help:      "the status of the proposal",
		}, []string{ProposalIDLabel}),
		Vote: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: config.Namespace,
			Subsystem: MetricsSubsystem,
			Name:      "vote",
			Help:      "validator vote the proposal",
		}, []string{ValidatorLabel, ProposalIDLabel}),
	}
}

func NopMetrics() *Metrics {
	return &Metrics{
		ProposalStatus: discard.NewGauge(),
		Vote:           discard.NewGauge(),
	}
}
