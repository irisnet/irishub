package gov

import (
	sdk "github.com/irisnet/irishub/types"

)

const StartingProposalID = 1

// GenesisState - all gov state that must be provided at genesis
type GenesisState struct {
	SystemHaltPeriod  int64                       `json:"terminator_period"`
	Params GovParams `json:"params"` // inflation params
}

func NewGenesisState(systemHaltPeriod int64, params GovParams) GenesisState {
	return GenesisState{
		SystemHaltPeriod:systemHaltPeriod,
		Params:params,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		SystemHaltPeriod:  20000,
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

	k.SetSystemHaltPeriod(ctx, data.SystemHaltPeriod)
	k.SetSystemHaltHeight(ctx, -1)
    k.SetParamSet(ctx,data.Params)
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {

	return GenesisState{
		SystemHaltPeriod:k.GetSystemHaltHeight(ctx),
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
		SystemHaltPeriod:  20,
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
