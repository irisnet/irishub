package bind

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

type SvcBindKeeper struct {
	ck       bank.CoinKeeper
	storeKey sdk.StoreKey // The (unexposed) key used to access the store from the Context.
}

// NewKeeper - Returns the Keeper
func NewSvcBindKeeper(key sdk.StoreKey, bankKeeper bank.CoinKeeper) SvcBindKeeper {
	return SvcBindKeeper{bankKeeper, key}
}

func (k SvcBindKeeper) Has(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(key)
}

func (k SvcBindKeeper) Set(ctx sdk.Context, key, value []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(key, value)
}
