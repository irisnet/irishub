package gov

import (
	distr "github.com/irisnet/irishub/app/v1/distribution/types"
	"github.com/irisnet/irishub/app/v1/mint"
	promutil "github.com/irisnet/irishub/tools/prometheus"
	sdk "github.com/irisnet/irishub/types"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	cfg "github.com/tendermint/tendermint/config"
	"strconv"
)

const (
	MetricsSubsystem = "module_gov"
	ValidatorLabel   = "validator_address"
	ProposalIDLabel  = "proposal_id"
	ParamKeyLabel    = "parameter_key"
)

type Label = stdprometheus.Labels

type Metrics struct {
	ProposalStatus *stdprometheus.GaugeVec // 1:DepositPeriod 2:VotingPeriod 3:Passed 4:Rejected
	Vote           *stdprometheus.GaugeVec // 1:Yes 2:Abstain 3:No 4:NoWithVeto
	Param          *stdprometheus.GaugeVec
	enabled        bool
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
		ProposalStatus: proposalStatusVec,
		Vote:           voteVec,
		Param:          paramVec,
		enabled:        config.Prometheus,
	}
}

func NopMetrics() *Metrics {
	return &Metrics{
		ProposalStatus: promutil.EmptyGaugeVec(),
		Vote:           promutil.EmptyGaugeVec(),
		Param:          promutil.EmptyGaugeVec(),
	}
}

func (metrics *Metrics) SetProposalStatus(proposalID uint64, status ProposalStatus) {
	promutil.SafeExec(func() {
		s := float64(status)
		metrics.ProposalStatus.WithLabelValues(strconv.FormatUint(proposalID, 10)).Set(s)
	}, metrics.enabled)
}

func (metrics *Metrics) DeleteProposalStatus(proposalID uint64) {
	promutil.SafeExec(func() {
		metrics.ProposalStatus.DeleteLabelValues(strconv.FormatUint(proposalID, 10))
	}, metrics.enabled)
}

func (metrics *Metrics) AddParameter(key string, value interface{}) {
	promutil.SafeExec(func() {
		switch key {
		case string(mint.KeyInflation), string(distr.KeyBaseProposerReward), string(distr.KeyBonusProposerReward), string(distr.KeyCommunityTax):
			valueFloat64, err := strconv.ParseFloat(value.(sdk.Dec).String(), 64)
			if err == nil {
				metrics.Param.WithLabelValues(key).Set(valueFloat64)
			}
		default:
		}
	}, metrics.enabled)
}
func (metrics *Metrics) AddVote(consAddr string, proposalID uint64, option VoteOption) {
	promutil.SafeExec(func() {
		labels := Label{
			ValidatorLabel:  consAddr,
			ProposalIDLabel: strconv.FormatUint(proposalID, 10),
		}
		metrics.Vote.With(labels).Set(float64(option))
	}, metrics.enabled)

}

func (metrics *Metrics) DeleteVote(valAddr string, proposalID uint64) {
	promutil.SafeExec(func() {
		labels := Label{
			ValidatorLabel:  valAddr,
			ProposalIDLabel: strconv.FormatUint(proposalID, 10),
		}
		metrics.Vote.Delete(labels)
	}, metrics.enabled)
}
