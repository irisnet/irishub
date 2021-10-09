package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/irisnet/irismod/modules/farm/types"
)

const (
	CreatePoolFee      = "create_pool_fee"
	MaxRewardCategoryN = "max_reward_category_n"
)

// RandomizedGenState generates a random GenesisState for farm
func RandomizedGenState(simState *module.SimulationState) {
	var (
		createPoolFee      sdk.Int
		maxRewardCategoryN uint32
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, CreatePoolFee, &createPoolFee, simState.Rand,
		func(r *rand.Rand) { createPoolFee = sdk.NewInt(5000) },
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, MaxRewardCategoryN, &maxRewardCategoryN, simState.Rand,
		func(r *rand.Rand) { maxRewardCategoryN = 2 },
	)

	farmPoolGenesis := types.NewGenesisState(
		types.NewParams(sdk.NewCoin(sdk.DefaultBondDenom, createPoolFee), maxRewardCategoryN),
		nil, nil,
	)

	bz, err := json.MarshalIndent(&farmPoolGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", types.ModuleName, bz)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(farmPoolGenesis)
}
