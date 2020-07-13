package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/irisnet/irishub/modules/random/types"
)

// GenPendingRandomRequests gets the randomized requests
func GenPendingRandomRequests(r *rand.Rand) map[string][]types.Request {
	pendingRequestNum := simtypes.RandIntBetween(r, 10, 50)
	pendingRequests := make(map[string][]types.Request)

	for i := 0; i < pendingRequestNum; i++ {
		height := simtypes.RandIntBetween(r, 100, 100000)
		consumerAccounts := simtypes.RandomAccounts(r, 1)
		txHash := []byte(simtypes.RandStringOfLength(r, 32))

		request := types.NewRequest(int64(height), consumerAccounts[0].Address, txHash)
		leftHeight := fmt.Sprintf("%d", simtypes.RandIntBetween(r, 100, 100000))

		pendingRequests[leftHeight] = append(pendingRequests[leftHeight], request)
	}

	return pendingRequests
}

// RandomomizedGenState generates a random GenesisState for rand
func RandomomizedGenState(simState *module.SimulationState) {
	pendingRequests := GenPendingRandomRequests(simState.Rand)
	randGenesis := types.NewGenesisState(pendingRequests)

	fmt.Printf("Selected randomly generated rand genesis:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, randGenesis))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(randGenesis)
}
