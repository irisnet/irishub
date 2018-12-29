package gov

import (
	sdk "github.com/irisnet/irishub/types"

)

const StartingProposalID = 1

// GenesisState - all gov state that must be provided at genesis
type GenesisState struct {
	Params GovParams `json:"params"` // inflation params
}

func NewGenesisState(systemHaltPeriod int64, params GovParams) GenesisState {
	return GenesisState{
		Params:params,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	err := ValidateGenesis(data)
	if err != nil {
		// TODO: Handle this with #870
		panic(err)
	}

	err = k.setInitialProposalID(ctx, StartingProposalID)
	if err != nil {
		// TODO: Handle this with #870
		panic(err)
	}

	k.SetSystemHaltHeight(ctx, -1)
    k.SetParamSet(ctx,data.Params)
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {

	return GenesisState{
		Params:k.GetParamSet(ctx),
	}
}

func ValidateGenesis(data GenesisState) error {
	err := validateParams(data.Params)
	if err != nil {
		return err
	}
	return nil
}

// get raw genesis raw message for testing
func DefaultGenesisStateForCliTest() GenesisState {

	return GenesisState{
		Params:DefaultParams(),
	}
}

func PrepForZeroHeightGenesis(ctx sdk.Context, k Keeper) {
	proposals := k.GetProposalsFiltered(ctx, nil, nil, StatusDepositPeriod, 0)
	for _, proposal := range proposals {
		proposalID := proposal.GetProposalID()
		k.RefundDeposits(ctx, proposalID)
	}

	proposals = k.GetProposalsFiltered(ctx, nil, nil, StatusVotingPeriod, 0)
	for _, proposal := range proposals {
		proposalID := proposal.GetProposalID()
		k.RefundDeposits(ctx, proposalID)
	}
}
