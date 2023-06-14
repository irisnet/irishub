package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/service/types"
)

// MaxRequestTimeout returns the maximum request timeout
func (k Keeper) MaxRequestTimeout(ctx sdk.Context) int64 {
	return k.GetParams(ctx).MaxRequestTimeout
}

// MinDepositMultiple returns the minimum deposit multiple
func (k Keeper) MinDepositMultiple(ctx sdk.Context) int64 {
	return k.GetParams(ctx).MinDepositMultiple
}

// MinDeposit returns the minimum deposit
func (k Keeper) MinDeposit(ctx sdk.Context) sdk.Coins {
	return k.GetParams(ctx).MinDeposit
}

// ServiceFeeTax returns the service fee tax
func (k Keeper) ServiceFeeTax(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).ServiceFeeTax
}

// SlashFraction returns the slashing fraction
func (k Keeper) SlashFraction(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).SlashFraction
}

// ComplaintRetrospect returns the complaint retrospect duration
func (k Keeper) ComplaintRetrospect(ctx sdk.Context) time.Duration {
	return k.GetParams(ctx).ComplaintRetrospect
}

// ArbitrationTimeLimit returns the arbitration time limit
func (k Keeper) ArbitrationTimeLimit(ctx sdk.Context) time.Duration {
	return k.GetParams(ctx).ArbitrationTimeLimit
}

// TxSizeLimit returns the tx size limit
func (k Keeper) TxSizeLimit(ctx sdk.Context) uint64 {
	return k.GetParams(ctx).TxSizeLimit
}

// BaseDenom returns the base denom of service module
func (k Keeper) BaseDenom(ctx sdk.Context) string {
	return k.GetParams(ctx).BaseDenom
}

// RestrictedServiceFeeDenom returns the boolean value which
// indicates if the service fee only accepts the base denom
func (k Keeper) RestrictedServiceFeeDenom(ctx sdk.Context) bool {
	return k.GetParams(ctx).RestrictedServiceFeeDenom
}

// GetParams sets the farm module parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(types.ParamsKey))
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams sets the farm module parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}
	store.Set(types.ParamsKey, bz)

	return nil
}
