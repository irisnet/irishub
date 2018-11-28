package service

import (
	"github.com/irisnet/irishub/modules/service"
	sdk "github.com/irisnet/irishub/types"
)

type DefOutput struct {
	service.SvcDef
	Methods []service.MethodProperty `json:"methods"`
}

type FeesOutput struct {
	ReturnedFee sdk.Coins `json:"returned_fee"`
	IncomingFee sdk.Coins `json:"incoming_fee"`
}
