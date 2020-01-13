package service

import (
	sdk "github.com/irisnet/irishub/types"
)

type FeesOutput struct {
	ReturnedFee sdk.Coins `json:"returned_fee"`
	IncomingFee sdk.Coins `json:"incoming_fee"`
}
