package keeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v3/oracle/internal/types"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/guardian"
	"github.com/irisnet/irishub/store"
	sdk "github.com/irisnet/irishub/types"
)

var (
	_ types.ServiceKeeper = MockServiceKeeper{}

	responses = []string{
		`{"last":100,"high":100,"low":50}`,
		`{"last":100,"high":200,"low":50}`,
		`{"last":100,"high":300,"low":50}`,
		`{"last":100,"high":400,"low":50}`,
	}
)

// create a codec used only for testing
func makeTestCodec() *codec.Codec {
	var cdc = codec.New()
	bank.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	guardian.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

func createTestInput(t *testing.T, amt sdk.Int, nAccs int64) (sdk.Context, Keeper, []auth.Account) {
	keyAcc := protocol.KeyAccount
	keyOracle := sdk.NewKVStoreKey("oracle")

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyOracle, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	cdc := makeTestCodec()
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-oracle", Time: time.Now()}, false, log.NewNopLogger())

	ak := auth.NewAccountKeeper(cdc, keyAcc, auth.ProtoBaseAccount)
	initialCoins := sdk.Coins{
		sdk.NewCoin(sdk.IrisAtto, amt),
	}
	initialCoins = initialCoins.Sort()
	accs := createTestAccs(ctx, int(nAccs), initialCoins, &ak)

	mockServiceKeeper := NewMockServiceKeeper()

	keeper := NewKeeper(cdc, keyOracle, types.DefaultCodespace, mockServiceKeeper)

	return ctx, keeper, accs
}

func createTestAccs(ctx sdk.Context, numAccs int, initialCoins sdk.Coins, ak *auth.AccountKeeper) (accs []auth.Account) {
	for i := 0; i < numAccs; i++ {
		privKey := secp256k1.GenPrivKey()
		pubKey := privKey.PubKey()
		addr := sdk.AccAddress(pubKey.Address())
		acc := auth.NewBaseAccountWithAddress(addr)
		acc.Coins = initialCoins
		acc.PubKey = pubKey
		acc.AccountNumber = uint64(i)
		ak.SetAccount(ctx, &acc)
		accs = append(accs, &acc)
	}

	return
}

type Callback func(ctx sdk.Context, requestContextID []byte, responseOutput []string)

type MockServiceKeeper struct {
	cxtMap      map[string]types.RequestContext
	callbackMap map[string]Callback
}

func NewMockServiceKeeper() MockServiceKeeper {
	cxtMap := make(map[string]types.RequestContext)
	callbackMap := make(map[string]Callback)
	return MockServiceKeeper{
		cxtMap:      cxtMap,
		callbackMap: callbackMap,
	}
}

func (m MockServiceKeeper) RegisterResponseHandler(moduleName string,
	callback func(ctx sdk.Context, requestContextID []byte, responseOutput []string)) error {
	m.callbackMap[moduleName] = callback
	return nil
}

func (m MockServiceKeeper) GetRequestContext(ctx sdk.Context,
	requestContextID []byte) (types.RequestContext, bool) {
	reqCtx, ok := m.cxtMap[string(requestContextID)]
	return reqCtx, ok
}

func (m MockServiceKeeper) CreateRequestContext(ctx sdk.Context,
	serviceName string,
	providers []sdk.AccAddress,
	consumer sdk.AccAddress,
	input string,
	serviceFeeCap sdk.Coins,
	timeout int64,
	repeated bool,
	repeatedFrequency uint64,
	repeatedTotal int64,
	state types.RequestContextState,
	respThreshold uint16,
	moduleName string) ([]byte, sdk.Error) {

	var reqCtxID = "mockRequest"
	reqCtx := types.RequestContext{
		ServiceName:       serviceName,
		Providers:         providers,
		Consumer:          consumer,
		Input:             input,
		ServiceFeeCap:     serviceFeeCap,
		Timeout:           timeout,
		Repeated:          repeated,
		RepeatedFrequency: repeatedFrequency,
		RepeatedTotal:     repeatedTotal,
		BatchCounter:      0,
		State:             state,
		ResponseThreshold: respThreshold,
		ResponseHandler:   moduleName,
	}
	m.cxtMap[reqCtxID] = reqCtx
	return []byte(reqCtxID), nil
}

func (m MockServiceKeeper) UpdateRequestContext(ctx sdk.Context, requestContextID []byte) error {
	return nil
}

func (m MockServiceKeeper) StartRequestContext(ctx sdk.Context, requestContextID []byte) error {
	reqCtx := m.cxtMap[string(requestContextID)]
	for i := int64(1); i <= reqCtx.RepeatedTotal; i++ {
		reqCtx.BatchCounter = uint64(i)
		m.cxtMap[string(requestContextID)] = reqCtx
		ctx = ctx.WithBlockHeader(abci.Header{
			ChainID: ctx.BlockHeader().ChainID,
			Height:  ctx.BlockHeight() + 1,
			Time:    ctx.BlockTime().Add(2 * time.Minute),
		})
		callback := m.callbackMap[reqCtx.ResponseHandler]
		callback(ctx, requestContextID, responses)
	}
	return nil
}

func (m MockServiceKeeper) PauseRequestContext(ctx sdk.Context, requestContextID []byte) error {
	panic("implement me")
}

func (m MockServiceKeeper) KillRequestContext(ctx sdk.Context, requestContextID []byte) error {
	panic("implement me")
}
