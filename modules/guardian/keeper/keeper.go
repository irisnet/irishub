package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/guardian/types"
)

// Keeper of the guardian store
type Keeper struct {
	cdc      codec.Marshaler
	storeKey sdk.StoreKey
}

// NewKeeper returns a guardian keeper
func NewKeeper(cdc codec.Marshaler, key sdk.StoreKey) Keeper {
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
	bz := k.cdc.MustMarshalBinaryBare(&guardian)
	store.Set(types.GetProfilerKey(guardian.Address), bz)
}

// DeleteProfiler delete the stored profiler
func (k Keeper) DeleteProfiler(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetProfilerKey(address))
}

// GetProfiler retrieves the profiler by specified address
func (k Keeper) GetProfiler(ctx sdk.Context, addr sdk.AccAddress) (guardian types.Guardian, found bool) {
	store := ctx.KVStore(k.storeKey)
	if bz := store.Get(types.GetProfilerKey(addr)); bz != nil {
		k.cdc.MustUnmarshalBinaryBare(bz, &guardian)
		return guardian, true
	}
	return guardian, false
}

// IterateProfilers iterates through all profilers
func (k Keeper) IterateProfilers(
	ctx sdk.Context,
	op func(profiler types.Guardian) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.GetProfilersSubspaceKey())
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var profiler types.Guardian
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &profiler)

		if stop := op(profiler); stop {
			break
		}
	}
}

// AddTrustee add a trustee
func (k Keeper) AddTrustee(ctx sdk.Context, guardian types.Guardian) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&guardian)
	store.Set(types.GetTrusteeKey(guardian.GetAddress()), bz)
}

// DeleteTrustee delete the stored trustee
func (k Keeper) DeleteTrustee(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetTrusteeKey(address))
}

// GetTrustee retrieves the trustee by specified address
func (k Keeper) GetTrustee(ctx sdk.Context, addr sdk.AccAddress) (guardian types.Guardian, found bool) {
	store := ctx.KVStore(k.storeKey)
	if bz := store.Get(types.GetTrusteeKey(addr)); bz != nil {
		k.cdc.MustUnmarshalBinaryBare(bz, &guardian)
		return guardian, true
	}
	return guardian, false
}

// IterateTrustees iterates through all trustees
func (k Keeper) IterateTrustees(
	ctx sdk.Context,
	op func(trustee types.Guardian) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.GetTrusteesSubspaceKey())
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var trustee types.Guardian
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &trustee)

		if stop := op(trustee); stop {
			break
		}
	}
}

func (k Keeper) Authorized(ctx sdk.Context, addr sdk.AccAddress) bool {
	_, found := k.GetProfiler(ctx, addr)
	return found
}
