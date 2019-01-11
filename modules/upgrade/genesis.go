package upgrade

import (
	"github.com/irisnet/irishub/modules/params"
	"github.com/irisnet/irishub/modules/upgrade/params"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/version"
)

// GenesisState - all upgrade state that must be provided at genesis
type GenesisState struct {
	GenesisVersion VersionInfo          `json:genesis_version`
	UpgradeParams  upgradeparams.Params `json:upgrade_govparams`
}

// InitGenesis - build the genesis version For first Version
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	genesisVersion := data.GenesisVersion

	k.AddNewVersionInfo(ctx, genesisVersion)
	k.protocolKeeper.ClearUpgradeConfig(ctx)
	k.protocolKeeper.SetCurrentVersion(ctx, genesisVersion.UpgradeInfo.Protocol.Version)
	params.InitGenesisParameter(&upgradeparams.UpgradeParameter, ctx, data.UpgradeParams)
}

// WriteGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context) GenesisState {
	return GenesisState{
		NewVersionInfo(sdk.NewUpgradeConfig(0, sdk.NewProtocolDefinition(uint64(0), "", uint64(1))), true),
		upgradeparams.Params{},
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		NewVersionInfo(sdk.NewUpgradeConfig(0, sdk.NewProtocolDefinition(uint64(0), "https://github.com/irisnet/irishub/releases/tag/v"+version.Version, uint64(1))), true),
		upgradeparams.NewUpgradeParams(),
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	return GenesisState{
		NewVersionInfo(sdk.NewUpgradeConfig(0, sdk.NewProtocolDefinition(uint64(0), "https://github.com/irisnet/irishub/releases/tag/v"+version.Version, uint64(1))), true),
		upgradeparams.NewUpgradeParams(),
	}
}
