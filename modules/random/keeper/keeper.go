package keeper

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/irisnet/irishub/modules/random/types"
)

// Keeper defines the random module Keeper
type Keeper struct {
	cdc           codec.Marshaler
	storeKey      sdk.StoreKey
	bankKeeper    types.BankKeeper
	serviceKeeper types.ServiceKeeper
}

// NewKeeper returns a new random keeper
func NewKeeper(cdc codec.Marshaler, key sdk.StoreKey, bankKeeper types.BankKeeper, serviceKeeper types.ServiceKeeper) Keeper {
	keeper := Keeper{
		cdc:           cdc,
		storeKey:      key,
		bankKeeper:    bankKeeper,
		serviceKeeper: serviceKeeper,
	}

	_ = serviceKeeper.RegisterResponseCallback(types.ModuleName, keeper.HandlerResponse)
	_ = serviceKeeper.RegisterStateCallback(types.ModuleName, keeper.HandlerStateChanged)
	return keeper
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
func (k Keeper) RequestRandom(ctx sdk.Context, consumer sdk.AccAddress,
	blockInterval uint64, oracle bool,
	serviceFeeCap sdk.Coins,
) (types.Request, error) {
	currentHeight := ctx.BlockHeight()
	destHeight := currentHeight + int64(blockInterval)

	// get tx hash
	txHash := types.SHA256(ctx.TxBytes())

	var request types.Request
	if oracle {
		// create paused request context
		requestContextID, err := k.RequestService(ctx, consumer, serviceFeeCap)
		if err != nil {
			return request, err
		}

		// build request
		request = types.NewRequest(currentHeight, consumer, txHash, oracle, serviceFeeCap, requestContextID)
	} else {
		// build request
		request = types.NewRequest(currentHeight, consumer, txHash, oracle, nil, nil)
	}

	// generate the request id
	reqID := types.GenerateRequestID(request)

	// add to the queue
	k.EnqueueRandomRequest(ctx, destHeight, reqID, request)

	return request, nil
}

// SetRandom stores the random number
func (k Keeper) SetRandom(ctx sdk.Context, reqID []byte, random types.Random) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&random)
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
	store.Delete(types.KeyRandomRequestQueue(height, reqID))
}

// SetOracleRandRequest stores the oracle random request
func (k Keeper) SetOracleRandRequest(ctx sdk.Context, requestContextID []byte, request types.Request) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&request)
	store.Set(types.KeyOracleRandomRequest(requestContextID), bz)
}

// GetOracleRandRequest retrieves the oracle random request by the specified request id
func (k Keeper) GetOracleRandRequest(ctx sdk.Context, requestContextID []byte) (types.Request, error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyOracleRandomRequest(requestContextID))
	if bz == nil {
		return types.Request{}, sdkerrors.Wrap(types.ErrInvalidRequestContextID, hex.EncodeToString(requestContextID))
	}

	var request types.Request
	k.cdc.MustUnmarshalBinaryBare(bz, &request)

	return request, nil
}

// DeleteOracleRandRequest delete an oracle random request
func (k Keeper) DeleteOracleRandRequest(ctx sdk.Context, requestContextID []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyOracleRandomRequest(requestContextID))
}

// GetRand retrieves the random number by the specified request id
func (k Keeper) GetRandom(ctx sdk.Context, reqID []byte) (types.Random, error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyRandom(reqID))
	if bz == nil {
		return types.Random{}, sdkerrors.Wrap(types.ErrInvalidReqID, hex.EncodeToString(reqID))
	}

	var random types.Random
	k.cdc.MustUnmarshalBinaryBare(bz, &random)

	return random, nil
}

// IterateRandoms iterates through all the random numbers
func (k Keeper) IterateRandoms(ctx sdk.Context, op func(r types.Random) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.PrefixRandom)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var random types.Random
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &random)

		if stop := op(random); stop {
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
		height, _ := strconv.ParseInt(string(keyParts[1]), 10, 64)

		var request types.Request
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &request)

		if stop := op(height, request); stop {
			break
		}
	}
}
