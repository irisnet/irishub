package keeper

import (
	sdk "github.com/irisnet/irishub/types"
)

var (
	prefixGateway      = []byte("gateways:")
	prefixOwnerGateway = []byte("ownerGateways:")
	prefixOwnerToken   = []byte("ownerTokens:")
)

//Init Initialize module parameters during network upgrade
func (k Keeper) Init(ctx sdk.Context) {
	logger := k.Logger(ctx).With("handler", "Init")

	logger.Info("Begin execute upgrade method")
	store := ctx.KVStore(k.storeKey)

	// delete gateways
	k.iterateWithPrefix(ctx, prefixGateway, func(key []byte) {
		logger.Info("Delete gateway information", "key", string(key))
		store.Delete(key)
	})

	// delete gateway owner
	k.iterateWithPrefix(ctx, prefixOwnerGateway, func(key []byte) {
		logger.Info("Delete gateway owner", "key", string(key))
		store.Delete(key)
	})

	// delete tokens
	k.iterateWithPrefix(ctx, PrefixToken, func(key []byte) {
		logger.Info("Delete token", "key", string(key))
		store.Delete(key)
	})

	// delete token owner
	k.iterateWithPrefix(ctx, prefixOwnerToken, func(key []byte) {
		logger.Info("Delete token owner", "key", string(key))
		store.Delete(key)
	})

	//reset params
	param := k.GetParamSet(ctx)
	k.SetParamSet(ctx, param)

	logger.Info("End execute upgrade method")
}

func (k Keeper) iterateWithPrefix(ctx sdk.Context, prefix []byte, op func(key []byte)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		op(iterator.Key())
	}
}
