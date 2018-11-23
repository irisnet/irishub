package service

import (
	"github.com/irisnet/irishub/modules/service"
)

type DefOutput struct {
	service.SvcDef
	Methods []service.MethodProperty `json:"methods"`
}

type FeesOutput struct {
	ReturnedFee service.ReturnedFee `json:"returned_fee"`
	IncomingFee service.IncomingFee `json:"incoming_fee"`
}
