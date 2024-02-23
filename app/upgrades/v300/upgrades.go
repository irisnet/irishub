package v300

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"

	"github.com/irisnet/irishub/v2/app/upgrades"
)

// Upgrade defines a struct containing necessary fields that a SoftwareUpgradeProposal
var Upgrade = upgrades.Upgrade{
	UpgradeName:               "v3.0",
	UpgradeHandlerConstructor: upgradeHandlerConstructor,
	StoreUpgrades: &storetypes.StoreUpgrades{
		Added: []string{icahosttypes.StoreKey},
	},
}

func upgradeHandlerConstructor(
	m *module.Manager,
	c module.Configurator,
	app upgrades.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// initialize ICS27 module
		initICAModule(ctx,m, fromVM)
		return app.ModuleManager.RunMigrations(ctx, c, fromVM)
	}
}

func initICAModule(ctx sdk.Context,m *module.Manager, fromVM module.VersionMap) {
	icaModule := m.Modules[icatypes.ModuleName].(ica.AppModule)
	fromVM[icatypes.ModuleName] = icaModule.ConsensusVersion()
	controllerParams := icacontrollertypes.Params{}
	hostParams := icahosttypes.Params{
		HostEnabled: true,
		AllowMessages: []string{
			authzMsgExec,
			authzMsgGrant,
			authzMsgRevoke,
			bankMsgSend,
			bankMsgMultiSend,
			distrMsgSetWithdrawAddr,
			distrMsgWithdrawValidatorCommission,
			distrMsgFundCommunityPool,
			distrMsgWithdrawDelegatorReward,
			feegrantMsgGrantAllowance,
			feegrantMsgRevokeAllowance,
			legacyGovMsgVoteWeighted,
			legacyGovMsgSubmitProposal,
			legacyGovMsgDeposit,
			legacyGovMsgVote,
			govMsgVoteWeighted,
			govMsgSubmitProposal,
			govMsgDeposit,
			govMsgVote,
			stakingMsgEditValidator,
			stakingMsgDelegate,
			stakingMsgUndelegate,
			stakingMsgBeginRedelegate,
			stakingMsgCreateValidator,
			vestingMsgCreateVestingAccount,
			ibcMsgTransfer,
			nftMsgIssueDenom,
			nftMsgTransferDenom,
			nftMsgMintNFT,
			nftMsgEditNFT,
			nftMsgTransferNFT,
			nftMsgBurnNFT,
			mtMsgIssueDenom,
			mtMsgTransferDenom,
			mtMsgMintMT,
			mtMsgEditMT,
			mtMsgTransferMT,
			mtMsgBurnMT,
		},
	}

	ctx.Logger().Info("start to init interchainaccount module...")
	icaModule.InitModule(ctx, controllerParams, hostParams)
	ctx.Logger().Info("start to run module migrations...")
}
