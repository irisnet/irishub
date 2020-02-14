package nft_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/modules/incubator/nft"
	"github.com/irisnet/modules/incubator/nft/internal/types"
)

const (
	module    = "module"
	denom     = "denom"
	nftID     = "nft-id"
	sender    = "sender"
	recipient = "recipient"
	tokenURI  = "token-uri"
)

func TestInvalidMsg(t *testing.T) {
	app, ctx := createTestApp(false)
	h := nft.GenericHandler(app.NFTKeeper)
	res, err := h(ctx, sdk.NewTestMsg())
	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, strings.Contains(err.Error(), "unrecognized nft message type"))
}

func TestTransferNFTMsg(t *testing.T) {
	app, ctx := createTestApp(false)
	h := nft.GenericHandler(app.NFTKeeper)

	// An NFT to be transferred
	nft := types.NewBaseNFT(id, address, "TokenURI")

	// Define MsgTransferNft
	transferNftMsg := types.NewMsgTransferNFT(address, address2, denom, id, tokenURI)

	// handle should fail trying to transfer NFT that doesn't exist
	res, err := h(ctx, transferNftMsg)
	require.Error(t, err)
	require.Nil(t, res)

	// Create token (collection and owner)
	err = app.NFTKeeper.MintNFT(ctx, denom, &nft)
	require.Nil(t, err)
	require.True(t, CheckInvariants(app.NFTKeeper, ctx))

	// handle should succeed when nft exists and is transferred by owner
	res, err = h(ctx, transferNftMsg)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.True(t, CheckInvariants(app.NFTKeeper, ctx))

	// event events should be emitted correctly
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case module:
				require.Equal(t, value, types.ModuleName)
			case denom:
				require.Equal(t, value, denom)
			case nftID:
				require.Equal(t, value, id)
			case sender:
				require.Equal(t, value, address.String())
			case recipient:
				require.Equal(t, value, address2.String())
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}

	// nft should have been transferred as a result of the message
	nftAfterwards, err := app.NFTKeeper.GetNFT(ctx, denom, id)
	require.NoError(t, err)
	require.True(t, nftAfterwards.GetOwner().Equals(address2))

	transferNftMsg = types.NewMsgTransferNFT(address2, address3, denom, id, tokenURI)

	// handle should succeed when nft exists and is transferred by owner
	res, err = h(ctx, transferNftMsg)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.True(t, CheckInvariants(app.NFTKeeper, ctx))

	// Create token (collection and owner)
	err = app.NFTKeeper.MintNFT(ctx, denom2, &nft)
	require.Nil(t, err)
	require.True(t, CheckInvariants(app.NFTKeeper, ctx))

	transferNftMsg = types.NewMsgTransferNFT(address2, address3, denom2, id, tokenURI)

	// handle should succeed when nft exists and is transferred by owner
	res, err = h(ctx, transferNftMsg)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.True(t, CheckInvariants(app.NFTKeeper, ctx))
}

func TestEditNFTMetadataMsg(t *testing.T) {
	app, ctx := createTestApp(false)
	h := nft.GenericHandler(app.NFTKeeper)

	// An NFT to be edited
	nft := types.NewBaseNFT(id, address, tokenURI)

	// Create token (collection and address)
	err := app.NFTKeeper.MintNFT(ctx, denom, &nft)
	require.Nil(t, err)

	// Define MsgTransferNft
	failingEditNFTMetadata := types.NewMsgEditNFTMetadata(address, id, denom2, tokenURI2)

	res, err := h(ctx, failingEditNFTMetadata)
	require.Error(t, err)
	require.Nil(t, res)

	// Define MsgTransferNft
	editNFTMetadata := types.NewMsgEditNFTMetadata(address, id, denom, tokenURI2)

	res, err = h(ctx, editNFTMetadata)
	require.NoError(t, err)
	require.NotNil(t, res)

	// event events should be emitted correctly
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case module:
				require.Equal(t, value, types.ModuleName)
			case denom:
				require.Equal(t, value, denom)
			case nftID:
				require.Equal(t, value, id)
			case sender:
				require.Equal(t, value, address.String())
			case tokenURI:
				require.Equal(t, value, tokenURI2)
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}

	nftAfterwards, err := app.NFTKeeper.GetNFT(ctx, denom, id)
	require.NoError(t, err)
	require.Equal(t, tokenURI2, nftAfterwards.GetTokenURI())
}

func TestMintNFTMsg(t *testing.T) {
	app, ctx := createTestApp(false)
	h := nft.GenericHandler(app.NFTKeeper)

	// Define MsgMintNFT
	mintNFT := types.NewMsgMintNFT(address, address, id, denom, tokenURI)

	// minting a token should succeed
	res, err := h(ctx, mintNFT)
	require.NoError(t, err)
	require.NotNil(t, res)

	// event events should be emitted correctly
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case module:
				require.Equal(t, value, types.ModuleName)
			case denom:
				require.Equal(t, value, denom)
			case nftID:
				require.Equal(t, value, id)
			case sender:
				require.Equal(t, value, address.String())
			case recipient:
				require.Equal(t, value, address.String())
			case tokenURI:
				require.Equal(t, value, tokenURI)
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}

	nftAfterwards, err := app.NFTKeeper.GetNFT(ctx, denom, id)

	require.NoError(t, err)
	require.Equal(t, tokenURI, nftAfterwards.GetTokenURI())

	// minting the same token should fail
	res, err = h(ctx, mintNFT)
	require.Error(t, err)
	require.Nil(t, res)

	require.True(t, CheckInvariants(app.NFTKeeper, ctx))
}

func TestBurnNFTMsg(t *testing.T) {
	app, ctx := createTestApp(false)
	h := nft.GenericHandler(app.NFTKeeper)

	// An NFT to be burned
	nft := types.NewBaseNFT(id, address, tokenURI)

	// Create token (collection and address)
	err := app.NFTKeeper.MintNFT(ctx, denom, &nft)
	require.Nil(t, err)

	exists := app.NFTKeeper.IsNFT(ctx, denom, id)
	require.True(t, exists)

	// burning a non-existent NFT should fail
	failBurnNFT := types.NewMsgBurnNFT(address, id2, denom)
	res, err := h(ctx, failBurnNFT)
	require.Error(t, err)
	require.Nil(t, res)

	// NFT should still exist
	exists = app.NFTKeeper.IsNFT(ctx, denom, id)
	require.True(t, exists)

	// burning the NFt should succeed
	burnNFT := types.NewMsgBurnNFT(address, id, denom)

	res, err = h(ctx, burnNFT)
	require.NoError(t, err)
	require.NotNil(t, res)

	// event events should be emitted correctly
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case module:
				require.Equal(t, value, types.ModuleName)
			case denom:
				require.Equal(t, value, denom)
			case nftID:
				require.Equal(t, value, id)
			case sender:
				require.Equal(t, value, address.String())
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}

	// the NFT should not exist after burn
	exists = app.NFTKeeper.IsNFT(ctx, denom, id)
	require.False(t, exists)

	ownerReturned := app.NFTKeeper.GetOwner(ctx, address)
	require.Equal(t, 0, ownerReturned.Supply())

	require.True(t, CheckInvariants(app.NFTKeeper, ctx))
}
