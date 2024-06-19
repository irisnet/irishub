package keeper

// DONTCOVER

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"irismod.io/mt/types"
)

// RegisterInvariants registers all supply invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "supply", SupplyInvariant(k))
}

// AllInvariants runs all invariants of the MT module.
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		return SupplyInvariant(k)(ctx)
	}
}

// SupplyInvariant checks that the total amount of MTs on collections matches the total amount owned by addresses
func SupplyInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		err := types.ValidateGenesis(*k.ExportGenesisState(ctx))
		if err != nil {
			return sdk.FormatInvariant(
				types.ModuleName, "supply",
				fmt.Sprintf("MT supply invariants check failed, %s", err.Error()),
			), true
		}

		return sdk.FormatInvariant(
			types.ModuleName, "supply",
			"MT supply invariants check passed",
		), false
	}
}
