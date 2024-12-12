package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	mathsdk "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/irisnet/irishub/v4/modules/mint/types"
)

// Simulation parameter constants
const (
	Inflation = "inflation"
)

// GenInflation randomized Inflation
func GenInflation(r *rand.Rand) mathsdk.LegacyDec {
	return mathsdk.LegacyNewDecWithPrec(int64(r.Intn(99)), 2)
}

// RandomizedGenState generates a random GenesisState for mint
func RandomizedGenState(simState *module.SimulationState) {
	// minter
	var inflation = mathsdk.LegacyZeroDec()
	simState.AppParams.GetOrGenerate(
		Inflation, &inflation, simState.Rand,
		func(r *rand.Rand) { inflation = GenInflation(r) },
	)

	params := types.Params{Inflation: inflation, MintDenom: sdk.DefaultBondDenom}
	mintGenesis := types.NewGenesisState(types.DefaultMinter(), params)

	bz, err := json.MarshalIndent(&mintGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", types.ModuleName, bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(mintGenesis)
}
