package protocol

import (
	"encoding/json"

	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type Protocol interface {
	GetVersion() uint64
	GetRouter() Router
	GetQueryRouter() QueryRouter
	GetAnteHandler() sdk.AnteHandler                   // ante handler for fee and auth
	GetFeeRefundHandler() sdk.FeeRefundHandler         // fee handler for fee refund
	GetFeePreprocessHandler() sdk.FeePreprocessHandler // fee handler for fee preprocessor
	ExportAppStateAndValidators(ctx sdk.Context, forZeroHeight bool) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error)
	ValidateTx(ctx sdk.Context, txBytes []byte, msgs []sdk.Msg) sdk.Error

	// may be nil
	GetInitChainer() sdk.InitChainer1  // initialize state with validators and state blob
	GetBeginBlocker() sdk.BeginBlocker // logic to run before any txs
	GetEndBlocker() sdk.EndBlocker     // logic to run after all txs, and to determine valset changes

	GetKVStoreKeyList() []*sdk.KVStoreKey
	Load()
	Init()
	GetCodec() *codec.Codec
	InitMetrics(store sdk.CommitMultiStore) // init metrics
}
