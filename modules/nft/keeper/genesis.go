package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/nft/types"
)

// InitGenesis stores the NFT genesis.
func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	if err := types.ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	for _, c := range data.Collections {
		creator, err := sdk.AccAddressFromBech32(c.Denom.Creator)
		if err != nil {
			panic(err)
		}
		if err := k.SaveDenom(ctx,
			c.Denom.Id,
			c.Denom.Name,
			c.Denom.Schema,
			c.Denom.Symbol,
			creator,
			c.Denom.MintRestricted,
			c.Denom.UpdateRestricted,
			c.Denom.Description,
			c.Denom.Uri,
			c.Denom.UriHash,
			c.Denom.Data,
		); err != nil {
			panic(err)
		}

		if err := k.SaveCollection(ctx, c); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	collections, err := k.GetCollections(ctx)
	if err != nil {
		panic(err)
	}
	return types.NewGenesisState(collections)
}
