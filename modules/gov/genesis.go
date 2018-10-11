package gov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/gov/params"
	"github.com/irisnet/irishub/modules/iparam"
	"github.com/irisnet/irishub/types"
	"fmt"
)

// GenesisState - all staking state that must be provided at genesis
type GenesisState struct {
	StartingProposalID int64                        `json:"starting_proposalID"`
	DepositProcedure   govparams.DepositProcedure   `json:"deposit_period"`
	VotingProcedure    govparams.VotingProcedure    `json:"voting_period"`
	TallyingProcedure  govparams.TallyingProcedure  `json:"tallying_procedure"`
}

func NewGenesisState(startingProposalID int64, dp govparams.DepositProcedure, vp govparams.VotingProcedure, tp govparams.TallyingProcedure) GenesisState {
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
	//k.setDepositProcedure(ctx, data.DepositProcedure)
	iparam.InitGenesisParameter(&govparams.DepositProcedureParameter, ctx, data.DepositProcedure)
	iparam.InitGenesisParameter(&govparams.VotingProcedureParameter, ctx, data.VotingProcedure)
	iparam.InitGenesisParameter(&govparams.TallyingProcedureParameter, ctx, data.TallyingProcedure)

}

// WriteGenesis - output genesis parameters
func WriteGenesis(ctx sdk.Context, k Keeper) GenesisState {
	startingProposalID, _ := k.getNewProposalID(ctx)
	depositProcedure := govparams.GetDepositProcedure(ctx)
	votingProcedure := govparams.GetVotingProcedure(ctx)
	tallyingProcedure := govparams.GetTallyingProcedure(ctx)

	return GenesisState{
		StartingProposalID: startingProposalID,
		DepositProcedure:   depositProcedure,
		VotingProcedure:    votingProcedure,
		TallyingProcedure:  tallyingProcedure,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	Denom  := "iris"
	IrisCt := types.NewDefaultCoinType(Denom)
	minDeposit, err := IrisCt.ConvertToMinCoin(fmt.Sprintf("%d%s", 1000, Denom))
	if err != nil {
		panic(err)
	}
	return GenesisState{
		StartingProposalID: 1,
		DepositProcedure: govparams.DepositProcedure{
			MinDeposit:       sdk.Coins{minDeposit},
			MaxDepositPeriod: 20000,
		},
		VotingProcedure: govparams.VotingProcedure{
			VotingPeriod: 20000,
		},
		TallyingProcedure: govparams.TallyingProcedure{
			Threshold:         sdk.NewRat(1, 2),
			Veto:              sdk.NewRat(1, 3),
			GovernancePenalty: sdk.NewRat(1, 100),
		},
	}
}
