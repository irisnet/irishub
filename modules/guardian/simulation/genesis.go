package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/irisnet/irishub/modules/guardian/internal/types"
)

// Simulation parameter constants
const (
	ServiceFee = "fee"
)

// GenServicefee randomized Servicefee
func GenServicefee(r *rand.Rand) sdk.Dec {
	return simulation.RandomDecAmount(r, sdk.NewDecWithPrec(1, 0))
}

// RandomizedGenState generates a random GenesisState for bank
func RandomizedGenState(simState *module.SimulationState) {
	var serviceFee sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, ServiceFee, &serviceFee, simState.Rand,
		func(r *rand.Rand) { serviceFee = GenServicefee(r) },
	)

	var (
		profilers types.Profilers
		trustees  types.Trustees
	)

	profilers = make([]types.Guardian, simState.NumBonded)
	trustees = make([]types.Guardian, simState.NumBonded)
	for i := 0; i < int(simState.NumBonded); i++ {
		profilers = append(profilers, types.NewGuardian("", types.Genesis, simState.Accounts[i].Address, simState.Accounts[i].Address))
		trustees = append(trustees, types.NewGuardian("", types.Genesis, simState.Accounts[i].Address, simState.Accounts[i].Address))
	}
	swapGenesis := types.NewGenesisState(profilers, trustees)

	fmt.Printf("Selected randomly generated guardian parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, swapGenesis))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(swapGenesis)
}
