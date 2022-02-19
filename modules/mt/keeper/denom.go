package keeper

import (
	"crypto/sha256"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/mt/types"
)

const denomIdPrefix = "mt-denom-%d"

// genDenomID generate a denom ID by auto increment sequence
func (k Keeper) genDenomID(ctx sdk.Context) string {
	sequence := k.getDenomSequence(ctx)
	denomID := fmt.Sprintf(denomIdPrefix, sequence)
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(denomID)))
	k.setDenomSequence(ctx, sequence+1)
	return hash
}

// HasDenomID returns whether the specified denom ID exists
func (k Keeper) HasDenomID(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyDenomID(id))
}

// SetDenom is responsible for saving the definition of denom
func (k Keeper) SetDenom(ctx sdk.Context, denom types.Denom) error {
	if k.HasDenomID(ctx, denom.Id) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "Denom already exists: %s", denom.Id)
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&denom)
	store.Set(types.KeyDenomID(denom.Id), bz)
	store.Set(types.KeyDenomName(denom.Name), []byte(denom.Id))
	return nil
}

// GetDenom returns the denom by id
func (k Keeper) GetDenom(ctx sdk.Context, id string) (denom types.Denom, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyDenomID(id))
	if len(bz) == 0 {
		return denom, false
	}

	k.cdc.MustUnmarshal(bz, &denom)
	return denom, true
}

// GetDenoms returns all the denoms
func (k Keeper) GetDenoms(ctx sdk.Context) (denoms []types.Denom) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyDenomID(""))
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
	if !k.HasDenomID(ctx, denom.Id) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "Denom not found: %s", denom.Id)
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&denom)
	store.Set(types.KeyDenomID(denom.Id), bz)
	return nil
}

// getDenomSequence gets the next denom sequence from the store.
func (k Keeper) getDenomSequence(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(types.KeyNextDenomSequence))
	if bz == nil {
		return 1
	}
	return sdk.BigEndianToUint64(bz)
}

// setDenomSequence sets the next denom sequence to the store.
func (k Keeper) setDenomSequence(ctx sdk.Context, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(sequence)
	store.Set([]byte(types.KeyNextDenomSequence), bz)
}
