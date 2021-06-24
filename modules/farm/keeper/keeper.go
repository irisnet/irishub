package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/irisnet/irismod/modules/farm/types"
)

// Keeper of the farm store
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      codec.Marshaler

	paramSpace paramstypes.Subspace
	// name of the fee collector
	feeCollectorName string
	validateLPToken  types.ValidateLPToken
	bk               types.BankKeeper
	ak               types.AccountKeeper
}

func NewKeeper(cdc codec.Marshaler,
	storeKey sdk.StoreKey,
	bk types.BankKeeper,
	ak types.AccountKeeper,
	validateLPToken types.ValidateLPToken,
	paramSpace paramstypes.Subspace,
	feeCollectorName string,
) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(ParamKeyTable())
	}

	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// ensure farm module accounts are set
	if addr := ak.GetModuleAddress(types.RewardCollector); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.RewardCollector))
	}

	return Keeper{
		storeKey:         storeKey,
		cdc:              cdc,
		bk:               bk,
		ak:               ak,
		validateLPToken:  validateLPToken,
		paramSpace:       paramSpace,
		feeCollectorName: feeCollectorName,
	}
}

// CreatePool creates an new farm pool
func (k Keeper) SetPool(ctx sdk.Context, pool types.FarmPool) {
	pool.Rules = nil
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&pool)
	store.Set(types.KeyFarmPool(pool.Name), bz)
}

// GetPool return the specified farm pool
func (k Keeper) GetPool(ctx sdk.Context, poolName string) (types.FarmPool, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyFarmPool(poolName))
	if len(bz) == 0 {
		return types.FarmPool{}, false
	}

	var pool types.FarmPool
	k.cdc.MustUnmarshalBinaryBare(bz, &pool)
	return pool, true
}

func (k Keeper) SetRewardRules(ctx sdk.Context, poolName string, rules types.RewardRules) {
	for _, r := range rules {
		k.SetRewardRule(ctx, poolName, r)
	}
}

func (k Keeper) SetRewardRule(ctx sdk.Context, poolName string, rule types.RewardRule) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&rule)
	store.Set(types.KeyRewardRule(poolName, rule.Reward), bz)
}

func (k Keeper) GetRewardRules(ctx sdk.Context, poolName string) (rules types.RewardRules) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PrefixRewardRule(poolName))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var r types.RewardRule
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &r)
		rules = append(rules, r)
	}
	return
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "irismod/farm")
}

func (k Keeper) IteratorRewardRules(ctx sdk.Context, poolName string, fun func(r types.RewardRule)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PrefixRewardRule(poolName))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var r types.RewardRule
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &r)
		fun(r)
	}
}
