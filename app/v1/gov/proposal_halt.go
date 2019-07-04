package gov

import sdk "github.com/irisnet/irishub/types"

var _ Proposal = (*SystemHaltProposal)(nil)

type SystemHaltProposal struct {
	BasicProposal
}

func (pp *SystemHaltProposal) Validate(ctx sdk.Context, k Keeper) sdk.Error {
	if err := pp.BasicProposal.Validate(ctx, k); err != nil {
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

	if err := pp.Validate(ctx, gk); err != nil {
		logger.Error("Execute SystemHaltProposal failed", "height", ctx.BlockHeight(), "proposalId", pp.ProposalID, "err", err.Error())
		return err
	}
	if gk.GetSystemHaltHeight(ctx) == -1 {
		gk.SetSystemHaltHeight(ctx, ctx.BlockHeight()+gk.GetSystemHaltPeriod(ctx))
		logger.Info("Execute SystemHaltProposal begin", "SystemHaltHeight", gk.GetSystemHaltHeight(ctx))
	} else {
		logger.Info("SystemHalt Period is in process", "SystemHaltHeight", gk.GetSystemHaltHeight(ctx))

	}
	return nil
}
