// nolint
package types

import (
	"context"
	"time"

	"github.com/gogo/protobuf/proto"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

/*
The intent of Context is for it to be an immutable object that can be
cloned and updated cheaply with WithValue() and passed forward to the
next decorator or handler. For example,

 func MsgHandler(ctx Context, tx Tx) Result {
 	...
 	ctx = ctx.WithValue(key, value)
 	...
 }
*/
type Context struct {
	ctx             context.Context
	ms              MultiStore
	header          abci.Header
	chainID         string
	txBytes         []byte
	logger          log.Logger
	voteInfo        []abci.VoteInfo
	gasMeter        GasMeter
	blockGasMeter   GasMeter
	checkTx         bool
	minimumFee      Coins
	consParams      *abci.ConsensusParams
	coinFlowTrigger string
	coinFlowTags    CoinFlowTags
	validTxCounter  *ValidTxCounter
}

// Read-only accessors
func (c Context) Context() context.Context        { return c.ctx }
func (c Context) MultiStore() MultiStore          { return c.ms }
func (c Context) BlockHeight() int64              { return c.header.Height }
func (c Context) BlockTime() time.Time            { return c.header.Time }
func (c Context) ChainID() string                 { return c.chainID }
func (c Context) TxBytes() []byte                 { return c.txBytes }
func (c Context) Logger() log.Logger              { return c.logger }
func (c Context) VoteInfos() []abci.VoteInfo      { return c.voteInfo }
func (c Context) GasMeter() GasMeter              { return c.gasMeter }
func (c Context) BlockGasMeter() GasMeter         { return c.blockGasMeter }
func (c Context) IsCheckTx() bool                 { return c.checkTx }
func (c Context) MinimumFees() Coins              { return c.minimumFee }
func (c Context) CoinFlowTags() CoinFlowTags      { return c.coinFlowTags }
func (c Context) CoinFlowTrigger() string         { return c.coinFlowTrigger }
func (c Context) ValidTxCounter() *ValidTxCounter { return c.validTxCounter }

// clone the header before returning
func (c Context) BlockHeader() abci.Header {
	var msg = proto.Clone(&c.header).(*abci.Header)
	return *msg
}

func (c Context) ConsensusParams() *abci.ConsensusParams {
	return proto.Clone(c.consParams).(*abci.ConsensusParams)
}

// NewContext create a new context
func NewContext(ms MultiStore, header abci.Header, isCheckTx bool, logger log.Logger) Context {
	// https://github.com/gogo/protobuf/issues/519
	header.Time = header.Time.UTC()
	return Context{
		ctx:             context.Background(),
		ms:              ms,
		header:          header,
		chainID:         header.ChainID,
		checkTx:         isCheckTx,
		logger:          logger,
		gasMeter:        NewInfiniteGasMeter(),
		minimumFee:      Coins{},
		coinFlowTrigger: "",
		coinFlowTags:    NewCoinFlowRecord(false),
		validTxCounter:  NewValidTxCounter(),
	}
}

func (c Context) WithContext(ctx context.Context) Context {
	c.ctx = ctx
	return c
}

func (c Context) WithMultiStore(ms MultiStore) Context {
	c.ms = ms
	return c
}

func (c Context) WithBlockHeader(header abci.Header) Context {
	// https://github.com/gogo/protobuf/issues/519
	header.Time = header.Time.UTC()
	c.header = header
	return c
}

func (c Context) WithBlockTime(newTime time.Time) Context {
	newHeader := c.BlockHeader()
	// https://github.com/gogo/protobuf/issues/519
	newHeader.Time = newTime.UTC()
	return c.WithBlockHeader(newHeader)
}

func (c Context) WithProposer(addr ConsAddress) Context {
	newHeader := c.BlockHeader()
	newHeader.ProposerAddress = addr.Bytes()
	return c.WithBlockHeader(newHeader)
}

func (c Context) WithBlockHeight(height int64) Context {
	newHeader := c.BlockHeader()
	newHeader.Height = height
	return c.WithBlockHeader(newHeader)
}

func (c Context) WithChainID(chainID string) Context {
	c.chainID = chainID
	return c
}

func (c Context) WithCheckValidNum(txCounter *ValidTxCounter) Context {
	c.validTxCounter = txCounter
	return c
}

func (c Context) WithTxBytes(txBytes []byte) Context {
	c.txBytes = txBytes
	return c
}

func (c Context) WithLogger(logger log.Logger) Context {
	c.logger = logger
	return c
}

func (c Context) WithVoteInfos(voteInfo []abci.VoteInfo) Context {
	c.voteInfo = voteInfo
	return c
}

func (c Context) WithGasMeter(meter GasMeter) Context {
	c.gasMeter = meter
	return c
}

func (c Context) WithBlockGasMeter(meter GasMeter) Context {
	c.blockGasMeter = meter
	return c
}

func (c Context) WithIsCheckTx(isCheckTx bool) Context {
	c.checkTx = isCheckTx
	return c
}

func (c Context) WithMinimumFees(minimumFees Coins) Context {
	c.minimumFee = minimumFees
	return c
}

func (c Context) WithConsensusParams(params *abci.ConsensusParams) Context {
	c.consParams = params
	return c
}

func (c Context) WithCoinFlowTags(cTag CoinFlowTags) Context {
	c.coinFlowTags = cTag
	return c
}

// WithCoinFlowTrigger set the coinFlowTrigger for context
// in handler, coinFlowTrigger = tx.hash
// in begin/end block, coinFlowTrigger = {modules-name}{begin/end}Blocker
func (c Context) WithCoinFlowTrigger(coinFlowTrigger string) Context {
	c.coinFlowTrigger = coinFlowTrigger
	return c
}

// is context nil
func (c Context) IsZero() bool {
	return c.ms == nil
}

// WithValue is deprecated, provided for backwards compatibility
// Please use
//     ctx = ctx.WithContext(context.WithValue(ctx.Context(), key, false))
// instead of
//     ctx = ctx.WithValue(key, false)
func (c Context) WithValue(key, value interface{}) Context {
	c.ctx = context.WithValue(c.ctx, key, value)
	return c
}

// Value is deprecated, provided for backwards compatibility
// Please use
//     ctx.Context().Value(key)
// instead of
//     ctx.Value(key)
func (c Context) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}

// ----------------------------------------------------------------------------
// Store / Caching
// ----------------------------------------------------------------------------

// KVStore fetches a KVStore from the MultiStore.
func (c Context) KVStore(key StoreKey) KVStore {
	return c.MultiStore().GetKVStore(key).Gas(c.GasMeter(), cachedKVGasConfig)
}

// TransientStore fetches a TransientStore from the MultiStore.
func (c Context) TransientStore(key StoreKey) KVStore {
	return c.MultiStore().GetKVStore(key).Gas(c.GasMeter(), cachedTransientGasConfig)
}

// CacheContext returns a new cached context. The cached context is
// written to the context when writeCache is called.
func (c Context) CacheContext() (cc Context, writeCache func()) {
	cms := c.MultiStore().CacheMultiStore()
	cc = c.WithMultiStore(cms)
	return cc, cms.Write
}
