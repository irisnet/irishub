package wrapper

import (
	"github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	farmkeeper "github.com/irisnet/irismod/modules/farm/keeper"
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
func (f farmGovHook) AfterProposalDeposit(ctx types.Context, proposalID uint64, depositorAddr types.AccAddress) error {
	f.gh.AfterProposalDeposit(ctx, proposalID, depositorAddr)
	return nil
}

// AfterProposalFailedMinDeposit implements types.GovHooks.
func (f farmGovHook) AfterProposalFailedMinDeposit(ctx types.Context, proposalID uint64) error {
	f.gh.AfterProposalFailedMinDeposit(ctx, proposalID)
	return nil
}

// AfterProposalSubmission implements types.GovHooks.
func (f farmGovHook) AfterProposalSubmission(ctx types.Context, proposalID uint64) error {
	f.gh.AfterProposalSubmission(ctx, proposalID)
	return nil
}

// AfterProposalVote implements types.GovHooks.
func (f farmGovHook) AfterProposalVote(ctx types.Context, proposalID uint64, voterAddr types.AccAddress) error {
	f.gh.AfterProposalVote(ctx, proposalID,voterAddr)
	return nil
}

// AfterProposalVotingPeriodEnded implements types.GovHooks.
func (f farmGovHook) AfterProposalVotingPeriodEnded(ctx types.Context, proposalID uint64) error {
	f.gh.AfterProposalVotingPeriodEnded(ctx, proposalID)
	return nil
}
