package keeper

import (
	"fmt"

	"github.com/irisnet/irishub/app/v1/rand/internal/types"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec

	// codespace
	codespace sdk.CodespaceType
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		storeKey:  key,
		cdc:       cdc,
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

// RequestRand requests a random number
func (k Keeper) RequestRand(ctx sdk.Context, consumer sdk.AccAddress, blockInterval uint64) (sdk.Tags, sdk.Error) {
	currentHeight := ctx.BlockHeight()
	destHeight := currentHeight + int64(blockInterval)

	// get tx hash
	txHash := sdk.SHA256(ctx.TxBytes())

	// build request
	request := types.NewRequest(currentHeight, consumer, txHash)

	// generate the request id
	reqID := types.GenerateRequestID(request)

	// add to the queue
	k.EnqueueRandRequest(ctx, destHeight, reqID, request)

	reqTags := sdk.NewTags(
		types.TagReqID, []byte(reqID),
	)

	return reqTags, nil
}

// SetRand stores the random number
func (k Keeper) SetRand(ctx sdk.Context, reqID string, rand types.Rand) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(rand)
	store.Set(KeyRand(reqID), bz)
}

// EnqueueRandRequest enqueue the random number request
func (k Keeper) EnqueueRandRequest(ctx sdk.Context, height int64, reqID string, request types.Request) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(reqID)
	store.Set(KeyRandRequestQueue(height, reqID), bz)
}

// DequeueRandRequest removes the random number request by the specified height and request id
func (k Keeper) DequeueRandRequest(ctx sdk.Context, height int64, reqID string) {
	store := ctx.KVStore(k.storeKey)

	// delete the key
	store.Delete(KeyRandRequestQueue(height, reqID))
}

// GetRand retrieves the random number by the specified request id
func (k Keeper) GetRand(ctx sdk.Context, reqID string) (types.Rand, sdk.Error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(KeyRand(reqID))
	if bz == nil {
		return types.Rand{}, types.ErrInvalidReqID(k.codespace, fmt.Sprintf("the request id does not exist: %s", reqID))
	}

	var rand types.Rand
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &rand)

	return rand, nil
}

// IterateRandRequestQueueByHeight iterates the random number request queue by the specified height
func (k Keeper) IterateRandRequestQueueByHeight(ctx sdk.Context, height int64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, KeyRandRequestQueueSubspace(height))
}

// IterateRands iterates through all the random numbers
func (k Keeper) IterateRands(ctx sdk.Context, op func(r types.Rand) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, PrefixRand)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var rand types.Rand
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &rand)

		if stop := op(rand); stop {
			break
		}
	}
}

// IterateRandRequestQueue iterates through the random number request queue
func (k Keeper) IterateRandRequestQueue(ctx sdk.Context, op func(r types.Request) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, PrefixRandRequestQueue)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var request types.Request
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &request)

		if stop := op(request); stop {
			break
		}
	}
}
