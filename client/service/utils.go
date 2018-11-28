package service

import (
	"github.com/irisnet/irishub/modules/service"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DefOutput struct {
	Definition service.SvcDef           `json:"definition"`
	Methods    []service.MethodProperty `json:"methods"`
}

type FeesOutput struct {
	ReturnedFee sdk.Coins `json:"returned_fee"`
	IncomingFee sdk.Coins `json:"incoming_fee"`
}
