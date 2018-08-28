package consensus

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"log"
	"time"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/irisnet/irishub/app"
	"fmt"
	"github.com/spf13/viper"
)

// TODO
const keyStoreStake = "KeyActiveProposalQueue"


// Metrics contains metrics exposed by this package.
type GovMetrics struct {
	//number of active proposals
	activeProposalNum metrics.Gauge
	// number of active proposals that voted by the validator with given address
	needToVoteActiveProposalNum metrics.Gauge
	Address     sdk.AccAddress
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
		Address:     make([]byte, 0),
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

func (gov *Metrics) Start(ctx app.Context) {

	account_addr := viper.GetString("account-address")
	gov.SetAddress(account_addr)

	go func() {
		for {
			gov.RecordMetrics(ctx, ctx.Cdc)
			time.Sleep(time.Second*1)
		}
	}()
}

func (gov *Metrics) RecordMetrics(ctx app.Context, cdc *wire.Codec) {
	count := 0
	need_to_vote := 0
	if active_proposals, err :=  getAllActiveProposalsID(cdc, ctx);err != nil{
		log.Println(err.Error())
	}else{
		count = len(active_proposals)
		for _, proposal_id := range active_proposals{
			if _, err := getVote(proposal_id, gov.govMetrics.Address, cdc, ctx);
			err != nil{
				need_to_vote++
			}
		}
	}

	gov.govMetrics.activeProposalNum.Set(float64(count))
	gov.govMetrics.needToVoteActiveProposalNum.Set(float64(need_to_vote))
}
//-------------------------help functions--------------------------------------

func getAllInactiveProposalsID(cdc *wire.Codec, ctx app.Context) (proposals gov.ProposalQueue, err error){
	if res, err := ctx.Ctx.QueryStore(gov.KeyInactiveProposalQueue, "gov"); err != nil{
		return gov.ProposalQueue{}, err
	}else {
		err = cdc.UnmarshalBinary(res, &proposals)
		return proposals, err
	}
}

func getAllActiveProposalsID(cdc *wire.Codec, ctx app.Context) (proposals gov.ProposalQueue, err error){
	if res, err := ctx.Ctx.QueryStore(gov.KeyActiveProposalQueue, "gov"); err != nil{
		return gov.ProposalQueue{}, err
	}else {
		err = cdc.UnmarshalBinary(res, &proposals)
		return proposals, err
	}

}

func getProposal(ID int64, cdc *wire.Codec, ctx app.Context) ( *gov.Proposal, error){
	if res, err := ctx.Ctx.QueryStore(gov.KeyProposal(ID), "gov"); err != nil{
		return nil, err
	}else {
		var proposal *gov.Proposal
		err = cdc.UnmarshalBinary(res, proposal)
		return proposal, err
	}
}


func getVote(proposalID int64, voterAddr sdk.AccAddress, cdc *wire.Codec, ctx app.Context) (vote gov.Vote, err error){
	if res, err := ctx.Ctx.QueryStore(gov.KeyVote(proposalID, voterAddr), "gov"); err != nil{
		return gov.Vote{}, err
	}else {
		if len(res) == 0{
			return gov.Vote{}, fmt.Errorf("cannot find the vote that %s vote for proposal %d", voterAddr.String(), proposalID)
		}
		err = cdc.UnmarshalBinary(res, &vote)
		return vote, err
	}
}
