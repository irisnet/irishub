package gov

import (
	"fmt"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/gov/params"
	"github.com/irisnet/irishub/types"
	"time"
	"github.com/irisnet/irishub/modules/params"
)

// GenesisState - all gov state that must be provided at genesis
type GenesisState struct {
	TerminatorPeriod   int64                       `json:"terminator_period"`
	StartingProposalID uint64                      `json:"starting_proposalID"`
	Deposits           []DepositWithMetadata       `json:"deposits"`
	Votes              []VoteWithMetadata          `json:"votes"`
	Proposals          []Proposal                  `json:"proposals"`
	DepositProcedure   govparams.DepositProcedure  `json:"deposit_period"`
	VotingProcedure    govparams.VotingProcedure   `json:"voting_period"`
	TallyingProcedure  govparams.TallyingProcedure `json:"tallying_procedure"`
}

type DepositWithMetadata struct {
	ProposalID uint64  `json:"proposal_id"`
	Deposit    Deposit `json:"deposit"`
}

// VoteWithMetadata (just for genesis)
type VoteWithMetadata struct {
	ProposalID uint64 `json:"proposal_id"`
	Vote       Vote   `json:"vote"`
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

	//k.setDepositProcedure(ctx, data.DepositProcedure)
	////////////////////  iris begin  ///////////////////////////
	params.InitGenesisParameter(&govparams.DepositProcedureParameter, ctx, data.DepositProcedure)
	params.InitGenesisParameter(&govparams.VotingProcedureParameter, ctx, data.VotingProcedure)
	params.InitGenesisParameter(&govparams.TallyingProcedureParameter, ctx, data.TallyingProcedure)
	////////////////////  iris end  /////////////////////////////
	for _, deposit := range data.Deposits {
		k.setDeposit(ctx, deposit.ProposalID, deposit.Deposit.Depositor, deposit.Deposit)
	}
	for _, vote := range data.Votes {
		k.setVote(ctx, vote.ProposalID, vote.Vote.Voter, vote.Vote)
	}
	for _, proposal := range data.Proposals {
		k.SetProposal(ctx, proposal)
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	startingProposalID, _ := k.peekCurrentProposalID(ctx)
	////////////////////  iris begin  ///////////////////////////
	depositProcedure := govparams.GetDepositProcedure(ctx)
	votingProcedure := govparams.GetVotingProcedure(ctx)
	tallyingProcedure := govparams.GetTallyingProcedure(ctx)
	////////////////////  iris end  /////////////////////////////

	var deposits []DepositWithMetadata
	var votes []VoteWithMetadata
	proposals := k.GetProposalsFiltered(ctx, nil, nil, StatusNil, 0)
	for _, proposal := range proposals {
		proposalID := proposal.GetProposalID()
		depositsIterator := k.GetDeposits(ctx, proposalID)
		for ; depositsIterator.Valid(); depositsIterator.Next() {
			var deposit Deposit
			k.cdc.MustUnmarshalBinaryLengthPrefixed(depositsIterator.Value(), &deposit)
			deposits = append(deposits, DepositWithMetadata{proposalID, deposit})
		}
		votesIterator := k.GetVotes(ctx, proposalID)
		for ; votesIterator.Valid(); votesIterator.Next() {
			var vote Vote
			k.cdc.MustUnmarshalBinaryLengthPrefixed(votesIterator.Value(), &vote)
			votes = append(votes, VoteWithMetadata{proposalID, vote})
		}
	}
	return GenesisState{
		StartingProposalID: startingProposalID,
		Deposits:           deposits,
		Votes:              votes,
		Proposals:          proposals,
		DepositProcedure:   depositProcedure,
		VotingProcedure:    votingProcedure,
		TallyingProcedure:  tallyingProcedure,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	Denom := "iris"
	IrisCt := types.NewDefaultCoinType(Denom)
	minDeposit, err := IrisCt.ConvertToMinCoin(fmt.Sprintf("%d%s", 1000, Denom))
	if err != nil {
		panic(err)
	}
	return GenesisState{
		TerminatorPeriod:20000,
		StartingProposalID: 1,
		DepositProcedure: govparams.DepositProcedure{
			MinDeposit:       sdk.Coins{minDeposit},
			MaxDepositPeriod: time.Duration(172800) * time.Second,
		},
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
	Denom := "iris"
	IrisCt := types.NewDefaultCoinType(Denom)
	minDeposit, err := IrisCt.ConvertToMinCoin(fmt.Sprintf("%d%s", 10, Denom))
	if err != nil {
		panic(err)
	}
	return GenesisState{
		TerminatorPeriod:20,
		StartingProposalID: 1,
		DepositProcedure: govparams.DepositProcedure{
			MinDeposit:       sdk.Coins{minDeposit},
			MaxDepositPeriod: time.Duration(60) * time.Second,
		},
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

// get raw genesis raw message for testing
func DefaultGenesisStateForLCDTest() GenesisState {
	Denom := "iris"
	IrisCt := types.NewDefaultCoinType(Denom)
	minDeposit, err := IrisCt.ConvertToMinCoin(fmt.Sprintf("%d%s", 10, Denom))
	if err != nil {
		panic(err)
	}
	return GenesisState{
		TerminatorPeriod:20,
		StartingProposalID: 1,
		DepositProcedure: govparams.DepositProcedure{
			MinDeposit:       sdk.Coins{minDeposit},
			MaxDepositPeriod: time.Duration(172800) * time.Second,
		},
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
