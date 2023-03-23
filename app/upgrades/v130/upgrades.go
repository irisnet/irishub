package v130

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	tibcmttypes "github.com/bianjieai/tibc-go/modules/tibc/apps/mt_transfer/types"

	mttypes "github.com/irisnet/irismod/modules/mt/types"

	"github.com/irisnet/irishub/app/upgrades"
	"github.com/irisnet/irishub/app/upgrades/v130/tibc"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:               "v1.3",
	UpgradeHandlerConstructor: upgradeHandlerConstructor,
	StoreUpgrades: &storetypes.StoreUpgrades{
		Added: []string{tibcmttypes.StoreKey, mttypes.StoreKey},
	},
}

func upgradeHandlerConstructor(m *module.Manager, c module.Configurator, app upgrades.AppKeepers) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if err := upgrades.CreateClient(ctx,
			app.AppCodec,
			tibc.ClientData,
			app.TIBCkeeper.ClientKeeper,
		); err != nil {
			return nil, err
		}
		return app.ModuleManager.RunMigrations(ctx, c, fromVM)
	}
}
