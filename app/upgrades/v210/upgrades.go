package v210

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

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
		// this migration is optional
		// add proposal ids with proposers which are active (deposit or voting period)
		// proposals := make(map[uint64]string)
		// proposals[1] = "cosmos1luyncewxk4lm24k6gqy8y5dxkj0klr4tu0lmnj"
		// v4.AddProposerAddressToProposal(
		// 	ctx,
		// 	sdk.NewKVStoreKey(v4.ModuleName),
		// 	app.AppCodec,
		// 	proposals,
		// )

		// Migrate Tendermint consensus parameters from x/params module to a
		// dedicated x/consensus module.
		baseAppLegacySS, ok := app.ParamsKeeper.GetSubspace(baseapp.Paramspace)
		if !ok {
			panic("failed to get legacy param subspace")
		}
		baseapp.MigrateParams(ctx, baseAppLegacySS, &app.ConsensusParamsKeeper)
		return app.ModuleManager.RunMigrations(ctx, c, fromVM)
	}
}
