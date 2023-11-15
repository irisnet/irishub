package app

import (
	"fmt"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/irisnet/irishub/v2/app/upgrades"
	v200 "github.com/irisnet/irishub/v2/app/upgrades/v200"
	v210 "github.com/irisnet/irishub/v2/app/upgrades/v210"
)

var (
	router = upgrades.NewUpgradeRouter().
		Register(v200.Upgrade).
		Register(v210.Upgrade)
)

// RegisterUpgradePlans register a handler of upgrade plan
func (app *IrisApp) RegisterUpgradePlans() {
	app.setupUpgradeStoreLoaders()
	app.setupUpgradeHandlers()
}

func (app *IrisApp) appKeepers() upgrades.AppKeepers {
	return upgrades.AppKeepers{
		AppCodec:              app.AppCodec(),
		HTLCKeeper:            app.HTLCKeeper,
		BankKeeper:            app.BankKeeper,
		AccountKeeper:         app.AccountKeeper,
		ServiceKeeper:         app.ServiceKeeper,
		GetKey:                app.GetKey,
		ModuleManager:         app.mm,
		TIBCkeeper:            app.TIBCKeeper,
		IBCKeeper:             app.IBCKeeper,
		EvmKeeper:             app.EvmKeeper,
		FeeMarketKeeper:       app.FeeMarketKeeper,
		TokenKeeper:           app.TokenKeeper,
		ReaderWriter:          app,
		ConsensusParamsKeeper: app.ConsensusParamsKeeper,
		ParamsKeeper:          app.ParamsKeeper,
	}
}

// configure store loader that checks if version == upgradeHeight and applies store upgrades
func (app *IrisApp) setupUpgradeStoreLoaders() {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	app.SetStoreLoader(
		upgradetypes.UpgradeStoreLoader(
			upgradeInfo.Height,
			router.UpgradeInfo(upgradeInfo.Name).StoreUpgrades,
		),
	)
}

func (app *IrisApp) setupUpgradeHandlers() {
	for upgradeName, upgrade := range router.Routers() {
		app.UpgradeKeeper.SetUpgradeHandler(
			upgradeName,
			upgrade.UpgradeHandlerConstructor(
				app.mm,
				app.configurator,
				app.appKeepers(),
			),
		)
	}
}
