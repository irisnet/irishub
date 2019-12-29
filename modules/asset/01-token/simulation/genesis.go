package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/irisnet/irishub/modules/asset/01-token/internal/types"
)

// Simulation parameter constants
const (
	AssetTaxRate      = "asset_tax_rate"
	IssueTokenBaseFee = "issue_token_base_fee"
	MintTokenFeeRatio = "mint_token_fee_ratio"
)

// RandomDec randomized sdk.RandomDec
func RandomDec(r *rand.Rand) sdk.Dec {
	return sdk.NewDec(r.Int63())
}

// RandomInt randomized sdk.Int
func RandomInt(r *rand.Rand) sdk.Int {
	return sdk.NewInt(r.Int63())
}

// RandomizedGenState generates a random GenesisState for bank
func RandomizedGenState(simState *module.SimulationState) {

	var assetTaxRate sdk.Dec
	var issueTokenBaseFee sdk.Int
	var mintTokenFeeRatio sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, AssetTaxRate, &assetTaxRate, simState.Rand,
		func(r *rand.Rand) { assetTaxRate = RandomDec(r) },
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, IssueTokenBaseFee, &issueTokenBaseFee, simState.Rand,
		func(r *rand.Rand) { issueTokenBaseFee = RandomInt(r) },
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, MintTokenFeeRatio, &mintTokenFeeRatio, simState.Rand,
		func(r *rand.Rand) { mintTokenFeeRatio = RandomDec(r) },
	)

	assetGenesis := types.NewGenesisState(
		types.NewParams(assetTaxRate, sdk.NewCoin(sdk.DefaultBondDenom, issueTokenBaseFee), mintTokenFeeRatio),
		types.Tokens{},
	)

	fmt.Printf("Selected randomly generated bank parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, assetGenesis))
	simState.GenState[types.SubModuleName] = simState.Cdc.MustMarshalJSON(assetGenesis)
}
