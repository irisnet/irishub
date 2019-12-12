package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/guardian/internal/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec

	// codespace
	codespace sdk.CodespaceType
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, codespace sdk.CodespaceType) Keeper {
	keeper := Keeper{
		storeKey:  key,
		cdc:       cdc,
		codespace: codespace,
	}
	return keeper
}

// Add a profiler, only a existing profiler can add a new and the profiler is not existed
func (k Keeper) AddProfiler(ctx sdk.Context, guardian types.Guardian) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(guardian)
	store.Set(types.GetProfilerKey(guardian.Address), bz)
}

func (k Keeper) DeleteProfiler(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetProfilerKey(address))
}

func (k Keeper) GetProfiler(ctx sdk.Context, addr sdk.AccAddress) (guardian types.Guardian, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetProfilerKey(addr))
	if bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &guardian)
		return guardian, true
	}
	return guardian, false
}

// Gets all profilers
func (k Keeper) ProfilersIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GetProfilersSubspaceKey())
}

// Add a trustee
func (k Keeper) AddTrustee(ctx sdk.Context, guardian types.Guardian) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(guardian)
	store.Set(types.GetTrusteeKey(guardian.Address), bz)
}

func (k Keeper) DeleteTrustee(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetTrusteeKey(address))
}

func (k Keeper) GetTrustee(ctx sdk.Context, addr sdk.AccAddress) (guardian types.Guardian, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetTrusteeKey(addr))
	if bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &guardian)
		return guardian, true
	}
	return guardian, false
}

// Gets all trustees
func (k Keeper) TrusteesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GetTrusteesSubspaceKey())
}
