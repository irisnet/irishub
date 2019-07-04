package gov

import sdk "github.com/irisnet/irishub/types"

var _ Proposal = (*PlainTextProposal)(nil)

type PlainTextProposal struct {
	BasicProposal
}

func (pp *PlainTextProposal) Validate(ctx sdk.Context, k Keeper) sdk.Error {
	return pp.BasicProposal.Validate(ctx, k)
}

func (pp *PlainTextProposal) Execute(ctx sdk.Context, gk Keeper) sdk.Error {
	logger := ctx.Logger()
	if err := pp.Validate(ctx, gk); err != nil {
		logger.Error("Execute PlainTextProposal failed", "height", ctx.BlockHeight(), "proposalId", pp.ProposalID, "err", err.Error())
		return err
	}
	return nil
}
