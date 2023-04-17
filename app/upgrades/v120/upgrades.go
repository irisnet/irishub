package v120

import (
	sdkmath "cosmossdk.io/math"
	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	ibchost "github.com/cosmos/ibc-go/v5/modules/core/24-host"

	coinswaptypes "github.com/irisnet/irismod/modules/coinswap/types"
	farmtypes "github.com/irisnet/irismod/modules/farm/types"
	"github.com/irisnet/irismod/modules/htlc"
	htlctypes "github.com/irisnet/irismod/modules/htlc/types"
	nftmodule "github.com/irisnet/irismod/modules/nft/module"
	nfttypes "github.com/irisnet/irismod/modules/nft/types"
	"github.com/irisnet/irismod/modules/oracle"
	oracletypes "github.com/irisnet/irismod/modules/oracle/types"
	"github.com/irisnet/irismod/modules/random"
	randomtypes "github.com/irisnet/irismod/modules/random/types"
	"github.com/irisnet/irismod/modules/record"
	recordtypes "github.com/irisnet/irismod/modules/record/types"
	"github.com/irisnet/irismod/modules/service"
	servicetypes "github.com/irisnet/irismod/modules/service/types"
	"github.com/irisnet/irismod/modules/token"
	tokentypes "github.com/irisnet/irismod/modules/token/types"

	tibcnfttypes "github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
	tibcclienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	tibchost "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"

	"github.com/irisnet/irishub/app/upgrades"
	"github.com/irisnet/irishub/app/upgrades/v120/tibc"
	"github.com/irisnet/irishub/modules/guardian"
	guardiantypes "github.com/irisnet/irishub/modules/guardian/types"
	"github.com/irisnet/irishub/modules/mint"
	minttypes "github.com/irisnet/irishub/modules/mint/types"
	"github.com/irisnet/irishub/types"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:               "v1.2",
	UpgradeHandlerConstructor: upgradeHandlerConstructor,
	StoreUpgrades: &store.StoreUpgrades{
		Added: []string{farmtypes.StoreKey, feegrant.StoreKey, tibchost.StoreKey, tibcnfttypes.StoreKey},
	},
}

func upgradeHandlerConstructor(m *module.Manager, c module.Configurator, app upgrades.AppKeepers) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// init farm params
		amount := sdkmath.NewIntWithDecimal(1000, int(types.NativeToken.Scale))
		farmtypes.SetDefaultGenesisState(farmtypes.GenesisState{
			Params: farmtypes.Params{
				PoolCreationFee:     sdk.NewCoin(types.NativeToken.MinUnit, amount),
				MaxRewardCategories: 2,
			}},
		)
		tibcclienttypes.SetDefaultGenesisState(tibcclienttypes.GenesisState{
			NativeChainName: "irishub-mainnet",
		})

		if err := upgrades.CreateClient(ctx,
			app.AppCodec,
			tibc.ClientData,
			app.TIBCkeeper.ClientKeeper,
		); err != nil {
			return nil, err
		}
		fromVM[authtypes.ModuleName] = 1
		fromVM[banktypes.ModuleName] = 1
		fromVM[stakingtypes.ModuleName] = 1
		fromVM[govtypes.ModuleName] = 1
		fromVM[distrtypes.ModuleName] = 1
		fromVM[slashingtypes.ModuleName] = 1
		fromVM[coinswaptypes.ModuleName] = 1
		fromVM[ibchost.ModuleName] = 1
		fromVM[capabilitytypes.ModuleName] = capability.AppModule{}.ConsensusVersion()
		fromVM[genutiltypes.ModuleName] = genutil.AppModule{}.ConsensusVersion()
		fromVM[minttypes.ModuleName] = mint.AppModule{}.ConsensusVersion()
		fromVM[paramstypes.ModuleName] = params.AppModule{}.ConsensusVersion()
		fromVM[crisistypes.ModuleName] = crisis.AppModule{}.ConsensusVersion()
		fromVM[upgradetypes.ModuleName] = crisis.AppModule{}.ConsensusVersion()
		fromVM[evidencetypes.ModuleName] = evidence.AppModule{}.ConsensusVersion()
		fromVM[feegrant.ModuleName] = feegrantmodule.AppModule{}.ConsensusVersion()
		fromVM[guardiantypes.ModuleName] = guardian.AppModule{}.ConsensusVersion()
		fromVM[tokentypes.ModuleName] = token.AppModule{}.ConsensusVersion()
		fromVM[recordtypes.ModuleName] = record.AppModule{}.ConsensusVersion()
		fromVM[nfttypes.ModuleName] = nftmodule.AppModule{}.ConsensusVersion()
		fromVM[htlctypes.ModuleName] = htlc.AppModule{}.ConsensusVersion()
		fromVM[servicetypes.ModuleName] = service.AppModule{}.ConsensusVersion()
		fromVM[oracletypes.ModuleName] = oracle.AppModule{}.ConsensusVersion()
		fromVM[randomtypes.ModuleName] = random.AppModule{}.ConsensusVersion()
		return app.ModuleManager.RunMigrations(ctx, c, fromVM)
	}
}
