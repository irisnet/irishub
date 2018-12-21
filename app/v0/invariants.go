package v0

import (
	"fmt"
	sdk "github.com/irisnet/irishub/types"
	banksim "github.com/irisnet/irishub/modules/bank/simulation"
	distrsim "github.com/irisnet/irishub/modules/distribution/simulation"
	"github.com/irisnet/irishub/modules/mock/simulation"
	stakesim "github.com/irisnet/irishub/modules/stake/simulation"
)

func (p *ProtocolVersion0) runtimeInvariants() []simulation.Invariant {
	return []simulation.Invariant{
		banksim.NonnegativeBalanceInvariant(p.accountMapper),
		distrsim.ValAccumInvariants(p.distrKeeper, p.StakeKeeper),
		stakesim.SupplyInvariants(p.bankKeeper, p.StakeKeeper,
			p.feeCollectionKeeper, p.distrKeeper, p.accountMapper),
		stakesim.PositivePowerInvariant(p.StakeKeeper),
	}
}

func (p *ProtocolVersion0) assertRuntimeInvariants(ctx sdk.Context) {
	invariants := p.runtimeInvariants()
	for _, inv := range invariants {
		if err := inv(ctx); err != nil {
			panic(fmt.Errorf("invariant broken: %s", err))
		}
	}
}
