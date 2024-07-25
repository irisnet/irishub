package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	"mods.irisnet.org/modules/coinswap/types"
)

const (
	keyFee                    = "swap_fee"
	keyPoolCreationFee        = "pool_creation_fee"
	keyTaxRate                = "tax_rate"
	keyUnilateralLiquidityFee = "unilateral_liquidity_fee"
)

// RandomizedGenState generates a random GenesisState for coinswap
func RandomizedGenState(simState *module.SimulationState) {
	var (
		fee                    sdk.Dec
		poolCreationFee        sdk.Coin
		taxRate                sdk.Dec
		unilateralLiquidityFee sdk.Dec
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, keyFee, &fee, simState.Rand,
		func(r *rand.Rand) {
			num := simulation.RandIntBetween(simState.Rand, 1, 9)
			fee = sdk.NewDecWithPrec(int64(num), 3)
		},
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, keyPoolCreationFee, &poolCreationFee, simState.Rand,
		func(r *rand.Rand) { poolCreationFee = sdk.NewInt64Coin(sdk.DefaultBondDenom, r.Int63n(100)) },
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, keyTaxRate, &taxRate, simState.Rand,
		func(r *rand.Rand) {
			num := simulation.RandIntBetween(simState.Rand, 1, 5)
			taxRate = sdk.NewDecWithPrec(int64(num), 1)
		},
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, keyUnilateralLiquidityFee, &unilateralLiquidityFee, simState.Rand,
		func(r *rand.Rand) {
			num := simulation.RandIntBetween(simState.Rand, 1, 3)
			unilateralLiquidityFee = sdk.NewDecWithPrec(int64(num), 3)
		},
	)

	params := types.NewParams(fee, taxRate, unilateralLiquidityFee, poolCreationFee)
	genesis := &types.GenesisState{
		Params:        params,
		StandardDenom: sdk.DefaultBondDenom,
		Sequence:      1,
	}

	bz, err := json.MarshalIndent(&genesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", types.ModuleName, bz)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(genesis)
}
