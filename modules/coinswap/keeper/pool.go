package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/irisnet/irismod/modules/coinswap/types"
)

// CreatePool create a liquidity that saves relevant information about popular pool tokens
func (k Keeper) CreatePool(ctx sdk.Context, counterpartyDenom string) types.Pool {
	sequence := k.getPoolSequence(ctx)
	lptDenom := types.GetLptDenom(sequence)
	pool := &types.Pool{
		Id:                types.GetPoolId(counterpartyDenom),
		StandardDenom:     k.GetStandardDenom(ctx),
		CounterpartyDenom: counterpartyDenom,
		EscrowAddress:     types.GetReservePoolAddr(lptDenom).String(),
		LptDenom:          lptDenom,
	}
	k.increasePoolSequence(ctx, sequence)
	k.setPool(ctx, pool)
	return *pool
}

// GetPool return the liquidity pool by the specified anotherCoinDenom
func (k Keeper) GetPool(ctx sdk.Context, poolId string) (types.Pool, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPoolKey(poolId))
	if bz == nil {
		return types.Pool{}, false
	}

	pool := &types.Pool{}
	k.cdc.MustUnmarshalBinaryBare(bz, pool)
	return *pool, true
}

// GetPoolByLptDenom return the liquidity pool by the specified anotherCoinDenom
func (k Keeper) GetPoolByLptDenom(ctx sdk.Context, lptDenom string) (types.Pool, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLptDenomKey(lptDenom))
	if bz == nil {
		return types.Pool{}, false
	}

	poolId := &gogotypes.StringValue{}
	k.cdc.MustUnmarshalBinaryBare(bz, poolId)
	return k.GetPool(ctx, poolId.Value)
}

// GetPoolBalances return the liquidity pool by the specified anotherCoinDenom
func (k Keeper) GetPoolBalances(ctx sdk.Context, escrowAddress string) (coins sdk.Coins, err error) {
	address, err := sdk.AccAddressFromBech32(escrowAddress)
	if err != nil {
		return coins, err
	}
	acc := k.ak.GetAccount(ctx, address)
	if acc == nil {
		return nil, sdkerrors.Wrap(types.ErrReservePoolNotExists, escrowAddress)
	}
	return k.bk.GetAllBalances(ctx, acc.GetAddress()), nil
}

func (k Keeper) GetPoolBalancesByLptDenom(ctx sdk.Context, lptDenom string) (coins sdk.Coins, err error) {
	address := types.GetReservePoolAddr(lptDenom)
	acc := k.ak.GetAccount(ctx, address)
	if acc == nil {
		return nil, sdkerrors.Wrap(types.ErrReservePoolNotExists, address.String())
	}
	return k.bk.GetAllBalances(ctx, acc.GetAddress()), nil
}

// GetLptDenomFromDenoms returns the liquidity pool token denom for the provided denominations.
func (k Keeper) GetLptDenomFromDenoms(ctx sdk.Context, denom1, denom2 string) (string, error) {
	if denom1 == denom2 {
		return "", types.ErrEqualDenom
	}

	standardDenom := k.GetStandardDenom(ctx)
	if denom1 != standardDenom && denom2 != standardDenom {
		return "", sdkerrors.Wrap(types.ErrNotContainStandardDenom, fmt.Sprintf("standard denom: %s, denom1: %s, denom2: %s", standardDenom, denom1, denom2))
	}

	counterpartyDenom := denom1
	if counterpartyDenom == standardDenom {
		counterpartyDenom = denom2
	}
	poolId := types.GetPoolId(counterpartyDenom)
	pool, has := k.GetPool(ctx, poolId)
	if !has {
		return "", sdkerrors.Wrapf(types.ErrReservePoolNotExists, "liquidity pool token: %s", counterpartyDenom)
	}
	return pool.LptDenom, nil
}

// ValidatePool Verify the legitimacy of the liquidity pool
func (k Keeper) ValidatePool(ctx sdk.Context, lptDenom string) error {
	if err := types.ValidateLptDenom(lptDenom); err != nil {
		return err
	}

	pool, has := k.GetPoolByLptDenom(ctx, lptDenom)
	if !has {
		return sdkerrors.Wrapf(types.ErrReservePoolNotExists, "liquidity pool token: %s", lptDenom)
	}

	_, err := k.GetPoolBalances(ctx, pool.EscrowAddress)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) setPool(ctx sdk.Context, pool *types.Pool) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(pool)
	store.Set(types.GetPoolKey(pool.Id), bz)

	// save by lpt denom
	poolId := &gogotypes.StringValue{Value: pool.Id}
	poolIdBz := k.cdc.MustMarshalBinaryBare(poolId)
	store.Set(types.GetLptDenomKey(pool.LptDenom), poolIdBz)
}

// getPoolSequence gets the next pool sequence from the store.
func (k Keeper) getPoolSequence(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(types.KeyNextPoolSequence))
	if bz == nil {
		return 1
	}
	return sdk.BigEndianToUint64(bz)
}

// increasePoolSequence sets the next pool sequence to the store.
func (k Keeper) increasePoolSequence(ctx sdk.Context, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(sequence + 1)
	store.Set([]byte(types.KeyNextPoolSequence), bz)
}
