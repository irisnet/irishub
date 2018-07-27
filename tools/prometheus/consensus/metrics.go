package consensus

import (
	cctx "context"
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/wire" // XXX fix
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/irisnet/irishub/tools"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/tendermint/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"log"
	"strings"
	"time"
	"container/list"
	"bytes"
	"github.com/spf13/viper"
)

type BlockInfo struct {
	Height  int64
	Time    time.Time
	signed	int // whether the given address signed the block, 1 if signed, 0 else.
}

// Metrics contains metrics exposed by this package.
type Metrics struct {
	// Height of the chain.
	Height metrics.Gauge

	// Number of rounds.
	Rounds metrics.Gauge

	// Number of validators.
	Validators metrics.Gauge
	// Number of Candidates
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

	// Time between this and the last block.
	BlockIntervalSeconds metrics.Histogram
	//average block interval in last 100 blocks (in seconds)
	AvgBlockIntervalSeconds metrics.Gauge
	//block info
	blockInfo *list.List// queue of BlockInfo

	//Voting Power of the validator
	VotingPower metrics.Gauge
	//ratio of Voting Power of the validator to total voting power
	VotingPowerRatio metrics.Gauge
	// ratio of signed blocks in last 100 blocks
	UpTime metrics.Gauge
	// missed precommited since monitor up
	MissedPrecommits metrics.Gauge
	// given address
	Address types.Address
	SignedCount int
	MissedCount int


	// Number of transactions.
	NumTxs metrics.Gauge
	// Size of the block.
	BlockSizeBytes metrics.Gauge
	// Total number of transactions.
	TotalTxs metrics.Gauge
}

// PrometheusMetrics returns Metrics build using Prometheus client library.
func PrometheusMetrics() *Metrics {
	return &Metrics{
		Height: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "height",
			Help:      "Height of the chain.",
		}, []string{}),
		Rounds: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "rounds",
			Help:      "Number of rounds.",
		}, []string{}),

		Validators: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "validators",
			Help:      "Number of validators.",
		}, []string{}),

		Candidates: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "candidates",
			Help:      "Number of Candidates.",
		}, []string{}),

		ValidatorsPower: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "validators_power",
			Help:      "Total power of all validators.",
		}, []string{}),
		MissingValidators: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "missing_validators",
			Help:      "Number of validators who did not sign.",
		}, []string{}),
		MissingValidatorsPower: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "missing_validators_power",
			Help:      "Total power of the missing validators.",
		}, []string{}),
		ByzantineValidators: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "byzantine_validators",
			Help:      "Number of validators who tried to double sign.",
		}, []string{}),
		ByzantineValidatorsPower: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "byzantine_validators_power",
			Help:      "Total power of the byzantine validators.",
		}, []string{}),

		BlockIntervalSeconds: prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
			Subsystem: "consensus",
			Name:      "block_interval_seconds",
			Help:      "Time between this and the last block.",
			Buckets:   []float64{1, 2.5, 5, 10, 60},
		}, []string{}),
		AvgBlockIntervalSeconds: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "avg_block_interval_seconds",
			Help:      "average block interval of last 100 blocks (in seconds).",
		}, []string{}),
		blockInfo:list.New(),
		UpTime:prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "missed_up_time",
			Help:      "ratio of signed blocks in last 100 blocks.",
		}, []string{}),
		MissedPrecommits:prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "missed_precommits_count",
			Help:      "missed precommited since monitor up.",
		}, []string{}),
		VotingPower:prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "voting_power",
			Help:      "voting power of the validator",
		}, []string{}),
		VotingPowerRatio:prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "voting_power_ratio",
			Help:      "ratio of voting power of the validator to total voting power",
		}, []string{}),
		Address:make([]byte, 0),
		SignedCount:0,
		MissedCount:0,

		NumTxs: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "num_txs",
			Help:      "Number of transactions.",
		}, []string{}),
		BlockSizeBytes: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "block_size_bytes",
			Help:      "Size of the block.",
		}, []string{}),
		TotalTxs: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "total_txs",
			Help:      "Total number of transactions.",
		}, []string{}),
	}
}

func (cs *Metrics) SetAddress(addr_str string){
	if addr, err := sdk.GetValAddressHex(addr_str);err != nil{
		log.Println("parse address falid ", err)
	}else {
		cs.Address = addr
	}
}

func (cs *Metrics) Start(ctx tools.Context) {
	context, _ := cctx.WithTimeout(cctx.Background(), 10*time.Second)
	validaor_addr := viper.GetString("address")
	cs.SetAddress(validaor_addr)
	var client = ctx.Client
	//开启监听事件
	client.Start()

	blockC := make(chan interface{})

	err := client.Subscribe(context, "monitor", types.EventQueryNewBlock, blockC)

	if err != nil {
		log.Println("got ", err)
	}

	go func() {
		for e := range blockC {
			block := e.(types.TMEventData).(types.EventDataNewBlock)
			cs.RecordMetrics(ctx, ctx.Cdc, block.Block, ctx.StoreName)
		}
	}()
}

func (cs *Metrics) RecordMetrics(ctx tools.Context, cdc *wire.Codec, block *types.Block, storeName string) {
	cs.Height.Set(float64(block.Height))
	if len(block.Evidence.Evidence) != 0{
		for _, evidence := range block.Evidence.Evidence{
			fmt.Println()
			fmt.Println(evidence)
			fmt.Println()
		}
	}

	cs.ByzantineValidators.Set(float64(len(block.Evidence.Evidence)))

	missingValidators := 0
	missingValidatorsPower := int64(0)
	validatorsPower := int64(0)
	resultValidators, err := ctx.Client.Validators(&block.Height)
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
		if bytes.Equal(cs.Address, val.Address){
			votingPower = val.VotingPower
		}
		valMap[val.Address.String()] = *val
		validatorsPower += val.VotingPower
	}
	cs.Candidates.Set(float64(getCandidatesNum(cdc, ctx, storeName)))
	cs.MissingValidators.Set(float64(missingValidators))
	cs.MissingValidatorsPower.Set(float64(missingValidatorsPower))
	cs.ValidatorsPower.Set(float64(validatorsPower))
	cs.Validators.Set(float64(len(validators)))

	cs.VotingPower.Set(float64(votingPower))
	cs.VotingPowerRatio.Set(float64(votingPower)/float64(validatorsPower))

	byzantineValidatorsPower := int64(0)
	for _, ev := range block.Evidence.Evidence {
		addr := strings.ToUpper(hex.EncodeToString(ev.Address()))
		if val, ok := valMap[addr]; ok {
			byzantineValidatorsPower += val.VotingPower
		}
	}
	cs.ByzantineValidatorsPower.Set(float64(byzantineValidatorsPower))

	if block.Height > 1 {
		lastBlockHight := block.Height - 1
		lastBlock, _ := ctx.Client.Block(&lastBlockHight)
		interval := block.Time.Sub(lastBlock.BlockMeta.Header.Time).Seconds()
		cs.BlockIntervalSeconds.Observe(interval)
	}

	if block.Height > 0{
		signed := 0
		for _, vote := range block.LastCommit.Precommits{
			if bytes.Equal(vote.ValidatorAddress.Bytes(), cs.Address.Bytes()){
				signed = 1
				break
			}
		}
		cs.MissedCount += 1- signed
		cs.SignedCount += signed

		cs.blockInfo.PushBack(BlockInfo{Height:block.Height, Time:block.Time, signed:signed})
		firstBlock := cs.blockInfo.Front().Value.(BlockInfo)
		if cs.blockInfo.Len() > 100{
			cs.blockInfo.Remove(cs.blockInfo.Front())
			cs.SignedCount -= firstBlock.signed
		}

		avgInterval := time.Now().Sub(firstBlock.Time).Seconds()/float64(cs.blockInfo.Len())
		cs.AvgBlockIntervalSeconds.Set(avgInterval)
		cs.UpTime.Set(float64(cs.SignedCount)/float64(cs.blockInfo.Len()))
		cs.MissedPrecommits.Set(float64(cs.MissedCount))
	}

	cs.NumTxs.Set(float64(block.NumTxs))
	cs.TotalTxs.Set(float64(block.TotalTxs))

	bz, _ := cdc.MarshalBinaryBare(block)
	cs.BlockSizeBytes.Set(float64(len(bz)))
}

func getCandidatesNum(cdc *wire.Codec, ctx tools.Context, storeName string) int {
	key := stake.ValidatorsKey
	resKVs, err := ctx.QuerySubspace(cdc, key, storeName)
	if err != nil {
		fmt.Println(err)
	}
	return len(resKVs)
}
