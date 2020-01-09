package keeper

import (
	"testing"

	"github.com/irisnet/irishub/app/v3/asset/internal/types"
	"github.com/irisnet/irishub/tests"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

// TestAssetAnteHandler tests the ante handler of asset
func TestAssetAnteHandler(t *testing.T) {
	ms, accountKey, assetKey, paramskey, paramsTkey := tests.SetupMultiStore()

	cdc := codec.New()
	types.RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	paramsKeeper := params.NewKeeper(cdc, paramskey, paramsTkey)
	ak := auth.NewAccountKeeper(cdc, accountKey, auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(cdc, ak)
	keeper := NewKeeper(cdc, assetKey, bk, types.DefaultCodespace, paramsKeeper.Subspace(types.DefaultParamSpace))

	// set test accounts
	addr1 := sdk.AccAddress([]byte("addr1"))
	addr2 := sdk.AccAddress([]byte("addr2"))
	acc1 := ak.NewAccountWithAddress(ctx, addr1)
	acc2 := ak.NewAccountWithAddress(ctx, addr2)

	// get asset fees
	nativeTokenIssueFee := keeper.getTokenIssueFee(ctx, "sym")
	nativeTokenMintFee := keeper.getTokenMintFee(ctx, "sym")

	// construct msgs
	msgIssueNativeToken := types.MsgIssueToken{Source: types.AssetSource(0x00), Symbol: "sym"}
	msgMintNativeToken := types.MsgMintToken{TokenId: "i.sym"}
	msgNonAsset1 := sdk.NewTestMsg(addr1)
	msgNonAsset2 := sdk.NewTestMsg(addr2)

	// construct test txs
	tx1 := auth.StdTx{Msgs: []sdk.Msg{msgIssueNativeToken, msgMintNativeToken}}
	tx2 := auth.StdTx{Msgs: []sdk.Msg{msgIssueNativeToken, msgNonAsset1, msgMintNativeToken}}
	tx3 := auth.StdTx{Msgs: []sdk.Msg{msgNonAsset2, msgIssueNativeToken, msgMintNativeToken}}

	// set signers and construct an ante handler
	newCtx := auth.WithSigners(ctx, []auth.Account{acc1, acc2})
	anteHandler := NewAnteHandler(keeper)

	// assert that the ante handler will return with `abort` set to false
	acc1.SetCoins(acc1.GetCoins().Add(sdk.Coins{nativeTokenMintFee}))
	_, res, abort := anteHandler(newCtx, tx1, false)
	require.Equal(t, false, abort)
	require.Equal(t, true, res.IsOK())

	// assert that the ante handler will return with `abort` set to false
	acc1.SetCoins(sdk.Coins{nativeTokenIssueFee})
	_, res, abort = anteHandler(newCtx, tx2, false)
	require.Equal(t, false, abort)
	require.Equal(t, true, res.IsOK())

	// assert that the ante handler will return with `abort` set to false
	acc1.SetCoins(sdk.Coins{})
	_, res, abort = anteHandler(newCtx, tx3, false)
	require.Equal(t, false, abort)
	require.Equal(t, true, res.IsOK())

	// assert that the ante handler will return with `abort` set to true
	newCtx = auth.WithSigners(ctx, []auth.Account{})
	_, res, abort = anteHandler(newCtx, tx3, false)
	require.Equal(t, true, abort)
	require.Equal(t, false, res.IsOK())
}
