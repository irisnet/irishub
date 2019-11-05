package keeper

import (
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/distribution/types"
	"github.com/irisnet/irishub/modules/params"
	sdk "github.com/irisnet/irishub/types"
)

// keeper of the stake store
type Keeper struct {
	storeKey    sdk.StoreKey
	cdc         *codec.Codec
	paramSpace  params.Subspace
	bankKeeper  types.BankKeeper
	stakeKeeper types.StakeKeeper
	feeKeeper   types.FeeKeeper

	// codespace
	codespace sdk.CodespaceType
	// metrics
	metrics *Metrics
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace, ck types.BankKeeper,
	sk types.StakeKeeper, fk types.FeeKeeper, codespace sdk.CodespaceType, metrics *Metrics) Keeper {

	keeper := Keeper{
		storeKey:    key,
		cdc:         cdc,
		paramSpace:  paramSpace.WithTypeTable(ParamTypeTable()),
		bankKeeper:  ck,
		stakeKeeper: sk,
		feeKeeper:   fk,
		codespace:   codespace,
		metrics:     metrics,
	}
	return keeper
}

//______________________________________________________________________

// get the global fee pool distribution info
func (k Keeper) GetFeePool(ctx sdk.Context) (feePool types.FeePool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(FeePoolKey)
	if b == nil {
		panic("Stored fee pool should not have been nil")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &feePool)
	return
}

// set the global fee pool distribution info
func (k Keeper) SetGenesisFeePool(ctx sdk.Context, feePool types.FeePool) {
	coins, _ := feePool.CommunityPool.TruncateDecimal()
	k.bankKeeper.IncreaseLoosenToken(ctx, coins)
	feePool.CommunityPool = types.NewDecCoins(coins)
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(feePool)
	store.Set(FeePoolKey, b)
}

// set the global fee pool distribution info
func (k Keeper) SetFeePool(ctx sdk.Context, feePool types.FeePool) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(feePool)
	store.Set(FeePoolKey, b)
}

// get the total validator accum for the ctx height
// in the fee pool
func (k Keeper) GetFeePoolValAccum(ctx sdk.Context) sdk.Dec {

	// withdraw self-delegation
	height := ctx.BlockHeight()
	totalPower := sdk.NewDecFromInt(k.stakeKeeper.GetLastTotalPower(ctx))
	fp := k.GetFeePool(ctx)
	return fp.GetTotalValAccum(height, totalPower)
}

//______________________________________________________________________

// set the proposer public key for this block
func (k Keeper) GetPreviousProposerConsAddr(ctx sdk.Context) (consAddr sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)

	b := store.Get(ProposerKey)
	if b == nil {
		panic("Previous proposer not set")
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &consAddr)
	return
}

// get the proposer public key for this block
func (k Keeper) SetPreviousProposerConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(consAddr)
	store.Set(ProposerKey, b)
}

//______________________________________________________________________

// get context required for withdraw operations
func (k Keeper) GetWithdrawContext(ctx sdk.Context,
	valOperatorAddr sdk.ValAddress) types.WithdrawContext {

	feePool := k.GetFeePool(ctx)
	height := ctx.BlockHeight()
	validator := k.stakeKeeper.Validator(ctx, valOperatorAddr)
	lastValPower := k.stakeKeeper.GetLastValidatorPower(ctx, valOperatorAddr)
	lastTotalPower := sdk.NewDecFromInt(k.stakeKeeper.GetLastTotalPower(ctx))

	return types.NewWithdrawContext(
		feePool, height, lastTotalPower, sdk.NewDecFromInt(lastValPower),
		validator.GetCommission())
}

//__________________________________________________________________________________
// used in simulation

// iterate over all the validator distribution infos (inefficient, just used to check invariants)
func (k Keeper) IterateValidatorDistInfos(ctx sdk.Context,
	fn func(index int64, distInfo types.ValidatorDistInfo) (stop bool)) {

	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, ValidatorDistInfoKey)
	defer iter.Close()
	index := int64(0)
	for ; iter.Valid(); iter.Next() {
		var vdi types.ValidatorDistInfo
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &vdi)
		if fn(index, vdi) {
			return
		}
		index++
	}
}
