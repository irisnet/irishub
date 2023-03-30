package v200

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	icahosttypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/host/types"

	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/evmos/ethermint/x/feemarket"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	irisevm "github.com/irisnet/irishub/modules/evm"

	"github.com/irisnet/irishub/app/upgrades"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:               "v2.0",
	UpgradeHandlerConstructor: upgradeHandlerConstructor,
	StoreUpgrades: &storetypes.StoreUpgrades{
		Added:   []string{evmtypes.StoreKey, feemarkettypes.StoreKey},
		Deleted: []string{icahosttypes.StoreKey},
	},
}

func upgradeHandlerConstructor(m *module.Manager, c module.Configurator, app upgrades.AppKeepers) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		fromVM[evmtypes.ModuleName] = irisevm.AppModule{}.ConsensusVersion()
		fromVM[feemarkettypes.ModuleName] = feemarket.AppModule{}.ConsensusVersion()

		app.EvmKeeper.SetParams(ctx, evmParams)
		app.FeeMarketKeeper.SetParams(ctx, feemarketParams)

		if err := evmToken.Validate(); err != nil {
			return nil, err
		}

		if err := app.TokenKeeper.AddToken(ctx, evmToken); err != nil {
			return nil, err
		}

		consensusParams := app.ReaderWriter.GetConsensusParams(ctx)
		consensusParams.Block.MaxGas = 20000000
		app.ReaderWriter.StoreConsensusParams(ctx, consensusParams)
		return app.ModuleManager.RunMigrations(ctx, c, fromVM)
	}
}
