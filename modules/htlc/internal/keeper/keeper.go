package keeper

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/htlc/internal/types"
)

var (
	// HTLCLockedCoinsAccAddr store All HTLC locked coins
	HTLCLockedCoinsAccAddr = sdk.AccAddress(crypto.AddressHash([]byte("HTLCLockedCoins")))
)

// Keeper defines the HTLC module Keeper
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	bk       types.BankKeeper

	// codespace
	codespace sdk.CodespaceType
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bk types.BankKeeper, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		storeKey:  key,
		cdc:       cdc,
		bk:        bk,
		codespace: codespace,
	}
}

// Codespace returns the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// GetCdc returns the cdc
func (k Keeper) GetCdc() *codec.Codec {
	return k.cdc
}

// CreateHTLC creates an HTLC
func (k Keeper) CreateHTLC(ctx sdk.Context, htlc types.HTLC, hashLock []byte) (sdk.Event, sdk.Error) {
	// check if the hash lock already exists
	if k.HasHashLock(ctx, hashLock) {
		return sdk.Event{}, types.ErrHashLockAlreadyExists(types.DefaultCodespace, fmt.Sprintf("the hash lock already exists: %s", hex.EncodeToString(hashLock)))
	}

	// transfer the specified tokens to a dedicated HTLC Address
	if err := k.bk.SendCoins(ctx, htlc.Sender, HTLCLockedCoinsAccAddr, htlc.Amount); err != nil {
		return sdk.Event{}, err
	}

	// set the HTLC
	k.SetHTLC(ctx, htlc, hashLock)

	// add to the expiration queue
	k.AddHTLCToExpireQueue(ctx, htlc.ExpireHeight, hashLock)

	createEvent := sdk.NewEvent(
		types.EventTypeCreateHTLC,
		sdk.NewAttribute(types.AttributeValueSender, htlc.Sender.String()),
		sdk.NewAttribute(types.AttributeValueReceiver, htlc.To.String()),
		sdk.NewAttribute(types.AttributeValueReceiverOnOtherChain, htlc.ReceiverOnOtherChain),
		sdk.NewAttribute(types.AttributeValueAmount, htlc.Amount.String()),
		sdk.NewAttribute(types.AttributeValueHashLock, hex.EncodeToString(hashLock)),
	)

	return createEvent, nil
}

// ClaimHTLC claim an HTLC
func (k Keeper) ClaimHTLC(ctx sdk.Context, hashLock []byte, secret []byte) (sdk.Event, sdk.Error) {
	// get the HTLC
	htlc, err := k.GetHTLC(ctx, hashLock)
	if err != nil {
		return sdk.Event{}, err
	}

	// check if the HTLC is open
	if htlc.State != types.OPEN {
		return sdk.Event{}, types.ErrStateIsNotOpen(k.codespace, fmt.Sprintf("the HTLC is not open"))
	}

	// check if the secret matches with the hash lock
	if !bytes.Equal(GetHashLock(secret, htlc.Timestamp), hashLock) {
		return sdk.Event{}, types.ErrInvalidSecret(k.codespace, fmt.Sprintf("invalid secret: %s", hex.EncodeToString(secret)))
	}

	// do the claim
	if err := k.bk.SendCoins(ctx, HTLCLockedCoinsAccAddr, htlc.To, htlc.Amount); err != nil {
		return sdk.Event{}, err
	}

	// update the secret and state in HTLC
	htlc.Secret = secret
	htlc.State = types.COMPLETED
	k.SetHTLC(ctx, htlc, hashLock)

	// delete from the expiration queue
	k.DeleteHTLCFromExpireQueue(ctx, htlc.ExpireHeight, hashLock)

	claimEvent := sdk.NewEvent(
		types.EventTypeClaimHTLC,
		sdk.NewAttribute(types.AttributeValueSender, htlc.Sender.String()),
		sdk.NewAttribute(types.AttributeValueReceiver, htlc.To.String()),
		sdk.NewAttribute(types.AttributeValueHashLock, hex.EncodeToString(hashLock)),
		sdk.NewAttribute(types.AttributeValueSecret, hex.EncodeToString(secret)),
	)

	return claimEvent, nil
}

// RefundHTLC refund an HTLC
func (k Keeper) RefundHTLC(ctx sdk.Context, hashLock []byte) (sdk.Event, sdk.Error) {
	// get the HTLC
	htlc, err := k.GetHTLC(ctx, hashLock)
	if err != nil {
		return sdk.Event{}, err
	}

	// check if the HTLC is expired
	if htlc.State != types.EXPIRED {
		return sdk.Event{}, types.ErrStateIsNotOpen(k.codespace, fmt.Sprintf("the htlc is not expired"))
	}

	// do the refund
	if err := k.bk.SendCoins(ctx, HTLCLockedCoinsAccAddr, htlc.Sender, htlc.Amount); err != nil {
		return sdk.Event{}, err
	}

	// update the state in HTLC
	htlc.State = types.REFUNDED
	k.SetHTLC(ctx, htlc, hashLock)

	refundEvent := sdk.NewEvent(
		types.EventTypeRefundHTLC,
		sdk.NewAttribute(types.AttributeValueSender, htlc.Sender.String()),
		sdk.NewAttribute(types.AttributeValueHashLock, hex.EncodeToString(hashLock)),
	)

	return refundEvent, nil
}

// HasHashLock returns whether the hashlock already exists
func (k Keeper) HasHashLock(ctx sdk.Context, hashLock []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(KeyHTLC(hashLock))
}

// SetHTLC stores the HTLC
func (k Keeper) SetHTLC(ctx sdk.Context, htlc types.HTLC, hashLock []byte) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(htlc)
	store.Set(KeyHTLC(hashLock), bz)
}

// GetHTLC retrieves the HTLC by the specified hash lock
func (k Keeper) GetHTLC(ctx sdk.Context, hashLock []byte) (types.HTLC, sdk.Error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(KeyHTLC(hashLock))
	if bz == nil {
		return types.HTLC{}, types.ErrInvalidHashLock(k.codespace, fmt.Sprintf("invalid hash lock: %s", hex.EncodeToString(hashLock)))
	}

	var htlc types.HTLC
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &htlc)

	return htlc, nil
}

// AddHTLCToExpireQueue adds the specified HTLC to the expiration queue
func (k Keeper) AddHTLCToExpireQueue(ctx sdk.Context, expireHeight uint64, hashLock []byte) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(hashLock)
	store.Set(KeyHTLCExpireQueue(expireHeight, hashLock), bz)
}

// DeleteHTLCFromExpireQueue removes the specified HTLC from the expiration queue
func (k Keeper) DeleteHTLCFromExpireQueue(ctx sdk.Context, expireHeight uint64, hashLock []byte) {
	store := ctx.KVStore(k.storeKey)

	// delete the key
	store.Delete(KeyHTLCExpireQueue(expireHeight, hashLock))
}

// GetHashLock calculates the hash lock from the given secret and timestamp
func GetHashLock(secret []byte, timestamp uint64) []byte {
	if timestamp > 0 {
		return types.SHA256(append(secret, sdk.Uint64ToBigEndian(timestamp)...))
	}

	return types.SHA256(secret)
}

// IterateHTLCs iterates through the HTLCs
func (k Keeper) IterateHTLCs(ctx sdk.Context, op func(hlock []byte, h types.HTLC) (stop bool)) {
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
