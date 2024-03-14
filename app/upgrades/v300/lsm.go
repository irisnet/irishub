package v300

import (
	"sort"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// keeper contains the staking keeper functions required
// for the migration
type keeper interface {
	GetAllDelegations(ctx sdk.Context) []types.Delegation
	GetAllValidators(ctx sdk.Context) []types.Validator
	SetDelegation(ctx sdk.Context, delegation types.Delegation)
	SetValidator(ctx sdk.Context, validator types.Validator)
	RefreshTotalLiquidStaked(ctx sdk.Context) error
	GetParams(ctx sdk.Context) (params types.Params)
	SetParams(ctx sdk.Context, params types.Params) error
}

// MigrateParamsStore migrates the params store to the latest version.
//
// ctx - sdk context
// k - keeper
func MigrateParamsStore(ctx sdk.Context, k keeper) {
	params := k.GetParams(ctx)
	params.ValidatorBondFactor = ValidatorBondFactor
	params.ValidatorLiquidStakingCap = ValidatorLiquidStakingCap
	params.GlobalLiquidStakingCap = GlobalLiquidStakingCap
	k.SetParams(ctx, params)
}

// MigrateValidators Set each validator's ValidatorBondShares and LiquidShares to 0
func MigrateValidators(ctx sdk.Context, k keeper) {
	for _, validator := range k.GetAllValidators(ctx) {
		validator.ValidatorBondShares = sdk.ZeroDec()
		validator.LiquidShares = sdk.ZeroDec()
		k.SetValidator(ctx, validator)
	}
}

// MigrateDelegations Set each delegation's ValidatorBond field to false
func MigrateDelegations(ctx sdk.Context, k keeper) {
	for _, delegation := range k.GetAllDelegations(ctx) {
		delegation.ValidatorBond = false
		k.SetDelegation(ctx, delegation)
	}
}

// MigrateUBDEntries will remove the ubdEntries with same creation_height
// and create a new ubdEntry with updated balance and initial_balance
func MigrateUBDEntries(ctx sdk.Context, store storetypes.KVStore, cdc codec.BinaryCodec) error {
	iterator := sdk.KVStorePrefixIterator(store, types.UnbondingDelegationKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		ubd := types.MustUnmarshalUBD(cdc, iterator.Value())

		entriesAtSameCreationHeight := make(map[int64][]types.UnbondingDelegationEntry)
		for _, ubdEntry := range ubd.Entries {
			entriesAtSameCreationHeight[ubdEntry.CreationHeight] = append(entriesAtSameCreationHeight[ubdEntry.CreationHeight], ubdEntry)
		}

		creationHeights := make([]int64, 0, len(entriesAtSameCreationHeight))
		for k := range entriesAtSameCreationHeight {
			creationHeights = append(creationHeights, k)
		}

		sort.Slice(creationHeights, func(i, j int) bool { return creationHeights[i] < creationHeights[j] })

		ubd.Entries = make([]types.UnbondingDelegationEntry, 0, len(creationHeights))

		for _, h := range creationHeights {
			ubdEntry := types.UnbondingDelegationEntry{
				Balance:        sdk.ZeroInt(),
				InitialBalance: sdk.ZeroInt(),
			}
			for _, entry := range entriesAtSameCreationHeight[h] {
				ubdEntry.Balance = ubdEntry.Balance.Add(entry.Balance)
				ubdEntry.InitialBalance = ubdEntry.InitialBalance.Add(entry.InitialBalance)
				ubdEntry.CreationHeight = entry.CreationHeight
				ubdEntry.CompletionTime = entry.CompletionTime
			}
			ubd.Entries = append(ubd.Entries, ubdEntry)
		}

		// set the new ubd to the store
		setUBDToStore(ctx, store, cdc, ubd)
	}
	return nil
}

func setUBDToStore(_ sdk.Context, store storetypes.KVStore, cdc codec.BinaryCodec, ubd types.UnbondingDelegation) {
	delegatorAddress := sdk.MustAccAddressFromBech32(ubd.DelegatorAddress)

	bz := types.MustMarshalUBD(cdc, ubd)

	addr, err := sdk.ValAddressFromBech32(ubd.ValidatorAddress)
	if err != nil {
		panic(err)
	}

	key := types.GetUBDKey(delegatorAddress, addr)

	store.Set(key, bz)
}

// MigrateStore performs the in-place store migration for adding LSM support to v0.45.16-ics, including:
//   - Adding params ValidatorBondFactor, GlobalLiquidStakingCap, ValidatorLiquidStakingCap
//   - Setting each validator's ValidatorBondShares and LiquidShares to 0
//   - Setting each delegation's ValidatorBond field to false
//   - Calculating the total liquid staked by summing the delegations from ICA accounts
func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec, k keeper) error {
	store := ctx.KVStore(storeKey)

	ctx.Logger().Info("Staking LSM Migration: Migrating param store")
	MigrateParamsStore(ctx, k)

	ctx.Logger().Info("Staking LSM Migration: Migrating validators")
	MigrateValidators(ctx, k)

	ctx.Logger().Info("Staking LSM Migration: Migrating delegations")
	MigrateDelegations(ctx, k)

	ctx.Logger().Info("Staking LSM Migration: Migrating UBD entries")
	if err := MigrateUBDEntries(ctx, store, cdc); err != nil {
		return err
	}

	ctx.Logger().Info("Staking LSM Migration: Calculating total liquid staked")
	return k.RefreshTotalLiquidStaked(ctx)
}
