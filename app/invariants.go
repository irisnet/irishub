package app

import (
	"fmt"
	"time"

	sdk "github.com/irisnet/irishub/types"
	banksim "github.com/irisnet/irishub/modules/bank/simulation"
	distrsim "github.com/irisnet/irishub/modules/distribution/simulation"
	"github.com/irisnet/irishub/modules/mock/simulation"
	stakesim "github.com/irisnet/irishub/modules/stake/simulation"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (app *IrisApp) runtimeInvariants() []simulation.Invariant {
	return []simulation.Invariant{
		banksim.NonnegativeBalanceInvariant(app.accountMapper),
		distrsim.ValAccumInvariants(app.distrKeeper, app.stakeKeeper),
		stakesim.SupplyInvariants(app.bankKeeper, app.stakeKeeper,
			app.feeCollectionKeeper, app.distrKeeper, app.serviceKeeper, app.accountMapper),
		stakesim.PositivePowerInvariant(app.stakeKeeper),
	}
}

func (app *IrisApp) assertRuntimeInvariants() {
	ctx := app.NewContext(false, abci.Header{Height: app.LastBlockHeight() + 1})
	app.assertRuntimeInvariantsOnContext(ctx)
}

func (app *IrisApp) assertRuntimeInvariantsOnContext(ctx sdk.Context) {
	start := time.Now()
	invariants := app.runtimeInvariants()
	for _, inv := range invariants {
		if err := inv(ctx); err != nil {
			panic(fmt.Errorf("invariant broken: %s", err))
		}
	}
	end := time.Now()
	diff := end.Sub(start)
	app.BaseApp.Logger.With("module", "invariants").Info("Asserted all invariants", "duration", diff)
}
