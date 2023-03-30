package app

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/irisnet/irishub/app/upgrades"
	v110 "github.com/irisnet/irishub/app/upgrades/v110"
	v120 "github.com/irisnet/irishub/app/upgrades/v120"
	v130 "github.com/irisnet/irishub/app/upgrades/v130"
	v140 "github.com/irisnet/irishub/app/upgrades/v140"
)

var (
	plans = []upgrades.Upgrade{
		v110.Upgrade,
		v120.Upgrade,
		v130.Upgrade,
		v140.Upgrade,
	}
)

// RegisterUpgradePlans register a handler of upgrade plan
func (app *IrisApp) RegisterUpgradePlans() {
	for _, u := range plans {
		app.registerUpgradeHandler(u.UpgradeName,
			u.StoreUpgrades,
			u.UpgradeHandlerConstructor(
				app.mm,
				app.configurator,
				app.appKeepers(),
			),
		)
	}
}

func (app *IrisApp) appKeepers() upgrades.AppKeepers {
	return upgrades.AppKeepers{
		AppCodec:      app.AppCodec(),
		HTLCKeeper:    app.HTLCKeeper,
		BankKeeper:    app.BankKeeper,
		ServiceKeeper: app.ServiceKeeper,
		GetKey:        app.GetKey,
		ModuleManager: app.mm,
		TIBCkeeper:    app.TIBCKeeper,
		ReaderWriter:  app,
	}
}

// registerUpgradeHandler implements the upgrade execution logic of the upgrade module
func (app *IrisApp) registerUpgradeHandler(
	planName string,
	upgrades *storetypes.StoreUpgrades,
	upgradeHandler upgradetypes.UpgradeHandler,
) {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		app.Logger().Info("not found upgrade plan", "planName", planName, "err", err.Error())
		return
	}

	if upgradeInfo.Name == planName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		// this configures a no-op upgrade handler for the planName upgrade
		app.UpgradeKeeper.SetUpgradeHandler(planName, upgradeHandler)
		// configure store loader that checks if version+1 == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, upgrades))
	}
}
