package v300

import (
	"context"
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	ica "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"

	"github.com/irisnet/irishub/v4/app/upgrades"
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
	return func(context context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(context)
		if err := mergeEVM(ctx, box); err != nil {
			return nil, err
		}

		if err := mergeFeeMarket(ctx, box); err != nil {
			return nil, err
		}

		if err := mergeToken(ctx, box); err != nil {
			return nil, err
		}

		if err := mergeGov(ctx, box); err != nil {
			return nil, err
		}
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

func mergeEVM(ctx sdk.Context, box upgrades.Toolbox) error {
	ctx.Logger().Info("start to run evm module migrations...")

	params := box.EvmKeeper.GetParams(ctx)
	params.AllowUnprotectedTxs = true
	return box.EvmKeeper.SetParams(ctx, params)
}

func mergeFeeMarket(ctx sdk.Context, box upgrades.Toolbox) error {
	ctx.Logger().Info("start to run feeMarket module migrations...")

	params := box.FeeMarketKeeper.GetParams(ctx)
	params.MinGasPrice = EvmMinGasPrice
	return box.FeeMarketKeeper.SetParams(ctx, params)
}

func mergeToken(ctx sdk.Context, box upgrades.Toolbox) error {
	ctx.Logger().Info("start to run token module migrations...")

	params := box.TokenKeeper.GetParams(ctx)
	params.EnableErc20 = true
	params.Beacon = BeaconContractAddress
	return box.TokenKeeper.SetParams(ctx, params)
}

func mergeGov(ctx sdk.Context, box upgrades.Toolbox) error {
	ctx.Logger().Info("start to run gov module migrations...")

	params, err := box.GovKeeper.Params.Get(ctx)
	if err != nil {
		return err
	}
	params.MinDepositRatio = MinDepositRatio.String()
	return box.GovKeeper.Params.Set(ctx, params)
}
