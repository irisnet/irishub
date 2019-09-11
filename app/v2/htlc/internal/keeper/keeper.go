package keeper

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/app/v2/htlc/internal/types"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/tendermint/tendermint/crypto"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	bk       types.BankKeeper

	// codespace
	codespace sdk.CodespaceType
	// params subspace
	paramSpace params.Subspace
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bk types.BankKeeper, codespace sdk.CodespaceType, paramSpace params.Subspace) Keeper {
	return Keeper{
		storeKey:   key,
		cdc:        cdc,
		bk:         bk,
		codespace:  codespace,
		paramSpace: paramSpace.WithTypeTable(types.ParamTypeTable()),
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

// CreateHTLC creates a HTLC
func (k Keeper) CreateHTLC(ctx sdk.Context, htlc types.HTLC, hashLock []byte) (sdk.Tags, sdk.Error) {
	// check if the hash lock already exists
	if k.HasHashLock(ctx, hashLock) {
		return nil, types.ErrHashLockAlreadyExists(types.DefaultCodespace, fmt.Sprintf("the hash lock already exists: %s", hex.EncodeToString(hashLock)))
	}

	// transfer the specified tokens to a dedicated HTLC Address
	htlcAddr := getHTLCAddress(htlc.Amount.Denom)
	if _, err := k.bk.SendCoins(ctx, htlc.Sender, htlcAddr, sdk.Coins{htlc.Amount}); err != nil {
		return nil, err
	}

	// add to coinflow
	ctx.CoinFlowTags().AppendCoinFlowTag(ctx, htlc.Sender.String(), htlcAddr.String(), htlc.Amount.String(), sdk.CoinHTLCCreateFlow, "")

	// set the htlc
	k.SetHTLC(ctx, htlc, hashLock)

	// add to the expiration queue
	k.AddHTLCToExpireQueue(ctx, htlc.ExpireHeight, hashLock)

	createTags := sdk.NewTags(
		types.TagSender, []byte(htlc.Sender.String()),
		types.TagReceiver, []byte(htlc.Receiver.String()),
		types.TagReceiverOnOtherChain, htlc.ReceiverOnOtherChain,
		types.TagHashLock, []byte(hex.EncodeToString(hashLock)),
	)

	return createTags, nil
}

func (k Keeper) ClaimHTLC(ctx sdk.Context, secret []byte, hashLock []byte) (sdk.Tags, sdk.Error) {
	// get the htlc
	htlc, err := k.GetHTLC(ctx, hashLock)
	if err != nil {
		return nil, err
	}

	// check if the htlc is open
	if htlc.State != types.StateOpen {
		return nil, types.ErrStateIsNotOpen(k.codespace, fmt.Sprintf("the htlc is not open"))
	}

	// check if the secret matches with the hash lock
	if !bytes.Equal(getHashLock(secret, htlc.Timestamp), hashLock) {
		return nil, types.ErrInvalidSecret(k.codespace, fmt.Sprintf("invalid secret: %s", hex.EncodeToString(secret)))
	}

	// do claim
	htlcAddr := getHTLCAddress(htlc.Amount.Denom)
	if _, err := k.bk.SendCoins(ctx, htlcAddr, htlc.Receiver, sdk.Coins{htlc.Amount}); err != nil {
		return nil, err
	}

	// update the secret and state in HTLC
	htlc.Secret = secret
	htlc.State = types.StateCompleted
	k.SetHTLC(ctx, htlc, hashLock)

	// delete from the expiration queue
	k.DeleteHTLCFromExpireQueue(ctx, uint64(ctx.BlockHeight()), hashLock)

	// add to coinflow
	ctx.CoinFlowTags().AppendCoinFlowTag(ctx, htlcAddr.String(), htlc.Receiver.String(), htlc.Amount.String(), sdk.CoinHTLCClaimFlow, "")

	calimTags := sdk.NewTags(
		types.TagSender, []byte(htlc.Sender.String()),
		types.TagReceiver, []byte(htlc.Receiver.String()),
		types.TagHashLock, []byte(hex.EncodeToString(hashLock)),
		types.TagSecret, []byte(hex.EncodeToString(secret)),
	)

	return calimTags, nil
}

func (k Keeper) RefundHTLC(ctx sdk.Context, hashLock []byte) (sdk.Tags, sdk.Error) {
	// get the htlc
	htlc, err := k.GetHTLC(ctx, hashLock)
	if err != nil {
		return nil, err
	}

	// check if the htlc is expired
	if htlc.State != types.StateExpired {
		return nil, types.ErrStateIsNotOpen(k.codespace, fmt.Sprintf("the htlc is not expired"))
	}

	// do refund
	htlcAddr := getHTLCAddress(htlc.Amount.Denom)
	if _, err := k.bk.SendCoins(ctx, htlcAddr, htlc.Sender, sdk.Coins{htlc.Amount}); err != nil {
		return nil, err
	}

	// update the state in HTLC
	htlc.State = types.StateRefunded
	k.SetHTLC(ctx, htlc, hashLock)

	// add to coinflow
	ctx.CoinFlowTags().AppendCoinFlowTag(ctx, htlcAddr.String(), htlc.Sender.String(), htlc.Amount.String(), sdk.CoinHTLCRefundFlow, "")

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

// SetHTLC stores the htlc
func (k Keeper) SetHTLC(ctx sdk.Context, htlc types.HTLC, hashLock []byte) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(htlc)
	store.Set(KeyHTLC(hashLock), bz)
}

// GetHTLC retrieves the htlc by the specified hash lock
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

// AddHTLCToExpireQueue adds the htlc to the expiration queue
func (k Keeper) AddHTLCToExpireQueue(ctx sdk.Context, expireHeight uint64, hashLock []byte) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(hashLock)
	store.Set(KeyHTLCExpireQueue(expireHeight, hashLock), bz)
}

// DeleteHTLCFromExpireQueue removes the htlc from the expiration queue
func (k Keeper) DeleteHTLCFromExpireQueue(ctx sdk.Context, expireHeight uint64, hashLock []byte) {
	store := ctx.KVStore(k.storeKey)

	// delete the key
	store.Delete(KeyHTLCExpireQueue(expireHeight, hashLock))
}

// getHTLCAddress returns a dedicated address for locking tokens by the specified denom
func getHTLCAddress(denom string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(denom)))
}

// getHashLock calculates the hash lock from the given secret and timestamp
func getHashLock(secret []byte, timestamp uint64) []byte {
	if timestamp > 0 {
		return sdk.SHA256(append(secret, sdk.Uint64ToBigEndian(timestamp)...))
	}

	return sdk.SHA256(secret)
}

// IterateHTLCExpireQueueByHeight iterates the HTLC expiration queue by the specified height
func (k Keeper) IterateHTLCExpireQueueByHeight(ctx sdk.Context, height uint64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, KeyHTLCExpireQueueSubspace(height))
}
