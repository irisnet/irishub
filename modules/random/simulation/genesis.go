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

// GenPendingRandomRequests gets the randomized requests
func GenPendingRandomRequests(r *rand.Random) map[string][]types.Request {
	pendingRequestNum := simulation.RandomIntBetween(r, 10, 50)
	pendingRequests := make(map[string][]types.Request)

	for i := 0; i < pendingRequestNum; i++ {
		height := simulation.RandomIntBetween(r, 100, 100000)
		consumerAccounts := simulation.RandomomAccounts(r, 1)
		txHash := []byte(simulation.RandomStringOfLength(r, 32))

		request := types.NewRequest(int64(height), consumerAccounts[0].Address, txHash)
		leftHeight := fmt.Sprintf("%d", simulation.RandomIntBetween(r, 100, 100000))

		pendingRequests[leftHeight] = append(pendingRequests[leftHeight], request)
	}

	return pendingRequests
}

// RandomomizedGenState generates a random GenesisState for rand
func RandomomizedGenState(simState *module.SimulationState) {
	pendingRequests := GenPendingRandomRequests(simState.Random)
	randGenesis := types.NewGenesisState(pendingRequests)

	fmt.Printf("Selected randomly generated rand genesis:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, randGenesis))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(randGenesis)
}
