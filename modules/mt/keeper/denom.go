package keeper

import (
	"crypto/sha256"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"mods.irisnet.org/modules/mt/types"
)

const denomIdPrefix = "mt-denom-%d"

// genDenomID generate a denom ID by auto increment sequence
func (k Keeper) genDenomID(ctx sdk.Context) string {
	sequence := k.GetDenomSequence(ctx)
	denomID := fmt.Sprintf(denomIdPrefix, sequence)
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(denomID)))
	k.SetDenomSequence(ctx, sequence+1)
	return hash
}

// HasDenom returns whether the specified denom ID exists
func (k Keeper) HasDenom(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyDenom(id))
}

// SetDenom is responsible for saving the definition of denom
func (k Keeper) SetDenom(ctx sdk.Context, denom types.Denom) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&denom)
	store.Set(types.KeyDenom(denom.Id), bz)
}

// GetDenom returns the denom by id
func (k Keeper) GetDenom(ctx sdk.Context, id string) (denom types.Denom, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyDenom(id))
	if len(bz) == 0 {
		return denom, false
	}

	k.cdc.MustUnmarshal(bz, &denom)
	return denom, true
}

// GetDenoms returns all the denoms
func (k Keeper) GetDenoms(ctx sdk.Context) (denoms []types.Denom) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.KeyDenom(""))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var denom types.Denom
		k.cdc.MustUnmarshal(iterator.Value(), &denom)
		denoms = append(denoms, denom)
	}
	return denoms
}

// UpdateDenom is responsible for updating the definition of denom
func (k Keeper) UpdateDenom(ctx sdk.Context, denom types.Denom) error {
	if !k.HasDenom(ctx, denom.Id) {
		return errorsmod.Wrapf(sdkerrors.ErrNotFound, "denom not found (%s)", denom.Id)
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&denom)
	store.Set(types.KeyDenom(denom.Id), bz)
	return nil
}

// GetDenomSequence gets the next denom sequence from the store.
func (k Keeper) GetDenomSequence(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(types.KeyNextDenomSequence))
	if bz == nil {
		return 1
	}
	return sdk.BigEndianToUint64(bz)
}

// SetDenomSequence sets the next denom sequence to the store.
func (k Keeper) SetDenomSequence(ctx sdk.Context, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(sequence)
	store.Set([]byte(types.KeyNextDenomSequence), bz)
}
