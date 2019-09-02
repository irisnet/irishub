package keeper

import (
	"testing"

	"github.com/irisnet/irishub/app/v1/rand/internal/types"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/store"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	randKey := sdk.NewKVStoreKey("randkey")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(randKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return ms, randKey
}

func TestRequestRandKeeper(t *testing.T) {
	ms, randKey := setupMultiStore()

	cdc := codec.New()
	types.RegisterCodec(cdc)

	keeper := NewKeeper(cdc, randKey, types.DefaultCodespace)

	// define variables
	txBytes := []byte("testtx")
	txHeight := int64(10000)
	blockInterval := uint64(100)
	destHeight := txHeight + int64(blockInterval)
	consumer := sdk.AccAddress([]byte("consumer"))

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

	// request a rand
	_, err := keeper.RequestRand(ctx, consumer, blockInterval)
	require.Nil(t, err)

	// get request id
	reqID := types.GenerateRequestID(types.NewRequest(txHeight, consumer, sdk.SHA256(txBytes)))

	// get the pending request and assert the result is not nil
	store := ctx.KVStore(randKey)
	bz := store.Get(KeyRandRequestQueue(destHeight, reqID))
	require.NotNil(t, bz)

	// decode the request
	var request types.Request
	cdc.MustUnmarshalBinaryLengthPrefixed(bz, &request)
	require.Equal(t, txHeight, request.Height)
	require.Equal(t, consumer, request.Consumer)
	require.Equal(t, sdk.SHA256(txBytes), request.TxHash)

	// get the rand and assert the result is nil
	bz = store.Get(KeyRand(reqID))
	require.Nil(t, bz)
}
