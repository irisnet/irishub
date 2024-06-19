package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"irismod.io/nft/types"
)

// SaveCollection saves all NFTs and returns an error if there already exists
func (k Keeper) SaveCollection(ctx sdk.Context, collection types.Collection) error {
	for _, nft := range collection.NFTs {
		if err := k.SaveNFT(
			ctx,
			collection.Denom.Id,
			nft.GetID(),
			nft.GetName(),
			nft.GetURI(),
			nft.GetURIHash(),
			nft.GetData(),
			nft.GetOwner(),
		); err != nil {
			return err
		}
	}
	return nil
}

// GetCollections returns all the collections
func (k Keeper) GetCollections(ctx sdk.Context) (cs []types.Collection, err error) {
	for _, class := range k.nk.GetClasses(ctx) {
		nfts, err := k.GetNFTs(ctx, class.Id)
		if err != nil {
			return nil, err
		}

		denom, err := k.GetDenomInfo(ctx, class.Id)
		if err != nil {
			return nil, err
		}

		cs = append(cs, types.NewCollection(*denom, nfts))
	}
	return cs, nil
}

// GetTotalSupply returns the number of NFTs by the specified denom ID
func (k Keeper) GetTotalSupply(ctx sdk.Context, denomID string) uint64 {
	return k.nk.GetTotalSupply(ctx, denomID)
}

// GetBalance returns the amount of NFTs by the specified conditions
func (k Keeper) GetBalance(ctx sdk.Context, id string, owner sdk.AccAddress) (supply uint64) {
	return k.nk.GetBalance(ctx, id, owner)
}
