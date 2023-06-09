package simulation

import (
	"encoding/json"
	"math/rand"
	"testing"

	mt "github.com/irisnet/irismod/modules/mt/types"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/simapp"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func TestRandomizedGenState(t *testing.T) {
	app := simapp.Setup(t, false)

	s := rand.NewSource(1)
	r := rand.New(s)

	simState := module.SimulationState{
		AppParams:    make(simtypes.AppParams),
		Cdc:          app.AppCodec(),
		Rand:         r,
		NumBonded:    3,
		Accounts:     simtypes.RandomAccounts(r, 3),
		InitialStake: sdkmath.NewInt(1000),
		GenState:     make(map[string]json.RawMessage),
	}

	RandomizedGenState(&simState)
	var mtGenesis mt.GenesisState
	simState.Cdc.MustUnmarshalJSON(simState.GenState[mt.ModuleName], &mtGenesis)

	require.Len(t, mtGenesis.Collections, len(simState.Accounts))
	require.Len(t, mtGenesis.Owners, len(simState.Accounts))
}
