package v0

import (
	"fmt"

	"github.com/irisnet/irishub/modules/bank"
	distr "github.com/irisnet/irishub/modules/distribution"
	"github.com/irisnet/irishub/modules/stake"
	sdk "github.com/irisnet/irishub/types"
)

func (p *ProtocolV0) runtimeInvariants() []sdk.Invariant {
	return []sdk.Invariant{
		bank.NonnegativeBalanceInvariant(p.accountMapper),

		distr.ValAccumInvariants(p.distrKeeper, p.StakeKeeper),
		distr.DelAccumInvariants(p.distrKeeper, p.StakeKeeper),
		distr.CanWithdrawInvariant(p.distrKeeper, p.StakeKeeper),

		stake.SupplyInvariants(p.bankKeeper, p.StakeKeeper,
			p.feeKeeper, p.distrKeeper, p.accountMapper),
		stake.NonNegativePowerInvariant(p.StakeKeeper),
		stake.PositiveDelegationInvariant(p.StakeKeeper),
		stake.DelegatorSharesInvariant(p.StakeKeeper),
	}
}

func (p *ProtocolV0) assertRuntimeInvariants(ctx sdk.Context) {
	if p.invariantLevel != sdk.InvariantError && p.invariantLevel != sdk.InvariantPanic {
		return
	}
	if p.invariantLevel == sdk.InvariantError && !p.checkInvariant {
		return
	}
	invariants := p.runtimeInvariants()
	ctx = ctx.WithLogger(ctx.Logger().With("module", "iris/invariant"))
	for _, inv := range invariants {
		if err := inv(ctx); err != nil {
			if p.invariantLevel == sdk.InvariantPanic {
				panic(fmt.Errorf("invariant broken: %s", err))
			} else {
				p.metrics.InvariantFailure.With("error", err.Error()).Add(float64(1))
				p.logger.Error(fmt.Sprintf("Invariant broken: height %d, reason %s", ctx.BlockHeight(), err.Error()))
			}
		}
	}
}
