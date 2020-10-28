package keeper

import (
	gogotypes "github.com/gogo/protobuf/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/service/types"
)

// RefundServiceFee refunds the service fee to the specified consumer
func (k Keeper) RefundServiceFee(ctx sdk.Context, consumer sdk.AccAddress, serviceFee sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.RequestAccName, consumer, serviceFee)
}

// AddEarnedFee adds the earned fee for the given provider
func (k Keeper) AddEarnedFee(ctx sdk.Context, provider sdk.AccAddress, fee sdk.Coins) error {
	taxRate := k.ServiceFeeTax(ctx)

	taxCoins := sdk.Coins{}
	for _, coin := range fee {
		taxAmount := sdk.NewDecFromInt(coin.Amount).Mul(taxRate).TruncateInt()
		taxCoins = taxCoins.Add(sdk.NewCoin(coin.Denom, taxAmount))
	}

	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.RequestAccName, k.feeCollectorName, taxCoins); err != nil {
		return err
	}

	earnedFee, hasNeg := fee.SafeSub(taxCoins)
	if hasNeg {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "%s is less than %s", fee, taxCoins)
	}

	// add the provider's earned fees
	earnedFees, _ := k.GetEarnedFees(ctx, provider)
	k.SetEarnedFees(ctx, provider, earnedFees.Add(earnedFee...))

	// add the owner's earned fees
	owner, _ := k.GetOwner(ctx, provider)
	ownerEarnedFees, _ := k.GetOwnerEarnedFees(ctx, owner)
	k.SetOwnerEarnedFees(ctx, owner, ownerEarnedFees.Add(earnedFee...))

	return nil
}

// SetEarnedFees sets the earned fees for the specified provider
func (k Keeper) SetEarnedFees(ctx sdk.Context, provider sdk.AccAddress, fees sdk.Coins) {
	store := ctx.KVStore(k.storeKey)

	for i := range fees {
		bz := k.cdc.MustMarshalBinaryBare(&fees[i])
		store.Set(types.GetEarnedFeesKey(provider, fees[i].Denom), bz)
	}
}

// GetEarnedFees retrieves the earned fees of the specified provider
func (k Keeper) GetEarnedFees(ctx sdk.Context, provider sdk.AccAddress) (fees sdk.Coins, found bool) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.GetEarnedFeesSubspace(provider))

	fees = sdk.NewCoins()
	for ; iterator.Valid(); iterator.Next() {
		var balance sdk.Coin
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &balance)
		fees = fees.Add(balance)
	}

	return fees, true
}

// DeleteEarnedFees removes the earned fees of the specified provider
func (k Keeper) DeleteEarnedFees(ctx sdk.Context, provider sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetEarnedFeesSubspace(provider))

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}

// SetOwnerEarnedFees sets the earned fees for the specified owner
func (k Keeper) SetOwnerEarnedFees(ctx sdk.Context, owner sdk.AccAddress, fees sdk.Coins) {
	store := ctx.KVStore(k.storeKey)

	for i := range fees {
		bz := k.cdc.MustMarshalBinaryBare(&fees[i])
		store.Set(types.GetOwnerEarnedFeesKey(owner, fees[i].Denom), bz)
	}
}

// GetOwnerEarnedFees retrieves the earned fees of the specified owner
func (k Keeper) GetOwnerEarnedFees(ctx sdk.Context, owner sdk.AccAddress) (fees sdk.Coins, found bool) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.GetOwnerEarnedFeesSubspace(owner))

	fees = sdk.NewCoins()
	for ; iterator.Valid(); iterator.Next() {
		var balance sdk.Coin
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &balance)
		fees = fees.Add(balance)
	}

	return fees, true
}

// DeleteOwnerEarnedFees removes the earned fees of the specified owner
func (k Keeper) DeleteOwnerEarnedFees(ctx sdk.Context, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetOwnerEarnedFeesSubspace(owner))

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}

// WithdrawEarnedFees withdraws the earned fees of the specified provider or owner
func (k Keeper) WithdrawEarnedFees(ctx sdk.Context, owner, provider sdk.AccAddress) error {
	if !provider.Empty() {
		providerOwner, _ := k.GetOwner(ctx, provider)
		if !owner.Equals(providerOwner) {
			return sdkerrors.Wrap(types.ErrNotAuthorized, "owner not matching")
		}
	}

	ownerEarnedFees, found := k.GetOwnerEarnedFees(ctx, owner)
	if !found {
		return sdkerrors.Wrap(types.ErrNoEarnedFees, owner.String())
	}

	var withdrawFees sdk.Coins

	if !provider.Empty() {
		earnedFees, found := k.GetEarnedFees(ctx, provider)
		if !found {
			return sdkerrors.Wrap(types.ErrNoEarnedFees, provider.String())
		}

		k.DeleteEarnedFees(ctx, provider)

		if earnedFees.IsEqual(ownerEarnedFees) {
			k.DeleteOwnerEarnedFees(ctx, owner)
		} else {
			k.SetOwnerEarnedFees(ctx, owner, ownerEarnedFees.Sub(earnedFees))
		}

		withdrawFees = earnedFees
	} else {
		iterator := k.OwnerProvidersIterator(ctx, owner)
		defer iterator.Close()

		for ; iterator.Valid(); iterator.Next() {
			provider := sdk.AccAddress(iterator.Key()[sdk.AddrLen+1:])
			k.DeleteEarnedFees(ctx, provider)
		}

		k.DeleteOwnerEarnedFees(ctx, owner)
		withdrawFees = ownerEarnedFees
	}

	withdrawAddr := k.GetWithdrawAddress(ctx, owner)

	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.RequestAccName, withdrawAddr, withdrawFees)
}

// AllEarnedFeesIterator returns an iterator for all the earned fees
func (k Keeper) AllEarnedFeesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.EarnedFeesKey)
}

// RefundEarnedFees refunds all the earned fees
func (k Keeper) RefundEarnedFees(ctx sdk.Context) error {
	iterator := k.AllEarnedFeesIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		provider := iterator.Key()[1:]

		var earnedFee sdk.Coin
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &earnedFee)

		if err := k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx, types.RequestAccName, provider, sdk.NewCoins(earnedFee),
		); err != nil {
			return err
		}
	}

	return nil
}

// RefundServiceFees refunds the service fees of all the active requests
func (k Keeper) RefundServiceFees(ctx sdk.Context) error {
	iterator := k.AllActiveRequestsIterator(ctx.KVStore(k.storeKey))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var requestID gogotypes.BytesValue
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &requestID)

		request, _ := k.GetRequest(ctx, requestID.Value)

		consumer, err := sdk.AccAddressFromBech32(request.Consumer)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid consumer address (%s)", err)
		}

		if err := k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx, types.RequestAccName, consumer, request.ServiceFee,
		); err != nil {
			return err
		}
	}

	return nil
}
