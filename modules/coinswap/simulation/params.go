package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/irisnet/irishub/modules/coinswap/internal/types"
)

const keyServiceFee = "Fee"

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simulation.ParamChange {
	return []simulation.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, keyServiceFee,
			func(r *rand.Rand) string {
				return fmt.Sprintf("%v", GenServicefee(r))
			},
		),
	}
}
