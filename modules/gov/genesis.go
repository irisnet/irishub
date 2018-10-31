
package gov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/gov/params"
	"fmt"
	"github.com/irisnet/irishub/types"
	"time"
	"github.com/irisnet/irishub/iparam"
)

// GenesisState - all gov state that must be provided at genesis
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
	////////////////////  iris begin  ///////////////////////////
	iparam.InitGenesisParameter(&govparams.DepositProcedureParameter, ctx, data.DepositProcedure)
	iparam.InitGenesisParameter(&govparams.VotingProcedureParameter, ctx, data.VotingProcedure)
	iparam.InitGenesisParameter(&govparams.TallyingProcedureParameter, ctx, data.TallyingProcedure)
	////////////////////  iris end  /////////////////////////////
}

// WriteGenesis - output genesis parameters
func WriteGenesis(ctx sdk.Context, k Keeper) GenesisState {
	startingProposalID, _ := k.getNewProposalID(ctx)

	////////////////////  iris begin  ///////////////////////////
	depositProcedure := govparams.GetDepositProcedure(ctx)
	votingProcedure := govparams.GetVotingProcedure(ctx)
	tallyingProcedure := govparams.GetTallyingProcedure(ctx)
	////////////////////  iris end  /////////////////////////////


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
			MaxDepositPeriod: time.Duration(172800) * time.Second,
		},
		VotingProcedure: govparams.VotingProcedure{
			VotingPeriod: time.Duration(172800) * time.Second,
		},
		TallyingProcedure: govparams.TallyingProcedure{
			Threshold:         sdk.NewDecWithPrec(5, 1),
			Veto:              sdk.NewDecWithPrec(334, 3),
			GovernancePenalty: sdk.NewDecWithPrec(1, 2),
		},
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForCliTest() GenesisState {
	Denom  := "iris"
	IrisCt := types.NewDefaultCoinType(Denom)
	minDeposit, err := IrisCt.ConvertToMinCoin(fmt.Sprintf("%d%s", 10, Denom))
	if err != nil {
		panic(err)
	}
	return GenesisState{
		StartingProposalID: 1,
		DepositProcedure: govparams.DepositProcedure{
			MinDeposit:       sdk.Coins{minDeposit},
			MaxDepositPeriod: time.Duration(172800) * time.Second,
		},
		VotingProcedure: govparams.VotingProcedure{
			VotingPeriod: time.Duration(172800) * time.Second,
		},
		TallyingProcedure: govparams.TallyingProcedure{
			Threshold:         sdk.NewDecWithPrec(5, 1),
			Veto:              sdk.NewDecWithPrec(334, 3),
			GovernancePenalty: sdk.NewDecWithPrec(1, 2),
		},
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForLCDTest() GenesisState {
	Denom  := "iris"
	IrisCt := types.NewDefaultCoinType(Denom)
	minDeposit, err := IrisCt.ConvertToMinCoin(fmt.Sprintf("%d%s", 10, Denom))
	if err != nil {
		panic(err)
	}
	return GenesisState{
		StartingProposalID: 1,
		DepositProcedure: govparams.DepositProcedure{
			MinDeposit:       sdk.Coins{minDeposit},
			MaxDepositPeriod: time.Duration(172800) * time.Second,
		},
		VotingProcedure: govparams.VotingProcedure{
			VotingPeriod: time.Duration(172800) * time.Second,
		},
		TallyingProcedure: govparams.TallyingProcedure{
			Threshold:         sdk.NewDecWithPrec(5, 1),
			Veto:              sdk.NewDecWithPrec(334, 3),
			GovernancePenalty: sdk.NewDecWithPrec(1, 2),
		},
	}
}
