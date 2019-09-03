package keeper

import (
	"encoding/hex"
	"fmt"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/params"
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

// CreateHTLC creates a HTLC
func (k Keeper) CreateHTLC(ctx sdk.Context, htlc types.HTLC, secretHashLock []byte) (sdk.Tags, sdk.Error) {
	// check if the secret hash lock already exists
	if k.HasSecretHashLock(ctx, secretHashLock) {
		return nil, types.ErrSecretHashLockAlreadyExists(types.DefaultCodespace, fmt.Sprintf("the secret hash lock already exists: %s", hex.EncodeToString(secretHashLock)))
	}

	// transfer the specified tokens to HTLCCoinsAccAddr
	_, err := k.bk.SendCoins(ctx, htlc.Sender, auth.HTLCCoinsAccAddr, sdk.Coins{htlc.OutAmount})
	if err != nil {
		return nil, err
	}

	// set the htlc
	k.SetHTLC(ctx, htlc, secretHashLock)

	// add to the expiration queue
	k.AddHTLCToExpireQueue(ctx, htlc.ExpireHeight, secretHashLock)

	createTags := sdk.NewTags(
		types.TagSender, []byte(htlc.Sender),
		types.TagReceiver, []byte(htlc.Sender),
		types.TagReceiverOnOtherChain, []byte(htlc.ReceiverOnOtherChain),
		types.TagOutAmount, []byte(htlc.OutAmount.String()),
		types.TagInAmount, sdk.Uint64ToBigEndian(htlc.InAmount),
		types.TagSecretHashLock, []byte(hex.EncodeToString(secretHashLock)),
		types.TagTimestamp, sdk.Uint64ToBigEndian(htlc.Timestamp),
		types.TagExpireHeight, sdk.Uint64ToBigEndian(htlc.ExpireHeight),
	)

	return createTags, nil
}

func (k Keeper) HasSecretHashLock(ctx sdk.Context, secretHashLock []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(KeyHTLC(secretHashLock))
}

// SetHTLC stores the htlc
func (k Keeper) SetHTLC(ctx sdk.Context, htlc types.HTLC, secretHashLock []byte) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(htlc)
	store.Set(KeyHTLC(secretHashLock), bz)
}

// GetHTLC retrieves the htlc by the specified secret hash lock
func (k Keeper) GetHTLC(ctx sdk.Context, secretHashLock []byte) (types.HTLC, sdk.Error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(KeyHTLC(secretHashLock))
	if bz == nil {
		return types.HTLC{}, types.ErrInvalidSecretHashLock(k.codespace, fmt.Sprintf("invalid secret hash lock: %s", hex.EncodeToString(secretHashLock)))
	}

	var htlc types.HTLC
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &htlc)

	return htlc, nil
}

// AddHTLCToExpireQueue adds the htlc to the expiration queue
func (k Keeper) AddHTLCToExpireQueue(ctx sdk.Context, expireHeight uint64, secretHashLock []byte) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(secretHashLock)
	store.Set(KeyHTLCExpireQueue(expireHeight, secretHashLock), bz)
}

// DeleteHTLCFromExpireQueue removes the htlc from the expiration queue
func (k Keeper) DeleteHTLCFromExpireQueue(ctx sdk.Context, expireHeight uint64, secretHashLock []byte) {
	store := ctx.KVStore(k.storeKey)

	// delete the key
	store.Delete(KeyHTLCExpireQueue(expireHeight, secretHashLock))
}
