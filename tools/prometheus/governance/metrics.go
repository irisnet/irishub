package governance

import (
	"fmt"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/modules/gov"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"log"
	"time"
	govtypes "github.com/irisnet/irishub/types/gov"
)

// TODO
const (
	storeName = "gov"
)

// Metrics contains metrics exposed by this package.
type GovMetrics struct {
	//number of active proposals
	activeProposalNum metrics.Gauge
	// number of active proposals that voted by the validator with given address
	needToVoteActiveProposalNum metrics.Gauge
	Address                     sdk.AccAddress
}

func NewGovMetrics() GovMetrics {
	return GovMetrics{
		activeProposalNum: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "governance",
			Name:      "active_proposal_num",
			Help:      "Number of active proposals.",
		}, []string{}),

		needToVoteActiveProposalNum: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "governance",
			Name:      "need_voted_active_proposal_num",
			Help:      "Number of active proposals that need to vote.",
		}, []string{}),
		Address: make([]byte, 0),
	}
}

type Metrics struct {
	govMetrics GovMetrics
}

// PrometheusMetrics returns Metrics build using Prometheus client library.
func PrometheusMetrics() *Metrics {
	govMetrics := NewGovMetrics()
	return &Metrics{
		govMetrics: govMetrics,
	}
}

func (gov *Metrics) SetAddress(addr_str string) {
	if addr, err := sdk.AccAddressFromBech32(addr_str); err != nil {
		log.Println("parse account address falid ", err)
	} else {
		if len(addr) == 0 {
			log.Println("account address is null ")
		} else {
			gov.govMetrics.Address = addr
		}
	}
}

func (gov *Metrics) Start(ctx context.CLIContext) {

	accountAddr := viper.GetString("account-address")
	gov.SetAddress(accountAddr)

	go func() {
		for {
			gov.RecordMetrics(ctx, ctx.Codec)
			time.Sleep(time.Second * 1)
		}
	}()
}

func (gov *Metrics) RecordMetrics(ctx context.CLIContext, cdc *codec.Codec) {
	count := 0
	needToVote := 0
	if activeProposals, err := getAllActiveProposalsID(cdc, ctx); err != nil {
		log.Println(err.Error())
	} else {
		count = len(activeProposals)
		for _, proposalId := range activeProposals {
			if _, err := getVote(int64(proposalId), gov.govMetrics.Address, cdc, ctx); err != nil {
				needToVote++
			}
		}
	}

	gov.govMetrics.activeProposalNum.Set(float64(count))
	gov.govMetrics.needToVoteActiveProposalNum.Set(float64(needToVote))
}

//-------------------------help functions--------------------------------------

func getAllInactiveProposalsID(cdc *codec.Codec, ctx context.CLIContext) (proposals govtypes.ProposalQueue, err error) {
	if res, err := ctx.QueryStore(gov.PrefixInactiveProposalQueue, storeName); err != nil {
		return govtypes.ProposalQueue{}, err
	} else {
		err = cdc.UnmarshalBinaryLengthPrefixed(res, &proposals)
		return proposals, err
	}
}

func getAllActiveProposalsID(cdc *codec.Codec, ctx context.CLIContext) (proposals govtypes.ProposalQueue, err error) {
	if res, err := ctx.QueryStore(gov.PrefixActiveProposalQueue, storeName); len(res) == 0 || err != nil {
		return govtypes.ProposalQueue{}, err
	} else {
		err = cdc.UnmarshalBinaryLengthPrefixed(res, &proposals)
		return proposals, err
	}

}

func getProposal(ID int64, cdc *codec.Codec, ctx context.CLIContext) (*govtypes.Proposal, error) {
	if res, err := ctx.QueryStore(gov.KeyProposal(uint64(ID)), storeName); err != nil {
		return nil, err
	} else {
		var proposal *govtypes.Proposal
		err = cdc.UnmarshalBinaryLengthPrefixed(res, proposal)
		return proposal, err
	}
}

func getVote(proposalID int64, voterAddr sdk.AccAddress, cdc *codec.Codec, ctx context.CLIContext) (vote govtypes.Vote, err error) {
	if res, err := ctx.QueryStore(gov.KeyVote(uint64(proposalID), voterAddr), storeName); err != nil {
		return govtypes.Vote{}, err
	} else {
		if len(res) == 0 {
			return govtypes.Vote{}, fmt.Errorf("cannot find the vote that %s vote for proposal %d", voterAddr.String(), proposalID)
		}
		err = cdc.UnmarshalBinaryLengthPrefixed(res, &vote)
		return vote, err
	}
}
