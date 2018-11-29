package gov

import (
	"fmt"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/upgrade/params"
)

var _ Proposal = (*SoftwareUpgradeProposal)(nil)

type SoftwareUpgradeProposal struct {
	TextProposal

}

func (sp *SoftwareUpgradeProposal) Execute(ctx sdk.Context, k Keeper) error {
	logger := ctx.Logger().With("module", "x/gov")

	if upgradeparams.GetCurrentUpgradeProposalId(ctx) == 0 {
		upgradeparams.SetCurrentUpgradeProposalId(ctx,sp.ProposalID)
		upgradeparams.SetProposalAcceptHeight(ctx,ctx.BlockHeight())
		logger.Info("Execute SoftwareProposal begin", "info", fmt.Sprintf("current height:%d", ctx.BlockHeight()))

	} else {
		logger.Info("Software Upgrade Switch Period is in process.", "info", fmt.Sprintf("current height:%d", ctx.BlockHeight()))

	}

	return nil
}
