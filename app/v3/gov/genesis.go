package gov

import (
	sdk "github.com/irisnet/irishub/types"
)

const StartingProposalID = 1

// GenesisState - all gov state that must be provided at genesis
type GenesisState struct {
	Params         GovParams `json:"params"` // inflation params
	NextProposalID uint64    `json:"next_proposal_id"`
	Deposits       Deposits  `json:"deposits"`
	Votes          Votes     `json:"votes"`
	Proposals      Proposals `json:"proposals"`
}

func NewGenesisState(systemHaltPeriod int64, params GovParams) GenesisState {
	return GenesisState{
		Params: params,
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
	if err := ValidateGenesis(data); err != nil {
		// TODO: Handle this with #870
		panic(err)
	}

	if err := k.setInitialProposalID(ctx, StartingProposalID); err != nil {
		// TODO: Handle this with #870
		panic(err)
	}

	k.SetSystemHaltHeight(ctx, -1)
	k.SetParamSet(ctx, data.Params)
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	proposals := k.GetProposals(ctx)

	var deposits Deposits
	var votes Votes
	for _, proposal := range proposals {
		depositsIterator := k.GetDeposits(ctx, proposal.GetProposalID())
		defer depositsIterator.Close()
		for ; depositsIterator.Valid(); depositsIterator.Next() {
			var deposit Deposit
			k.cdc.MustUnmarshalBinaryLengthPrefixed(depositsIterator.Value(), deposit)
			deposits = append(deposits, deposit)
		}

		vitesIterator := k.GetVotes(ctx, proposal.GetProposalID())
		defer vitesIterator.Close()
		for ; vitesIterator.Valid(); vitesIterator.Next() {
			var vote Vote
			k.cdc.MustUnmarshalBinaryLengthPrefixed(vitesIterator.Value(), vote)
			votes = append(votes, vote)
		}
	}

	nextProposalID, err := k.peekCurrentProposalID(ctx)
	if err != nil {
		panic(err)
	}

	return GenesisState{
		Params:         k.GetParamSet(ctx),
		NextProposalID: nextProposalID,
		Deposits:       deposits,
		Votes:          votes,
		Proposals:      proposals,
	}
}

func ValidateGenesis(data GenesisState) error {
	if err := validateParams(data.Params); err != nil {
		return err
	}
	return nil
}

// get raw genesis raw message for testing
func DefaultGenesisStateForCliTest() GenesisState {
	return GenesisState{
		Params: DefaultParamsForTest(),
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
