package service

import (
	"github.com/irisnet/irishub/app/v1/service"
	sdk "github.com/irisnet/irishub/types"
)

type DefOutput struct {
	Definition service.SvcDef           `json:"definition"`
	Methods    []service.MethodProperty `json:"methods"`
}

type FeesOutput struct {
	ReturnedFee sdk.Coins `json:"returned_fee"`
	IncomingFee sdk.Coins `json:"incoming_fee"`
}
