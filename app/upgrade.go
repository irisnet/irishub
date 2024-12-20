package app

import (
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/irisnet/irishub/v4/app/upgrades"
	v200 "github.com/irisnet/irishub/v4/app/upgrades/v200"
	v210 "github.com/irisnet/irishub/v4/app/upgrades/v210"
	v300 "github.com/irisnet/irishub/v4/app/upgrades/v300"
	v400 "github.com/irisnet/irishub/v4/app/upgrades/v400"
)

var (
	router = upgrades.NewUpgradeRouter().
		Register(v200.Upgrade).
		Register(v210.Upgrade).
		Register(v300.Upgrade).
		Register(v400.Upgrade)
)

// RegisterUpgradePlans register a handler of upgrade plan
func (app *IrisApp) RegisterUpgradePlans() {
	app.setupUpgradeStoreLoaders()
	app.setupUpgradeHandlers()
}

func (app *IrisApp) toolbox() upgrades.Toolbox {
	return upgrades.Toolbox{
		AppCodec:      app.AppCodec(),
		ModuleManager: app.mm,
		ReaderWriter:  app,
		AppKeepers:    app.AppKeepers,
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
	box := app.toolbox()
	for upgradeName, upgrade := range router.Routers() {
		app.UpgradeKeeper.SetUpgradeHandler(
			upgradeName,
			upgrade.UpgradeHandlerConstructor(
				app.mm,
				app.configurator,
				box,
			),
		)
	}
}
