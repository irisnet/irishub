package keeper

import (
	"fmt"

	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v3/service/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

// RefundServiceFee refunds the service fee to the specified consumer
func (k Keeper) RefundServiceFee(ctx sdk.Context, consumer sdk.AccAddress, serviceFee sdk.Coins) sdk.Error {
	_, err := k.bk.SendCoins(ctx, auth.ServiceRequestCoinsAccAddr, consumer, serviceFee)

	if !serviceFee.IsZero() {
		ctx.CoinFlowTags().AppendCoinFlowTag(ctx, auth.ServiceRequestCoinsAccAddr.String(),
			consumer.String(), serviceFee.String(), sdk.ServiceFeeRefundFlow, "")
	}

	if err != nil {
		return err
	}

	return nil
}

// AddEarnedFee adds the earned fee for the given provider
func (k Keeper) AddEarnedFee(ctx sdk.Context, provider sdk.AccAddress, fee sdk.Coins) sdk.Error {
	params := k.GetParamSet(ctx)
	taxRate := params.ServiceFeeTax

	taxCoins := sdk.Coins{}
	for _, coin := range fee {
		taxAmount := sdk.NewDecFromInt(coin.Amount).Mul(taxRate).TruncateInt()
		taxCoins = taxCoins.Add(sdk.NewCoins(sdk.NewCoin(coin.Denom, taxAmount)))
	}

	_, err := k.bk.SendCoins(ctx, auth.ServiceRequestCoinsAccAddr, auth.ServiceTaxCoinsAccAddr, taxCoins)
	if err != nil {
		return err
	}

	earnedFee, hasNeg := fee.SafeSub(taxCoins)
	if hasNeg {
		errMsg := fmt.Sprintf("%s is less than %s", fee, taxCoins)
		return sdk.ErrInsufficientFunds(errMsg)
	}

	// add the provider's earned fees
	earnedFees, _ := k.GetEarnedFees(ctx, provider)
	k.SetEarnedFees(ctx, provider, earnedFees.Add(earnedFee))

	// add the owner's earned fees
	owner, _ := k.GetOwner(ctx, provider)
	ownerEarnedFees, _ := k.GetOwnerEarnedFees(ctx, owner)
	k.SetOwnerEarnedFees(ctx, owner, ownerEarnedFees.Add(earnedFee))

	return nil
}

// SetEarnedFees sets the earned fees for the specified provider
func (k Keeper) SetEarnedFees(ctx sdk.Context, provider sdk.AccAddress, fees sdk.Coins) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(fees)
	store.Set(GetEarnedFeesKey(provider), bz)
}

// GetEarnedFees retrieves the earned fees of the specified provider
func (k Keeper) GetEarnedFees(ctx sdk.Context, provider sdk.AccAddress) (fees sdk.Coins, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(GetEarnedFeesKey(provider))
	if bz == nil {
		return fees, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &fees)
	return fees, true
}

// DeleteEarnedFees removes the earned fees of the specified provider
func (k Keeper) DeleteEarnedFees(ctx sdk.Context, provider sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetEarnedFeesKey(provider))
}

// SetOwnerEarnedFees sets the earned fees for the specified owner
func (k Keeper) SetOwnerEarnedFees(ctx sdk.Context, owner sdk.AccAddress, fees sdk.Coins) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(fees)
	store.Set(GetOwnerEarnedFeesKey(owner), bz)
}

// GetOwnerEarnedFees retrieves the earned fees of the specified owner
func (k Keeper) GetOwnerEarnedFees(ctx sdk.Context, owner sdk.AccAddress) (fees sdk.Coins, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(GetOwnerEarnedFeesKey(owner))
	if bz == nil {
		return fees, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &fees)
	return fees, true
}

// DeleteOwnerEarnedFees removes the earned fees of the specified owner
func (k Keeper) DeleteOwnerEarnedFees(ctx sdk.Context, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetOwnerEarnedFeesKey(owner))
}

// WithdrawEarnedFees withdraws the earned fees of the specified provider or owner
func (k Keeper) WithdrawEarnedFees(ctx sdk.Context, owner, provider sdk.AccAddress) sdk.Error {
	if !provider.Empty() {
		providerOwner, _ := k.GetOwner(ctx, provider)
		if !owner.Equals(providerOwner) {
			return types.ErrNotAuthorized(k.codespace, "owner not matching")
		}
	}

	ownerEarnedFees, found := k.GetOwnerEarnedFees(ctx, owner)
	if !found {
		return types.ErrNoEarnedFees(k.codespace, owner)
	}

	var withdrawFees sdk.Coins

	if !provider.Empty() {
		earnedFees, found := k.GetEarnedFees(ctx, provider)
		if !found {
			return types.ErrNoEarnedFees(k.codespace, provider)
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

	_, err := k.bk.SendCoins(ctx, auth.ServiceRequestCoinsAccAddr, withdrawAddr, withdrawFees)
	if err != nil {
		return err
	}

	return nil
}

// WithdrawTax withdraws the service tax to the speicified destination address by the trustee
func (k Keeper) WithdrawTax(ctx sdk.Context, trustee sdk.AccAddress, destAddress sdk.AccAddress, amt sdk.Coins) sdk.Error {
	if _, found := k.gk.GetTrustee(ctx, trustee); !found {
		return types.ErrInvalidTrustee(k.codespace, trustee)
	}

	_, err := k.bk.SendCoins(ctx, auth.ServiceTaxCoinsAccAddr, destAddress, amt)
	if err != nil {
		return err
	}

	return nil
}

// AllEarnedFeesIterator returns an iterator for all the earned fees
func (k Keeper) AllEarnedFeesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, earnedFeesKey)
}

// RefundEarnedFees refunds all the earned fees
func (k Keeper) RefundEarnedFees(ctx sdk.Context) sdk.Error {
	iterator := k.AllEarnedFeesIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		provider := iterator.Key()[1:]

		var earnedFees sdk.Coins
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &earnedFees)

		_, err := k.bk.SendCoins(ctx, auth.ServiceRequestCoinsAccAddr, provider, earnedFees)
		if err != nil {
			return err
		}
	}

	return nil
}

// RefundServiceFees refunds the service fees of all the active requests
func (k Keeper) RefundServiceFees(ctx sdk.Context) sdk.Error {
	iterator := k.AllActiveRequestsIterator(ctx.KVStore(k.storeKey))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var requestID cmn.HexBytes
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &requestID)

		request, _ := k.GetRequest(ctx, requestID)

		_, err := k.bk.SendCoins(ctx, auth.ServiceRequestCoinsAccAddr, request.Consumer, request.ServiceFee)
		if err != nil {
			return err
		}
	}

	return nil
}
