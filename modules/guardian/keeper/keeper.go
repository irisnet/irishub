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
	cdc      codec.Codec
	storeKey sdk.StoreKey
}

// NewKeeper returns a guardian keeper
func NewKeeper(cdc codec.Codec, key sdk.StoreKey) Keeper {
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

// Add a super, only a existing super can add a new and the super is not existed
func (k Keeper) AddSuper(ctx sdk.Context, super types.Super) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&super)
	address, _ := sdk.AccAddressFromBech32(super.Address)
	store.Set(types.GetSuperKey(address), bz)
}

// DeleteSuper delete the stored super
func (k Keeper) DeleteSuper(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetSuperKey(address))
}

// GetSuper retrieves the super by specified address
func (k Keeper) GetSuper(ctx sdk.Context, addr sdk.AccAddress) (super types.Super, found bool) {
	store := ctx.KVStore(k.storeKey)
	if bz := store.Get(types.GetSuperKey(addr)); bz != nil {
		k.cdc.MustUnmarshal(bz, &super)
		return super, true
	}
	return super, false
}

// IterateSupers iterates through all supers
func (k Keeper) IterateSupers(
	ctx sdk.Context,
	op func(super types.Super) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.GetSupersSubspaceKey())
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var super types.Super
		k.cdc.MustUnmarshal(iterator.Value(), &super)

		if stop := op(super); stop {
			break
		}
	}
}

func (k Keeper) Authorized(ctx sdk.Context, addr sdk.AccAddress) bool {
	_, found := k.GetSuper(ctx, addr)
	return found
}
