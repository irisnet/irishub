package keeper

import (
	"bytes"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/irisnet/irismod/modules/mt/types"
)

// AddBalance adds amounts to an account
func (k Keeper) AddBalance(ctx sdk.Context,
	denomID, mtID string,
	amount uint64,
	addr sdk.AccAddress) error {

	balance := k.GetBalance(ctx, denomID, mtID, addr)
	if math.MaxUint64-balance < amount {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "overflow: max %d, got %d", math.MaxUint64-balance, amount)
	}
	balance += amount

	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalAmount(k.cdc, balance)
	store.Set(types.KeyBalance(addr, denomID, mtID), bz)
	return nil
}

// SubBalance subs amounts from an account
func (k Keeper) SubBalance(ctx sdk.Context,
	denomID, mtID string,
	amount uint64,
	addr sdk.AccAddress) {

	store := ctx.KVStore(k.storeKey)
	balance := k.GetBalance(ctx, denomID, mtID, addr)
	balance -= amount

	bz := types.MustMarshalAmount(k.cdc, balance)
	store.Set(types.KeyBalance(addr, denomID, mtID), bz)
}

// GetBalance gets balance of an MT
func (k Keeper) GetBalance(ctx sdk.Context,
	denomID, mtID string,
	addr sdk.AccAddress) uint64 {

	store := ctx.KVStore(k.storeKey)

	amount := store.Get(types.KeyBalance(addr, denomID, mtID))
	if len(amount) == 0 {
		return 0
	}

	return types.MustUnMarshalAmount(k.cdc, amount)
}

// getBalances gets balances of all accounts, should only be used in exporting genesis states
func (k Keeper) getBalances(ctx sdk.Context) []types.Owner {

	store := ctx.KVStore(k.storeKey)

	it := sdk.KVStorePrefixIterator(store, types.PrefixBalance)
	defer it.Close()

	var ownerMap map[string]map[string]map[string]uint64
	ownerMap = make(map[string]map[string]map[string]uint64)

	for ; it.Valid(); it.Next() {
		keys := bytes.Split(it.Key(), types.Delimiter)

		address := string(keys[1])
		denomID := string(keys[2])
		mtID := string(keys[3])
		amount := types.MustUnMarshalAmount(k.cdc, it.Value())

		if _, ok := ownerMap[address]; !ok {
			ownerMap[address] = make(map[string]map[string]uint64)
		}

		if _, ok := ownerMap[address][denomID]; !ok {
			ownerMap[address][denomID] = make(map[string]uint64)
		}

		ownerMap[address][denomID][mtID] = amount
	}

	var owners []types.Owner
	for addr, denomMap := range ownerMap {
		var denomBalances []types.DenomBalance
		for denomID, mtMap := range denomMap {
			var balances []types.Balance
			for mtID, amount := range mtMap {
				balance := types.NewBalance(mtID, amount)
				balances = append(balances, balance)
			}
			denomBalance := types.NewDenomBalance(denomID, balances)
			denomBalances = append(denomBalances, denomBalance)
		}

		owner := types.NewOwner(addr, denomBalances)
		owners = append(owners, owner)
	}

	return owners
}

// Transfer transfers mts
func (k Keeper) Transfer(ctx sdk.Context,
	denomID, mtID string,
	amount uint64,
	from, to sdk.AccAddress) error {

	k.SubBalance(ctx, denomID, mtID, amount, from)

	return k.AddBalance(ctx, denomID, mtID, amount, to)
}

// GetDenomSupply returns the number of Mts by the specified denom ID
func (k Keeper) GetDenomSupply(ctx sdk.Context, denomID string) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeySupply(denomID, ""))
	if len(bz) == 0 {
		return 0
	}
	return types.MustUnMarshalSupply(k.cdc, bz)
}

// GetMTSupply returns the supply of a specified MT
func (k Keeper) GetMTSupply(ctx sdk.Context, denomID, mtID string) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeySupply(denomID, mtID))
	if len(bz) == 0 {
		return 0
	}
	return types.MustUnMarshalSupply(k.cdc, bz)
}

// IncreaseDenomSupply increase total supply (count of MTs) of a denom
func (k Keeper) IncreaseDenomSupply(ctx sdk.Context, denomID string) {
	supply := k.GetDenomSupply(ctx, denomID)
	supply++

	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalSupply(k.cdc, supply)
	store.Set(types.KeySupply(denomID, ""), bz)
}

// IncreaseMTSupply increase total supply of an MT
func (k Keeper) IncreaseMTSupply(ctx sdk.Context, denomID, mtID string, amount uint64) error {
	supply := k.GetMTSupply(ctx, denomID, mtID)
	supply += amount
	if math.MaxUint64-supply < amount {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "overflow: max %d, got %d", math.MaxUint64-supply, amount)
	}

	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalSupply(k.cdc, supply)
	store.Set(types.KeySupply(denomID, mtID), bz)
	return nil
}

// decreaseMTSupply decrease total supply of an MT
func (k Keeper) decreaseMTSupply(ctx sdk.Context, denomID, mtID string, amount uint64) {
	supply := k.GetMTSupply(ctx, denomID, mtID)
	supply -= amount

	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalSupply(k.cdc, supply)
	store.Set(types.KeySupply(denomID, mtID), bz)
}
