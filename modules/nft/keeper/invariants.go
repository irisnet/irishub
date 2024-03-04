package keeper

// DONTCOVER

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/nft/types"
)

// RegisterInvariants registers all supply invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "supply", SupplyInvariant(k))
}

// SupplyInvariant checks that the total amount of NFTs on collections matches the total amount owned by addresses
func SupplyInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		ownersCollectionsSupply := make(map[string]uint64)
		var msg string
		count := 0

		collections, err := k.GetCollections(ctx)
		if err != nil {
			panic(err)
		}

		for _, collection := range collections {
			ownersCollectionsSupply[collection.Denom.Id] = uint64(len(collection.NFTs))
		}

		for denom, supply := range ownersCollectionsSupply {
			if supply != k.GetTotalSupply(ctx, denom) {
				count++
				msg += fmt.Sprintf(
					"total %s NFTs supply invariance:\n"+
						"\ttotal %s NFTs supply: %d\n"+
						"\tsum of %s NFTs by owner: %d\n",
					denom, denom, supply, denom, ownersCollectionsSupply[denom],
				)
			}
		}
		broken := count != 0

		return sdk.FormatInvariant(
			types.ModuleName, "supply",
			fmt.Sprintf("%d NFT supply invariants found\n%s", count, msg),
		), broken
	}
}
