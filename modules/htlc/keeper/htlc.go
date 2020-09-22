package keeper

import (
	"bytes"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/htlc/types"
)

// CreateHTLC creates an HTLC
func (k Keeper) CreateHTLC(
	ctx sdk.Context,
	sender,
	to sdk.AccAddress,
	receiverOnOtherChain string,
	amount sdk.Coins,
	hashLock tmbytes.HexBytes,
	timestamp,
	timeLock uint64,
) error {
	// check if the HTLC already exists
	if k.HasHTLC(ctx, hashLock) {
		return sdkerrors.Wrap(types.ErrHTLCExists, hashLock.String())
	}

	// transfer the specified tokens to the HTLC module account
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, amount)
	if err != nil {
		return err
	}

	expirationHeight := uint64(ctx.BlockHeight()) + timeLock
	state := types.Open
	htlc := types.NewHTLC(sender, to, receiverOnOtherChain, amount, nil, timestamp, expirationHeight, state)

	// set the HTLC
	k.SetHTLC(ctx, htlc, hashLock)

	// add to the expiration queue
	k.AddHTLCToExpiredQueue(ctx, htlc.ExpirationHeight, hashLock)

	return nil
}

// ClaimHTLC claims the specified HTLC with the given secret
func (k Keeper) ClaimHTLC(ctx sdk.Context, hashLock tmbytes.HexBytes, secret tmbytes.HexBytes) error {
	// query the HTLC
	htlc, found := k.GetHTLC(ctx, hashLock)
	if !found {
		return sdkerrors.Wrap(types.ErrUnknownHTLC, hashLock.String())
	}

	// check if the HTLC is open
	if htlc.State != types.Open {
		return sdkerrors.Wrap(types.ErrHTLCNotOpen, hashLock.String())
	}

	// check if the secret matches with the hash lock
	if !bytes.Equal(types.GetHashLock(secret, htlc.Timestamp), hashLock) {
		return sdkerrors.Wrap(types.ErrInvalidSecret, secret.String())
	}

	// do the claim
	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, htlc.To, htlc.Amount)
	if err != nil {
		return err
	}

	// update the secret and state of the HTLC
	htlc.Secret = secret
	htlc.State = types.Completed
	k.SetHTLC(ctx, htlc, hashLock)

	// delete from the expiration queue
	k.DeleteHTLCFromExpiredQueue(ctx, htlc.ExpirationHeight, hashLock)

	return nil
}

// RefundHTLC refunds the specified HTLC
func (k Keeper) RefundHTLC(ctx sdk.Context, hashLock tmbytes.HexBytes) error {
	// query the HTLC
	htlc, found := k.GetHTLC(ctx, hashLock)
	if !found {
		return sdkerrors.Wrap(types.ErrUnknownHTLC, hashLock.String())
	}

	// check if the HTLC is expired
	if htlc.State != types.Expired {
		return sdkerrors.Wrap(types.ErrHTLCNotExpired, hashLock.String())
	}

	// do the refund
	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, htlc.Sender, htlc.Amount)
	if err != nil {
		return err
	}

	// update the state of the HTLC
	htlc.State = types.Refunded
	k.SetHTLC(ctx, htlc, hashLock)

	return nil
}

// HasHTLC checks if the given HTLC exists
func (k Keeper) HasHTLC(ctx sdk.Context, hashLock tmbytes.HexBytes) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetHTLCKey(hashLock))
}

// SetHTLC sets the given HTLC
func (k Keeper) SetHTLC(ctx sdk.Context, htlc types.HTLC, hashLock tmbytes.HexBytes) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&htlc)
	store.Set(types.GetHTLCKey(hashLock), bz)
}

// GetHTLC retrieves the specified HTLC
func (k Keeper) GetHTLC(ctx sdk.Context, hashLock tmbytes.HexBytes) (htlc types.HTLC, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetHTLCKey(hashLock))
	if bz == nil {
		return htlc, false
	}

	k.cdc.MustUnmarshalBinaryBare(bz, &htlc)

	return htlc, true
}

// AddHTLCToExpiredQueue adds the specified HTLC to the expiration queue
func (k Keeper) AddHTLCToExpiredQueue(ctx sdk.Context, expirationHeight uint64, hashLock tmbytes.HexBytes) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetHTLCExpiredQueueKey(expirationHeight, hashLock), []byte{})
}

// DeleteHTLCFromExpiredQueue removes the specified HTLC from the expiration queue
func (k Keeper) DeleteHTLCFromExpiredQueue(ctx sdk.Context, expirationHeight uint64, hashLock tmbytes.HexBytes) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetHTLCExpiredQueueKey(expirationHeight, hashLock))
}

// IterateHTLCs iterates through the HTLCs
func (k Keeper) IterateHTLCs(
	ctx sdk.Context,
	op func(hlock tmbytes.HexBytes, h types.HTLC) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.HTLCKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		hashLock := tmbytes.HexBytes(iterator.Key()[1:])

		var htlc types.HTLC
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &htlc)

		if stop := op(hashLock, htlc); stop {
			break
		}
	}
}

// IterateHTLCExpiredQueueByHeight iterates through the HTLC expiration queue by the specified height
func (k Keeper) IterateHTLCExpiredQueueByHeight(
	ctx sdk.Context,
	height uint64,
	op func(hlock tmbytes.HexBytes, h types.HTLC) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.GetHTLCExpiredQueueSubspace(height))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		hashLock := tmbytes.HexBytes(iterator.Key()[9:])
		htlc, _ := k.GetHTLC(ctx, hashLock)

		if stop := op(hashLock, htlc); stop {
			break
		}
	}
}
