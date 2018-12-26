package serviceparams

import (
	sdk "github.com/irisnet/irishub/types"
)

func GetSericeParams(ctx sdk.Context) Params {
	ServiceParameter.LoadValue(ctx)
	return ServiceParameter.Value
}

func GetMaxRequestTimeout(ctx sdk.Context) int64 {
	ServiceParameter.LoadValue(ctx)
	return ServiceParameter.Value.MaxRequestTimeout
}

func GetMinDepositMultiple(ctx sdk.Context) int64 {
	ServiceParameter.LoadValue(ctx)
	return ServiceParameter.Value.MinDepositMultiple
}

func GetServiceFeeTax(ctx sdk.Context) sdk.Dec {
	ServiceParameter.LoadValue(ctx)
	return ServiceParameter.Value.ServiceFeeTax
}

func GetSlashFraction(ctx sdk.Context) sdk.Dec {
	ServiceParameter.LoadValue(ctx)
	return ServiceParameter.Value.SlashFraction
}

func SetMinProviderDeposit(ctx sdk.Context, i int64) {
	ServiceParameter.Value.MinDepositMultiple = i
	ServiceParameter.SaveValue(ctx)
}

func SetMaxRequestTimeout(ctx sdk.Context, i int64) {
	ServiceParameter.Value.MaxRequestTimeout = i
	ServiceParameter.SaveValue(ctx)
}

func SetServiceFeeTax(ctx sdk.Context, i sdk.Dec) {
	ServiceParameter.Value.ServiceFeeTax = i
	ServiceParameter.SaveValue(ctx)
}

func SetSlashFraction(ctx sdk.Context, i sdk.Dec) {
	ServiceParameter.Value.SlashFraction = i
	ServiceParameter.SaveValue(ctx)
}
