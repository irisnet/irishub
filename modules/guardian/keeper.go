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
func (k Keeper) AddProfiler(ctx sdk.Context, profiler Profiler) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(profiler)
	store.Set(GetProfilerKey(profiler.Addr), bz)
	return nil
}

func (k Keeper) GetProfiler(ctx sdk.Context, addr sdk.AccAddress) (profiler Profiler, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(GetProfilerKey(addr))
	if bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &profiler)
		return profiler, true
	}
	return profiler, false
}

// Gets all profilers
func (k Keeper) GetProfilers(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, GetProfilersSubspaceKey())
}

// Add a trustee
func (k Keeper) AddTrustee(ctx sdk.Context, trustee Trustee) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(trustee)
	store.Set(GetTrusteeKey(trustee.Addr), bz)
	return nil
}

func (k Keeper) GetTrustee(ctx sdk.Context, addr sdk.AccAddress) (trustee Trustee, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(GetTrusteeKey(addr))
	if bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &trustee)
		return trustee, true
	}
	return trustee, false
}

// Gets all trustees
func (k Keeper) GetTrustees(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, GetTrusteesSubspaceKey())
}
