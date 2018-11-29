package serviceparams

import (
	sdk "github.com/irisnet/irishub/types"
)

func GetMaxRequestTimeout(ctx sdk.Context) int64 {
	MaxRequestTimeoutParameter.LoadValue(ctx)
	return MaxRequestTimeoutParameter.Value
}

func GetMinDepositMultiple(ctx sdk.Context) int64 {
	MinDepositMultipleParameter.LoadValue(ctx)
	return MinDepositMultipleParameter.Value
}

func SetMinProviderDeposit(ctx sdk.Context, i int64) {
	MinDepositMultipleParameter.Value = i
	MinDepositMultipleParameter.SaveValue(ctx)
}

func SetMaxRequestTimeout(ctx sdk.Context, i int64) {
	MaxRequestTimeoutParameter.Value = i
	MaxRequestTimeoutParameter.SaveValue(ctx)
}
