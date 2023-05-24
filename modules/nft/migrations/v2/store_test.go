package v2_test

import (
	"fmt"
	"math/rand"
	"testing"

	gogotypes "github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/nft/keeper"
	v2 "github.com/irisnet/irismod/modules/nft/migrations/v2"
	"github.com/irisnet/irismod/modules/nft/types"
	"github.com/irisnet/irismod/simapp"
)

func TestMigrate(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	storeKey := app.GetKey(types.StoreKey)
	cdc := app.AppCodec()

	collections := prepareData(ctx, storeKey, cdc)
	require.NoError(t, v2.Migrate(ctx, storeKey, cdc, app.NFTKeeper.Logger(ctx), app.NFTKeeper.SaveDenom))
	check(t, ctx, app.NFTKeeper, collections)

}

func prepareData(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.Codec) (collection []types.Collection) {
	addrs := simapp.CreateTestAddrs(10)
	for i := 1; i <= 10; i++ {
		denom := types.Denom{
			Id:               fmt.Sprintf("denom%d", i),
			Name:             fmt.Sprintf("denomName%d", i),
			Schema:           fmt.Sprintf("denomSchema%d", i),
			Creator:          addrs[rand.Intn(len(addrs))].String(),
			Symbol:           fmt.Sprintf("denomSymbol%d", i),
			MintRestricted:   false,
			UpdateRestricted: true,
			Description:      fmt.Sprintf("denomDescription%d", i),
			Uri:              fmt.Sprintf("denomUri%d", i),
			UriHash:          fmt.Sprintf("denomUriHash%d", i),
			Data:             fmt.Sprintf("denomData%d", i),
		}
		setDenom(ctx, storeKey, cdc, denom)

		var tokens []types.BaseNFT
		for j := 1; j <= 100; j++ {
			token := types.BaseNFT{
				Id:      fmt.Sprintf("nft%d", j),
				Name:    fmt.Sprintf("nftName%d", j),
				URI:     fmt.Sprintf("nftURI%d", j),
				Data:    fmt.Sprintf("nftData%d", j),
				Owner:   addrs[rand.Intn(len(addrs))].String(),
				UriHash: fmt.Sprintf("nftUriHash%d", j),
			}
			tokens = append(tokens, token)
			mintNFT(ctx, storeKey, cdc, denom.Id, token)
		}
		collection = append(collection, types.Collection{
			Denom: denom,
			NFTs:  tokens,
		})
	}
	return
}

func check(t *testing.T, ctx sdk.Context, k keeper.Keeper, collections []types.Collection) {
	for _, collection := range collections {
		denom := collection.Denom
		d, err := k.GetDenomInfo(ctx, denom.Id)
		require.NoError(t, err)
		require.EqualValues(t, denom, *d)

		for _, token := range collection.NFTs {
			nft, err := k.GetNFT(ctx, denom.Id, token.Id)
			require.NoError(t, err)
			require.EqualValues(t, token, nft)
		}
	}
	keeper.SupplyInvariant(k)
}

// SetDenom is responsible for saving the definition of denom
func setDenom(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.Codec, denom types.Denom) {
	store := ctx.KVStore(storeKey)
	bz := cdc.MustMarshal(&denom)
	store.Set(v2.KeyDenom(denom.Id), bz)
	store.Set(v2.KeyDenomName(denom.Name), []byte(denom.Id))
}

// MintNFT mints an NFT and manages the NFT's existence within Collections and Owners
func mintNFT(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.Codec, denomID string, baseToken types.BaseNFT) {
	setNFT(ctx, storeKey, cdc, denomID, baseToken)
	setOwner(ctx, storeKey, cdc, denomID, baseToken.Id, baseToken.Owner)
	increaseSupply(ctx, storeKey, cdc, denomID)
}

func setNFT(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.Codec, denomID string, baseToken types.BaseNFT) {
	store := ctx.KVStore(storeKey)

	bz := cdc.MustMarshal(&baseToken)
	store.Set(v2.KeyNFT(denomID, baseToken.Id), bz)
}

func setOwner(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.Codec, denomID, tokenID, owner string) {
	store := ctx.KVStore(storeKey)
	bz := mustMarshalTokenID(cdc, tokenID)
	ownerAddr := sdk.MustAccAddressFromBech32(owner)
	store.Set(v2.KeyOwner(ownerAddr, denomID, tokenID), bz)
}

func increaseSupply(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.Codec, denomID string) {
	supply := getTotalSupply(ctx, storeKey, cdc, denomID)
	supply++

	store := ctx.KVStore(storeKey)
	bz := mustMarshalSupply(cdc, supply)
	store.Set(v2.KeyCollection(denomID), bz)
}

func getTotalSupply(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.Codec, denomID string) uint64 {
	store := ctx.KVStore(storeKey)
	bz := store.Get(v2.KeyCollection(denomID))
	if len(bz) == 0 {
		return 0
	}
	return mustUnMarshalSupply(cdc, bz)
}

func mustMarshalSupply(cdc codec.Codec, supply uint64) []byte {
	supplyWrap := gogotypes.UInt64Value{Value: supply}
	return cdc.MustMarshal(&supplyWrap)
}

func mustUnMarshalSupply(cdc codec.Codec, value []byte) uint64 {
	var supplyWrap gogotypes.UInt64Value
	cdc.MustUnmarshal(value, &supplyWrap)
	return supplyWrap.Value
}

func mustMarshalTokenID(cdc codec.Codec, tokenID string) []byte {
	tokenIDWrap := gogotypes.StringValue{Value: tokenID}
	return cdc.MustMarshal(&tokenIDWrap)
}
