package serviceparams

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetMaxRequestTimeout(ctx sdk.Context) int64 {
	MaxRequestTimeoutParameter.LoadValue(ctx)
	return MaxRequestTimeoutParameter.Value
}

func GetMinProviderDeposit(ctx sdk.Context) sdk.Coins {
	MinProviderDepositParameter.LoadValue(ctx)
	return MinProviderDepositParameter.Value
}

func SetMinProviderDeposit(ctx sdk.Context, i sdk.Coins) {
	MinProviderDepositParameter.Value = i
	MinProviderDepositParameter.SaveValue(ctx)
}

func SetMaxRequestTimeout(ctx sdk.Context, i int64) {
	MaxRequestTimeoutParameter.Value = i
	MaxRequestTimeoutParameter.SaveValue(ctx)
}
