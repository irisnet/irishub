package gov

import (
	"strconv"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	distr "github.com/irisnet/irishub/modules/distribution/types"
	"github.com/irisnet/irishub/modules/mint"
	promutil "github.com/irisnet/irishub/tools/prometheus"
	sdk "github.com/irisnet/irishub/types"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	cfg "github.com/tendermint/tendermint/config"
)

const (
	MetricsSubsystem = "module_gov"
	ValidatorLabel   = "validator_address"
	ProposalIDLabel  = "proposal_id"
	ParamKeyLabel    = "parameter_key"
)

type Metrics struct {
	ProposalStatus metrics.Gauge // 0:DepositPeriod 1:VotingPeriod 2:Pass 3:Reject 4:Other
	Vote           metrics.Gauge // 0:Yes 1:No 2:NoWithVeto 3:Abstain
	Param          metrics.Gauge
}

// PrometheusMetrics returns Metrics build using Prometheus client library.
func PrometheusMetrics(config *cfg.InstrumentationConfig) *Metrics {
	if !config.Prometheus {
		return NopMetrics()
	}

	proposalStatusVec := stdprometheus.NewGaugeVec(stdprometheus.GaugeOpts{
		Namespace: config.Namespace,
		Subsystem: MetricsSubsystem,
		Name:      "proposal_status",
		Help:      "the status of the proposal",
	}, []string{ProposalIDLabel})

	voteVec := stdprometheus.NewGaugeVec(stdprometheus.GaugeOpts{
		Namespace: config.Namespace,
		Subsystem: MetricsSubsystem,
		Name:      "vote",
		Help:      "validator vote the proposal",
	}, []string{ValidatorLabel, ProposalIDLabel})

	paramVec := stdprometheus.NewGaugeVec(stdprometheus.GaugeOpts{
		Namespace: config.Namespace,
		Subsystem: MetricsSubsystem,
		Name:      "parameter",
		Help:      "parameter changes",
	}, []string{ParamKeyLabel})

	promutil.RegisterMetrics(proposalStatusVec, voteVec, paramVec)

	return &Metrics{
		ProposalStatus: prometheus.NewGauge(proposalStatusVec),
		Vote:           prometheus.NewGauge(voteVec),
		Param:          prometheus.NewGauge(paramVec),
	}
}

func NopMetrics() *Metrics {
	return &Metrics{
		ProposalStatus: discard.NewGauge(),
		Vote:           discard.NewGauge(),
		Param:          discard.NewGauge(),
	}
}

func SetParameterMetrics(metrics *Metrics, key string, value interface{}) {
	switch key {
	case string(mint.KeyInflation), string(distr.KeyBaseProposerReward), string(distr.KeyBonusProposerReward), string(distr.KeyCommunityTax):
		valueFloat64, err := strconv.ParseFloat(value.(sdk.Dec).String(), 64)
		if err == nil {
			metrics.Param.With(ParamKeyLabel, key).Set(valueFloat64)
		}
	default:
	}
}
