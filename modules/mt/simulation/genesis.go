package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/irisnet/irismod/modules/mt/types"
)

const (
	kitties = "kitties"
	doggos  = "doggos"
)

// RandomizedGenState generates a random GenesisState for mt
func RandomizedGenState(simState *module.SimulationState) {
	collections := types.NewCollections(
		types.NewCollection(
			types.Denom{
				Id:      doggos,
				Name:    doggos,
			},
			types.MTs{},
		),
		types.NewCollection(
			types.Denom{
				Id:      kitties,
				Name:    kitties,
			},
			types.MTs{}),
	)

	mtGenesis := types.NewGenesisState(collections)

	bz, err := json.MarshalIndent(mtGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", types.ModuleName, bz)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(mtGenesis)
}

func RandnMTID(r *rand.Rand, min, max int) string {
	n := simtypes.RandIntBetween(r, min, max)
	id := simtypes.RandStringOfLength(r, n)
	return strings.ToLower(id)
}
