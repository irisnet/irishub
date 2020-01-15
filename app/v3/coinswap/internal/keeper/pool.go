package keeper

import (
	"fmt"

	"github.com/irisnet/irishub/app/v3/coinswap/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

const KeyPool = "pool"

type Pool struct {
	sdk.Coins
	Name string
}

func NewPool(name string, coins sdk.Coins) Pool {
	return Pool{
		Coins: coins,
		Name:  name,
	}
}

// GetPool returns the total balance of an reserve pool at the
// provided denomination.
func (k Keeper) GetPool(ctx sdk.Context, uniID string) (pool Pool, existed bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(keyPool(uniID))

	if err := k.cdc.UnmarshalBinaryLengthPrefixed(bz, &pool); err != nil {
		return pool, false
	}
	return pool, true
}

// GetPool returns all pools.
func (k Keeper) GetPools(ctx sdk.Context) (pools []Pool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(KeyPool))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var pool Pool
		if err := k.cdc.UnmarshalBinaryLengthPrefixed(iterator.Value(), &pool); err == nil {
			pools = append(pools, pool)
		}
	}
	return
}

// SetPool is responsible for storing the poll to database
func (k Keeper) SetPool(ctx sdk.Context, pool Pool) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.MarshalBinaryLengthPrefixed(pool)
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}
	store.Set(keyPool(pool.Name), bz)
	return nil
}

// SendCoinsFromAccountToPool is responsible for deducting some coins of the account and adding it to the pool
func (k Keeper) SendCoinsFromAccountToPool(ctx sdk.Context, from sdk.AccAddress, uniID string, amount sdk.Coins) sdk.Error {
	if _, _, err := k.bk.SubtractCoins(ctx, from, amount); err != nil {
		return err
	}
	pool, existed := k.GetPool(ctx, uniID)
	if !existed {
		return types.ErrReservePoolNotExists(fmt.Sprintf("reserve pool for %s not found", uniID))
	}
	pool.Coins = pool.Add(amount)
	return k.SetPool(ctx, pool)
}

// SendCoinsFromPoolToAccount is responsible for deducting some coins of the pool and adding it to the account
func (k Keeper) SendCoinsFromPoolToAccount(ctx sdk.Context, receiver sdk.AccAddress, uniID string, amount sdk.Coins) sdk.Error {
	if _, _, err := k.bk.AddCoins(ctx, receiver, amount); err != nil {
		return err
	}
	pool, existed := k.GetPool(ctx, uniID)
	if !existed {
		return types.ErrReservePoolNotExists(fmt.Sprintf("reserve pool for %s not found", uniID))
	}
	pool.Coins = pool.Coins.Sub(amount)
	return k.SetPool(ctx, pool)
}

//MintCoins is responsible for minting some coins and adding it to the account/pool
func (k Keeper) MintCoins(ctx sdk.Context, receiver sdk.AccAddress, uniID string, mintAmount sdk.Int) sdk.Error {
	pool, existed := k.GetPool(ctx, uniID)
	if !existed {
		return types.ErrReservePoolNotExists(fmt.Sprintf("reserve pool for %s not found", uniID))
	}

	uniDenom, err := types.GetUniDenom(uniID)
	if err != nil {
		return err
	}
	// mint liquidity vouchers for reserve Pool
	mintCoins := sdk.NewCoins(sdk.NewCoin(uniDenom, mintAmount))
	pool.Coins = pool.Coins.Add(mintCoins)
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

//BurnCoins is responsible for burning some coins from the account/pool
func (k Keeper) BurnCoins(ctx sdk.Context, from sdk.AccAddress, uniID string, burnCoin sdk.Coin) sdk.Error {
	pool, existed := k.GetPool(ctx, uniID)
	if !existed {
		return types.ErrReservePoolNotExists(fmt.Sprintf("reserve pool for %s not found", uniID))
	}
	// burn liquidity from pool
	burnCoins := sdk.NewCoins(burnCoin)
	pool.Coins = pool.Coins.Sub(burnCoins)
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

func keyPool(uniID string) []byte {
	return []byte(fmt.Sprintf("%s:%s", KeyPool, uniID))
}
