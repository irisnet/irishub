package rand

import (
	"fmt"

	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/app/v1/rand/tags"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec

	// codespace
	codespace sdk.CodespaceType

	// params subspace
	paramSpace params.Subspace
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, codespace sdk.CodespaceType, paramSpace params.Subspace) Keeper {
	return Keeper{
		storeKey:   key,
		cdc:        cdc,
		codespace:  codespace,
		paramSpace: paramSpace.WithTypeTable(ParamTypeTable()),
	}
}

// Codespace returns the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// RequestRand requests a random number
func (k Keeper) RequestRand(ctx sdk.Context, consumer sdk.AccAddress) (sdk.Tags, sdk.Error) {
	currentHeight := ctx.BlockHeight()
	destHeight := currentHeight + BlockNumAfter

	// build a request
	request := NewRequest(currentHeight, consumer)
	// generate the request id
	reqID := generateRequestID(request)

	// set the initial rand
	k.SetRand(ctx, reqID, NewRand(destHeight, consumer, sdk.ZeroDec()))

	// set the request
	k.SetRandRequest(ctx, reqID, request)

	// add to the queue
	k.EnqueueRandRequest(ctx, destHeight, reqID)

	reqTags := sdk.NewTags(
		tags.ReqID, []byte(reqID),
		tags.Consumer, consumer,
	)

	return reqTags, nil
}

// SetRand stores the random number
func (k Keeper) SetRand(ctx sdk.Context, reqID string, rand Rand) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(rand)
	store.Set(KeyRand(reqID), bz)
}

// SetRandRequest stores the random number request
func (k Keeper) SetRandRequest(ctx sdk.Context, reqID string, request Request) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(request)
	store.Set(KeyRandRequest(reqID), bz)
}

// EnqueueRandRequest enqueue the random number request
func (k Keeper) EnqueueRandRequest(ctx sdk.Context, height int64, reqID string) {
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
func (k Keeper) GetRand(ctx sdk.Context, reqID string) (Rand, sdk.Error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(KeyRand(reqID))
	if bz == nil {
		return Rand{}, ErrInvalidReqID(k.codespace, fmt.Sprintf("the request id does not exist: %s", reqID))
	}

	var rand Rand
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &rand)

	return rand, nil
}

// GetRandRequest retrieves the random number request by the specified request id
func (k Keeper) GetRandRequest(ctx sdk.Context, reqID string) (Request, sdk.Error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(KeyRandRequest(reqID))
	if bz == nil {
		return Request{}, ErrInvalidReqID(k.codespace, fmt.Sprintf("the request id does not exist: %s", reqID))
	}

	var request Request
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &request)

	return request, nil
}

// IterateRandRequestQueueByHeight iterates the random number request queue by the specified height
func (k Keeper) IterateRandRequestQueueByHeight(ctx sdk.Context, height int64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, KeyRandRequestQueueSubspace(height))
}

// IterateRands iterates through all the random numbers
func (k Keeper) IterateRands(ctx sdk.Context, op func(r Rand) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, PrefixRand)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var rand Rand
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &rand)

		if stop := op(rand); stop {
			break
		}
	}
}

// IterateRandRequests iterates through all the random number requests
func (k Keeper) IterateRandRequests(ctx sdk.Context, op func(r Request) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, PrefixRandRequest)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var request Request
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &request)

		if stop := op(request); stop {
			break
		}
	}
}

// IterateRandRequestQueue iterates through the random number request queue
func (k Keeper) IterateRandRequestQueue(ctx sdk.Context, op func(r Request) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, PrefixRandRequestQueue)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var reqID string
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &reqID)

		request, err := k.GetRandRequest(ctx, reqID)
		if err != nil {
			continue
		}

		if stop := op(request); stop {
			break
		}
	}
}

// Init initializes the keeper
func (k Keeper) Init(ctx sdk.Context) {
	k.SetParamSet(ctx, DefaultParams())
}
