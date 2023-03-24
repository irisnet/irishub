package v140

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ica "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts"
	icacontrollertypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/types"

	"github.com/irisnet/irishub/app/upgrades"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:               "v1.4",
	UpgradeHandlerConstructor: upgradeHandlerConstructor,
	StoreUpgrades: &storetypes.StoreUpgrades{
		Added: []string{authzkeeper.StoreKey},
	},
}

func upgradeHandlerConstructor(m *module.Manager, c module.Configurator, app upgrades.AppKeepers) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// version upgrade:
		//	nft :    1 -> 2
		// 	auth:    2 -> 3
		// 	bank:    2 -> 3
		//	coinswap 3 -> 4
		// 	feegrant 1 -> 2
		// 	gov      2 -> 3
		// 	staking  2 -> 3
		// 	upgrade  2 -> 3

		// added module:
		//  authz

		// ibc application:
		//  27-interchain-accounts
		icaModule := app.ModuleManager.Modules[icatypes.ModuleName].(ica.AppModule)
		fromVM[icatypes.ModuleName] = icaModule.ConsensusVersion()
		// create ICS27 Controller submodule params
		controllerParams := icacontrollertypes.Params{}
		// create ICS27 Host submodule params
		hostParams := icahosttypes.Params{
			HostEnabled:   true,
			AllowMessages: msgTypes,
		}

		ctx.Logger().Info("start to init interchainaccount module...")
		// initialize ICS27 module
		icaModule.InitModule(ctx, controllerParams, hostParams)
		ctx.Logger().Info("start to run module migrations...")
		return app.ModuleManager.RunMigrations(ctx, c, fromVM)
	}
}
