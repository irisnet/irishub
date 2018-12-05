package fee

import (
	sdk "github.com/irisnet/irishub/types"
)

var IrisCt = sdk.NewDefaultCoinType("iris")

// GenesisState - all auth state that must be provided at genesis
type GenesisState struct {
	CollectedFees     sdk.Coins `json:"collected_fees"` // collected fees
	FeeTokenNative    string    `json:"fee_token_native"`
	GasPriceThreshold int64     `json:"gas_price_threshold"`
}

// Create a new genesis state
func NewGenesisState(collectedFees sdk.Coins) GenesisState {
	return GenesisState{
		CollectedFees:     collectedFees,
		FeeTokenNative:    IrisCt.MinUnit.Denom,
		GasPriceThreshold: 20000000000, // 20(glue), 20*10^9, 1 glue = 10^9 lue/gas, 1 iris = 10^18 lue
	}
}

// Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(sdk.Coins{})
}

// Init store state from genesis data
func InitGenesis(ctx sdk.Context, keeper FeeCollectionKeeper, ps FeeManager, data GenesisState) {
	keeper.setCollectedFees(ctx, data.CollectedFees)
	ps.paramSpace.Set(ctx, nativeFeeTokenKey, data.FeeTokenNative)
	ps.paramSpace.Set(ctx, nativeGasPriceThresholdKey, sdk.NewInt(data.GasPriceThreshold).String())
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, keeper FeeCollectionKeeper) GenesisState {
	collectedFees := keeper.GetCollectedFees(ctx)
	return NewGenesisState(collectedFees)
}
