package keeper

import (
	"encoding/json"
	"testing"

	"github.com/irisnet/irishub/tests"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/app/v3/asset/internal/types"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

func TestKeeper_IssueToken(t *testing.T) {
	ms, accountKey, assetKey, paramskey, paramsTkey := tests.SetupMultiStore()

	cdc := codec.New()
	types.RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	pk := params.NewKeeper(cdc, paramskey, paramsTkey)
	ak := auth.NewAccountKeeper(cdc, accountKey, auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(cdc, ak)
	keeper := NewKeeper(cdc, assetKey, bk, types.DefaultCodespace, pk.Subspace(types.DefaultParamSpace))
	addr := sdk.AccAddress([]byte("addr1"))

	acc := ak.NewAccountWithAddress(ctx, addr)

	ft := types.NewFungibleToken(types.NATIVE, "", "btc", "btc", 1, "", "satoshi", sdk.NewIntWithDecimal(1, 0), sdk.NewIntWithDecimal(1, 0), true, acc.GetAddress())
	_, err := keeper.IssueToken(ctx, ft)
	assert.NoError(t, err)

	assert.True(t, keeper.HasToken(ctx, "btc"))

	token, found := keeper.getToken(ctx, "btc")
	assert.True(t, found)

	assert.Equal(t, ft.GetDenom(), token.GetDenom())
	assert.Equal(t, ft.Owner, ft.Owner)

	msgJson, _ := json.Marshal(ft)
	assetJson, _ := json.Marshal(token)
	assert.Equal(t, msgJson, assetJson)
}

//TODO:finish the test
func TestKeeper_EditToken(t *testing.T) {
	ms, accountKey, assetKey, paramskey, paramsTkey := tests.SetupMultiStore()

	cdc := codec.New()
	types.RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	pk := params.NewKeeper(cdc, paramskey, paramsTkey)
	ak := auth.NewAccountKeeper(cdc, accountKey, auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(cdc, ak)
	keeper := NewKeeper(cdc, assetKey, bk, types.DefaultCodespace, pk.Subspace(types.DefaultParamSpace))
	addr := sdk.AccAddress([]byte("addr1"))

	acc := ak.NewAccountWithAddress(ctx, addr)

	ft := types.NewFungibleToken(types.NATIVE, "", "btc", "btc", 1, "", "satoshi", sdk.NewIntWithDecimal(1, 0), sdk.NewIntWithDecimal(21000000, 0), true, acc.GetAddress())

	_, err := keeper.IssueToken(ctx, ft)
	assert.NoError(t, err)

	assert.True(t, keeper.HasToken(ctx, "i.btc"))

	token, found := keeper.getToken(ctx, "i.btc")
	assert.True(t, found)

	assert.Equal(t, ft.GetDenom(), token.GetDenom())
	assert.Equal(t, ft.Owner, token.Owner)

	msgJson, _ := json.Marshal(ft)
	assetJson, _ := json.Marshal(token)
	assert.Equal(t, msgJson, assetJson)

	//TODO:finish the edit token
	mintable := types.False
	msgEditToken := types.NewMsgEditToken("BTC Token", "btc", "btc", "btc", 0, mintable, acc.GetAddress())
	_, err = keeper.EditToken(ctx, msgEditToken)
	assert.NoError(t, err)
}

func TestMintTokenKeeper(t *testing.T) {
	ms, accountKey, assetKey, paramskey, paramsTkey := tests.SetupMultiStore()

	cdc := codec.New()
	types.RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	pk := params.NewKeeper(cdc, paramskey, paramsTkey)
	ak := auth.NewAccountKeeper(cdc, accountKey, auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(cdc, ak)
	keeper := NewKeeper(cdc, assetKey, bk, types.DefaultCodespace, pk.Subspace(types.DefaultParamSpace))
	//keeper.Init(ctx)

	addr := sdk.AccAddress([]byte("addr1"))

	acc := ak.NewAccountWithAddress(ctx, addr)
	amtCoin, _ := sdk.NewIntFromString("1000000000000000000000000000")
	coin := sdk.Coins{sdk.NewCoin("iris-atto", amtCoin)}
	bk.AddCoins(ctx, addr, coin)
	ak.IncreaseTotalLoosenToken(ctx, coin)

	ft := types.NewFungibleToken(types.NATIVE, "", "btc", "btc", 0, "", "satoshi", sdk.NewIntWithDecimal(1000, 0), sdk.NewIntWithDecimal(10000, 0), true, acc.GetAddress())
	_, err := keeper.IssueToken(ctx, ft)
	assert.NoError(t, err)

	assert.True(t, keeper.HasToken(ctx, "btc"))

	token, found := keeper.getToken(ctx, "btc")
	assert.True(t, found)

	assert.Equal(t, ft.GetDenom(), token.GetDenom())
	assert.Equal(t, ft.Owner, ft.Owner)

	msgJson, _ := json.Marshal(ft)
	assetJson, _ := json.Marshal(token)
	assert.Equal(t, msgJson, assetJson)

	msgMintToken := types.NewMsgMintToken("btc", acc.GetAddress(), nil, 1000)
	_, err = keeper.MintToken(ctx, msgMintToken)
	assert.NoError(t, err)

	balance := bk.GetCoins(ctx, acc.GetAddress())
	amt := balance.AmountOf("btc-min")
	assert.Equal(t, "2000", amt.String())
}

func TestTransferOwnerKeeper(t *testing.T) {
	ms, accountKey, assetKey, paramskey, paramsTkey := tests.SetupMultiStore()

	cdc := codec.New()
	types.RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	pk := params.NewKeeper(cdc, paramskey, paramsTkey)
	ak := auth.NewAccountKeeper(cdc, accountKey, auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(cdc, ak)
	keeper := NewKeeper(cdc, assetKey, bk, types.DefaultCodespace, pk.Subspace(types.DefaultParamSpace))

	srcOwner := sdk.AccAddress([]byte("TokenSrcOwner"))

	acc := ak.NewAccountWithAddress(ctx, srcOwner)

	ft := types.NewFungibleToken(types.NATIVE, "", "btc", "btc", 1, "", "satoshi", sdk.NewIntWithDecimal(1, 0), sdk.NewIntWithDecimal(21000000, 0), true, acc.GetAddress())

	_, err := keeper.IssueToken(ctx, ft)
	assert.NoError(t, err)

	assert.True(t, keeper.HasToken(ctx, "i.btc"))

	token, found := keeper.getToken(ctx, "i.btc")
	assert.True(t, found)

	assert.Equal(t, ft.GetDenom(), token.GetDenom())
	assert.Equal(t, ft.Owner, token.Owner)

	msgJson, _ := json.Marshal(ft)
	assetJson, _ := json.Marshal(token)
	assert.Equal(t, msgJson, assetJson)

	dstOwner := sdk.AccAddress([]byte("TokenDstOwner"))
	msg := types.MsgTransferTokenOwner{
		SrcOwner: srcOwner,
		DstOwner: dstOwner,
		TokenId:  "btc",
	}
	_, err = keeper.TransferTokenOwner(ctx, msg)
	assert.NoError(t, err)

	token, found = keeper.getToken(ctx, "i.btc")
	assert.True(t, found)
	assert.Equal(t, dstOwner, token.Owner)
}
