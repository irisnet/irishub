package stake

import (
	"bytes"
	"fmt"
	"runtime/debug"

	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/modules/bank"
	"github.com/irisnet/irishub/modules/distribution"
	"github.com/irisnet/irishub/modules/stake/keeper"
	"github.com/irisnet/irishub/modules/stake/types"
	sdk "github.com/irisnet/irishub/types"
)

// AllInvariants runs all invariants of the stake module.
// Currently: total supply, positive power
func AllInvariants(ck bank.Keeper, k Keeper,
	f auth.FeeKeeper, d distribution.Keeper,
	am auth.AccountKeeper) sdk.Invariant {

	return func(ctx sdk.Context) error {
		err := SupplyInvariants(ck, k, f, d, am)(ctx)
		if err != nil {
			return err
		}

		err = NonNegativePowerInvariant(k)(ctx)
		if err != nil {
			return err
		}

		err = PositiveDelegationInvariant(k)(ctx)
		if err != nil {
			return err
		}

		err = DelegatorSharesInvariant(k)(ctx)
		if err != nil {
			return err
		}

		return nil
	}
}

// SupplyInvariants checks that the total supply reflects all held loose tokens, bonded tokens, and unbonding delegations
// nolint: unparam
func SupplyInvariants(ck bank.Keeper, k Keeper,
	f auth.FeeKeeper, d distribution.Keeper, am auth.AccountKeeper) sdk.Invariant {
	return func(ctx sdk.Context) (err error) {

		defer func() {
			if r := recover(); r != nil {
				switch rType := r.(type) {
				case error:
					err = rType
				default:
					err = fmt.Errorf(string(debug.Stack()))
				}
			}
		}()

		pool := k.GetPool(ctx)

		loose := sdk.ZeroDec()
		bonded := sdk.ZeroDec()
		am.IterateAccounts(ctx, func(acc auth.Account) bool {
			loose = loose.Add(sdk.NewDecFromInt(acc.GetCoins().AmountOf(types.StakeDenom)))
			return false
		})
		k.IterateUnbondingDelegations(ctx, func(_ int64, ubd UnbondingDelegation) bool {
			if ubd.Balance.Amount.LT(sdk.ZeroInt()) {
				panic(fmt.Errorf("found negative balance in unbonding delegation"))
			}
			loose = loose.Add(sdk.NewDecFromInt(ubd.Balance.Amount))
			return false
		})
		k.IterateValidators(ctx, func(_ int64, validator sdk.Validator) bool {
			switch validator.GetStatus() {
			case sdk.Bonded:
				bonded = bonded.Add(validator.GetTokens())
			case sdk.Unbonding:
				loose = loose.Add(validator.GetTokens())
			case sdk.Unbonded:
				loose = loose.Add(validator.GetTokens())
			}
			return false
		})

		feePool := d.GetFeePool(ctx)

		// add outstanding fees
		loose = loose.Add(sdk.NewDecFromInt(f.GetCollectedFees(ctx).AmountOf(types.StakeDenom)))

		// add community pool
		loose = loose.Add(feePool.CommunityPool.AmountOf(types.StakeDenom))

		// add validator distribution pool
		loose = loose.Add(feePool.ValPool.AmountOf(types.StakeDenom))

		// add validator distribution commission and yet-to-be-withdrawn-by-delegators
		d.IterateValidatorDistInfos(ctx,
			func(_ int64, distInfo distribution.ValidatorDistInfo) (stop bool) {
				loose = loose.Add(distInfo.DelPool.AmountOf(types.StakeDenom))
				loose = loose.Add(distInfo.ValCommission.AmountOf(types.StakeDenom))
				return false
			},
		)

		// Loose tokens should equal coin supply plus unbonding delegations
		// plus tokens on unbonded validators
		if !pool.GetLoosenTokenAmount(ctx).Equal(loose) {
			return fmt.Errorf("loose token invariance:\n\tbank.LooseTokens: %v"+
				"\n\tsum of account tokens: %v", pool.GetLoosenTokenAmount(ctx).TruncateInt(), loose.TruncateInt())
		}

		// Bonded tokens should equal sum of tokens with bonded validators
		if !pool.BondedPool.BondedTokens.Equal(bonded) {
			return fmt.Errorf("bonded token invariance:\n\tpool.BondedTokens: %v"+
				"\n\tsum of account tokens: %v", pool.BondedPool.BondedTokens, bonded)
		}

		return nil
	}
}

// NonNegativePowerInvariant checks that all stored validators have >= 0 power.
func NonNegativePowerInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (err error) {

		defer func() {
			if r := recover(); r != nil {
				switch rType := r.(type) {
				case error:
					err = rType
				default:
					err = fmt.Errorf(string(debug.Stack()))
				}
			}
		}()

		iterator := k.ValidatorsPowerStoreIterator(ctx)

		for ; iterator.Valid(); iterator.Next() {
			validator, found := k.GetValidator(ctx, iterator.Value())
			if !found {
				panic(fmt.Sprintf("validator record not found for address: %X\n", iterator.Value()))
			}

			pool := k.GetPool(ctx)
			powerKey := keeper.GetValidatorsByPowerIndexKey(validator, pool)

			validatorInfo := fmt.Sprintf("\n\tOperator address: %s\n\tValidator name: %s\n\tValidator Token: %s\n\tValidator Shares: %s\n",
				validator.GetOperator().String(), validator.GetMoniker(), validator.GetTokens().String(), validator.GetDelegatorShares().String())

			if !bytes.Equal(iterator.Key(), powerKey) {
				return fmt.Errorf("Power store invariance:\n%s", validatorInfo)
			}

			if validator.GetTokens().IsNegative() {
				return fmt.Errorf("Validator token is negative!\n%s", validatorInfo)
			}
			// if validator delegator shares is zero, validator will be deleted once its status becomes unbonded, thus validator tokens will be lost
			if !validator.GetTokens().IsZero() && validator.GetDelegatorShares().IsZero() {
				return fmt.Errorf("Validator token is not zero but delegation shares is zero!\n%s", validatorInfo)
			}
		}

		iterator.Close()
		return nil
	}
}

// PositiveDelegationInvariant checks that all stored delegations have > 0 shares.
func PositiveDelegationInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (err error) {

		defer func() {
			if r := recover(); r != nil {
				switch rType := r.(type) {
				case error:
					err = rType
				default:
					err = fmt.Errorf(string(debug.Stack()))
				}
			}
		}()

		delegations := k.GetAllDelegations(ctx)
		for _, delegation := range delegations {
			if delegation.Shares.IsNegative() {
				return fmt.Errorf("delegation with negative shares: %+v", delegation)
			}
			if delegation.Shares.IsZero() {
				return fmt.Errorf("delegation with zero shares: %+v", delegation)
			}
		}

		return nil
	}
}

// DelegatorSharesInvariant checks whether all the delegator shares which persist
// in the delegator object add up to the correct total delegator shares
// amount stored in each validator
func DelegatorSharesInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (err error) {

		defer func() {
			if r := recover(); r != nil {
				switch rType := r.(type) {
				case error:
					err = rType
				default:
					err = fmt.Errorf(string(debug.Stack()))
				}
			}
		}()

		validators := k.GetAllValidators(ctx)
		for _, validator := range validators {

			valTotalDelShares := validator.GetDelegatorShares()

			totalDelShares := sdk.ZeroDec()
			delegations := k.GetValidatorDelegations(ctx, validator.GetOperator())
			for _, delegation := range delegations {
				totalDelShares = totalDelShares.Add(delegation.Shares)
			}

			if !valTotalDelShares.Equal(totalDelShares) {
				return fmt.Errorf("broken delegator shares invariance:\n"+
					"\tvalidator.DelegatorShares: %v\n"+
					"\tsum of Delegator.Shares: %v", valTotalDelShares, totalDelShares)
			}
		}
		return nil
	}
}
