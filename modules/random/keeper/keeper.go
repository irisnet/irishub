package keeper

import (
	"encoding/hex"

	"github.com/cometbft/cometbft/libs/log"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"mods.irisnet.org/modules/random/types"
)

// Keeper defines the random module Keeper
type Keeper struct {
	cdc           codec.Codec
	storeKey      storetypes.StoreKey
	bankKeeper    types.BankKeeper
	serviceKeeper types.ServiceKeeper
}

// NewKeeper returns a new random keeper
func NewKeeper(
	cdc codec.Codec,
	key storetypes.StoreKey,
	bankKeeper types.BankKeeper,
	serviceKeeper types.ServiceKeeper,
) Keeper {
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
	return ctx.Logger().With("module", types.ModuleName)
}

// GetCdc returns the cdc
func (k Keeper) GetCdc() codec.BinaryCodec {
	return k.cdc
}

// RequestRandom requests a random number
func (k Keeper) RequestRandom(
	ctx sdk.Context, consumer sdk.AccAddress,
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
		request = types.NewRequest(
			currentHeight,
			consumer.String(),
			hex.EncodeToString(txHash),
			oracle,
			serviceFeeCap,
			requestContextID.String(),
		)
	} else {
		// build request
		request = types.NewRequest(currentHeight, consumer.String(), hex.EncodeToString(txHash), oracle, nil, "")
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
	bz := k.cdc.MustMarshal(&random)
	store.Set(types.KeyRandom(reqID), bz)
}

// EnqueueRandomRequest enqueues the random number request
func (k Keeper) EnqueueRandomRequest(
	ctx sdk.Context,
	height int64,
	reqID []byte,
	request types.Request,
) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&request)
	store.Set(types.KeyRandomRequestQueue(height, reqID), bz)
}

// DequeueRandomRequest removes the random number request by the specified height and request id
func (k Keeper) DequeueRandomRequest(ctx sdk.Context, height int64, reqID []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyRandomRequestQueue(height, reqID))
}

// SetOracleRandRequest stores the oracle random number request
func (k Keeper) SetOracleRandRequest(
	ctx sdk.Context,
	requestContextID []byte,
	request types.Request,
) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&request)
	store.Set(types.KeyOracleRandomRequest(requestContextID), bz)
}

// GetOracleRandRequest retrieves the oracle random number request by the specified request id
func (k Keeper) GetOracleRandRequest(
	ctx sdk.Context,
	requestContextID []byte,
) (types.Request, error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyOracleRandomRequest(requestContextID))
	if bz == nil {
		return types.Request{}, errorsmod.Wrap(
			types.ErrInvalidRequestContextID,
			hex.EncodeToString(requestContextID),
		)
	}

	var request types.Request
	k.cdc.MustUnmarshal(bz, &request)

	return request, nil
}

// DeleteOracleRandRequest deletes an oracle random number request
func (k Keeper) DeleteOracleRandRequest(ctx sdk.Context, requestContextID []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyOracleRandomRequest(requestContextID))
}

// GetRandom retrieves the random number by the specified request id
func (k Keeper) GetRandom(ctx sdk.Context, reqID []byte) (types.Random, error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyRandom(reqID))
	if bz == nil {
		return types.Random{}, errorsmod.Wrap(types.ErrInvalidReqID, hex.EncodeToString(reqID))
	}

	var random types.Random
	k.cdc.MustUnmarshal(bz, &random)

	return random, nil
}

// IterateRandoms iterates through all the random numbers
func (k Keeper) IterateRandoms(ctx sdk.Context, op func(r types.Random) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.RandomKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var random types.Random
		k.cdc.MustUnmarshal(iterator.Value(), &random)

		if stop := op(random); stop {
			break
		}
	}
}

// IterateRandomRequestQueueByHeight iterates through the random number request queue by the specified height
func (k Keeper) IterateRandomRequestQueueByHeight(ctx sdk.Context, height int64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.KeyRandomRequestQueueSubspace(height))
}

// IterateRandomRequestQueue iterates through the random number request queue
func (k Keeper) IterateRandomRequestQueue(
	ctx sdk.Context,
	op func(h int64, reqID []byte, r types.Request) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.RandomRequestQueueKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		height := sdk.BigEndianToUint64(iterator.Key()[1:9])
		reqID := iterator.Key()[9:]

		var request types.Request
		k.cdc.MustUnmarshal(iterator.Value(), &request)

		if stop := op(int64(height), reqID, request); stop {
			break
		}
	}
}
