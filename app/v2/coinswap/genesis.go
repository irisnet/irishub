package coinswap

import (
	"fmt"
	"github.com/irisnet/irishub/app/v2/coinswap/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

// TODO: ...

// GenesisState - coinswap genesis state
type GenesisState struct {
	Params types.Params `json:"params"`
}

// NewGenesisState is the constructor function for GenesisState
func NewGenesisState(params types.Params) GenesisState {
	return GenesisState{
		Params: params,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() GenesisState {
	return NewGenesisState(types.DefaultParams())
}

// InitGenesis new coinswap genesis
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	//
	if err := ValidateGenesis(data); err != nil {
		panic(fmt.Errorf("panic for ValidateGenesis,%v", err))
	}
	k.SetParams(ctx, data.Params)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return NewGenesisState(keeper.GetParams(ctx))
}

// ValidateGenesis - placeholder function
func ValidateGenesis(data GenesisState) error {
	if err := types.ValidateParams(data.Params); err != nil {
		return err
	}
	return nil
}
