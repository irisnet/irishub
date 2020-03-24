package keeper

import (
	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/modules/guardian"
	sdk "github.com/irisnet/irishub/types"
)

var (
	prefixGateway      = []byte("gateways:")
	prefixOwnerGateway = []byte("ownerGateways:")
	prefixOwnerToken   = []byte("ownerTokens:")

	prefixTotalSupply = []byte("totalSupply:")
)

//Init Initializes module parameters during network upgrade
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

	// delete non-iris coins from profilers
	k.deleteCoinsFromAccounts(ctx, k.getAllProfilers(ctx))

	// delete total supplies
	k.deleteTotalSupplies(ctx)

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

func (k Keeper) getAllProfilers(ctx sdk.Context) []sdk.AccAddress {
	var profilers []sdk.AccAddress

	iterator := k.gk.ProfilersIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		var profiler guardian.Guardian
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &profiler)

		profilers = append(profilers, profiler.Address)
	}

	return profilers
}

func (k Keeper) deleteCoinsFromAccounts(ctx sdk.Context, addrs []sdk.AccAddress) {
	for _, addr := range addrs {
		coins := k.bk.GetCoins(ctx, addr)
		nonIrisCoins, _ := coins.SafeSub(sdk.NewCoins(sdk.NewCoin(sdk.IrisAtto, coins.AmountOf(sdk.IrisAtto))))

		if !nonIrisCoins.IsZero() {
			_, _, _ = k.bk.SubtractCoins(ctx, addr, nonIrisCoins)
		}
	}
}

func (k Keeper) deleteTotalSupplies(ctx sdk.Context) {
	store := ctx.KVStore(protocol.KeyAccount)

	iterator := sdk.KVStorePrefixIterator(store, prefixTotalSupply)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}
