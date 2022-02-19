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
				Schema:  "",
				Creator: "",
				Symbol:  "dog",
			},
			types.MTs{},
		),
		types.NewCollection(
			types.Denom{
				Id:      kitties,
				Name:    kitties,
				Schema:  "",
				Creator: "",
				Symbol:  "kit",
			},
			types.MTs{}),
	)
	for _, acc := range simState.Accounts {
		// 10% of accounts own an MT
		if simState.Rand.Intn(100) < 10 {
			baseMT := types.NewMT(
				RandnMTID(simState.Rand, types.MinDenomLen, types.MaxDenomLen), // id
				simtypes.RandStringOfLength(simState.Rand, 10),
				acc.Address,
				simtypes.RandStringOfLength(simState.Rand, 45), // tokenURI
				simtypes.RandStringOfLength(simState.Rand, 32), // tokenURIHash
				simtypes.RandStringOfLength(simState.Rand, 10),
			)

			// 50% doggos and 50% kitties
			if simState.Rand.Intn(100) < 50 {
				collections[0].Denom.Creator = baseMT.Owner
				collections[0] = collections[0].AddMT(baseMT)
			} else {
				collections[1].Denom.Creator = baseMT.Owner
				collections[1] = collections[1].AddMT(baseMT)
			}
		}
	}

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
