package simulation

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/irisnet/irishub/modules/guardian/internal/keeper"
)

// Simulation operation weights constants
const ()

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(appParams simulation.AppParams, cdc *codec.Codec, k keeper.Keeper) simulation.WeightedOperations {

	return simulation.WeightedOperations{}
}
