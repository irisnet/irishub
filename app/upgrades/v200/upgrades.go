package v200

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"

	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/evmos/ethermint/x/feemarket"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	"github.com/irisnet/irishub/v2/app/upgrades"
	irisevm "github.com/irisnet/irishub/v2/modules/evm"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:               "v2.0",
	UpgradeHandlerConstructor: upgradeHandlerConstructor,
	StoreUpgrades: &storetypes.StoreUpgrades{
		Added:   []string{evmtypes.StoreKey, feemarkettypes.StoreKey},
		Deleted: []string{icahosttypes.StoreKey},
	},
}

func upgradeHandlerConstructor(
	m *module.Manager,
	c module.Configurator,
	app upgrades.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		fromVM[evmtypes.ModuleName] = irisevm.AppModule{}.ConsensusVersion()
		fromVM[feemarkettypes.ModuleName] = feemarket.AppModule{}.ConsensusVersion()

		if err := app.EvmKeeper.SetParams(ctx, evmParams); err != nil {
			return nil, err
		}

		if err := app.FeeMarketKeeper.SetParams(ctx, generateFeemarketParams(ctx.BlockHeight())); err != nil {
			return nil, err
		}

		//transfer token ownership
		owner, err := sdk.AccAddressFromBech32(evmToken.Owner)
		if err != nil {
			return nil, err
		}
		if err := app.TokenKeeper.UnsafeTransferTokenOwner(ctx, evmToken.Symbol, owner); err != nil {
			return nil, err
		}

		//update consensusParams.Block.MaxGas
		consensusParams := app.ReaderWriter.GetConsensusParams(ctx)
		consensusParams.Block.MaxGas = maxBlockGas
		app.ReaderWriter.StoreConsensusParams(ctx, consensusParams)

		//add Burner Permission for authtypes.FeeCollectorName
		feeModuleAccount := app.AccountKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName)
		account, ok := feeModuleAccount.(*authtypes.ModuleAccount)
		if !ok {
			return nil, fmt.Errorf("feeCollector accountis not *authtypes.ModuleAccount")
		}
		account.Permissions = append(account.Permissions, authtypes.Burner)
		app.AccountKeeper.SetModuleAccount(ctx, account)

		// delete ica moudule version from upgrade moudule
		store := ctx.KVStore(app.GetKey(upgradetypes.StoreKey))
		versionStore := prefix.NewStore(store, []byte{types.VersionMapByte})
		versionStore.Delete([]byte(icatypes.ModuleName))

		return app.ModuleManager.RunMigrations(ctx, c, fromVM)
	}
}
