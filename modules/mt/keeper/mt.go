package keeper

import (
	"crypto/sha256"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"irismod.io/mt/exported"
	"irismod.io/mt/types"
)

const mtIdPrefix = "mt-%d"

// genMTID generate an MT ID by auto increment sequence
func (k Keeper) genMTID(ctx sdk.Context) string {
	sequence := k.GetMTSequence(ctx)
	mtID := fmt.Sprintf(mtIdPrefix, sequence)
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(mtID)))
	k.SetMTSequence(ctx, sequence+1)
	return hash
}

// GetMT gets the the specified MT
func (k Keeper) GetMT(ctx sdk.Context, denomID, mtID string) (mt exported.MT, err error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyMT(denomID, mtID))
	if bz == nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrNotFound, "mt not found (%s)", mtID)
	}

	var baseMT types.MT
	k.cdc.MustUnmarshal(bz, &baseMT)

	// get mt supply
	baseMT.Supply = k.GetMTSupply(ctx, denomID, mtID)

	return baseMT, nil
}

// GetMTs returns all MTs by the specified denom ID
func (k Keeper) GetMTs(ctx sdk.Context, denomID string) (mts []exported.MT) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.KeyMT(denomID, ""))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var baseMT types.MT
		k.cdc.MustUnmarshal(iterator.Value(), &baseMT)

		// get mt supply
		baseMT.Supply = k.GetMTSupply(ctx, denomID, baseMT.GetID())
		mts = append(mts, baseMT)
	}

	return mts
}

// Authorize checks if the sender is the owner of the given denom
func (k Keeper) Authorize(ctx sdk.Context, denomID string, owner sdk.AccAddress) error {
	denom, found := k.GetDenom(ctx, denomID)
	if !found {
		return errorsmod.Wrapf(sdkerrors.ErrNotFound, "denom not found (%s)", denomID)
	}

	if owner.String() != denom.Owner {
		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, owner.String())
	}

	return nil
}

// HasMT checks if the specified MT exists
func (k Keeper) HasMT(ctx sdk.Context, denomID, mtID string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyMT(denomID, mtID))
}

// SetMT set the MT to store
func (k Keeper) SetMT(ctx sdk.Context, denomID string, mt types.MT) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&mt)
	store.Set(types.KeyMT(denomID, mt.GetID()), bz)
}

// GetMTSequence gets the next MT sequence from the store.
func (k Keeper) GetMTSequence(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(types.KeyNextMTSequence))
	if bz == nil {
		return 1
	}
	return sdk.BigEndianToUint64(bz)
}

// SetMTSequence sets the next MT sequence to the store.
func (k Keeper) SetMTSequence(ctx sdk.Context, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(sequence)
	store.Set([]byte(types.KeyNextMTSequence), bz)
}
