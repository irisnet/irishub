package protocol

import (
	"github.com/irisnet/irishub/newapp/protocol/router"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/types/common"
)

type Protocol interface {
	GetDefinition() common.ProtocolDefinition
	GetRouter() router.Router
	GetQueryRouter() router.QueryRouter
	GetAnteHandler() sdk.AnteHandler                   // ante handler for fee and auth
	GetFeeRefundHandler() sdk.FeeRefundHandler         // fee handler for fee refund
	GetFeePreprocessHandler() sdk.FeePreprocessHandler // fee handler for fee preprocessor

	// may be nil
	GetInitChainer() sdk.InitChainer1  // initialize state with validators and state blob
	GetBeginBlocker() sdk.BeginBlocker // logic to run before any txs
	GetEndBlocker() sdk.EndBlocker     // logic to run after all txs, and to determine valset changes
	Load()
	Init()
}

type ProtocolBase struct {
	definition common.ProtocolDefinition
	//	engine 		*ProtocolEngine
}

func (pb ProtocolBase) GetDefinition() common.ProtocolDefinition {
	return pb.definition
}
