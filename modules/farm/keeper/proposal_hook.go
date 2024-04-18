package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

var _ govtypes.GovHooks = GovHook{}

// GovHook implements govtypes.GovHooks
type GovHook struct {
	k Keeper
}

// NewGovHook returns a new GovHook instance.
//
// It takes a parameter of type Keeper and returns a GovHook.
func NewGovHook(k Keeper) GovHook {
	return GovHook{k}
}

//AfterProposalFailedMinDeposit callback when the proposal is deleted due to insufficient collateral
func (h GovHook) AfterProposalFailedMinDeposit(ctx sdk.Context, proposalID uint64) {
	info, has := h.k.GetEscrowInfo(ctx, proposalID)
	if !has {
		return
	}
	//execute refund logic
	h.k.refundEscrow(ctx, info)
}

//AfterProposalVotingPeriodEnded callback when proposal voting is complete
func (h GovHook) AfterProposalVotingPeriodEnded(ctx sdk.Context, proposalID uint64) {
	info, has := h.k.GetEscrowInfo(ctx, proposalID)
	if !has {
		return
	}

	proposal, has := h.k.gk.GetProposal(ctx, proposalID)
	if !has {
		return
	}

	//when the proposal is passed, the content of the proposal is executed by the gov module, which is not directly processed here
	if proposal.Status == v1.StatusPassed {
		h.k.deleteEscrowInfo(ctx, proposalID)
		return
	}

	//when the proposal is not passed,execute refund logic
	h.k.refundEscrow(ctx, info)
}

// AfterProposalSubmission description of the Go function.
//
// Takes in sdk.Context and uint64 as parameters.
func (h GovHook) AfterProposalSubmission(ctx sdk.Context, proposalID uint64) {}

// AfterProposalDeposit is a function that...
//
// takes in ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress.
func (h GovHook) AfterProposalDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress) {}

// AfterProposalVote is a Go function.
//
// It takes parameters ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress.
func (h GovHook) AfterProposalVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) {}
