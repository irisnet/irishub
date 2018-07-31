package consensus

import (
	cctx "context"
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/wire" // XXX fix
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/irisnet/irishub/app"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/consensus"
  	//"github.com/irisnet/irishub/tools"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"log"
	"strings"
	"time"
	"container/list"
	"bytes"
	"github.com/spf13/viper"
)

// TODO
const keyStoreStake  = "stake"
type BlockInfo struct {
	Height  int64
	Time    time.Time
	signed	int // whether the given address signed the block, 1 if signed, 0 else.
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

	// Time between this and the last block.
	BlockIntervalSeconds metrics.Histogram
	//average block interval in last 100 blocks (in seconds)
	AvgBlockIntervalSeconds metrics.Gauge
	//block info
	blockInfo *list.List// queue of BlockInfo
	// missed prevote ratio in last 100 blocks (in seconds)
	MissedPrevotesRatio metrics.Gauge
	// given address
	Address types.Address
	SignedCount int


	// Number of transactions.
	NumTxs metrics.Gauge
	// Size of the block.
	BlockSizeBytes metrics.Gauge
	// Total number of transactions.
	TotalTxs metrics.Gauge
}

func NewIrisMetrics()  IrisMetrics{
	return IrisMetrics{
		Candidates:prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "candidates",
			Help:      "Number of Candidates.",
		}, []string{}),
	}
}

type Metrics struct {
	TmMetrics 	consensus.Metrics
	IrisMetrics IrisMetrics
}

// PrometheusMetrics returns Metrics build using Prometheus client library.
func PrometheusMetrics() *Metrics {
	tmMetrics := *consensus.PrometheusMetrics()
	irisMetrics := NewIrisMetrics()
	return &Metrics{
		TmMetrics	: tmMetrics,
		IrisMetrics	: irisMetrics,

/*
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
		MissedPrevotesRatio:prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "consensus",
			Name:      "missed_precommits_ratio",
			Help:      "missed precommits ratio of last 100 blocks.",
		}, []string{}),
		Address:make([]byte, 0),
		SignedCount:0,
    */
	}
}

func (cs *Metrics) SetAddress(addr_str string){
	if addr, err := sdk.GetValAddressHex(addr_str);err != nil{
		log.Println("parse address falid ", err)
	}else {
		cs.Address = addr
	}
}

func (cs *Metrics) Start(ctx app.Context) {
	context, _ := cctx.WithTimeout(cctx.Background(), 10*time.Second)

	var client = ctx.Ctx.Client
/*
func (cs *Metrics) Start(ctx tools.Context) {
	context, _ := cctx.WithTimeout(cctx.Background(), 10*time.Second)
	validaor_addr := viper.GetString("address")
	cs.SetAddress(validaor_addr)
	var client = ctx.Client
 */

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
			cs.RecordMetrics(ctx, ctx.Cdc, block.Block)
		}
	}()
}

func (cs *Metrics) RecordMetrics(ctx app.Context, cdc *wire.Codec, block *types.Block) {
	var client = ctx.Ctx.Client

	cs.TmMetrics.Height.Set(float64(block.Height))
	cs.TmMetrics.ByzantineValidators.Set(float64(len(block.Evidence.Evidence)))
/*
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
  */

	missingValidators := 0
	missingValidatorsPower := int64(0)
	validatorsPower := int64(0)
	resultValidators, err := client.Validators(&block.Height)
	if err != nil {
		panic(err)
	}
	validators := resultValidators.Validators
	valMap := make(map[string]types.Validator, len(validators))
	for i, val := range validators {
		var vote *types.Vote
		if i < len(block.LastCommit.Precommits) {
			vote = block.LastCommit.Precommits[i]
		}
		if vote == nil {
			missingValidators++
			missingValidatorsPower += val.VotingPower
		}

		valMap[val.Address.String()] = *val
		validatorsPower += val.VotingPower
	}
	cs.IrisMetrics.Candidates.Set(float64(getCandidatesNum(cdc, ctx)))
	cs.TmMetrics.MissingValidators.Set(float64(missingValidators))
	cs.TmMetrics.MissingValidatorsPower.Set(float64(missingValidatorsPower))
	cs.TmMetrics.ValidatorsPower.Set(float64(validatorsPower))
	cs.TmMetrics.Validators.Set(float64(len(validators)))

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
  /*
	if block.Height > 0{
		signed := 0
		for _, vote := range block.LastCommit.Precommits{
			if bytes.Equal(vote.ValidatorAddress.Bytes(), cs.Address.Bytes()){
				signed = 1
				break
			}
		}
		cs.blockInfo.PushBack(BlockInfo{Height:block.Height, Time:block.Time, signed:signed})
		firstBlock := cs.blockInfo.Front().Value.(BlockInfo)
		if cs.blockInfo.Len() > 100{
			cs.blockInfo.Remove(cs.blockInfo.Front())
		}
		cs.SignedCount += firstBlock.signed
		avgInterval := time.Now().Sub(firstBlock.Time).Seconds()/float64(cs.blockInfo.Len())
		cs.AvgBlockIntervalSeconds.Set(avgInterval)
		cs.MissedPrevotesRatio.Set(1 - float64(cs.SignedCount)/float64(cs.blockInfo.Len()))
	}

	cs.NumTxs.Set(float64(block.NumTxs))
	cs.TotalTxs.Set(float64(block.TotalTxs))
  */
	bz, _ := cdc.MarshalBinaryBare(block)
	cs.TmMetrics.BlockSizeBytes.Set(float64(len(bz)))
}

func getCandidatesNum(cdc *wire.Codec, ctx app.Context) int {
	key := stake.ValidatorsKey
	resKVs, err := ctx.Ctx.QuerySubspace(cdc, key, keyStoreStake)
	if err != nil {
		fmt.Println(err)
	}
	return len(resKVs)
}
