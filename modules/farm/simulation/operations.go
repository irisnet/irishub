package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/codec"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// Simulation operation weights constants
const (
	OpWeightMsgCreateRecord = "op_weight_msg_create_farm_pool"
)

//TODO
// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	cdc codec.JSONMarshaler) simulation.WeightedOperations {
	var weightCreate int
	appParams.GetOrGenerate(
		cdc, OpWeightMsgCreateRecord, &weightCreate, nil,
		func(_ *rand.Rand) {
			weightCreate = 50
		},
	)
	return simulation.WeightedOperations{}
}
