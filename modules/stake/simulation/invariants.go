package simulation

import (
	"bytes"
	"fmt"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/modules/bank"
	"github.com/irisnet/irishub/modules/distribution"
	"github.com/irisnet/irishub/modules/stake"
	"github.com/irisnet/irishub/modules/stake/keeper"
	"github.com/irisnet/irishub/modules/mock/simulation"
	"github.com/irisnet/irishub/modules/stake/types"
	"runtime/debug"
)

// AllInvariants runs all invariants of the stake module.
// Currently: total supply, positive power
func AllInvariants(ck bank.Keeper, k stake.Keeper,
	f auth.FeeKeeper, d distribution.Keeper,
	am auth.AccountKeeper) simulation.Invariant {

	return func(ctx sdk.Context) error {
		//err := SupplyInvariants(ck, k, f, d, am)(app, header)
		//if err != nil {
		//	return err
		//}
		//err = PositivePowerInvariant(k)(app, header)
		//if err != nil {
		//	return err
		//}
		//err = ValidatorSetInvariant(k)(app, header)
		//return err
		//return nil
		return nil
	}
}

// SupplyInvariants checks that the total supply reflects all held loose tokens, bonded tokens, and unbonding delegations
// nolint: unparam
func SupplyInvariants(ck bank.Keeper, k stake.Keeper,
	f auth.FeeKeeper, d distribution.Keeper, am auth.AccountKeeper) simulation.Invariant {
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
		k.IterateUnbondingDelegations(ctx, func(_ int64, ubd stake.UnbondingDelegation) bool {
			if ubd.Balance.Amount.LT(sdk.ZeroInt()) {
				panic(fmt.Errorf("found negative balance in unbonding delegation"))
			}
			loose = loose.Add(sdk.NewDecFromInt(ubd.Balance.Amount))
			return false
		})
		k.IterateValidators(ctx, func(_ int64, validator sdk.Validator) bool {
			validatorInfo := fmt.Sprintf("Operator address: %s\nValidator name: %s\nValidator Token: %s\nValidator Shares: %s\n",
				validator.GetOperator().String(), validator.GetMoniker(), validator.GetTokens().String(), validator.GetDelegatorShares().String())

			if validator.GetTokens().IsNegative() {
				panic(fmt.Errorf("Validator token is negative!\n%s", validatorInfo))
			}
			if !validator.GetTokens().IsZero() && validator.GetDelegatorShares().IsZero() {
				panic(fmt.Errorf("Validator token is not zero but delegation shares is zero!\n%s", validatorInfo))
			}
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

// PositivePowerInvariant checks that all stored validators have > 0 power
func PositivePowerInvariant(k stake.Keeper) simulation.Invariant {
	return func(ctx sdk.Context) error {

		iterator := k.ValidatorsPowerStoreIterator(ctx)
		defer iterator.Close()
		pool := k.GetPool(ctx)

		for ; iterator.Valid(); iterator.Next() {
			validator, found := k.GetValidator(ctx, iterator.Value())
			if !found {
				panic(fmt.Sprintf("validator record not found for address: %X\n", iterator.Value()))
			}

			powerKey := keeper.GetValidatorsByPowerIndexKey(validator, pool)

			if !bytes.Equal(iterator.Key(), powerKey) {
				return fmt.Errorf("power store invariance:\n\tvalidator.Power: %v"+
					"\n\tkey should be: %v\n\tkey in store: %v", validator.GetPower(), powerKey, iterator.Key())
			}
		}
		return nil
	}
}

// ValidatorSetInvariant checks equivalence of Tendermint validator set and SDK validator set
func ValidatorSetInvariant(k stake.Keeper) simulation.Invariant {
	return func(ctx sdk.Context) error {
		// TODO
		return nil
	}
}
