package lsm

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// SDK v0.53 no longer registers the old bank/staking/distribution crisis
// invariants. Check the accounting relationships touched by this migration
// explicitly, both before and after changing the initialized genesis.
func checkInvariants(ctx sdk.Context, keepers Keepers) error {
	bank := keepers.Bank.ExportGenesis(ctx)
	supply := sdk.NewMapCoins(sdk.NewCoins())
	for _, balance := range bank.Balances {
		if err := balance.Coins.Validate(); err != nil {
			return err
		}
		supply.Add(balance.Coins...)
	}
	if !supply.ToCoins().Equal(bank.Supply) {
		return fmt.Errorf("LSM genesis bank balances do not match supply")
	}
	validators, err := keepers.Staking.GetAllValidators(ctx)
	if err != nil {
		return err
	}
	delegations, err := keepers.Staking.GetAllDelegations(ctx)
	if err != nil {
		return err
	}
	shares := make(map[string]math.LegacyDec)
	delegationKeys := make(map[string]bool)
	for _, delegation := range delegations {
		if !delegation.Shares.IsPositive() {
			return fmt.Errorf("LSM genesis has nonpositive delegation shares")
		}
		sum, ok := shares[delegation.ValidatorAddress]
		if !ok {
			sum = math.LegacyZeroDec()
		}
		shares[delegation.ValidatorAddress] = sum.Add(delegation.Shares)
		delegationKeys[delegation.DelegatorAddress+"/"+delegation.ValidatorAddress] = true
	}
	for _, validator := range validators {
		sum, ok := shares[validator.OperatorAddress]
		if !ok {
			sum = math.LegacyZeroDec()
		}
		if !sum.Equal(validator.DelegatorShares) {
			return fmt.Errorf("LSM genesis delegation shares do not match validator %s", validator.OperatorAddress)
		}
		delete(shares, validator.OperatorAddress)
	}
	if len(shares) != 0 {
		return fmt.Errorf("LSM genesis delegation references a missing validator")
	}
	distribution := keepers.Distribution.ExportGenesis(ctx)
	references := make(map[string]uint32)
	periodKey := func(validator string, period uint64) string { return fmt.Sprintf("%s/%d", validator, period) }
	for _, info := range distribution.DelegatorStartingInfos {
		key := info.DelegatorAddress + "/" + info.ValidatorAddress
		if !delegationKeys[key] {
			return fmt.Errorf("LSM genesis reward starting info has no delegation: %s", key)
		}
		delete(delegationKeys, key)
		references[periodKey(info.ValidatorAddress, info.StartingInfo.PreviousPeriod)]++
	}
	if len(delegationKeys) != 0 {
		return fmt.Errorf("LSM genesis delegation has no reward starting info")
	}
	for _, current := range distribution.ValidatorCurrentRewards {
		if current.Rewards.Period == 0 {
			return fmt.Errorf("LSM genesis has zero current reward period")
		}
		references[periodKey(current.ValidatorAddress, current.Rewards.Period-1)]++
	}
	for _, event := range distribution.ValidatorSlashEvents {
		references[periodKey(event.ValidatorAddress, event.ValidatorSlashEvent.ValidatorPeriod)]++
	}
	for _, historical := range distribution.ValidatorHistoricalRewards {
		key := periodKey(historical.ValidatorAddress, historical.Period)
		if historical.Rewards.ReferenceCount != references[key] {
			return fmt.Errorf("LSM genesis reward reference count mismatch at %s", key)
		}
		delete(references, key)
	}
	if len(references) != 0 {
		return fmt.Errorf("LSM genesis references missing historical rewards")
	}
	liabilities := distribution.FeePool.CommunityPool
	for _, outstanding := range distribution.OutstandingRewards {
		liabilities = liabilities.Add(outstanding.OutstandingRewards...)
	}
	owed, _ := liabilities.TruncateDecimal()
	if !keepers.Bank.GetAllBalances(ctx, authtypes.NewModuleAddress(distrtypes.ModuleName)).IsAllGTE(owed) {
		return fmt.Errorf("LSM genesis distribution balance cannot cover rewards and community pool")
	}
	return nil
}
