package wrapper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	farmkeeper "mods.irisnet.org/modules/farm/keeper"
)

var _ govtypes.GovHooks = farmGovHook{}

type farmGovHook struct {
	gh farmkeeper.GovHook
}

// NewFarmGovHook creates a new farmGovHook instance.
//
// It takes a parameter of type farmkeeper.GovHook and returns a farmGovHook.
func NewFarmGovHook(gh farmkeeper.GovHook) govtypes.GovHooks {
	return farmGovHook{
		gh: gh,
	}
}

// AfterProposalDeposit implements types.GovHooks.
func (f farmGovHook) AfterProposalDeposit(c context.Context, proposalID uint64, depositorAddr types.AccAddress) error {
	ctx := types.UnwrapSDKContext(c)
	f.gh.AfterProposalDeposit(ctx, proposalID, depositorAddr)
	return nil
}

// AfterProposalFailedMinDeposit implements types.GovHooks.
func (f farmGovHook) AfterProposalFailedMinDeposit(c context.Context, proposalID uint64) error {
	ctx := types.UnwrapSDKContext(c)
	f.gh.AfterProposalFailedMinDeposit(ctx, proposalID)
	return nil
}

// AfterProposalSubmission implements types.GovHooks.
func (f farmGovHook) AfterProposalSubmission(c context.Context, proposalID uint64) error {
	ctx := types.UnwrapSDKContext(c)
	f.gh.AfterProposalSubmission(ctx, proposalID)
	return nil
}

// AfterProposalVote implements types.GovHooks.
func (f farmGovHook) AfterProposalVote(c context.Context, proposalID uint64, voterAddr types.AccAddress) error {
	ctx := types.UnwrapSDKContext(c)
	f.gh.AfterProposalVote(ctx, proposalID, voterAddr)
	return nil
}

// AfterProposalVotingPeriodEnded implements types.GovHooks.
func (f farmGovHook) AfterProposalVotingPeriodEnded(c context.Context, proposalID uint64) error {
	ctx := types.UnwrapSDKContext(c)
	f.gh.AfterProposalVotingPeriodEnded(ctx, proposalID)
	return nil
}
