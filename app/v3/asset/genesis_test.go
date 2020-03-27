package asset

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/params"
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
	_ = ms.LoadLatestVersion()

	return ms, accountKey, assetKey, guardianKey, paramskey, paramsTkey
}

func TestExportGenesis(t *testing.T) {
	ms, accountKey, assetKey, guardianKey, paramskey, paramsTkey := setupMultiStore()

	cdc := codec.New()
	RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	paramsKeeper := params.NewKeeper(cdc, paramskey, paramsTkey)
	ak := auth.NewAccountKeeper(cdc, accountKey, auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(cdc, ak)
	gk := guardian.NewKeeper(cdc, guardianKey, guardian.DefaultCodespace)
	keeper := NewKeeper(cdc, assetKey, bk, gk, DefaultCodespace, paramsKeeper.Subspace(DefaultParamSpace))

	// add token
	addr := sdk.AccAddress([]byte("addr1"))
	acc := ak.NewAccountWithAddress(ctx, addr)
	ft := NewFungibleToken("btc", "Bitcoin Network", "satoshi", 1, 1, 1, true, acc.GetAddress())

	genesis := GenesisState{
		Params: DefaultParams(),
		Tokens: Tokens{ft},
	}

	// initialize genesis
	InitGenesis(ctx, keeper, genesis)

	// query all tokens
	var tokens Tokens
	keeper.IterateTokens(ctx, func(token FungibleToken) (stop bool) {
		tokens = append(tokens, token)
		return false
	})

	require.Equal(t, len(tokens), 1)

	// export genesis
	genesisState := ExportGenesis(ctx, keeper)

	require.Equal(t, DefaultParams(), genesisState.Params)
	for _, token := range genesisState.Tokens {
		require.Equal(t, token, ft)
	}
}
