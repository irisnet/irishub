package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/app/v3/asset/internal/types"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/guardian"
	"github.com/irisnet/irishub/store"
	sdk "github.com/irisnet/irishub/types"
)

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey, *sdk.KVStoreKey, *sdk.KVStoreKey, *sdk.KVStoreKey, *sdk.TransientStoreKey) {
	db := dbm.NewMemDB()

	accountKey := sdk.NewKVStoreKey("accountKey")
	assetKey := sdk.NewKVStoreKey("assetKey")
	guardianKey := sdk.NewKVStoreKey("guardianKey")
	paramskey := sdk.NewKVStoreKey("params")
	paramsTkey := sdk.NewTransientStoreKey("transient_params")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(accountKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(assetKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramskey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramsTkey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return ms, accountKey, assetKey, guardianKey, paramskey, paramsTkey
}

// TestAssetAnteHandler tests the ante handler of asset
func TestAssetAnteHandler(t *testing.T) {
	ms, accountKey, assetKey, guardianKey, paramskey, paramsTkey := setupMultiStore()

	cdc := codec.New()
	types.RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	paramsKeeper := params.NewKeeper(cdc, paramskey, paramsTkey)
	ak := auth.NewAccountKeeper(cdc, accountKey, auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(cdc, ak)
	gk := guardian.NewKeeper(cdc, guardianKey, guardian.DefaultCodespace)
	keeper := NewKeeper(cdc, assetKey, bk, gk, types.DefaultCodespace, paramsKeeper.Subspace(types.DefaultParamSpace))

	//set params
	keeper.SetParamSet(ctx, types.DefaultParams())

	// set test accounts
	addr1 := sdk.AccAddress([]byte("addr1"))
	addr2 := sdk.AccAddress([]byte("addr2"))

	// get asset fees
	nativeTokenIssueFee := keeper.getTokenIssueFee(ctx, "sym")
	nativeTokenMintFee := keeper.getTokenMintFee(ctx, "sym")

	//msg
	msgIssueToken := types.MsgIssueToken{Symbol: "sym", Owner: addr1}
	msgIssueToken2 := types.MsgIssueToken{Symbol: "sym", Owner: addr2}
	msgMintToken := types.MsgMintToken{Symbol: "sym", Owner: addr1}

	//init account balance
	_, _, err := keeper.bk.AddCoins(ctx, addr1, sdk.Coins{nativeTokenIssueFee})
	require.NoError(t, err)

	//single msg
	tx := auth.StdTx{Msgs: []sdk.Msg{msgIssueToken}}
	_, res, abort := NewAnteHandler(keeper)(ctx, tx, false)
	require.Equal(t, false, abort)
	require.Equal(t, true, res.IsOK())

	//multiple msg, but insufficient coins
	tx = auth.StdTx{Msgs: []sdk.Msg{msgIssueToken, msgMintToken}}
	_, res, abort = NewAnteHandler(keeper)(ctx, tx, false)
	require.Equal(t, true, abort)
	require.Equal(t, false, res.IsOK())

	//multiple msg, success
	_, _, err = keeper.bk.AddCoins(ctx, addr1, sdk.Coins{nativeTokenMintFee})
	require.NoError(t, err)
	_, res, abort = NewAnteHandler(keeper)(ctx, tx, false)
	require.Equal(t, false, abort)
	require.Equal(t, true, res.IsOK())

	//multiple msg, but insufficient coins
	tx = auth.StdTx{Msgs: []sdk.Msg{msgIssueToken, msgIssueToken2, msgMintToken}}
	_, res, abort = NewAnteHandler(keeper)(ctx, tx, false)
	require.Equal(t, true, abort)
	require.Equal(t, false, res.IsOK())

	//multiple msg, success
	_, _, err = keeper.bk.AddCoins(ctx, addr2, sdk.Coins{nativeTokenIssueFee})
	require.NoError(t, err)
	tx = auth.StdTx{Msgs: []sdk.Msg{msgIssueToken, msgIssueToken2, msgMintToken}}
	_, res, abort = NewAnteHandler(keeper)(ctx, tx, false)
	require.Equal(t, false, abort)
	require.Equal(t, true, res.IsOK())
}
