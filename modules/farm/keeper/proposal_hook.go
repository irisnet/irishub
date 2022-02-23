package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

var _ govtypes.GovHooks = GovHook{}

type GovHook struct {
	k Keeper
}

func NewGovHook(k Keeper) GovHook {
	return GovHook{k}
}

//AfterProposalFailedMinDeposit callback when the proposal is deleted due to insufficient collateral
func (h GovHook) AfterProposalFailedMinDeposit(ctx sdk.Context, proposalID uint64) {
	info, has := h.k.getEscrowInfo(ctx, proposalID)
	if !has {
		return
	}
	//execute refund logic
	h.k.refundEscrow(ctx, proposalID, info)
}

//AfterProposalVotingPeriodEnded callback when proposal voting is complete
func (h GovHook) AfterProposalVotingPeriodEnded(ctx sdk.Context, proposalID uint64) {
	info, has := h.k.getEscrowInfo(ctx, proposalID)
	if !has {
		return
	}

	proposal, has := h.k.gk.GetProposal(ctx, proposalID)
	if !has {
		return
	}

	//when the proposal is passed, the content of the proposal is executed by the gov module, which is not directly processed here
	if proposal.Status == govtypes.StatusPassed {
		return
	}

	//when the proposal is not passed,execute refund logic
	h.k.refundEscrow(ctx, proposalID, info)
}

func (h GovHook) AfterProposalSubmission(ctx sdk.Context, proposalID uint64) {}
func (h GovHook) AfterProposalDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress) {
}
func (h GovHook) AfterProposalVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) {}
