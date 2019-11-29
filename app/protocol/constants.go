package protocol

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// all store name
	AccountStore         = "acc"
	ParamsStore          = "params"
	ParamsTransientStore = "transient_params"
	HtlcStore            = "htlc"

	// all route for query and handler
	AccountRoute = AccountStore
	ParamsRoute  = ParamsStore
	HtlcRoute    = HtlcStore
)

var (
	KeyAccount = sdk.NewKVStoreKey(AccountStore)
	KeyParams  = sdk.NewKVStoreKey(ParamsStore)
	TkeyParams = sdk.NewTransientStoreKey(ParamsTransientStore)
	KeyHtlc    = sdk.NewKVStoreKey(HtlcStore)
)
