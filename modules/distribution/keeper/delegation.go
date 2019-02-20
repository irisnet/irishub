package keeper

import (
	"fmt"
	"github.com/irisnet/irishub/modules/distribution/types"
	sdk "github.com/irisnet/irishub/types"
)

// check whether a delegator distribution info exists
func (k Keeper) HasDelegationDistInfo(ctx sdk.Context, delAddr sdk.AccAddress,
	valOperatorAddr sdk.ValAddress) (has bool) {
	store := ctx.KVStore(k.storeKey)
	return store.Has(GetDelegationDistInfoKey(delAddr, valOperatorAddr))
}

// get the delegator distribution info
func (k Keeper) GetDelegationDistInfo(ctx sdk.Context, delAddr sdk.AccAddress,
	valOperatorAddr sdk.ValAddress) (ddi types.DelegationDistInfo) {

	store := ctx.KVStore(k.storeKey)

	b := store.Get(GetDelegationDistInfoKey(delAddr, valOperatorAddr))
	if b == nil {
		panic("Stored delegation-distribution info should not have been nil")
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &ddi)
	return
}

// set the delegator distribution info
func (k Keeper) SetDelegationDistInfo(ctx sdk.Context, ddi types.DelegationDistInfo) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(ddi)
	store.Set(GetDelegationDistInfoKey(ddi.DelegatorAddr, ddi.ValOperatorAddr), b)
}

// remove a delegator distribution info
func (k Keeper) RemoveDelegationDistInfo(ctx sdk.Context, delAddr sdk.AccAddress,
	valOperatorAddr sdk.ValAddress) {

	store := ctx.KVStore(k.storeKey)
	store.Delete(GetDelegationDistInfoKey(delAddr, valOperatorAddr))
}

// remove all delegation distribution infos
func (k Keeper) RemoveDelegationDistInfos(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, DelegationDistInfoKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

// iterate over all the validator distribution infos
func (k Keeper) IterateDelegationDistInfos(ctx sdk.Context,
	fn func(index int64, distInfo types.DelegationDistInfo) (stop bool)) {

	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, DelegationDistInfoKey)
	defer iter.Close()
	index := int64(0)
	for ; iter.Valid(); iter.Next() {
		var ddi types.DelegationDistInfo
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &ddi)
		if fn(index, ddi) {
			return
		}
		index++
	}
}

//___________________________________________________________________________________________

// get the delegator withdraw address, return the delegator address if not set
func (k Keeper) GetDelegatorWithdrawAddr(ctx sdk.Context, delAddr sdk.AccAddress) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)

	b := store.Get(GetDelegatorWithdrawAddrKey(delAddr))
	if b == nil {
		return delAddr
	}
	return sdk.AccAddress(b)
}

// set the delegator withdraw address
func (k Keeper) SetDelegatorWithdrawAddr(ctx sdk.Context, delAddr, withdrawAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetDelegatorWithdrawAddrKey(delAddr), withdrawAddr.Bytes())
}

// remove a delegator withdraw info
func (k Keeper) RemoveDelegatorWithdrawAddr(ctx sdk.Context, delAddr, withdrawAddr sdk.AccAddress) {
	ctx.Logger().Debug("Remove delegator distribution info", "delegator", delAddr.String())
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetDelegatorWithdrawAddrKey(delAddr))
}

//___________________________________________________________________________________________

// return all rewards for a delegation
func (k Keeper) withdrawDelegationReward(ctx sdk.Context,
	delAddr sdk.AccAddress, valAddr sdk.ValAddress) (types.FeePool,
	types.ValidatorDistInfo, types.DelegationDistInfo, types.DecCoins) {
	logger := ctx.Logger()
	wc := k.GetWithdrawContext(ctx, valAddr)
	valInfo := k.GetValidatorDistInfo(ctx, valAddr)
	delInfo := k.GetDelegationDistInfo(ctx, delAddr, valAddr)
	validator := k.stakeKeeper.Validator(ctx, valAddr)
	delegation := k.stakeKeeper.Delegation(ctx, delAddr, valAddr)

	logger.Debug("Withdraw context", "commission_rate", wc.CommissionRate.String(), "total_power", wc.TotalPower, "validator_power", wc.ValPower, "community_pool", wc.FeePool.CommunityPool, "total_accum", wc.FeePool.TotalValAccum, "validator_pool", wc.FeePool.ValPool, "validator_total_delegation_shares", validator.GetDelegatorShares().String(), "delegation_shares", delegation.GetShares().String())
	logger.Debug("Before withdraw", "validator_distInfo", valInfo.String())
	delInfo, valInfo, feePool, withdraw := delInfo.WithdrawRewards(logger, wc, valInfo,
		validator.GetDelegatorShares(), delegation.GetShares())
	logger.Debug("After withdraw", "validator_distInfo", valInfo.String())

	recipient := k.GetDelegatorWithdrawAddr(ctx, delAddr)
	coins, _ := withdraw.TruncateDecimal()
	ctx.CoinFlowTags().AppendAddCoinSourceTag(ctx, recipient.String(), sdk.ValidatorDelegationReward, valAddr.String(), coins.String())
	return feePool, valInfo, delInfo, withdraw
}

// get all rewards for all delegations of a delegator
func (k Keeper) currentDelegationReward(ctx sdk.Context, delAddr sdk.AccAddress,
	valAddr sdk.ValAddress) types.DecCoins {

	wc := k.GetWithdrawContext(ctx, valAddr)

	valInfo := k.GetValidatorDistInfo(ctx, valAddr)
	delInfo := k.GetDelegationDistInfo(ctx, delAddr, valAddr)
	validator := k.stakeKeeper.Validator(ctx, valAddr)
	delegation := k.stakeKeeper.Delegation(ctx, delAddr, valAddr)

	estimation := delInfo.CurrentRewards(wc, valInfo,
		validator.GetDelegatorShares(), delegation.GetShares())

	return estimation
}

//___________________________________________________________________________________________

// withdraw all rewards for a single delegation
// NOTE: This gets called "onDelegationSharesModified",
// meaning any changes to bonded coins
func (k Keeper) WithdrawToDelegator(ctx sdk.Context, feePool types.FeePool,
	delAddr sdk.AccAddress, amount types.DecCoins) {

	withdrawAddr := k.GetDelegatorWithdrawAddr(ctx, delAddr)
	coinsToAdd, change := amount.TruncateDecimal()
	feePool.CommunityPool = feePool.CommunityPool.Plus(change)
	k.SetFeePool(ctx, feePool)

	ctx.Logger().Debug("Withdraw reward to delegator", "reward", coinsToAdd.String(), "change", change.ToString(), "delegator", delAddr.String())

	_, _, err := k.bankKeeper.AddCoins(ctx, withdrawAddr, coinsToAdd)
	if err != nil {
		panic(err)
	}
}

//___________________________________________________________________________________________

// withdraw all rewards for a single delegation
// NOTE: This gets called "onDelegationSharesModified",
// meaning any changes to bonded coins
func (k Keeper) WithdrawDelegationReward(ctx sdk.Context, delAddr sdk.AccAddress,
	valAddr sdk.ValAddress) (types.DecCoins, sdk.Error) {

	if !k.HasDelegationDistInfo(ctx, delAddr, valAddr) {
		ctx.Logger().Error("There is no delegation distribution information")
		return nil, types.ErrNoDelegationDistInfo(k.codespace)
	}

	feePool, valInfo, delInfo, withdraw :=
		k.withdrawDelegationReward(ctx, delAddr, valAddr)

	k.SetValidatorDistInfo(ctx, valInfo)
	k.SetDelegationDistInfo(ctx, delInfo)
	k.WithdrawToDelegator(ctx, feePool, delAddr, withdraw)
	return withdraw, nil
}

// current rewards for a single delegation
func (k Keeper) CurrentDelegationReward(ctx sdk.Context, delAddr sdk.AccAddress,
	valAddr sdk.ValAddress) (sdk.Coins, sdk.Error) {

	if !k.HasDelegationDistInfo(ctx, delAddr, valAddr) {
		return sdk.Coins{}, types.ErrNoDelegationDistInfo(k.codespace)
	}
	estCoins := k.currentDelegationReward(ctx, delAddr, valAddr)
	trucate, _ := estCoins.TruncateDecimal()
	return trucate, nil
}

//___________________________________________________________________________________________

// return all rewards for all delegations of a delegator
func (k Keeper) WithdrawDelegationRewardsAll(ctx sdk.Context, delAddr sdk.AccAddress) (types.DecCoins, sdk.Tags) {
	withdraw, tags := k.withdrawDelegationRewardsAll(ctx, delAddr)
	feePool := k.GetFeePool(ctx)
	k.WithdrawToDelegator(ctx, feePool, delAddr, withdraw)
	return withdraw, tags
}

func (k Keeper) withdrawDelegationRewardsAll(ctx sdk.Context,
	delAddr sdk.AccAddress) (types.DecCoins, sdk.Tags) {

	// iterate over all the delegations
	withdraw := types.DecCoins{}
	var tags sdk.Tags
	operationAtDelegation := func(_ int64, del sdk.Delegation) (stop bool) {

		valAddr := del.GetValidatorAddr()
		feePool, valInfo, delInfo, diWithdraw :=
			k.withdrawDelegationReward(ctx, delAddr, valAddr)
		withdraw = withdraw.Plus(diWithdraw)
		k.SetFeePool(ctx, feePool)
		k.SetValidatorDistInfo(ctx, valInfo)
		k.SetDelegationDistInfo(ctx, delInfo)
		diWithdrawTruncated, _ := diWithdraw.TruncateDecimal()
		ctx.Logger().Debug("Withdraw delegation reward", "reward", diWithdrawTruncated.String(), "delegator", del.GetDelegatorAddr().String(), "validator", del.GetValidatorAddr().String())
		tags = tags.AppendTag(fmt.Sprintf(sdk.TagRewardFromValidator, valAddr.String()), []byte(diWithdrawTruncated.String()))

		recipient := k.GetDelegatorWithdrawAddr(ctx, delAddr)
		coins, _ := diWithdraw.TruncateDecimal()
		ctx.CoinFlowTags().AppendAddCoinSourceTag(ctx, recipient.String(), sdk.ValidatorDelegationReward, valAddr.String(), coins.String())

		return false
	}
	k.stakeKeeper.IterateDelegations(ctx, delAddr, operationAtDelegation)
	return withdraw, tags
}

// get all rewards for all delegations of a delegator
func (k Keeper) CurrentDelegationRewardsAll(ctx sdk.Context,
	delAddr sdk.AccAddress) types.DecCoins {

	// iterate over all the delegations
	total := types.DecCoins{}
	operationAtDelegation := func(_ int64, del sdk.Delegation) (stop bool) {
		valAddr := del.GetValidatorAddr()
		est := k.currentDelegationReward(ctx, delAddr, valAddr)
		total = total.Plus(est)
		return false
	}
	k.stakeKeeper.IterateDelegations(ctx, delAddr, operationAtDelegation)
	return total
}
