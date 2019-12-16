package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/irisnet/modules/incubator/nft/internal/keeper"
	"github.com/irisnet/modules/incubator/nft/internal/types"
)

func TestSetCollection(t *testing.T) {
	app, ctx := createTestApp(false)

	// create a new nft with id = "id" and owner = "address"
	// MintNFT shouldn't fail when collection does not exist
	nft := types.NewBaseNFT(id, address, tokenURI)
	err := app.NFTKeeper.MintNFT(ctx, denom, &nft)
	require.NoError(t, err)

	// collection should exist
	collection, exists := app.NFTKeeper.GetCollection(ctx, denom)
	require.True(t, exists)

	// create a new NFT and add it to the collection created with the NFT mint
	nft2 := types.NewBaseNFT(id2, address, tokenURI)
	collection2, err2 := collection.AddNFT(&nft2)
	require.NoError(t, err2)
	app.NFTKeeper.SetCollection(ctx, denom, collection2)

	collection2, exists = app.NFTKeeper.GetCollection(ctx, denom)
	require.True(t, exists)
	require.Len(t, collection2.NFTs, 2)

	// reset collection for invariant sanity
	app.NFTKeeper.SetCollection(ctx, denom, collection)

	msg, fail := keeper.SupplyInvariant(app.NFTKeeper)(ctx)
	require.False(t, fail, msg)
}
func TestGetCollection(t *testing.T) {
	app, ctx := createTestApp(false)

	// collection shouldn't exist
	collection, exists := app.NFTKeeper.GetCollection(ctx, denom)
	require.Empty(t, collection)
	require.False(t, exists)

	// MintNFT shouldn't fail when collection does not exist
	nft := types.NewBaseNFT(id, address, tokenURI)
	err := app.NFTKeeper.MintNFT(ctx, denom, &nft)
	require.NoError(t, err)

	// collection should exist
	collection, exists = app.NFTKeeper.GetCollection(ctx, denom)
	require.True(t, exists)
	require.NotEmpty(t, collection)

	msg, fail := keeper.SupplyInvariant(app.NFTKeeper)(ctx)
	require.False(t, fail, msg)
}
func TestGetCollections(t *testing.T) {
	app, ctx := createTestApp(false)

	// collections should be empty
	collections := app.NFTKeeper.GetCollections(ctx)
	require.Empty(t, collections)

	// MintNFT shouldn't fail when collection does not exist
	nft := types.NewBaseNFT(id, address, tokenURI)
	err := app.NFTKeeper.MintNFT(ctx, denom, &nft)
	require.NoError(t, err)

	// collections should equal 1
	collections = app.NFTKeeper.GetCollections(ctx)
	require.NotEmpty(t, collections)
	require.Equal(t, len(collections), 1)

	msg, fail := keeper.SupplyInvariant(app.NFTKeeper)(ctx)
	require.False(t, fail, msg)
}
