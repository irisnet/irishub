package gov

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ Proposal = (*SoftwareUpgradeProposal)(nil)

type SoftwareUpgradeProposal struct {
	TextProposal
}

func (sp *SoftwareUpgradeProposal) Execute(ctx sdk.Context, k Keeper) error {
	logger := ctx.Logger().With("module", "x/gov")
	logger.Info("Execute SoftwareProposal begin", "info", fmt.Sprintf("current height:%d", ctx.BlockHeight()))


	bz := k.ps.GovSetter().GetRaw(ctx, "upgrade/proposalId")
	if bz == nil || len(bz) == 0 {
		logger.Error("Execute SoftwareProposal ", "err", "Parameter upgrade/proposalId is not exist")
	} else {
		err := k.ps.GovSetter().Set(ctx, "upgrade/proposalId", sp.ProposalID)
		if err != nil {
			return err
		}
	}

	bz = k.ps.GovSetter().GetRaw(ctx, "upgrade/proposalAcceptHeight")
	if bz == nil || len(bz) == 0 {
		logger.Error("Execute SoftwareProposal ", "err", "Parameter upgrade/proposalAcceptHeight is not exist")
	} else {
		err := k.ps.GovSetter().Set(ctx, "upgrade/proposalAcceptHeight", ctx.BlockHeight())
		if err != nil {
			return err
		}
	}
	return nil
}