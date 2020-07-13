package keeper

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irishub/modules/random/types"
)

// Keeper defines the rand module Keeper
type Keeper struct {
	cdc      codec.Marshaler
	storeKey sdk.StoreKey
}

// NewKeeper returns a new rand keeper
func NewKeeper(cdc codec.Marshaler, key sdk.StoreKey) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: key,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("%s", types.ModuleName))
}

// GetCdc returns the cdc
func (k Keeper) GetCdc() codec.Marshaler {
	return k.cdc
}

// RequestRandom requests a random number
func (k Keeper) RequestRandom(ctx sdk.Context, consumer sdk.AccAddress, blockInterval uint64) (types.Request, error) {
	currentHeight := ctx.BlockHeight()
	destHeight := currentHeight + int64(blockInterval)

	// get tx hash
	txHash := types.SHA256(ctx.TxBytes())

	// build request
	request := types.NewRequest(currentHeight, consumer, txHash)

	// generate the request id
	reqID := types.GenerateRequestID(request)

	// add to the queue
	k.EnqueueRandomRequest(ctx, destHeight, reqID, request)

	return request, nil
}

// SetRandom stores the random number
func (k Keeper) SetRandom(ctx sdk.Context, reqID []byte, rand types.Random) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&rand)
	store.Set(types.KeyRandom(reqID), bz)
}

// EnqueueRandomRequest enqueue the random number request
func (k Keeper) EnqueueRandomRequest(ctx sdk.Context, height int64, reqID []byte, request types.Request) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&request)
	store.Set(types.KeyRandomRequestQueue(height, reqID), bz)
}

// DequeueRandomRequest removes the random number request by the specified height and request id
func (k Keeper) DequeueRandomRequest(ctx sdk.Context, height int64, reqID []byte) {
	store := ctx.KVStore(k.storeKey)

	// delete the key
	store.Delete(types.KeyRandomRequestQueue(height, reqID))
}

// GetRandom retrieves the random number by the specified request id
func (k Keeper) GetRandom(ctx sdk.Context, reqID []byte) (types.Random, error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyRandom(reqID))
	if bz == nil {
		return types.Random{}, sdkerrors.Wrap(types.ErrInvalidReqID, fmt.Sprintf("request id does not exist: %s", hex.EncodeToString(reqID)))
	}

	var rand types.Random
	k.cdc.MustUnmarshalBinaryBare(bz, &rand)

	return rand, nil
}

// IterateRandoms iterates through all the random numbers
func (k Keeper) IterateRandoms(ctx sdk.Context, op func(r types.Random) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.PrefixRandom)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var rand types.Random
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &rand)

		if stop := op(rand); stop {
			break
		}
	}
}

// IterateRandomRequestQueueByHeight iterates the random number request queue by the specified height
func (k Keeper) IterateRandomRequestQueueByHeight(ctx sdk.Context, height int64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.KeyRandomRequestQueueSubspace(height))
}

// IterateRandomRequestQueue iterates through the random number request queue
func (k Keeper) IterateRandomRequestQueue(ctx sdk.Context, op func(h int64, r types.Request) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.PrefixRandomRequestQueue)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		keyParts := bytes.Split(iterator.Key(), types.KeyDelimiter)
		height := int64(binary.BigEndian.Uint64(keyParts[1]))

		var request types.Request
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &request)

		if stop := op(height, request); stop {
			break
		}
	}
}
