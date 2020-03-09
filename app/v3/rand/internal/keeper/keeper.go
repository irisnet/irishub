package keeper

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/irisnet/irishub/app/v3/rand/internal/types"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

type Keeper struct {
	storeKey  sdk.StoreKey
	cdc       *codec.Codec
	bk        types.BankKeeper
	sk        types.ServiceKeeper
	codespace sdk.CodespaceType
}

func NewKeeper(
	cdc *codec.Codec,
	key sdk.StoreKey,
	bk types.BankKeeper,
	sk types.ServiceKeeper,
	codespace sdk.CodespaceType,
) Keeper {
	keeper := Keeper{
		storeKey:  key,
		cdc:       cdc,
		bk:        bk,
		sk:        sk,
		codespace: codespace,
	}
	if err := sk.RegisterResponseCallback(types.ModuleName, keeper.HandlerResponse); err != nil {
		panic(err)
	}
	return keeper
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
func (k Keeper) RequestRand(
	ctx sdk.Context,
	consumer sdk.AccAddress,
	blockInterval uint64,
	oracle bool,
	serviceFeeCap sdk.Coins,
) (sdk.Tags, sdk.Error) {
	currentHeight := ctx.BlockHeight()
	destHeight := currentHeight + int64(blockInterval)

	// get tx hash
	txHash := sdk.SHA256(ctx.TxBytes())

	// build request
	request := types.NewRequest(currentHeight, consumer, txHash, oracle, serviceFeeCap)

	// generate the request id
	reqID := types.GenerateRequestID(request)

	// add to the queue
	k.EnqueueRandRequest(ctx, destHeight, reqID, request)

	reqTags := sdk.NewTags(
		types.TagReqID, []byte(hex.EncodeToString(reqID)),
		types.TagRandHeight, []byte(fmt.Sprintf("%d", destHeight)),
	)

	return reqTags, nil
}

// SetRand stores the random number
func (k Keeper) SetRand(ctx sdk.Context, reqID []byte, rand types.Rand) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(rand)
	store.Set(KeyRand(reqID), bz)
}

// EnqueueRandRequest enqueue the random number request
func (k Keeper) EnqueueRandRequest(ctx sdk.Context, height int64, reqID []byte, request types.Request) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(request)
	store.Set(KeyRandRequestQueue(height, reqID), bz)
}

// DequeueRandRequest removes the random number request by the specified height and request id
func (k Keeper) DequeueRandRequest(ctx sdk.Context, height int64, reqID []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(KeyRandRequestQueue(height, reqID))
}

// SetOracleRandRequest stores the oracle rand request
func (k Keeper) SetOracleRandRequest(ctx sdk.Context, requestContextID []byte, request types.Request) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(request)
	store.Set(KeyOracleRandRequest(requestContextID), bz)
}

// GetOracleRandRequest retrieves the oracle rand request by the specified request id
func (k Keeper) GetOracleRandRequest(ctx sdk.Context, requestContextID []byte) (types.Request, sdk.Error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(KeyOracleRandRequest(requestContextID))
	if bz == nil {
		return types.Request{}, types.ErrInvalidRequestContextID(k.codespace, fmt.Sprintf("invalid request context id: %s", hex.EncodeToString(requestContextID)))
	}

	var request types.Request
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &request)

	return request, nil
}

// DeleteOracleRandRequest delete an oracle rand request
func (k Keeper) DeleteOracleRandRequest(ctx sdk.Context, requestContextID []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(KeyOracleRandRequest(requestContextID))
}

// GetRand retrieves the random number by the specified request id
func (k Keeper) GetRand(ctx sdk.Context, reqID []byte) (types.Rand, sdk.Error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(KeyRand(reqID))
	if bz == nil {
		return types.Rand{}, types.ErrInvalidReqID(k.codespace, fmt.Sprintf("invalid request id: %s", hex.EncodeToString(reqID)))
	}

	var rand types.Rand
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &rand)

	return rand, nil
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

// IterateRandRequestQueueByHeight iterates the random number request queue by the specified height
func (k Keeper) IterateRandRequestQueueByHeight(ctx sdk.Context, height int64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, KeyRandRequestQueueSubspace(height))
}

// IterateRandRequestQueue iterates through the random number request queue
func (k Keeper) IterateRandRequestQueue(ctx sdk.Context, op func(h int64, r types.Request) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, PrefixRandRequestQueue)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		keyParts := bytes.Split(iterator.Key(), KeyDelimiter)
		height, _ := strconv.ParseInt(string(keyParts[1]), 10, 64)

		var request types.Request
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &request)

		if stop := op(height, request); stop {
			break
		}
	}
}
