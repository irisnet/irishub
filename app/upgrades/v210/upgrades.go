package v210

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"

	ibcnfttransfertypes "github.com/bianjieai/nft-transfer/types"

	"github.com/irisnet/irishub/v2/app/upgrades"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:               "v2.1",
	UpgradeHandlerConstructor: upgradeHandlerConstructor,
	StoreUpgrades: &storetypes.StoreUpgrades{
		Added: []string{crisistypes.StoreKey, consensustypes.StoreKey, ibcnfttransfertypes.StoreKey},
	},
}

func upgradeHandlerConstructor(
	m *module.Manager,
	c module.Configurator,
	app upgrades.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// Enable 09-localhost type in allowed clients according to
		// https://github.com/cosmos/ibc-go/blob/v7.3.0/docs/migrations/v7-to-v7_1.md
		params := app.IBCKeeper.ClientKeeper.GetParams(ctx)
		params.AllowedClients = append(params.AllowedClients, exported.Localhost)
		app.IBCKeeper.ClientKeeper.SetParams(ctx, params)

		// Migrate Tendermint consensus parameters from x/params module to a
		// dedicated x/consensus module.
		baseAppLegacySS := app.ParamsKeeper.Subspace(baseapp.Paramspace).
			WithKeyTable(paramstypes.ConsensusParamsKeyTable())
		baseapp.MigrateParams(ctx, baseAppLegacySS, &app.ConsensusParamsKeeper)
		return app.ModuleManager.RunMigrations(ctx, c, fromVM)
	}
}
