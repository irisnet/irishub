package keeper

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irismod/modules/mt/types"
)

func (k Keeper) deleteOwner(ctx sdk.Context, denomID, tokenID string, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyOwner(owner, denomID, tokenID))
}

func (k Keeper) getOwner(ctx sdk.Context,
	denomID, tokenID string,
	owner sdk.AccAddress,
) uint64 {
	store := ctx.KVStore(k.storeKey)

	ownerMt := store.Get(types.KeyOwner(owner, denomID, tokenID))

	return binary.BigEndian.Uint64(ownerMt)
}

func (k Keeper) setOwner(ctx sdk.Context,
	denomID, tokenID string,
	amount uint64,
	owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	bz := types.MustMarshalAmount(k.cdc, amount)
	store.Set(types.KeyOwner(owner, denomID, tokenID), bz)
}

func (k Keeper) swapOwner(ctx sdk.Context, denomID, tokenID string, srcOwner, dstOwner sdk.AccAddress) {

	// delete old owner key
	k.deleteOwner(ctx, denomID, tokenID, srcOwner)

	// set new owner key
	k.setOwner(ctx, denomID, tokenID, dstOwner)
}
