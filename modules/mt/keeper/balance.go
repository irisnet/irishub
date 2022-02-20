package keeper

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irismod/modules/mt/types"
)

func (k Keeper) addBalance(ctx sdk.Context,
	denomID, mtID string,
	amount uint64,
	addr sdk.AccAddress) {

	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalAmount(k.cdc, amount)
	store.Set(types.KeyOwner(addr, denomID, mtID), bz)
}

func (k Keeper) subBalance(ctx sdk.Context,
	denomID, mtID string,
	amount uint64,
	addr sdk.AccAddress) {

	store := ctx.KVStore(k.storeKey)
	balance := k.getBalance(ctx, denomID, mtID, addr) - amount

	bz := types.MustMarshalAmount(k.cdc, balance)
	store.Set(types.KeyOwner(addr, denomID, mtID), bz)
}

func (k Keeper) getBalance(ctx sdk.Context,
	denomID, mtID string,
	addr sdk.AccAddress) uint64 {

	store := ctx.KVStore(k.storeKey)

	ownerMt := store.Get(types.KeyOwner(addr, denomID, mtID))
	return binary.BigEndian.Uint64(ownerMt)
}

func (k Keeper) transfer(ctx sdk.Context,
	denomID, mtID string,
	amount uint64,
	from, to sdk.AccAddress) {

	k.subBalance(ctx, denomID, mtID, amount, from)

	k.addBalance(ctx, denomID, mtID, amount, to)
}
