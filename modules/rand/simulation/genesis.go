package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/irisnet/irishub/modules/rand/internal/types"
)

// GenPendingRandRequests gets the randomized requests
func GenPendingRandRequests(r *rand.Rand) map[string][]types.Request {
	pendingRequestNum := simulation.RandIntBetween(r, 10, 50)
	pendingRequests := make(map[string][]types.Request)

	for i := 0; i < pendingRequestNum; i++ {
		height := simulation.RandIntBetween(r, 100, 100000)
		consumerAccounts := simulation.RandomAccounts(r, 1)
		txHash := []byte(simulation.RandStringOfLength(r, 32))

		request := types.NewRequest(int64(height), consumerAccounts[0].Address, txHash)
		leftHeight := fmt.Sprintf("%d", simulation.RandIntBetween(r, 100, 100000))

		pendingRequests[leftHeight] = append(pendingRequests[leftHeight], request)
	}

	return pendingRequests
}

// RandomizedGenState generates a random GenesisState for rand
func RandomizedGenState(simState *module.SimulationState) {
	pendingRequests := GenPendingRandRequests(simState.Rand)
	randGenesis := types.NewGenesisState(pendingRequests)

	fmt.Printf("Selected randomly generated rand genesis:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, randGenesis))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(randGenesis)
}
