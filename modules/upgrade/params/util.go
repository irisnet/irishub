package upgradeparams

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

)

func GetCurrentUpgradeProposalId(ctx sdk.Context) uint64 {
	CurrentUpgradeProposalIdParameter.LoadValue(ctx)
	return  CurrentUpgradeProposalIdParameter.Value
}

func GetProposalAcceptHeight(ctx sdk.Context) int64 {
	ProposalAcceptHeightParameter.LoadValue(ctx)
	return  ProposalAcceptHeightParameter.Value
}

func GetSwitchPeriod(ctx sdk.Context) int64 {
	SwitchPeriodParameter.LoadValue(ctx)
	return  SwitchPeriodParameter.Value
}

func SetCurrentUpgradeProposalId(ctx sdk.Context, i uint64) {
	CurrentUpgradeProposalIdParameter.Value = i
	CurrentUpgradeProposalIdParameter.SaveValue(ctx)
}

func SetProposalAcceptHeight(ctx sdk.Context,i int64 ){
	ProposalAcceptHeightParameter.Value = i
	ProposalAcceptHeightParameter.SaveValue(ctx)
}

func SetSwitchPeriod(ctx sdk.Context,i int64 ){
	SwitchPeriodParameter.Value = i
	SwitchPeriodParameter.SaveValue(ctx)
}
