package simulation

// DONTCOVER

import (
	"github.com/cosmos/cosmos-sdk/types/module"

	token "github.com/irisnet/irishub/modules/asset/01-token"
)

// RandomizedGenState generates a random GenesisState for bank
func RandomizedGenState(simState *module.SimulationState) {
	token.RandomizedGenState(simState)
}
