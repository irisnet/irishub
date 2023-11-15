package guardian

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/v2/modules/guardian/keeper"
	"github.com/irisnet/irishub/v2/modules/guardian/types"
)

// InitGenesis stores genesis data
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) {
	if err := ValidateGenesis(data); err != nil {
		panic(fmt.Errorf("failed to initialize guardian genesis state: %s", err.Error()))
	}
	// Add supers
	for _, super := range data.Supers {
		keeper.AddSuper(ctx, super)
	}
}

// ExportGenesis outputs genesis data
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	var supers []types.Super
	k.IterateSupers(
		ctx,
		func(super types.Super) bool {
			supers = append(supers, super)
			return false
		},
	)

	return types.NewGenesisState(supers)
}

// ValidateGenesis performs basic validation of supply genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data types.GenesisState) error {
	for _, super := range data.Supers {
		if _, err := sdk.AccAddressFromBech32(super.Address); err != nil {
			return err
		}
		if _, err := sdk.AccAddressFromBech32(super.AddedBy); err != nil {
			return err
		}
	}
	return nil
}
