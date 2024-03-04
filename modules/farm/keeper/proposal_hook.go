package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/irisnet/irismod/modules/farm/types"
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
func (h GovHook) AfterProposalFailedMinDeposit(ctx sdk.Context, proposalID uint64) error {
	info, has := h.k.GetEscrowInfo(ctx, proposalID)
	if !has {
		return types.ErrEscrowInfoNotFound
	}
	//execute refund logic
	h.k.refundEscrow(ctx, info)
	return nil
}

//AfterProposalVotingPeriodEnded callback when proposal voting is complete
func (h GovHook) AfterProposalVotingPeriodEnded(ctx sdk.Context, proposalID uint64) error{
	info, has := h.k.GetEscrowInfo(ctx, proposalID)
	if !has {
		return types.ErrEscrowInfoNotFound
	}

	proposal, has := h.k.gk.GetProposal(ctx, proposalID)
	if !has {
		return types.ErrInvalidProposal
	}

	//when the proposal is passed, the content of the proposal is executed by the gov module, which is not directly processed here
	if proposal.Status == v1.StatusPassed {
		h.k.deleteEscrowInfo(ctx, proposalID)
		return types.ErrInvalidProposal
	}

	//when the proposal is not passed,execute refund logic
	h.k.refundEscrow(ctx, info)

	return nil
}

// AfterProposalSubmission description of the Go function.
//
// Takes in sdk.Context and uint64 as parameters.
// Returns an error.
func (h GovHook) AfterProposalSubmission(ctx sdk.Context, proposalID uint64) error{
	return nil
}
// AfterProposalDeposit is a function that...
//
// takes in ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress.
// Returns an error.
func (h GovHook) AfterProposalDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress) error {
	return nil
}
// AfterProposalVote is a Go function.
//
// It takes parameters ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress.
// It returns an error.
func (h GovHook) AfterProposalVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) error{
	return nil
}
