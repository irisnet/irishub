package upgrade

import (
	bam "github.com/irisnet/irishub/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"fmt"
	"github.com/irisnet/irishub/modules/iparam"
	"github.com/irisnet/irishub/modules/upgrade/params"
)

type GenesisState struct {
	SwitchPeriod int64    `json:"switch_period"`
}

// InitGenesis - build the genesis version For first Version
func InitGenesis(ctx sdk.Context, k Keeper, router bam.Router, data GenesisState) {

	RegisterModuleList(router)

	moduleList, found := GetModuleListFromBucket(0)
	fmt.Println(moduleList)
	if !found {
		panic("No module list info found for genesis version")
	}

	genesisVersion := NewVersion(0, 0, 0, moduleList)
	k.AddNewVersion(ctx, genesisVersion)

	iparam.InitGenesisParameter(&upgradeparams.ProposalAcceptHeightParameter, ctx, -1)
	iparam.InitGenesisParameter(&upgradeparams.CurrentUpgradeProposalIdParameter, ctx, -1)
	iparam.InitGenesisParameter(&upgradeparams.SwitchPeriodParameter, ctx, data.SwitchPeriod)
	InitGenesis_commitID(ctx, k)
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		SwitchPeriod: 57600,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	return GenesisState{
		SwitchPeriod: 40,
	}
}
