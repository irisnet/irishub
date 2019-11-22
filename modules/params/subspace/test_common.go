package subspace

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/store"
	sdk "github.com/irisnet/irishub/types"
)

// Keys for parameter access
const (
	TestParamStore = "ParamsTest"
)

// Returns components for testing
func DefaultTestComponents(t *testing.T, table TypeTable) (sdk.Context, Subspace, func([]*sdk.KVStoreKey) sdk.CommitID) {
	cdc := codec.New()
	key := sdk.NewKVStoreKey("params")
	tkey := sdk.NewTransientStoreKey("tparams")
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.WithTracer(os.Stdout)
	ms.WithTracingContext(sdk.TraceContext{})
	ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewTMLogger(os.Stdout))
	subspace := NewSubspace(cdc, key, tkey, TestParamStore).WithTypeTable(table)

	return ctx, subspace, ms.Commit
}
