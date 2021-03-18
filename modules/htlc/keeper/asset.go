package keeper

import (
	"time"

	gogotypes "github.com/gogo/protobuf/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/htlc/types"
)

// IncrementCurrentAssetSupply increments an asset's supply by the coin
func (k Keeper) IncrementCurrentAssetSupply(ctx sdk.Context, coin sdk.Coin) error {
	supply, found := k.GetAssetSupply(ctx, coin.Denom)
	if !found {
		return sdkerrors.Wrap(types.ErrAssetNotSupported, coin.Denom)
	}

	limit, err := k.GetSupplyLimit(ctx, coin.Denom)
	if err != nil {
		return err
	}
	supplyLimit := sdk.NewCoin(coin.Denom, limit.Limit)

	// Resulting current supply must be under asset's limit
	if supplyLimit.IsLT(supply.CurrentSupply.Add(coin)) {
		return sdkerrors.Wrapf(
			types.ErrExceedsSupplyLimit,
			"increase %s, asset supply %s, limit %s",
			coin, supply.CurrentSupply, supplyLimit,
		)
	}

	if limit.TimeLimited {
		timeBasedSupplyLimit := sdk.NewCoin(coin.Denom, limit.TimeBasedLimit)
		if timeBasedSupplyLimit.IsLT(supply.TimeLimitedCurrentSupply.Add(coin)) {
			return sdkerrors.Wrapf(
				types.ErrExceedsTimeBasedSupplyLimit,
				"increase %s, current time-based asset supply %s, limit %s",
				coin, supply.TimeLimitedCurrentSupply, timeBasedSupplyLimit,
			)
		}
		supply.TimeLimitedCurrentSupply = supply.TimeLimitedCurrentSupply.Add(coin)
	}

	supply.CurrentSupply = supply.CurrentSupply.Add(coin)
	k.SetAssetSupply(ctx, supply, coin.Denom)
	return nil
}

// DecrementCurrentAssetSupply decrement an asset's supply by the coin
func (k Keeper) DecrementCurrentAssetSupply(ctx sdk.Context, coin sdk.Coin) error {
	supply, found := k.GetAssetSupply(ctx, coin.Denom)
	if !found {
		return sdkerrors.Wrap(types.ErrAssetNotSupported, coin.Denom)
	}

	// Resulting current supply must be greater than or equal to 0
	// Use sdk.Int instead of sdk.Coin to prevent panic if true
	if supply.CurrentSupply.Amount.Sub(coin.Amount).IsNegative() {
		return sdkerrors.Wrapf(types.ErrInvalidCurrentSupply, "decrease %s, asset supply %s", coin, supply.CurrentSupply)
	}

	supply.CurrentSupply = supply.CurrentSupply.Sub(coin)
	k.SetAssetSupply(ctx, supply, coin.Denom)
	return nil
}

// IncrementIncomingAssetSupply increments an asset's incoming supply
func (k Keeper) IncrementIncomingAssetSupply(ctx sdk.Context, coin sdk.Coin) error {
	supply, found := k.GetAssetSupply(ctx, coin.Denom)
	if !found {
		return sdkerrors.Wrap(types.ErrAssetNotSupported, coin.Denom)
	}

	// 	Result of (current + incoming + amount) must be under asset's limit
	totalSupply := supply.CurrentSupply.Add(supply.IncomingSupply)

	limit, err := k.GetSupplyLimit(ctx, coin.Denom)
	if err != nil {
		return err
	}
	supplyLimit := sdk.NewCoin(coin.Denom, limit.Limit)
	if supplyLimit.IsLT(totalSupply.Add(coin)) {
		return sdkerrors.Wrapf(types.ErrExceedsSupplyLimit, "increase %s, asset supply %s, limit %s", coin, totalSupply, supplyLimit)
	}

	if limit.TimeLimited {
		timeLimitedTotalSupply := supply.TimeLimitedCurrentSupply.Add(supply.IncomingSupply)
		timeBasedSupplyLimit := sdk.NewCoin(coin.Denom, limit.TimeBasedLimit)
		if timeBasedSupplyLimit.IsLT(timeLimitedTotalSupply.Add(coin)) {
			return sdkerrors.Wrapf(
				types.ErrExceedsTimeBasedSupplyLimit,
				"increase %s, time-based asset supply %s, limit %s",
				coin, supply.TimeLimitedCurrentSupply, timeBasedSupplyLimit,
			)
		}
	}

	supply.IncomingSupply = supply.IncomingSupply.Add(coin)
	k.SetAssetSupply(ctx, supply, coin.Denom)
	return nil
}

// DecrementIncomingAssetSupply decrements an asset's incoming supply
func (k Keeper) DecrementIncomingAssetSupply(ctx sdk.Context, coin sdk.Coin) error {
	supply, found := k.GetAssetSupply(ctx, coin.Denom)
	if !found {
		return sdkerrors.Wrap(types.ErrAssetNotSupported, coin.Denom)
	}

	// Resulting incoming supply must be greater than or equal to 0
	// Use sdk.Int instead of sdk.Coin to prevent panic if true
	if supply.IncomingSupply.Amount.Sub(coin.Amount).IsNegative() {
		return sdkerrors.Wrapf(types.ErrInvalidIncomingSupply, "decrease %s, incoming supply %s", coin, supply.IncomingSupply)
	}

	supply.IncomingSupply = supply.IncomingSupply.Sub(coin)
	k.SetAssetSupply(ctx, supply, coin.Denom)
	return nil
}

// IncrementOutgoingAssetSupply increments an asset's outgoing supply
func (k Keeper) IncrementOutgoingAssetSupply(ctx sdk.Context, coin sdk.Coin) error {
	supply, found := k.GetAssetSupply(ctx, coin.Denom)
	if !found {
		return sdkerrors.Wrap(types.ErrAssetNotSupported, coin.Denom)
	}

	// Result of (outgoing + amount) must be less than current supply
	if supply.CurrentSupply.IsLT(supply.OutgoingSupply.Add(coin)) {
		return sdkerrors.Wrapf(
			types.ErrExceedsAvailableSupply,
			"swap amount %s, available supply %s",
			coin, supply.CurrentSupply.Amount.Sub(supply.OutgoingSupply.Amount),
		)
	}

	supply.OutgoingSupply = supply.OutgoingSupply.Add(coin)
	k.SetAssetSupply(ctx, supply, coin.Denom)
	return nil
}

// DecrementOutgoingAssetSupply decrements an asset's outgoing supply
func (k Keeper) DecrementOutgoingAssetSupply(ctx sdk.Context, coin sdk.Coin) error {
	supply, found := k.GetAssetSupply(ctx, coin.Denom)
	if !found {
		return sdkerrors.Wrap(types.ErrAssetNotSupported, coin.Denom)
	}

	// Resulting outgoing supply must be greater than or equal to 0
	// Use sdk.Int instead of sdk.Coin to prevent panic if true
	if supply.OutgoingSupply.Amount.Sub(coin.Amount).IsNegative() {
		return sdkerrors.Wrapf(types.ErrInvalidOutgoingSupply, "decrease %s, outgoing supply %s", coin, supply.OutgoingSupply)
	}

	supply.OutgoingSupply = supply.OutgoingSupply.Sub(coin)
	k.SetAssetSupply(ctx, supply, coin.Denom)
	return nil
}

// CreateNewAssetSupply creates a new AssetSupply in the store for the input denom
func (k Keeper) CreateNewAssetSupply(ctx sdk.Context, denom string) types.AssetSupply {
	supply := types.NewAssetSupply(
		sdk.NewCoin(denom, sdk.ZeroInt()), sdk.NewCoin(denom, sdk.ZeroInt()),
		sdk.NewCoin(denom, sdk.ZeroInt()), sdk.NewCoin(denom, sdk.ZeroInt()), time.Duration(0),
	)
	k.SetAssetSupply(ctx, supply, denom)
	return supply
}

// UpdateTimeBasedSupplyLimits updates the time based supply for each asset, resetting it if the current time window has elapsed.
func (k Keeper) UpdateTimeBasedSupplyLimits(ctx sdk.Context) {
	assets, found := k.GetAssets(ctx)
	if !found {
		return
	}
	previousBlockTime, found := k.GetPreviousBlockTime(ctx)
	if !found {
		previousBlockTime = ctx.BlockTime()
		k.SetPreviousBlockTime(ctx, previousBlockTime)
	}
	timeElapsed := ctx.BlockTime().Sub(previousBlockTime)
	for _, asset := range assets {
		supply, found := k.GetAssetSupply(ctx, asset.Denom)
		// if a new asset has been added by governance, create a new asset supply for it in the store
		if !found {
			supply = k.CreateNewAssetSupply(ctx, asset.Denom)
		}
		newTimeElapsed := supply.TimeElapsed + timeElapsed
		if asset.SupplyLimit.TimeLimited && newTimeElapsed < asset.SupplyLimit.TimePeriod {
			supply.TimeElapsed = newTimeElapsed
		} else {
			supply.TimeElapsed = time.Duration(0)
			supply.TimeLimitedCurrentSupply = sdk.NewCoin(asset.Denom, sdk.ZeroInt())
		}
		k.SetAssetSupply(ctx, supply, asset.Denom)
	}
	k.SetPreviousBlockTime(ctx, ctx.BlockTime())
}

// GetAssetSupply gets an asset's current supply from the store.
func (k Keeper) GetAssetSupply(ctx sdk.Context, denom string) (assetSupply types.AssetSupply, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetAssetSupplyKey(denom))
	if bz == nil {
		return types.AssetSupply{}, false
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &assetSupply)
	return assetSupply, true
}

// SetAssetSupply updates an asset's supply
func (k Keeper) SetAssetSupply(ctx sdk.Context, supply types.AssetSupply, denom string) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&supply)
	store.Set(types.GetAssetSupplyKey(denom), bz)
}

// IterateAssetSupplies provides an iterator over all stored AssetSupplies.
func (k Keeper) IterateAssetSupplies(
	ctx sdk.Context,
	cb func(supply types.AssetSupply) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.AssetSupplyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var supply types.AssetSupply
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &supply)

		if cb(supply) {
			break
		}
	}
}

// GetAllAssetSupplies returns all asset supplies from the store
func (k Keeper) GetAllAssetSupplies(ctx sdk.Context) (supplies []types.AssetSupply) {
	k.IterateAssetSupplies(
		ctx, func(supply types.AssetSupply) bool {
			supplies = append(supplies, supply)
			return false
		},
	)
	return
}

// GetPreviousBlockTime get the blocktime for the previous block
func (k Keeper) GetPreviousBlockTime(ctx sdk.Context) (blockTime time.Time, found bool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.PreviousBlockTimeKey)
	if b == nil {
		return time.Time{}, false
	}

	var timestamp gogotypes.Timestamp
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &timestamp)
	blockTime, _ = gogotypes.TimestampFromProto(&timestamp)
	return blockTime, true
}

// SetPreviousBlockTime set the time of the previous block
func (k Keeper) SetPreviousBlockTime(ctx sdk.Context, blockTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	timestamp, _ := gogotypes.TimestampProto(blockTime)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(timestamp)
	store.Set(types.PreviousBlockTimeKey, bz)
}
