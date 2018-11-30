package app

import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/bank"
	distr "github.com/irisnet/irishub/modules/distribution"
	"github.com/irisnet/irishub/modules/slashing"
	"github.com/irisnet/irishub/modules/stake"
	"github.com/irisnet/irishub/modules/gov"
	"github.com/irisnet/irishub/modules/record"
	"github.com/irisnet/irishub/modules/service"
	"github.com/irisnet/irishub/modules/upgrade"
	"github.com/irisnet/irishub/modules/profiling"
)

const LastProtocolVersion = 0

func (app *IrisApp) wireRouterForAllVersion() {
	for i := 0; i <= LastProtocolVersion; i++ {
		app.wireRouterForVerion(i)
	}
}

func (app *IrisApp) wireRouterForVerion(version int) {
	if version > LastProtocolVersion {
		panic("The protocal version is not valid!")
	}

	// register message routes
	// need to update each module's msg type
	switch version {
	case 0:
		app.Router().
			AddRoute("bank", []*sdk.KVStoreKey{app.keyAccount}, bank.NewHandler(app.bankKeeper)).
			AddRoute("stake", []*sdk.KVStoreKey{app.keyStake, app.keyAccount, app.keyMint, app.keyDistr}, stake.NewHandler(app.stakeKeeper)).
			AddRoute("slashing", []*sdk.KVStoreKey{app.keySlashing, app.keyStake}, slashing.NewHandler(app.slashingKeeper)).
			AddRoute("distr", []*sdk.KVStoreKey{app.keyDistr}, distr.NewHandler(app.distrKeeper)).
			AddRoute("gov", []*sdk.KVStoreKey{app.keyGov, app.keyAccount, app.keyStake, app.keyParams}, gov.NewHandler(app.govKeeper)).
			AddRoute("upgrade", []*sdk.KVStoreKey{app.keyUpgrade, app.keyStake}, upgrade.NewHandler(app.upgradeKeeper)).
			AddRoute("record", []*sdk.KVStoreKey{app.keyRecord}, record.NewHandler(app.recordKeeper)).
			AddRoute("service", []*sdk.KVStoreKey{app.keyService}, service.NewHandler(app.serviceKeeper)).
			AddRoute("profiling", []*sdk.KVStoreKey{app.KeyProfiling}, profiling.NewHandler(app.profilingKeeper))

		app.QueryRouter().
			AddRoute("gov", gov.NewQuerier(app.govKeeper)).
			AddRoute("stake", stake.NewQuerier(app.stakeKeeper, app.cdc))

		app.hookHub.
			AddHook(stakeTrigger, 0, app.distrKeeper.Hooks()).
			AddHook(stakeTrigger, 0, app.slashingKeeper.Hooks())

		break
	default:
		panic("The protocal version is not valid!")
	}
}
