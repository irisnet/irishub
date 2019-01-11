package keeper

import (
	"github.com/irisnet/irishub/modules/distribution/types"
	sdk "github.com/irisnet/irishub/types"
)

// check whether a validator has distribution info
func (k Keeper) HasValidatorDistInfo(ctx sdk.Context,
	operatorAddr sdk.ValAddress) (exists bool) {
	store := ctx.KVStore(k.storeKey)
	return store.Has(GetValidatorDistInfoKey(operatorAddr))
}

// get the validator distribution info
func (k Keeper) GetValidatorDistInfo(ctx sdk.Context,
	operatorAddr sdk.ValAddress) (vdi types.ValidatorDistInfo) {

	store := ctx.KVStore(k.storeKey)

	b := store.Get(GetValidatorDistInfoKey(operatorAddr))
	if b == nil {
		panic("Stored validator-distribution info should not have been nil")
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &vdi)
	return
}

// set the validator distribution info
func (k Keeper) SetValidatorDistInfo(ctx sdk.Context, vdi types.ValidatorDistInfo) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(vdi)
	store.Set(GetValidatorDistInfoKey(vdi.OperatorAddr), b)
}

// remove a validator distribution info
func (k Keeper) RemoveValidatorDistInfo(ctx sdk.Context, valAddr sdk.ValAddress) {
	// defensive check
	vdi := k.GetValidatorDistInfo(ctx, valAddr)
	if vdi.DelAccum.Accum.IsPositive() {
		panic("Should not delete validator with unwithdrawn delegator accum")
	}
	if !vdi.ValCommission.IsZero() {
		panic("Should not delete validator with unwithdrawn validator commission")
	}
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetValidatorDistInfoKey(valAddr))
}

// Get the calculated accum of a validator at the current block
// without affecting the state.
func (k Keeper) GetValidatorAccum(ctx sdk.Context, operatorAddr sdk.ValAddress) (sdk.Dec, sdk.Error) {
	if !k.HasValidatorDistInfo(ctx, operatorAddr) {
		return sdk.Dec{}, types.ErrNoValidatorDistInfo(k.codespace)
	}

	// withdraw self-delegation
	height := ctx.BlockHeight()
	lastValPower := k.stakeKeeper.GetLastValidatorPower(ctx, operatorAddr)
	valInfo := k.GetValidatorDistInfo(ctx, operatorAddr)
	accum := valInfo.GetValAccum(height, sdk.NewDecFromInt(lastValPower))

	return accum, nil
}

// takeValidatorFeePoolRewards updates the validator's distribution info
// from the global fee pool without withdrawing any rewards. This will be called
// from a onValidatorModified hook.
func (k Keeper) takeValidatorFeePoolRewards(ctx sdk.Context, operatorAddr sdk.ValAddress) sdk.Error {
	if !k.HasValidatorDistInfo(ctx, operatorAddr) {
		return types.ErrNoValidatorDistInfo(k.codespace)
	}
	accAddr := sdk.AccAddress(operatorAddr.Bytes())

	// withdraw reward for self-delegation
	if k.HasDelegationDistInfo(ctx, accAddr, operatorAddr) {
		fp, vi, di, withdraw :=
			k.withdrawDelegationReward(ctx, accAddr, operatorAddr)
		k.SetFeePool(ctx, fp)
		k.SetValidatorDistInfo(ctx, vi)
		k.SetDelegationDistInfo(ctx, di)
		k.WithdrawToDelegator(ctx, fp, accAddr, withdraw)
	}

	// withdraw validator commission rewards
	valInfo := k.GetValidatorDistInfo(ctx, operatorAddr)
	wc := k.GetWithdrawContext(ctx, operatorAddr)
	valInfo, feePool := valInfo.TakeFeePoolRewards(wc)

	k.SetFeePool(ctx, feePool)
	k.SetValidatorDistInfo(ctx, valInfo)

	return nil
}

// withdrawal all the validator rewards including the commission
func (k Keeper) WithdrawValidatorRewardsAll(ctx sdk.Context, operatorAddr sdk.ValAddress) (types.DecCoins, sdk.Tags, sdk.Error) {

	if !k.HasValidatorDistInfo(ctx, operatorAddr) {
		return nil, nil, types.ErrNoValidatorDistInfo(k.codespace)
	}

	// withdraw self-delegation
	accAddr := sdk.AccAddress(operatorAddr.Bytes())
	withdraw, resultTags := k.withdrawDelegationRewardsAll(ctx, accAddr)

	// withdraw validator commission rewards
	feePool, commission := k.withdrawValidatorCommission(ctx, operatorAddr)
	withdraw = withdraw.Plus(commission)
	commissionTruncated, _ := commission.TruncateDecimal()
	resultTags = resultTags.AppendTag(sdk.TagRewardCommission, []byte(commissionTruncated.String()))

	k.WithdrawToDelegator(ctx, feePool, accAddr, withdraw)
	return withdraw, resultTags, nil
}

func (k Keeper) withdrawValidatorCommission(ctx sdk.Context, operatorAddr sdk.ValAddress) (types.FeePool, types.DecCoins) {
	valInfo := k.GetValidatorDistInfo(ctx, operatorAddr)
	wc := k.GetWithdrawContext(ctx, operatorAddr)
	valInfo, feePool, commission := valInfo.WithdrawCommission(wc)
	k.SetValidatorDistInfo(ctx, valInfo)

	return feePool, commission
}

// get all the validator rewards including the commission
func (k Keeper) CurrentValidatorRewardsAll(ctx sdk.Context, operatorAddr sdk.ValAddress) (sdk.Coins, sdk.Error) {

	if !k.HasValidatorDistInfo(ctx, operatorAddr) {
		return sdk.Coins{}, types.ErrNoValidatorDistInfo(k.codespace)
	}

	// withdraw self-delegation
	accAddr := sdk.AccAddress(operatorAddr.Bytes())
	withdraw := k.CurrentDelegationRewardsAll(ctx, accAddr)

	// withdrawal validator commission rewards
	valInfo := k.GetValidatorDistInfo(ctx, operatorAddr)

	wc := k.GetWithdrawContext(ctx, operatorAddr)
	commission := valInfo.CurrentCommissionRewards(wc)
	withdraw = withdraw.Plus(commission)
	truncated, _ := withdraw.TruncateDecimal()
	return truncated, nil
}
