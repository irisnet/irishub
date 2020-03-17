package keeper

import (
	"fmt"

	"github.com/irisnet/irishub/app/v3/coinswap/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

const KeyPool = "pool"

// GetPool returns the total balance of an reserve pool at the
// provided denomination.
func (k Keeper) GetPool(ctx sdk.Context, uniID string) (pool types.Pool, existed bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(keyPool(uniID))

	if err := k.cdc.UnmarshalBinaryLengthPrefixed(bz, &pool); err != nil {
		return pool, false
	}
	return pool, true
}

// GetPool returns all pools.
func (k Keeper) GetPools(ctx sdk.Context) (pools []types.Pool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(KeyPool))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var pool types.Pool
		if err := k.cdc.UnmarshalBinaryLengthPrefixed(iterator.Value(), &pool); err == nil {
			pools = append(pools, pool)
		}
	}
	return
}

// SetPool is responsible for storing the pool
func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.MarshalBinaryLengthPrefixed(pool)
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}
	store.Set(keyPool(pool.Name), bz)
	return nil
}

// SendCoinsFromAccountToPool is responsible for deducting some coins of the account and adding it to the pool
func (k Keeper) SendCoinsFromAccountToPool(ctx sdk.Context, from sdk.AccAddress, voucherCoinName string, amount sdk.Coins) sdk.Error {
	if _, _, err := k.bk.SubtractCoins(ctx, from, amount); err != nil {
		return err
	}
	pool, existed := k.GetPool(ctx, voucherCoinName)
	if !existed {
		return types.ErrReservePoolNotExists(fmt.Sprintf("liquidity pool for %s not found", voucherCoinName))
	}
	pool.Add(amount)
	if amt := amount.AmountOf(sdk.IrisAtto); amt.GT(sdk.ZeroInt()) {
		k.bk.DecreaseLoosenToken(ctx, amount)
	}
	return k.SetPool(ctx, pool)
}

// SendCoinsFromPoolToAccount is responsible for deducting some coins of the pool and adding it to the account
func (k Keeper) SendCoinsFromPoolToAccount(ctx sdk.Context, receiver sdk.AccAddress, voucherCoinName string, amount sdk.Coins) sdk.Error {
	if _, _, err := k.bk.AddCoins(ctx, receiver, amount); err != nil {
		return err
	}
	pool, existed := k.GetPool(ctx, voucherCoinName)
	if !existed {
		return types.ErrReservePoolNotExists(fmt.Sprintf("liquidity pool for %s not found", voucherCoinName))
	}
	pool.Sub(amount)
	if amt := amount.AmountOf(sdk.IrisAtto); amt.GT(sdk.ZeroInt()) {
		k.bk.IncreaseLoosenToken(ctx, amount)
	}
	return k.SetPool(ctx, pool)
}

//MintLiquidity is responsible for minting some liquidity and adding it to the account/pool
func (k Keeper) MintLiquidity(ctx sdk.Context, receiver sdk.AccAddress, voucherCoinName string, amount sdk.Int) sdk.Error {
	pool, existed := k.GetPool(ctx, voucherCoinName)
	if !existed {
		return types.ErrReservePoolNotExists(fmt.Sprintf("liquidity pool for %s not found", voucherCoinName))
	}

	voucherDenom, err := types.GetVoucherDenom(voucherCoinName)
	if err != nil {
		return err
	}
	// mint liquidity vouchers for Pool
	mintCoins := sdk.NewCoins(sdk.NewCoin(voucherDenom, amount))
	pool.Add(mintCoins)
	if err := k.SetPool(ctx, pool); err != nil {
		return err
	}

	// mint liquidity vouchers for sender
	if _, _, err := k.bk.AddCoins(ctx, receiver, mintCoins); err != nil {
		return err
	}
	ctx.CoinFlowTags().AppendCoinFlowTag(ctx, "", receiver.String(), mintCoins.String(), sdk.MintTokenFlow, "")
	return nil
}

//BurnLiquidity is responsible for burning some liquidity from the account/pool
func (k Keeper) BurnLiquidity(ctx sdk.Context, from sdk.AccAddress, voucherCoinName string, amount sdk.Int) sdk.Error {
	pool, existed := k.GetPool(ctx, voucherCoinName)
	if !existed {
		return types.ErrReservePoolNotExists(fmt.Sprintf("liquidity pool for %s not found", voucherCoinName))
	}

	voucherDenom, err := types.GetVoucherDenom(voucherCoinName)
	if err != nil {
		return err
	}
	// burn liquidity from pool
	burnCoins := sdk.NewCoins(sdk.NewCoin(voucherDenom, amount))
	pool.Sub(burnCoins)
	if err := k.SetPool(ctx, pool); err != nil {
		return err
	}

	// burn liquidity from account
	if _, _, err := k.bk.SubtractCoins(ctx, from, burnCoins); err != nil {
		return err
	}
	ctx.CoinFlowTags().AppendCoinFlowTag(ctx, from.String(), "", burnCoins.String(), sdk.BurnFlow, "")
	return nil
}

func keyPool(voucherCoinName string) []byte {
	return []byte(fmt.Sprintf("%s:%s", KeyPool, voucherCoinName))
}
