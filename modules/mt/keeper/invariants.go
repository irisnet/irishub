package keeper

// DONTCOVER

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/mt/types"
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
		ownersCollectionsSupply := make(map[string]uint64)
		var msg string
		count := 0

		// TODO add invariant check
		//for _, owner := range k.GetOwners(ctx) {
		//	for _, idCollection := range owner.IDCollections {
		//		ownersCollectionsSupply[idCollection.DenomId] += uint64(idCollection.Supply())
		//	}
		//}

		for denom, supply := range ownersCollectionsSupply {
			if supply != k.GetTotalSupply(ctx, denom) {
				count++
				msg += fmt.Sprintf(
					"total %s MTs supply invariance:\n"+
						"\ttotal %s MTs supply: %d\n"+
						"\tsum of %s MTs by owner: %d\n",
					denom, denom, supply, denom, ownersCollectionsSupply[denom],
				)
			}
		}
		broken := count != 0

		return sdk.FormatInvariant(
			types.ModuleName, "supply",
			fmt.Sprintf("%d MT supply invariants found\n%s", count, msg),
		), broken
	}
}
