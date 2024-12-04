package v300

import (
	"context"
	"sort"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// keeper contains the staking keeper functions required
// for the migration
type keeper interface {
	GetAllDelegations(ctx context.Context) ([]types.Delegation, error)
	GetAllValidators(ctx context.Context) ([]types.Validator, error)
	SetDelegation(ctx context.Context, delegation types.Delegation) error
	SetValidator(ctx context.Context, validator types.Validator) error
	RefreshTotalLiquidStaked(ctx context.Context) error
	GetParams(ctx context.Context) (params types.Params, err error)
	SetParams(ctx context.Context, params types.Params) error
}

// migrateParamsStore migrates the params store to the latest version.
//
// ctx - sdk context
// k - keeper
func migrateParamsStore(ctx sdk.Context, k keeper) error {
	params, err := k.GetParams(ctx)
	if err != nil {
		return err
	}
	params.ValidatorBondFactor = ValidatorBondFactor
	params.ValidatorLiquidStakingCap = ValidatorLiquidStakingCap
	params.GlobalLiquidStakingCap = GlobalLiquidStakingCap
	return k.SetParams(ctx, params)
}

// migrateValidators Set each validator's ValidatorBondShares and LiquidShares to 0
func migrateValidators(ctx sdk.Context, k keeper) error {
	validators, err := k.GetAllValidators(ctx)
	if err != nil {
		return err
	}
	for _, validator := range validators {
		validator.ValidatorBondShares = math.LegacyZeroDec()
		validator.LiquidShares = math.LegacyZeroDec()
		if err := k.SetValidator(ctx, validator); err != nil {
			return err
		}
	}
	return nil
}

// migrateDelegations Set each delegation's ValidatorBond field to false
func migrateDelegations(ctx sdk.Context, k keeper) error {
	delegations, err := k.GetAllDelegations(ctx)
	if err != nil {
		return err
	}
	for _, delegation := range delegations {
		delegation.ValidatorBond = false
		if err := k.SetDelegation(ctx, delegation); err != nil {
			return err
		}
	}
	return nil
}

// MigrateUBDEntries will remove the ubdEntries with same creation_height
// and create a new ubdEntry with updated balance and initial_balance
func migrateUBDEntries(ctx sdk.Context, store storetypes.KVStore, cdc codec.BinaryCodec) error {
	iterator := storetypes.KVStorePrefixIterator(store, types.UnbondingDelegationKey)
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
				Balance:        math.ZeroInt(),
				InitialBalance: math.ZeroInt(),
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

// migrateStore performs the in-place store migration for adding LSM support to v0.45.16-ics, including:
//   - Adding params ValidatorBondFactor, GlobalLiquidStakingCap, ValidatorLiquidStakingCap
//   - Setting each validator's ValidatorBondShares and LiquidShares to 0
//   - Setting each delegation's ValidatorBond field to false
//   - Calculating the total liquid staked by summing the delegations from ICA accounts
func migrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec, k keeper) error {
	store := ctx.KVStore(storeKey)

	ctx.Logger().Info("Staking LSM Migration: Migrating param store")
	if err := migrateParamsStore(ctx, k); err != nil {
		return err
	}

	ctx.Logger().Info("Staking LSM Migration: Migrating validators")
	if err := migrateValidators(ctx, k); err != nil {
		return err
	}

	ctx.Logger().Info("Staking LSM Migration: Migrating delegations")
	migrateDelegations(ctx, k)

	ctx.Logger().Info("Staking LSM Migration: Migrating UBD entries")
	if err := migrateUBDEntries(ctx, store, cdc); err != nil {
		return err
	}

	ctx.Logger().Info("Staking LSM Migration: Calculating total liquid staked")
	return k.RefreshTotalLiquidStaked(ctx)
}
