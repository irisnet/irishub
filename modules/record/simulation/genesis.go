package simulation

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/irisnet/irismod/modules/record/types"
)

// RandomizedGenState generates a random GenesisState for record
func RandomizedGenState(simState *module.SimulationState) {
	records := make([]types.Record, simState.Rand.Intn(100))

	for i := 0; i < len(records); i++ {
		records[i], _ = genRecord(simState.Rand, simState.Accounts)
	}

	recordGenesis := types.NewGenesisState(records)

	bz, err := json.MarshalIndent(&recordGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", types.ModuleName, bz)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(recordGenesis)
}
