package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/irisnet/irismod/modules/service/types"
)

// ParamKeyTable for service module
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&types.Params{})
}

// MaxRequestTimeout returns the maximum request timeout
func (k Keeper) MaxRequestTimeout(ctx sdk.Context) (res int64) {
	k.paramSpace.Get(ctx, types.KeyMaxRequestTimeout, &res)
	return
}

// MinDepositMultiple returns the minimum deposit multiple
func (k Keeper) MinDepositMultiple(ctx sdk.Context) (res int64) {
	k.paramSpace.Get(ctx, types.KeyMinDepositMultiple, &res)
	return
}

// MinDeposit returns the minimum deposit
func (k Keeper) MinDeposit(ctx sdk.Context) (res sdk.Coins) {
	k.paramSpace.Get(ctx, types.KeyMinDeposit, &res)
	return
}

// ServiceFeeTax returns the service fee tax
func (k Keeper) ServiceFeeTax(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.KeyServiceFeeTax, &res)
	return
}

// SlashFraction returns the slashing fraction
func (k Keeper) SlashFraction(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.KeySlashFraction, &res)
	return
}

// ComplaintRetrospect returns the complaint retrospect duration
func (k Keeper) ComplaintRetrospect(ctx sdk.Context) (res time.Duration) {
	k.paramSpace.Get(ctx, types.KeyComplaintRetrospect, &res)
	return
}

// ArbitrationTimeLimit returns the arbitration time limit
func (k Keeper) ArbitrationTimeLimit(ctx sdk.Context) (res time.Duration) {
	k.paramSpace.Get(ctx, types.KeyArbitrationTimeLimit, &res)
	return
}

// TxSizeLimit returns the tx size limit
func (k Keeper) TxSizeLimit(ctx sdk.Context) (res uint64) {
	k.paramSpace.Get(ctx, types.KeyTxSizeLimit, &res)
	return
}

// BaseDenom returns the base denom of service module
func (k Keeper) BaseDenom(ctx sdk.Context) (res string) {
	k.paramSpace.Get(ctx, types.KeyBaseDenom, &res)
	return
}

// RestrictedServiceFeeDenom returns the boolean value which
// indicates if the service fee only accepts the base denom
func (k Keeper) RestrictedServiceFeeDenom(ctx sdk.Context) (res bool) {
	k.paramSpace.Get(ctx, types.KeyRestrictedServiceFeeDenom, &res)
	return
}

// GetParams gets all parameteras as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.MaxRequestTimeout(ctx),
		k.MinDepositMultiple(ctx),
		k.MinDeposit(ctx),
		k.ServiceFeeTax(ctx),
		k.SlashFraction(ctx),
		k.ComplaintRetrospect(ctx),
		k.ArbitrationTimeLimit(ctx),
		k.TxSizeLimit(ctx),
		k.BaseDenom(ctx),
		k.RestrictedServiceFeeDenom(ctx),
	)
}

// SetParams sets the params to the store
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
