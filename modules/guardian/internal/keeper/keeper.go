package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/guardian/exported"
	"github.com/irisnet/irishub/modules/guardian/internal/types"
)

// Keeper of the guardian store
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
}

// NewKeeper returns a guardian keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey) Keeper {
	keeper := Keeper{
		storeKey: key,
		cdc:      cdc,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("%s", types.ModuleName))
}

// Add a profiler, only a existing profiler can add a new and the profiler is not existed
func (k Keeper) AddProfiler(ctx sdk.Context, guardian types.Guardian) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(guardian)
	store.Set(types.GetProfilerKey(guardian.Address), bz)
}

// DeleteProfiler delete the stored profiler
func (k Keeper) DeleteProfiler(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetProfilerKey(address))
}

// GetProfiler retrieves the profiler by specified address
func (k Keeper) GetProfiler(ctx sdk.Context, addr sdk.AccAddress) (guardian exported.GuardianI, found bool) {
	store := ctx.KVStore(k.storeKey)
	if bz := store.Get(types.GetProfilerKey(addr)); bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &guardian)
		return guardian, true
	}
	return guardian, false
}

// ProfilersIterator gets all profilers
func (k Keeper) ProfilersIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GetProfilersSubspaceKey())
}

// AddTrustee add a trustee
func (k Keeper) AddTrustee(ctx sdk.Context, guardian exported.GuardianI) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(guardian)
	store.Set(types.GetTrusteeKey(guardian.GetAddress()), bz)
}

// DeleteTrustee delete the stored trustee
func (k Keeper) DeleteTrustee(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetTrusteeKey(address))
}

// GetTrustee retrieves the trustee by specified address
func (k Keeper) GetTrustee(ctx sdk.Context, addr sdk.AccAddress) (guardian exported.GuardianI, found bool) {
	store := ctx.KVStore(k.storeKey)
	if bz := store.Get(types.GetTrusteeKey(addr)); bz != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &guardian)
		return guardian, true
	}
	return guardian, false
}

// TrusteesIterator gets all trustees
func (k Keeper) TrusteesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GetTrusteesSubspaceKey())
}
