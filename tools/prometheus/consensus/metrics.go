package consensus

import (
	"bytes"
	"container/list"
	cctx "context"
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/consensus"
	"github.com/tendermint/tendermint/types"
	"log"
	"strings"
	"time"
	"github.com/irisnet/irishub/client/context"
)

// TODO
const keyStoreStake = "stake"

type BlockInfo struct {
	Height int64
	Time   time.Time
	signed int // whether the given address signed the block, 1 if signed, 0 else.
}

// Metrics contains metrics exposed by this package.

type IrisMetrics struct {
	Candidates metrics.Gauge
	// Total power of all validators.
	ValidatorsPower metrics.Gauge
	// Number of validators who did not sign.
	MissingValidators metrics.Gauge
	// Total power of the missing validators.
	MissingValidatorsPower metrics.Gauge
	// Number of validators who tried to double sign.
	ByzantineValidators metrics.Gauge
	// Total power of the byzantine validators.
	ByzantineValidatorsPower metrics.Gauge

	//average block interval in last 100 blocks (in seconds)
	AvgBlockIntervalSeconds metrics.Gauge
	//block info
	blockInfo *list.List // queue of BlockInfo

	//Voting Power of the validator
	VotingPower metrics.Gauge
	//ratio of Voting Power of the validator to total voting power
	VotingPowerRatio metrics.Gauge
	//Voting Power ratio of the validator who has maximum voting power
	MaximumVotingPowerRatio metrics.Gauge
	// ratio of signed blocks in last 100 blocks
	UpTime metrics.Gauge
	// missed precommited since monitor up
	MissedPrecommits metrics.Gauge
	// given address
	Address     types.Address
	SignedCount int
	MissedCount int
}

func NewIrisMetrics() IrisMetrics {
	return IrisMetrics{
		Candidates: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "candidates",
			Help:      "Number of Candidates.",
		}, []string{}),

		AvgBlockIntervalSeconds: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "avg_block_interval_seconds",
			Help:      "average block interval of last 100 blocks (in seconds).",
		}, []string{}),
		blockInfo: list.New(),
		UpTime: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "up_time",
			Help:      "ratio of signed blocks in last 100 blocks.",
		}, []string{}),
		MissedPrecommits: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "missed_precommits_count",
			Help:      "missed precommited since monitor up.",
		}, []string{}),
		VotingPower: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "voting_power",
			Help:      "voting power of the validator",
		}, []string{}),
		VotingPowerRatio: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "voting_power_ratio",
			Help:      "ratio of voting power of the validator to total voting power",
		}, []string{}),
		MaximumVotingPowerRatio: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "maximum_voting_power_ratio",
			Help:      "Voting Power ratio of the validator who has maximum voting power",
		}, []string{}),
		Address:     make([]byte, 0),
		SignedCount: 0,
		MissedCount: 0,
	}
}

type Metrics struct {
	TmMetrics   consensus.Metrics
	IrisMetrics IrisMetrics
}

// PrometheusMetrics returns Metrics build using Prometheus client library.
func PrometheusMetrics() *Metrics {
	tmMetrics := *consensus.PrometheusMetrics()
	irisMetrics := NewIrisMetrics()
	return &Metrics{
		TmMetrics:   tmMetrics,
		IrisMetrics: irisMetrics,
	}
}

func (cs *Metrics) SetAddress(addrStr string) {
	if addr, err := hex.DecodeString(addrStr); err != nil {
		log.Println("parse validator address falid ", err)
	} else {
		if len(addr) == 0 {
			log.Println("validator address is null ")
		} else {
			cs.IrisMetrics.Address = addr
		}
	}
}

func (cs *Metrics) Start(ctx context.CLIContext) {

	validatorAddr := viper.GetString("address")
	cs.SetAddress(validatorAddr)

	ct, _ := cctx.WithTimeout(cctx.Background(), 10*time.Second)
	var client = ctx.Client

	//开启监听事件
	client.Start()

	blockC := make(chan interface{})

	err := client.Subscribe(ct, "monitor", types.EventQueryNewBlock, blockC)

	if err != nil {
		log.Println("got ", err)
		return
	}

	go func() {
		for e := range blockC {
			block := e.(types.TMEventData).(types.EventDataNewBlock)
			cs.RecordMetrics(ctx, ctx.Codec, block.Block)
		}
	}()

	roundC := make(chan interface{})
	err = client.Subscribe(ct, "monitor", types.EventQueryNewRound, roundC)
	if err != nil {
		log.Println("got ", err)
		return
	}

	go func() {
		for e := range roundC {
			round := e.(types.TMEventData).(types.EventDataRoundState)
			cs.TmMetrics.Rounds.Set(float64(round.Round))
		}
	}()
}

func (cs *Metrics) RecordMetrics(ctx context.CLIContext, cdc *codec.Codec, block *types.Block) {
	var client = ctx.Client

	cs.TmMetrics.Height.Set(float64(block.Height))
	cs.TmMetrics.ByzantineValidators.Set(float64(len(block.Evidence.Evidence)))

	missingValidators := 0
	missingValidatorsPower := int64(0)
	validatorsPower := int64(0)
	mostPowerful := int64(0)

	resultValidators, err := client.Validators(&block.Height)
	if err != nil {
		panic(err)
	}
	validators := resultValidators.Validators
	valMap := make(map[string]types.Validator, len(validators))
	var votingPower int64
	for i, val := range validators {
		var vote *types.Vote
		if i < len(block.LastCommit.Precommits) {
			vote = block.LastCommit.Precommits[i]
		}
		if vote == nil {
			missingValidators++
			missingValidatorsPower += val.VotingPower
		}
		if bytes.Equal(cs.IrisMetrics.Address, val.Address) {
			votingPower = val.VotingPower
		}
		if val.VotingPower > mostPowerful {
			mostPowerful = val.VotingPower
		}
		valMap[val.Address.String()] = *val
		validatorsPower += val.VotingPower
	}
	cs.IrisMetrics.Candidates.Set(float64(getCandidatesNum(ctx)))
	cs.TmMetrics.MissingValidators.Set(float64(missingValidators))
	cs.TmMetrics.MissingValidatorsPower.Set(float64(missingValidatorsPower))
	cs.TmMetrics.ValidatorsPower.Set(float64(validatorsPower))
	cs.TmMetrics.Validators.Set(float64(len(validators)))

	cs.IrisMetrics.VotingPower.Set(float64(votingPower))
	cs.IrisMetrics.VotingPowerRatio.Set(float64(votingPower) / float64(validatorsPower))
	cs.IrisMetrics.MaximumVotingPowerRatio.Set(float64(mostPowerful) / float64(validatorsPower))

	byzantineValidatorsPower := int64(0)
	for _, ev := range block.Evidence.Evidence {
		addr := strings.ToUpper(hex.EncodeToString(ev.Address()))
		if val, ok := valMap[addr]; ok {
			byzantineValidatorsPower += val.VotingPower
		}
	}
	cs.TmMetrics.ByzantineValidatorsPower.Set(float64(byzantineValidatorsPower))

	if block.Height > 1 {
		lastBlockHight := block.Height - 1
		lastBlock, _ := client.Block(&lastBlockHight)
		interval := block.Time.Sub(lastBlock.BlockMeta.Header.Time).Seconds()
		cs.TmMetrics.BlockIntervalSeconds.Observe(interval)
	}

	cs.TmMetrics.NumTxs.Set(float64(block.NumTxs))
	cs.TmMetrics.TotalTxs.Set(float64(block.TotalTxs))

	if block.Height > 0 {
		signed := 0
		for _, vote := range block.LastCommit.Precommits {
			if vote != nil && bytes.Equal(vote.ValidatorAddress.Bytes(), cs.IrisMetrics.Address.Bytes()) {
				signed = 1
				break
			}
		}
		cs.IrisMetrics.MissedCount += 1 - signed
		cs.IrisMetrics.SignedCount += signed

		cs.IrisMetrics.blockInfo.PushBack(BlockInfo{Height: block.Height, Time: block.Time, signed: signed})
		firstBlock := cs.IrisMetrics.blockInfo.Front().Value.(BlockInfo)
		if cs.IrisMetrics.blockInfo.Len() > 100 {
			cs.IrisMetrics.blockInfo.Remove(cs.IrisMetrics.blockInfo.Front())
			cs.IrisMetrics.SignedCount -= firstBlock.signed
		}

		avgInterval := time.Now().Sub(firstBlock.Time).Seconds() / float64(cs.IrisMetrics.blockInfo.Len())
		cs.IrisMetrics.AvgBlockIntervalSeconds.Set(avgInterval)
		cs.IrisMetrics.UpTime.Set(float64(cs.IrisMetrics.SignedCount) / float64(cs.IrisMetrics.blockInfo.Len()))
		cs.IrisMetrics.MissedPrecommits.Set(float64(cs.IrisMetrics.MissedCount))
	}
	bz, _ := cdc.MarshalBinaryLengthPrefixedBare(block)
	cs.TmMetrics.BlockSizeBytes.Set(float64(len(bz)))
}

func getCandidatesNum(ctx context.CLIContext) int {
	key := stake.ValidatorsKey
	resKVs, err := ctx.QuerySubspace(key, keyStoreStake)
	if err != nil {
		fmt.Println(err)
	}
	return len(resKVs)
}
