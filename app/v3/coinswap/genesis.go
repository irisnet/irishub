package coinswap

import (
	"fmt"

	"github.com/irisnet/irishub/app/v3/coinswap/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

// GenesisState - coinswap genesis state
type GenesisState struct {
	Params types.Params `json:"params"`
	Pools  []types.Pool `json:"pool"`
}

// NewGenesisState is the constructor function for GenesisState
func NewGenesisState(params types.Params, pools []types.Pool) GenesisState {
	return GenesisState{
		Params: params,
		Pools:  pools,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() GenesisState {
	return NewGenesisState(types.DefaultParams(), nil)
}

// InitGenesis new coinswap genesis
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	if err := types.ValidateParams(data.Params); err != nil {
		panic(fmt.Errorf("panic for ValidateGenesis, %v", err))
	}
	k.SetParams(ctx, data.Params)

	for _, pool := range data.Pools {
		if err := CheckVoucherCoinName(pool.Name); err != nil {
			panic(fmt.Errorf("panic for ValidateGenesis, %v", err))
		}
		if pool.Balance().Empty() {
			panic(fmt.Sprintf("empty pool: %s", pool.Name))
		}
		if _, existed := k.GetPool(ctx, pool.Name); existed {
			panic(fmt.Sprintf("pool has existed: %s", pool.Name))
		}
		_ = k.SetPool(ctx, pool)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	params := keeper.GetParams(ctx)
	pools := keeper.GetPools(ctx)
	return NewGenesisState(params, pools)
}
