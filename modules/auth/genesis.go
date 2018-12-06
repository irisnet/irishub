package auth

import (
	sdk "github.com/irisnet/irishub/types"
)

// GenesisState - all auth state that must be provided at genesis
type GenesisState struct {
	CollectedFees sdk.Coins `json:"collected_fees"` // collected fees
}


type FeeGenesisStateConfig struct {
	FeeTokenNative    string `json:"fee_token_native"`
	GasPriceThreshold int64  `json:"gas_price_threshold"`
}


// Create a new genesis state
func NewGenesisState(collectedFees sdk.Coins) GenesisState {
	return GenesisState{
		CollectedFees: collectedFees,
	}
}

// Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(sdk.Coins{})
}

// Init store state from genesis data
func InitGenesis(ctx sdk.Context, keeper FeeCollectionKeeper, data GenesisState, ps FeeManager, params FeeGenesisStateConfig) {
	keeper.setCollectedFees(ctx, data.CollectedFees)

	ps.paramSpace.Set(ctx, nativeFeeTokenKey, params.FeeTokenNative)
	ps.paramSpace.Set(ctx, nativeGasPriceThresholdKey, sdk.NewInt(params.GasPriceThreshold).String())
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, keeper FeeCollectionKeeper) GenesisState {
	collectedFees := keeper.GetCollectedFees(ctx)
	return NewGenesisState(collectedFees)
}
