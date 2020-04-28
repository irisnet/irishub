package keeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/irisnet/irishub/app/v3/service"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/irisnet/irishub/app/v3/rand/internal/types"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/store"
	sdk "github.com/irisnet/irishub/types"
)

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	randKey := sdk.NewKVStoreKey("randkey")
	serviceKey := sdk.NewKVStoreKey("servicekey")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(randKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(serviceKey, sdk.StoreTypeIAVL, db)
	_ = ms.LoadLatestVersion()

	return ms, randKey, serviceKey
}

func TestRequestRandKeeper(t *testing.T) {
	ms, randKey, serviceKey := setupMultiStore()

	cdc := codec.New()
	types.RegisterCodec(cdc)
	service.RegisterCodec(cdc)

	mockServiceKeeper := NewMockServiceKeeper(serviceKey)
	mockBankKeeper := NewMockBankKeeper()

	keeper := NewKeeper(cdc, randKey, mockBankKeeper, mockServiceKeeper, types.DefaultCodespace)

	// define variables
	txBytes := []byte("testtx")
	txHeight := int64(10000)
	blockInterval := uint64(100)
	destHeight := txHeight + int64(blockInterval)
	consumer := sdk.AccAddress([]byte("consumer"))
	provider := sdk.AccAddress([]byte("provider"))

	// build context
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	ctx = ctx.WithBlockHeight(txHeight).WithTxBytes(txBytes)
	require.Equal(t, txHeight, ctx.BlockHeight())
	require.Equal(t, txBytes, ctx.TxBytes())

	// query rand request queue and assert that the result is empty
	var requests []types.Request
	keeper.IterateRandRequestQueue(ctx, func(h int64, r types.Request) bool {
		requests = append(requests, r)
		return false
	})
	require.True(t, len(requests) == 0)

	mockServiceKeeper.SetServiceBinding(ctx,
		service.NewServiceBinding(
			types.ServiceName, provider, sdk.NewCoins(), "", 0, true, time.Time{}, provider,
		))

	// request a rand
	_, err := keeper.RequestRand(ctx, consumer, blockInterval, false, sdk.NewCoins())
	require.Nil(t, err)

	// get request id
	reqID := types.GenerateRequestID(types.NewRequest(txHeight, consumer, sdk.SHA256(txBytes), false, sdk.NewCoins(), nil))

	// get the pending request and assert the result is not nil
	store := ctx.KVStore(randKey)
	bz := store.Get(KeyRandRequestQueue(destHeight, reqID))
	require.NotNil(t, bz)

	// decode the request
	var request types.Request
	cdc.MustUnmarshalBinaryLengthPrefixed(bz, &request)
	require.Equal(t, txHeight, request.Height)
	require.Equal(t, consumer, request.Consumer)
	require.Equal(t, cmn.HexBytes(sdk.SHA256(txBytes)), request.TxHash)

	// get the rand and assert the result is nil
	bz = store.Get(KeyRand(reqID))
	require.Nil(t, bz)
}
