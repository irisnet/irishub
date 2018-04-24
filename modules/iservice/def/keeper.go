package def

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

type SvcDefKeeper struct {
	ck       bank.CoinKeeper
	storeKey sdk.StoreKey // The (unexposed) key used to access the store from the Context.
}

// NewKeeper - Returns the Keeper
func NewSvcDefKeeper(key sdk.StoreKey, bankKeeper bank.CoinKeeper) SvcDefKeeper {
	return SvcDefKeeper{bankKeeper, key}
}

func (k SvcDefKeeper) Has(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(key)
}

func (k SvcDefKeeper) Set(ctx sdk.Context, key, value []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(key, value)
}
