package rand

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/irisnet/irishub/app/v3/rand/internal/keeper"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/store"
	sdk "github.com/irisnet/irishub/types"
)

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	randKey := sdk.NewKVStoreKey("randkey")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(randKey, sdk.StoreTypeIAVL, db)
	_ = ms.LoadLatestVersion()

	return ms, randKey
}

func TestExportRandGenesis(t *testing.T) {
	ms, randKey := setupMultiStore()

	cdc := codec.New()
	RegisterCodec(cdc)

	mockServiceKeeper := keeper.NewMockServiceKeeper()

	keeper := NewKeeper(cdc, randKey, mockServiceKeeper, DefaultCodespace)

	// define variables
	txBytes := []byte("testtx")
	txHeight := int64(10000)
	newBlockHeight := txHeight + 50
	consumer1 := sdk.AccAddress([]byte("consumer1"))
	consumer2 := sdk.AccAddress([]byte("consumer2"))
	blockInterval1 := uint64(100)
	blockInterval2 := uint64(200)

	// build context
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	ctx = ctx.WithBlockHeight(txHeight).WithTxBytes(txBytes)

	// request rands
	_, _ = keeper.RequestRand(ctx, consumer1, blockInterval1, false, sdk.NewCoins())
	_, _ = keeper.RequestRand(ctx, consumer2, blockInterval2, false, sdk.NewCoins())

	// get the pending requests from queue
	storedRequests := make(map[int64][]Request)
	keeper.IterateRandRequestQueue(ctx, func(h int64, r Request) bool {
		storedRequests[h] = append(storedRequests[h], r)
		return false
	})
	require.Equal(t, 2, len(storedRequests))

	// proceed to the new block
	ctx = ctx.WithBlockHeight(newBlockHeight)

	// export requests
	exportedGenesis := ExportGenesis(ctx, keeper)
	exportedRequests := exportedGenesis.PendingRandRequests
	require.Equal(t, 2, len(exportedRequests))

	// assert that exported requests are consistent with the requests in queue
	for height, requests := range exportedRequests {
		h, _ := strconv.ParseInt(height, 10, 64)
		storedHeight := h + newBlockHeight - 1

		require.Equal(t, storedRequests[storedHeight], requests)
	}
}
