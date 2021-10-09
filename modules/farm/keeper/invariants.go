package keeper

// DONTCOVER

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/farm/types"
)

// RegisterInvariants registers all invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "reward", RewardInvariant(k))
}

// AllInvariants runs all invariants of the farm module.
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		return RewardInvariant(k)(ctx)
	}
}

// RewardInvariant checks whether the amount of module account is consistent with the recorded in the farm
func RewardInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		expectedBalance := sdk.Coins{}
		balance := k.bk.GetAllBalances(ctx, k.ak.GetModuleAddress(types.ModuleName))

		k.IteratorAllPools(ctx, func(pool types.FarmPool) {
			expectedBalance = expectedBalance.Add(pool.TotalLptLocked)
			k.IteratorRewardRules(ctx, pool.Name, func(r types.RewardRule) {
				expectedBalance = expectedBalance.Add(sdk.NewCoin(r.Reward, r.RemainingReward))
			})
		})

		broken := !expectedBalance.IsEqual(balance)
		return sdk.FormatInvariant(
			types.ModuleName,
			"module account balance",
			fmt.Sprintf(
				"\tsum of accounts coins: %v\n"+
					"\tbalance:          %v\n",
				expectedBalance, balance,
			),
		), broken
	}
}
