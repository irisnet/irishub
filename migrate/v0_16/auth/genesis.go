package auth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	CollectedFees sdk.Coins `json:"collected_fee"`
	FeeAuth       FeeAuth   `json:"data"`
	Params        Params    `json:"params"`
}

type Params struct {
	GasPriceThreshold sdk.Int `json:"gas_price_threshold"`
	TxSizeLimit       uint64  `json:"tx_size"`
}

type FeeAuth struct {
	NativeFeeDenom string `json:"native_fee_denom"`
}
