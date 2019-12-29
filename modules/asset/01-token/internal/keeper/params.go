package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/irisnet/irishub/modules/asset/01-token/internal/types"
)

// ParamTable for staking module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

// SetParamSet sets the params
func (k Keeper) SetParamSet(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetParamSet returns the params
func (k Keeper) GetParamSet(ctx sdk.Context) types.Params {
	var params types.Params
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// AssetTaxRate returns parameter assetTaxRate
func (k Keeper) AssetTaxRate(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.KeyAssetTaxRate, &res)
	return
}

// IssueTokenBaseFee returns parameter issueTokenBaseFee
func (k Keeper) IssueTokenBaseFee(ctx sdk.Context) (res sdk.Coin) {
	k.paramSpace.Get(ctx, types.KeyIssueTokenBaseFee, &res)
	return
}

// MintTokenFeeRatio returns parameter mintTokenFeeRatio
func (k Keeper) MintTokenFeeRatio(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.KeyMintTokenFeeRatio, &res)
	return
}
