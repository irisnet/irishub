package auth

import (
	"fmt"
	sdk "github.com/irisnet/irishub/types"
)

// GenesisState - all auth state that must be provided at genesis
type GenesisState struct {
	CollectedFees sdk.Coins `json:"collected_fee"`
	FeeAuth       FeeAuth   `json:"data"`
	Params        Params    `json:"params"`
}

// Create a new genesis state
func NewGenesisState(collectedFees sdk.Coins, feeAuth FeeAuth, params Params) GenesisState {
	return GenesisState{
		CollectedFees: collectedFees,
		FeeAuth:       feeAuth,
		Params:        params,
	}
}

// Return a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		CollectedFees: nil,
		FeeAuth:       InitialFeeAuth(),
		Params:        DefaultParams(),
	}
}

// Init store state from genesis data
func InitGenesis(ctx sdk.Context, keeper FeeKeeper, ak AccountKeeper, data GenesisState) {
	if err := ValidateGenesis(data); err != nil {
		panic(err)
	}

	keeper.setCollectedFees(ctx, data.CollectedFees)
	ak.IncreaseTotalLoosenToken(ctx, data.CollectedFees)

	keeper.SetFeeAuth(ctx, data.FeeAuth)
	keeper.SetParamSet(ctx, data.Params)
	ak.InitTotalSupply(ctx)
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, keeper FeeKeeper, ak AccountKeeper) GenesisState {
	collectedFees := keeper.GetCollectedFees(ctx)
	feeAuth := keeper.GetFeeAuth(ctx)
	params := keeper.GetParamSet(ctx)
	return NewGenesisState(collectedFees, feeAuth, params)
}

func ValidateGenesis(data GenesisState) error {
	err := validateParams(data.Params)
	if err != nil {
		return err
	}
	err = ValidateFee(data.CollectedFees)
	if err != nil {
		return err
	}
	return nil
}

func ValidateFee(collectedFee sdk.Coins) error {
	if collectedFee == nil || collectedFee.Empty() {
		return nil
	}
	if !collectedFee.IsValidIrisAtto() {
		return sdk.ErrInvalidCoins(fmt.Sprintf("invalid collected fees [%s]", collectedFee))
	}
	return nil
}
