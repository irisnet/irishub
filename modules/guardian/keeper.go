package guardian

import (
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
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
func (k Keeper) AddProfiler(ctx sdk.Context, guardian Guardian) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(guardian)
	store.Set(GetProfilerKey(guardian.Address), bz)
	return nil
}

func (k Keeper) DeleteProfiler(ctx sdk.Context, address sdk.AccAddress) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetProfilerKey(address))
	return nil
}

func (k Keeper) GetProfiler(ctx sdk.Context, addr sdk.AccAddress) (guardian Guardian, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(GetProfilerKey(addr))
	if bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &guardian)
		return guardian, true
	}
	return guardian, false
}

// Gets all profilers
func (k Keeper) ProfilersIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, GetProfilersSubspaceKey())
}

// Add a trustee
func (k Keeper) AddTrustee(ctx sdk.Context, guardian Guardian) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(guardian)
	store.Set(GetTrusteeKey(guardian.Address), bz)
	return nil
}

func (k Keeper) DeleteTrustee(ctx sdk.Context, address sdk.AccAddress) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetProfilerKey(address))
	return nil
}

func (k Keeper) GetTrustee(ctx sdk.Context, addr sdk.AccAddress) (guardian Guardian, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(GetTrusteeKey(addr))
	if bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &guardian)
		return guardian, true
	}
	return guardian, false
}

// Gets all trustees
func (k Keeper) TrusteesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, GetTrusteesSubspaceKey())
}
