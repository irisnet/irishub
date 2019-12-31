package keeper

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irishub/modules/htlc/internal/types"
)

// Keeper defines the HTLC module Keeper
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	sk       types.SupplyKeeper
}

// NewKeeper returns a new HTLC keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, sk types.SupplyKeeper) Keeper {
	// ensure htlc module account is set
	if addr := sk.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		storeKey: key,
		cdc:      cdc,
		sk:       sk,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("%s", types.ModuleName))
}

// GetCdc returns the cdc
func (k Keeper) GetCdc() *codec.Codec {
	return k.cdc
}

// CreateHTLC creates an HTLC
func (k Keeper) CreateHTLC(ctx sdk.Context, htlc types.HTLC, hashLock types.HTLCHashLock) error {
	// check if the hash lock already exists
	if k.HasHashLock(ctx, hashLock) {
		return sdkerrors.Wrap(types.ErrHashLockAlreadyExists, hex.EncodeToString(hashLock))
	}

	// transfer the specified tokens to HTLC module address
	if err := k.sk.SendCoinsFromAccountToModule(ctx, htlc.Sender, types.ModuleName, htlc.Amount); err != nil {
		return err
	}

	// set the HTLC
	k.SetHTLC(ctx, htlc, hashLock)

	// add to the expiration queue
	k.AddHTLCToExpireQueue(ctx, htlc.ExpireHeight, hashLock)

	return nil
}

// ClaimHTLC claim an HTLC
func (k Keeper) ClaimHTLC(ctx sdk.Context, hashLock types.HTLCHashLock, secret types.HTLCSecret) (string, string, error) {
	// get the HTLC
	htlc, err := k.GetHTLC(ctx, hashLock)
	if err != nil {
		return "", "", types.ErrUnknownHTLC
	}

	// check if the HTLC is open
	if htlc.State != types.OPEN {
		return "", "", types.ErrStateIsNotOpen
	}

	// check if the secret matches with the hash lock
	if !bytes.Equal(types.GetHashLock(secret, htlc.Timestamp), hashLock) {
		return "", "", sdkerrors.Wrap(types.ErrInvalidSecret, hex.EncodeToString(secret))
	}

	// do the claim
	if err := k.sk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, htlc.To, htlc.Amount); err != nil {
		return "", "", err
	}

	// update the secret and state in HTLC
	htlc.Secret = secret
	htlc.State = types.COMPLETED
	k.SetHTLC(ctx, htlc, hashLock)

	// delete from the expiration queue
	k.DeleteHTLCFromExpireQueue(ctx, htlc.ExpireHeight, hashLock)

	return htlc.Sender.String(), htlc.To.String(), nil
}

// RefundHTLC refund an HTLC
func (k Keeper) RefundHTLC(ctx sdk.Context, hashLock types.HTLCHashLock) (string, error) {
	// get the HTLC
	htlc, err := k.GetHTLC(ctx, hashLock)
	if err != nil {
		return "", err
	}

	// check if the HTLC is expired
	if htlc.State != types.EXPIRED {
		return "", types.ErrStateIsNotExpired
	}

	// do the refund
	if err := k.sk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, htlc.Sender, htlc.Amount); err != nil {
		return "", err
	}

	k.DeleteHTLC(ctx, hashLock)

	return htlc.Sender.String(), nil
}

// HasHashLock returns whether the hashlock already exists
func (k Keeper) HasHashLock(ctx sdk.Context, hashLock types.HTLCHashLock) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(KeyHTLC(hashLock))
}

// SetHTLC stores the HTLC
func (k Keeper) SetHTLC(ctx sdk.Context, htlc types.HTLC, hashLock types.HTLCHashLock) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(htlc)
	store.Set(KeyHTLC(hashLock), bz)
}

// DeleteHTLC delete the stored HTLC
func (k Keeper) DeleteHTLC(ctx sdk.Context, hashLock types.HTLCHashLock) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(KeyHTLC(hashLock))
}

// GetHTLC retrieves the HTLC by the specified hash lock
func (k Keeper) GetHTLC(ctx sdk.Context, hashLock types.HTLCHashLock) (types.HTLC, error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(KeyHTLC(hashLock))
	if bz == nil {
		return types.HTLC{}, sdkerrors.Wrap(types.ErrInvalidHashLock, hex.EncodeToString(hashLock))
	}

	var htlc types.HTLC
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &htlc)

	return htlc, nil
}

// AddHTLCToExpireQueue adds the specified HTLC to the expiration queue
func (k Keeper) AddHTLCToExpireQueue(ctx sdk.Context, expireHeight uint64, hashLock types.HTLCHashLock) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(hashLock)
	store.Set(KeyHTLCExpireQueue(expireHeight, hashLock), bz)
}

// DeleteHTLCFromExpireQueue removes the specified HTLC from the expiration queue
func (k Keeper) DeleteHTLCFromExpireQueue(ctx sdk.Context, expireHeight uint64, hashLock types.HTLCHashLock) {
	store := ctx.KVStore(k.storeKey)

	// delete the key
	store.Delete(KeyHTLCExpireQueue(expireHeight, hashLock))
}

// IterateHTLCs iterates through the HTLCs
func (k Keeper) IterateHTLCs(ctx sdk.Context, op func(hlock types.HTLCHashLock, h types.HTLC) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, PrefixHTLC)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		keyParts := bytes.Split(iterator.Key(), KeyDelimiter)
		hashLock := keyParts[1]

		var htlc types.HTLC
		k.GetCdc().MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &htlc)

		if stop := op(hashLock, htlc); stop {
			break
		}
	}
}

// IterateHTLCExpireQueueByHeight iterates through the HTLC expiration queue by the specified height
func (k Keeper) IterateHTLCExpireQueueByHeight(ctx sdk.Context, height uint64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, KeyHTLCExpireQueueSubspace(height))
}
