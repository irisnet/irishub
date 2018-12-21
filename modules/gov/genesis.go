package gov

import (
	"github.com/irisnet/irishub/modules/params"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/types/gov/params"
	govtypes "github.com/irisnet/irishub/types/gov"
	"time"
)

// GenesisState - all gov state that must be provided at genesis
type GenesisState struct {
	TerminatorPeriod   int64                       `json:"terminator_period"`
	StartingProposalID uint64                      `json:"starting_proposalID"`
	DepositProcedure   govparams.DepositProcedure  `json:"deposit_period"`
	VotingProcedure    govparams.VotingProcedure   `json:"voting_period"`
	TallyingProcedure  govparams.TallyingProcedure `json:"tallying_procedure"`
}

func NewGenesisState(startingProposalID uint64, dp govparams.DepositProcedure, vp govparams.VotingProcedure, tp govparams.TallyingProcedure) GenesisState {
	return GenesisState{
		StartingProposalID: startingProposalID,
		DepositProcedure:   dp,
		VotingProcedure:    vp,
		TallyingProcedure:  tp,
	}
}

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {

	err := k.setInitialProposalID(ctx, data.StartingProposalID)
	if err != nil {
		// TODO: Handle this with #870
		panic(err)
	}

	k.SetTerminatorPeriod(ctx, data.TerminatorPeriod)
	k.SetTerminatorHeight(ctx, -1)

	params.InitGenesisParameter(&govparams.DepositProcedureParameter, ctx, data.DepositProcedure)
	params.InitGenesisParameter(&govparams.VotingProcedureParameter, ctx, data.VotingProcedure)
	params.InitGenesisParameter(&govparams.TallyingProcedureParameter, ctx, data.TallyingProcedure)
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	startingProposalID, _ := k.peekCurrentProposalID(ctx)
	depositProcedure := GetDepositProcedure(ctx)
	votingProcedure := GetVotingProcedure(ctx)
	tallyingProcedure := GetTallyingProcedure(ctx)

	return GenesisState{
		StartingProposalID: startingProposalID,
		DepositProcedure:   depositProcedure,
		VotingProcedure:    votingProcedure,
		TallyingProcedure:  tallyingProcedure,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		TerminatorPeriod:   20000,
		StartingProposalID: 1,
		DepositProcedure: govparams.NewDepositProcedure(),
		VotingProcedure: govparams.VotingProcedure{
			VotingPeriod: time.Duration(172800) * time.Second,
		},
		TallyingProcedure: govparams.TallyingProcedure{
			Threshold:     sdk.NewDecWithPrec(5, 1),
			Veto:          sdk.NewDecWithPrec(334, 3),
			Participation: sdk.NewDecWithPrec(667, 3),
		},
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForCliTest() GenesisState {

	depositProcedure := govparams.NewDepositProcedure()
	depositProcedure.MaxDepositPeriod = time.Duration(60) * time.Second
	return GenesisState{
		TerminatorPeriod:   20,
		StartingProposalID: 1,
		DepositProcedure: depositProcedure,
		VotingProcedure: govparams.VotingProcedure{
			VotingPeriod: time.Duration(60) * time.Second,
		},
		TallyingProcedure: govparams.TallyingProcedure{
			Threshold:     sdk.NewDecWithPrec(5, 1),
			Veto:          sdk.NewDecWithPrec(334, 3),
			Participation: sdk.NewDecWithPrec(667, 3),
		},
	}
}

func PrepForZeroHeightGenesis(ctx sdk.Context, k Keeper) {
	proposals := k.GetProposalsFiltered(ctx, nil, nil, govtypes.StatusDepositPeriod, 0)
	for _, proposal := range proposals {
		proposalID := proposal.GetProposalID()
		k.RefundDeposits(ctx, proposalID)
	}

	proposals = k.GetProposalsFiltered(ctx, nil, nil, govtypes.StatusVotingPeriod, 0)
	for _, proposal := range proposals {
		proposalID := proposal.GetProposalID()
		k.RefundDeposits(ctx, proposalID)
	}
}
