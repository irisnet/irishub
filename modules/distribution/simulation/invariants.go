package simulation

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
	distr "github.com/irisnet/irishub/modules/distribution"
	"github.com/irisnet/irishub/modules/mock/simulation"
)

// AllInvariants runs all invariants of the distribution module
// Currently: total supply, positive power
func AllInvariants(d distr.Keeper, sk distr.StakeKeeper) simulation.Invariant {
	return func(ctx sdk.Context) error {
		err := ValAccumInvariants(d, sk)(ctx)
		if err != nil {
			return err
		}
		return nil
	}
}

// ValAccumInvariants checks that the fee pool accum == sum all validators' accum
func ValAccumInvariants(k distr.Keeper, sk distr.StakeKeeper) simulation.Invariant {

	return func(ctx sdk.Context) error {
		height := ctx.BlockHeight()

		valAccum := sdk.ZeroDec()
		k.IterateValidatorDistInfos(ctx, func(_ int64, vdi distr.ValidatorDistInfo) bool {
			lastValPower := sk.GetLastValidatorPower(ctx, vdi.OperatorAddr)
			valAccum = valAccum.Add(vdi.GetValAccum(height, sdk.NewDecFromInt(lastValPower)))
			return false
		})

		lastTotalPower := sdk.NewDecFromInt(sk.GetLastTotalPower(ctx))
		totalAccum := k.GetFeePool(ctx).GetTotalValAccum(height, lastTotalPower)

		if !totalAccum.Equal(valAccum) {
			return fmt.Errorf("validator accum invariance: \n\tfee pool totalAccum: %v"+
				"\n\tvalidator accum \t%v\n", totalAccum.String(), valAccum.String())
		}

		return nil
	}
}
