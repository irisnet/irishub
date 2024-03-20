package v300

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"

	"github.com/irisnet/irishub/v3/app/upgrades"
)

// Upgrade defines a struct containing necessary fields that a SoftwareUpgradeProposal
var Upgrade = upgrades.Upgrade{
	UpgradeName:               "v3",
	UpgradeHandlerConstructor: upgradeHandlerConstructor,
	StoreUpgrades: &storetypes.StoreUpgrades{
		Added: []string{icahosttypes.StoreKey},
	},
}

func upgradeHandlerConstructor(
	m *module.Manager,
	c module.Configurator,
	box upgrades.Toolbox,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// initialize ICS27 module
		initICAModule(ctx, m, fromVM)

		// merge liquid staking module
		if err := mergeLSModule(ctx, box); err != nil {
			return nil, err
		}
		return box.ModuleManager.RunMigrations(ctx, c, fromVM)
	}
}

func initICAModule(ctx sdk.Context, m *module.Manager, fromVM module.VersionMap) {
	icaModule := m.Modules[icatypes.ModuleName].(ica.AppModule)
	fromVM[icatypes.ModuleName] = icaModule.ConsensusVersion()
	controllerParams := icacontrollertypes.Params{}
	hostParams := icahosttypes.Params{
		HostEnabled:   true,
		AllowMessages: allowMessages,
	}

	ctx.Logger().Info("start to run ica migrations...")
	icaModule.InitModule(ctx, controllerParams, hostParams)
}

func mergeLSModule(ctx sdk.Context, box upgrades.Toolbox) error {
	ctx.Logger().Info("start to run lsm module migrations...")

	storeKey := box.GetKey(stakingtypes.StoreKey)
	return migrateStore(ctx, storeKey, box.AppCodec, box.StakingKeeper)
}
