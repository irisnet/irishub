package v0

import (
	"fmt"

	banksim "github.com/irisnet/irishub/modules/bank/simulation"
	distrsim "github.com/irisnet/irishub/modules/distribution/simulation"
	"github.com/irisnet/irishub/modules/mock/simulation"
	stakesim "github.com/irisnet/irishub/modules/stake/simulation"
	sdk "github.com/irisnet/irishub/types"
)

func (p *ProtocolVersion0) runtimeInvariants() []simulation.Invariant {
	return []simulation.Invariant{
		banksim.NonnegativeBalanceInvariant(p.accountMapper),
		distrsim.ValAccumInvariants(p.distrKeeper, p.StakeKeeper),
		stakesim.SupplyInvariants(p.bankKeeper, p.StakeKeeper,
			p.feeKeeper, p.distrKeeper, p.accountMapper),
		stakesim.PositivePowerInvariant(p.StakeKeeper),
	}
}

func (p *ProtocolVersion0) assertRuntimeInvariants(ctx sdk.Context) {
	if p.invariantLevel != sdk.InvariantError && p.invariantLevel != sdk.InvariantPanic {
		return
	}
	invariants := p.runtimeInvariants()
	for _, inv := range invariants {
		if err := inv(ctx); err != nil {
			if p.invariantLevel == sdk.InvariantPanic {
				panic(fmt.Errorf("invariant broken: %s", err))
			} else {
				p.logger.Error(fmt.Sprintf("Invariant broken: height %d, reason %s", ctx.BlockHeight(), err.Error()))
			}
		}
	}
}
