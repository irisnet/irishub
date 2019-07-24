package keeper

import (
	"strconv"

	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"

	"github.com/irisnet/irishub/modules/bank"
	"github.com/irisnet/irishub/modules/params"
	"github.com/irisnet/irishub/modules/stake/types"
)

// keeper of the stake store
type Keeper struct {
	storeKey   sdk.StoreKey
	storeTKey  sdk.StoreKey
	cdc        *codec.Codec
	bankKeeper bank.Keeper
	hooks      sdk.StakingHooks
	paramstore params.Subspace

	// codespace
	codespace sdk.CodespaceType
	// metrics
	metrics *Metrics
}

func NewKeeper(cdc *codec.Codec, key, tkey sdk.StoreKey, ck bank.Keeper, paramstore params.Subspace, codespace sdk.CodespaceType, metrics *Metrics) Keeper {
	keeper := Keeper{
		storeKey:   key,
		storeTKey:  tkey,
		cdc:        cdc,
		bankKeeper: ck,
		paramstore: paramstore.WithTypeTable(ParamTypeTable()),
		hooks:      nil,
		codespace:  codespace,
		metrics:    metrics,
	}
	return keeper
}

// Set the validator hooks
func (k *Keeper) SetHooks(sh sdk.StakingHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set validator hooks twice")
	}
	k.hooks = sh
	return k
}

//_________________________________________________________________________

// return the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

//_______________________________________________________________________

// load the pool
func (k Keeper) GetPool(ctx sdk.Context) (pool types.Pool) {
	var bondedPool types.BondedPool
	store := ctx.KVStore(k.storeKey)
	b := store.Get(PoolKey)
	if b == nil {
		panic("stored pool should not have been nil")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &bondedPool)
	pool = types.Pool{
		BondedPool: bondedPool,
		BankKeeper: k.bankKeeper,
	}
	return
}

// set the pool
func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(pool.BondedPool)
	store.Set(PoolKey, b)
}

// load the pool status
func (k Keeper) GetPoolStatus(ctx sdk.Context) (poolStatus types.PoolStatus) {
	pool := k.GetPool(ctx)
	poolStatus = types.PoolStatus{
		LooseTokens:  pool.GetLoosenTokenAmount(ctx),
		BondedTokens: pool.BondedPool.BondedTokens,
	}
	return
}

//_______________________________________________________________________

// Load the last total validator power.
func (k Keeper) GetLastTotalPower(ctx sdk.Context) (power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(LastTotalPowerKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &power)
	return
}

// Set the last total validator power.
func (k Keeper) SetLastTotalPower(ctx sdk.Context, power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(power)
	store.Set(LastTotalPowerKey, b)
}

//_______________________________________________________________________

// Load the last validator power.
// Returns zero if the operator was not a validator last block.
func (k Keeper) GetLastValidatorPower(ctx sdk.Context, operator sdk.ValAddress) (power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(GetLastValidatorPowerKey(operator))
	if bz == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &power)
	return
}

// Set the last validator power.
func (k Keeper) SetLastValidatorPower(ctx sdk.Context, operator sdk.ValAddress, power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(power)
	store.Set(GetLastValidatorPowerKey(operator), bz)
}

// Iterate over last validator powers.
func (k Keeper) IterateLastValidatorPowers(ctx sdk.Context, handler func(operator sdk.ValAddress, power sdk.Int) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, LastValidatorPowerKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		addr := sdk.ValAddress(iter.Key()[len(LastValidatorPowerKey):])
		var power sdk.Int
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &power)
		if handler(addr, power) {
			break
		}
	}
}

// Delete the last validator power.
func (k Keeper) DeleteLastValidatorPower(ctx sdk.Context, operator sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetLastValidatorPowerKey(operator))
}

func (k Keeper) BondDenom() string {
	return types.StakeDenom
}

func (k Keeper) UpdateMetrics(ctx sdk.Context) {
	burnedToken, err := strconv.ParseFloat(sdk.NewDecFromInt(k.bankKeeper.GetBurnedCoins(ctx).AmountOf(types.StakeDenom)).QuoInt(sdk.AttoScaleFactor).String(), 64)
	if err == nil {
		k.metrics.BurnedToken.Set(burnedToken)
	}

	loosenToken, err := strconv.ParseFloat(sdk.NewDecFromInt(k.bankKeeper.GetLoosenCoins(ctx).AmountOf(types.StakeDenom)).QuoInt(sdk.AttoScaleFactor).String(), 64)
	if err == nil {
		k.metrics.LoosenToken.Set(loosenToken)
	}
}
