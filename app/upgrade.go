package app

import (
	"fmt"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/irisnet/irishub/app/upgrades"
	v110 "github.com/irisnet/irishub/app/upgrades/v110"
	v120 "github.com/irisnet/irishub/app/upgrades/v120"
	v130 "github.com/irisnet/irishub/app/upgrades/v130"
	v140 "github.com/irisnet/irishub/app/upgrades/v140"
	v200 "github.com/irisnet/irishub/app/upgrades/v200"
)

var (
	plans = []upgrades.Upgrade{
		v110.Upgrade,
		v120.Upgrade,
		v130.Upgrade,
		v140.Upgrade,
		v200.Upgrade,
	}
)

// RegisterUpgradePlans register a handler of upgrade plan
func (app *IrisApp) RegisterUpgradePlans() {
	app.setupUpgradeStoreLoaders()
	app.setupUpgradeHandlers()
}

func (app *IrisApp) appKeepers() upgrades.AppKeepers {
	return upgrades.AppKeepers{
		AppCodec:        app.AppCodec(),
		HTLCKeeper:      app.HTLCKeeper,
		BankKeeper:      app.BankKeeper,
		AccountKeeper:   app.AccountKeeper,
		ServiceKeeper:   app.ServiceKeeper,
		GetKey:          app.GetKey,
		ModuleManager:   app.mm,
		TIBCkeeper:      app.TIBCKeeper,
		EvmKeeper:       app.EvmKeeper,
		FeeMarketKeeper: app.FeeMarketKeeper,
		TokenKeeper:     app.TokenKeeper,
		ReaderWriter:    app,
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

	for _, upgrade := range plans {
		if upgradeInfo.Name == upgrade.UpgradeName {
			app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, upgrade.StoreUpgrades))
		}
	}
}

func (app *IrisApp) setupUpgradeHandlers() {
	for _, upgrade := range plans {
		app.UpgradeKeeper.SetUpgradeHandler(
			upgrade.UpgradeName,
			upgrade.UpgradeHandlerConstructor(
				app.mm,
				app.configurator,
				app.appKeepers(),
			),
		)
	}
}
