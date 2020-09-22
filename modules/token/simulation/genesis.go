package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/irisnet/irismod/modules/token/types"
)

// Simulation parameter constants
const (
	TokenTaxRate      = "token_tax_rate"
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

	var tokenTaxRate sdk.Dec
	var issueTokenBaseFee sdk.Int
	var mintTokenFeeRatio sdk.Dec
	var tokens []types.Token

	simState.AppParams.GetOrGenerate(
		simState.Cdc, TokenTaxRate, &tokenTaxRate, simState.Rand,
		func(r *rand.Rand) { tokenTaxRate = sdk.NewDecWithPrec(int64(r.Intn(5)), 1) },
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, IssueTokenBaseFee, &issueTokenBaseFee, simState.Rand,
		func(r *rand.Rand) {
			issueTokenBaseFee = sdk.NewInt(int64(10))

			for i := 0; i < 5; i++ {
				tokens = append(tokens, randToken(r, simState.Accounts))
			}
			tokens = append(tokens, types.GetNativeToken())
		},
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, MintTokenFeeRatio, &mintTokenFeeRatio, simState.Rand,
		func(r *rand.Rand) { mintTokenFeeRatio = sdk.NewDecWithPrec(int64(r.Intn(5)), 1) },
	)

	tokenGenesis := types.NewGenesisState(
		types.NewParams(tokenTaxRate, sdk.NewCoin(sdk.DefaultBondDenom, issueTokenBaseFee), mintTokenFeeRatio),
		tokens,
	)

	bz, err := json.MarshalIndent(&tokenGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", types.ModuleName, bz)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&tokenGenesis)
}
