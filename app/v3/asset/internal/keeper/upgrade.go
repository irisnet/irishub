package keeper

import (
	sdk "github.com/irisnet/irishub/types"
)

var (
	prefixGateway      = []byte("gateways:")
	prefixOwnerGateway = []byte("ownerGateways:")
)

//Init Initialize module parameters during network upgrade
func (k Keeper) Init(ctx sdk.Context) {
	logger := k.Logger(ctx).With("handler", "Init")

	logger.Info("Begin execute upgrade method")
	store := ctx.KVStore(k.storeKey)

	// delete gateways
	k.iterateGateways(ctx, prefixGateway, func(key []byte) {
		logger.Info("Delete gateway information", "key", string(key))
		store.Delete(key)
	})

	// delete gateway owner
	k.iterateGateways(ctx, prefixOwnerGateway, func(key []byte) {
		logger.Info("Delete gateway owner", "key", string(key))
		store.Delete(key)
	})

	// delete all tokens
	k.IterateTokensWithKeyOp(ctx, func(key []byte) {
		logger.Info("Delete token")
		store.Delete(key)
	})

	//reset params
	param := k.GetParamSet(ctx)
	k.SetParamSet(ctx, param)

	logger.Info("End execute upgrade method")
}

func (k Keeper) iterateGateways(ctx sdk.Context, prefix []byte, op func(key []byte)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		op(iterator.Key())
	}
}

// IterateTokensWithKeyOp iterates through all existing tokens with an operation on key
func (k Keeper) IterateTokensWithKeyOp(ctx sdk.Context, op func(key []byte)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, PrefixToken)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		op(iterator.Key())
	}
}
