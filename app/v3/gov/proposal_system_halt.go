package gov

import (
	sdk "github.com/irisnet/irishub/types"
)

var _ Proposal = (*SystemHaltProposal)(nil)

type SystemHaltProposal struct {
	BasicProposal
}

func (pp *SystemHaltProposal) Validate(ctx sdk.Context, k Keeper, verify bool) sdk.Error {
	if err := pp.BasicProposal.Validate(ctx, k, verify); err != nil {
		return err
	}

	_, found := k.guardianKeeper.GetProfiler(ctx, pp.GetProposer())
	if !found {
		return ErrNotProfiler(k.codespace, pp.GetProposer())
	}
	return nil
}

func (pp *SystemHaltProposal) Execute(ctx sdk.Context, gk Keeper) sdk.Error {
	logger := ctx.Logger()

	if err := pp.Validate(ctx, gk, false); err != nil {
		logger.Error("Execute SystemHaltProposal failed", "height", ctx.BlockHeight(), "proposalId", pp.ProposalID, "err", err.Error())
		return err
	}
	var height = gk.GetSystemHaltHeight(ctx)
	if height == -1 {
		haltHeight := gk.GetSystemHaltPeriod(ctx) + ctx.BlockHeight()
		gk.SetSystemHaltHeight(ctx, haltHeight)
		logger.Info("Execute SystemHaltProposal begin", "SystemHaltHeight", haltHeight)
	} else {
		logger.Info("SystemHalt Period is in process", "SystemHaltHeight", height)

	}
	return nil
}
