package keeper

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v2/htlc/internal/types"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

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
func (k Keeper) CreateHTLC(ctx sdk.Context, htlc types.HTLC, hashLock []byte) (sdk.Tags, sdk.Error) {
	// check if the hash lock already exists
	if k.HasHashLock(ctx, hashLock) {
		return nil, types.ErrHashLockAlreadyExists(types.DefaultCodespace, fmt.Sprintf("the hash lock already exists: %s", hex.EncodeToString(hashLock)))
	}

	// transfer the specified tokens to a dedicated HTLC Address
	if _, err := k.bk.SendCoins(ctx, htlc.Sender, auth.HTLCLockedCoinsAccAddr, htlc.Amount); err != nil {
		return nil, err
	}

	// add to coinflow
	ctx.CoinFlowTags().AppendCoinFlowTag(ctx, htlc.Sender.String(), auth.HTLCLockedCoinsAccAddr.String(), htlc.Amount.String(), sdk.CoinHTLCCreateFlow, "")

	// set the HTLC
	k.SetHTLC(ctx, htlc, hashLock)

	// add to the expiration queue
	k.AddHTLCToExpireQueue(ctx, htlc.ExpireHeight, hashLock)

	createTags := sdk.NewTags(
		types.TagSender, []byte(htlc.Sender.String()),
		types.TagReceiver, []byte(htlc.To.String()),
		types.TagReceiverOnOtherChain, []byte(htlc.ReceiverOnOtherChain),
		types.TagAmount, []byte(htlc.Amount.String()),
		types.TagHashLock, []byte(hex.EncodeToString(hashLock)),
	)

	return createTags, nil
}

func (k Keeper) ClaimHTLC(ctx sdk.Context, hashLock []byte, secret []byte) (sdk.Tags, sdk.Error) {
	// get the HTLC
	htlc, err := k.GetHTLC(ctx, hashLock)
	if err != nil {
		return nil, err
	}

	// check if the HTLC is open
	if htlc.State != types.OPEN {
		return nil, types.ErrStateIsNotOpen(k.codespace, fmt.Sprintf("the HTLC is not open"))
	}

	// check if the secret matches with the hash lock
	if !bytes.Equal(GetHashLock(secret, htlc.Timestamp), hashLock) {
		return nil, types.ErrInvalidSecret(k.codespace, fmt.Sprintf("invalid secret: %s", hex.EncodeToString(secret)))
	}

	// do the claim
	if _, err := k.bk.SendCoins(ctx, auth.HTLCLockedCoinsAccAddr, htlc.To, htlc.Amount); err != nil {
		return nil, err
	}

	// update the secret and state in HTLC
	htlc.Secret = secret
	htlc.State = types.COMPLETED
	k.SetHTLC(ctx, htlc, hashLock)

	// delete from the expiration queue
	k.DeleteHTLCFromExpireQueue(ctx, htlc.ExpireHeight, hashLock)

	// add to coinflow
	ctx.CoinFlowTags().AppendCoinFlowTag(ctx, auth.HTLCLockedCoinsAccAddr.String(), htlc.To.String(), htlc.Amount.String(), sdk.CoinHTLCClaimFlow, "")

	claimTags := sdk.NewTags(
		types.TagSender, []byte(htlc.Sender.String()),
		types.TagReceiver, []byte(htlc.To.String()),
		types.TagHashLock, []byte(hex.EncodeToString(hashLock)),
		types.TagSecret, []byte(hex.EncodeToString(secret)),
	)

	return claimTags, nil
}

func (k Keeper) RefundHTLC(ctx sdk.Context, hashLock []byte) (sdk.Tags, sdk.Error) {
	// get the HTLC
	htlc, err := k.GetHTLC(ctx, hashLock)
	if err != nil {
		return nil, err
	}

	// check if the HTLC is expired
	if htlc.State != types.EXPIRED {
		return nil, types.ErrStateIsNotOpen(k.codespace, fmt.Sprintf("the htlc is not expired"))
	}

	// do the refund
	if _, err := k.bk.SendCoins(ctx, auth.HTLCLockedCoinsAccAddr, htlc.Sender, htlc.Amount); err != nil {
		return nil, err
	}

	// update the state in HTLC
	htlc.State = types.REFUNDED
	k.SetHTLC(ctx, htlc, hashLock)

	// add to coinflow
	ctx.CoinFlowTags().AppendCoinFlowTag(ctx, auth.HTLCLockedCoinsAccAddr.String(), htlc.Sender.String(), htlc.Amount.String(), sdk.CoinHTLCRefundFlow, "")

	refundTags := sdk.NewTags(
		types.TagSender, []byte(htlc.Sender.String()),
		types.TagHashLock, []byte(hex.EncodeToString(hashLock)),
	)

	return refundTags, nil
}

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
		return sdk.SHA256(append(secret, sdk.Uint64ToBigEndian(timestamp)...))
	}

	return sdk.SHA256(secret)
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
