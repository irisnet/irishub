package iservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type IServiceKeeper interface {
	Has(ctx sdk.Context, key []byte) bool
	Set(ctx sdk.Context, key, value []byte)
}

type IService interface {
	CheckTx(ctx sdk.Context, msg sdk.Msg) sdk.Result
	DeliverTx(ctx sdk.Context, msg sdk.Msg) sdk.Result
}

type IServiceRoute = map[string]IService
