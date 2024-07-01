package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"mods.irisnet.org/modules/token/types"
	v1 "mods.irisnet.org/modules/token/types/v1"
)

// Simulation parameter constants
const (
	TokenTaxRate      = "token_tax_rate"
	IssueTokenBaseFee = "issue_token_base_fee"
	MintTokenFeeRatio = "mint_token_fee_ratio"
	EnableErc20       = "enable_erc20"
)

// RandomDec randomized sdk.RandomDec
func RandomDec(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(r.Int63n(9)+1, 1)
}

// RandomInt randomized math.Int
func RandomInt(r *rand.Rand) math.Int {
	return sdk.NewInt(r.Int63())
}

// RandomizedGenState generates a random GenesisState for bank
func RandomizedGenState(simState *module.SimulationState) {
	var tokenTaxRate sdk.Dec
	var issueTokenBaseFee math.Int
	var mintTokenFeeRatio sdk.Dec
	var enableErc20 bool
	var tokens []v1.Token

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
			tokens = append(tokens, v1.GetNativeToken())
		},
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, MintTokenFeeRatio, &mintTokenFeeRatio, simState.Rand,
		func(r *rand.Rand) { mintTokenFeeRatio = sdk.NewDecWithPrec(int64(r.Intn(5)), 1) },
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, EnableErc20, &enableErc20, simState.Rand,
		func(r *rand.Rand) {
			enableErc20 = true
		},
	)

	tokenGenesis := v1.NewGenesisState(
		v1.NewParams(tokenTaxRate, sdk.NewCoin(sdk.DefaultBondDenom, issueTokenBaseFee), mintTokenFeeRatio, enableErc20, ""),
		tokens,
	)

	bz, err := json.MarshalIndent(&tokenGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", types.ModuleName, bz)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&tokenGenesis)
}
