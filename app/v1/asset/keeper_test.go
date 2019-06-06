package asset

import (
	"testing"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/params"
	codec "github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/guardian"
	"github.com/irisnet/irishub/store"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	authKey := sdk.NewKVStoreKey("authkey")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, authKey
}

func TestCreateKeeper(t *testing.T) {
	ms, authKey := setupMultiStore()

	cdc := codec.New()
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	accountKeeper := auth.NewAccountKeeper(cdc, authKey, auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper)
	guardianKeeper := guardian.NewKeeper(cdc, protocol.KeyGuardian, guardian.DefaultCodespace)
	paramsKeeper := params.NewKeeper(cdc, protocol.KeyParams, protocol.TkeyParams)

	createKeeper := NewKeeper(cdc, protocol.KeyAsset, bankKeeper, guardianKeeper, DefaultCodespace, paramsKeeper.Subspace(DefaultParamSpace))

	// define variables
	owner := sdk.AccAddress([]byte("owner"))
	moniker := "moniker"
	identity := "identity"
	details := "details"
	website := "website"

	// construct a test gateway
	gateway := Gateway{
		Owner:    owner,
		Moniker:  moniker,
		Identity: identity,
		Details:  details,
		Website:  website,
	}

	// assert the gateway of the given moniker does not exist at the beginning
	require.False(t, createKeeper.HasGateway(ctx, moniker))

	// create a gateway and asset that the gateway exists now
	createKeeper.SetGateway(ctx, gateway)
	require.True(t, createKeeper.HasGateway(ctx, moniker))

	// asset GetGateway will return the previous gateway
	newGateway, _ := createKeeper.GetGateway(ctx, moniker)
	require.Equal(t, gateway, newGateway)
}
