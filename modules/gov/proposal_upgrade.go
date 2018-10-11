package gov

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/upgrade/params"
)

var _ Proposal = (*SoftwareUpgradeProposal)(nil)

type SoftwareUpgradeProposal struct {
	TextProposal
}

func (sp *SoftwareUpgradeProposal) Execute(ctx sdk.Context, k Keeper) error {
	logger := ctx.Logger().With("module", "x/gov")
	logger.Info("Execute SoftwareProposal begin", "info", fmt.Sprintf("current height:%d", ctx.BlockHeight()))

    upgradeparams.SetCurrentUpgradeProposalId(ctx,sp.ProposalID)
	upgradeparams.SetProposalAcceptHeight(ctx,ctx.BlockHeight())

	return nil
}
