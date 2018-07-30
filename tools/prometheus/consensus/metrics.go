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
	"log"
	"strings"
	"time"
)

// TODO
const keyStoreStake  = "stake"
// Metrics contains metrics exposed by this package.

type IrisMetrics struct {
	Candidates metrics.Gauge
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

	}
}

func (cs *Metrics) Start(ctx app.Context) {
	context, _ := cctx.WithTimeout(cctx.Background(), 10*time.Second)

	var client = ctx.Ctx.Client
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
