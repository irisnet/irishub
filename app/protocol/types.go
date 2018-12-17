package protocol

import (
	sdk "github.com/irisnet/irishub/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/irisnet/irishub/types/common"
	"encoding/json"
	protocolKeeper "github.com/irisnet/irishub/app/protocol/keeper"
)

type Protocol interface {
	GetDefinition() common.ProtocolDefinition
	GetRouter() Router
	GetQueryRouter() QueryRouter
	GetAnteHandler() sdk.AnteHandler                   // ante handler for fee and auth
	GetFeeRefundHandler() sdk.FeeRefundHandler         // fee handler for fee refund
	GetFeePreprocessHandler() sdk.FeePreprocessHandler // fee handler for fee preprocessor
	ExportAppStateAndValidators(ctx sdk.Context, forZeroHeight bool) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error)

	// may be nil
	GetInitChainer() sdk.InitChainer1  // initialize state with validators and state blob
	GetBeginBlocker() sdk.BeginBlocker // logic to run before any txs
	GetEndBlocker() sdk.EndBlocker     // logic to run after all txs, and to determine valset changes

	GetKVStoreKeyList() []*sdk.KVStoreKey
	Load(protocolKeeper.Keeper)
	Init()
}

type ProtocolBase struct {
	Definition common.ProtocolDefinition
	//	engine 		*ProtocolEngine
}

func (pb ProtocolBase) GetDefinition() common.ProtocolDefinition {
	return pb.Definition
}

