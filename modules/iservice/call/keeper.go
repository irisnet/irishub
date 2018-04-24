package call

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

type SvcCallKeeper struct {
	ck       bank.CoinKeeper
	storeKey sdk.StoreKey // The (unexposed) key used to access the store from the Context.
}

// NewKeeper - Returns the Keeper
func NewSvcCallKeeper(key sdk.StoreKey, bankKeeper bank.CoinKeeper) SvcCallKeeper {
	return SvcCallKeeper{bankKeeper, key}
}

func (k SvcCallKeeper) Has(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(key)
}

func (k SvcCallKeeper) Set(ctx sdk.Context, key, value []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(key, value)
}
