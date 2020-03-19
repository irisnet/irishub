package keeper

import (
	"testing"
	"time"

	"github.com/irisnet/irishub/app/v3/service/exported"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	cmn "github.com/tendermint/tendermint/libs/common"
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

	mockReqCtxID = []byte("mockRequest")
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
	keyGuardian := protocol.KeyGuardian
	keyOracle := sdk.NewKVStoreKey("oracle")

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyOracle, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyGuardian, sdk.StoreTypeIAVL, db)
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

	gk := guardian.NewKeeper(cdc, keyGuardian, guardian.DefaultCodespace)
	err = gk.AddProfiler(ctx, guardian.Guardian{
		Description: "oracle test",
		AccountType: guardian.Genesis,
		Address:     accs[0].GetAddress(),
		AddedBy:     nil,
	})
	require.Nil(t, err)

	keeper := NewKeeper(cdc, keyOracle, types.DefaultCodespace, gk, mockServiceKeeper)

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

type MockServiceKeeper struct {
	cxtMap      map[string]exported.RequestContext
	callbackMap map[string]exported.ResponseCallback
}

func NewMockServiceKeeper() MockServiceKeeper {
	cxtMap := make(map[string]exported.RequestContext)
	callbackMap := make(map[string]exported.ResponseCallback)
	return MockServiceKeeper{
		cxtMap:      cxtMap,
		callbackMap: callbackMap,
	}
}

func (m MockServiceKeeper) RegisterResponseCallback(moduleName string,
	respCallback exported.ResponseCallback) sdk.Error {
	m.callbackMap[moduleName] = respCallback
	return nil
}

func (m MockServiceKeeper) GetRequestContext(ctx sdk.Context,
	requestContextID cmn.HexBytes) (exported.RequestContext, bool) {
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
	superMode,
	repeated bool,
	repeatedFrequency uint64,
	repeatedTotal int64,
	state exported.RequestContextState,
	respThreshold uint16,
	moduleName string) (cmn.HexBytes, sdk.Tags, sdk.Error) {

	reqCtx := exported.RequestContext{
		ServiceName:       serviceName,
		Providers:         providers,
		Consumer:          consumer,
		Input:             input,
		ServiceFeeCap:     serviceFeeCap,
		Timeout:           timeout,
		SuperMode:         superMode,
		Repeated:          repeated,
		RepeatedFrequency: repeatedFrequency,
		RepeatedTotal:     repeatedTotal,
		BatchCounter:      0,
		State:             state,
		ResponseThreshold: respThreshold,
		ModuleName:        moduleName,
	}
	m.cxtMap[string(mockReqCtxID)] = reqCtx
	return mockReqCtxID, sdk.NewTags(), nil
}

func (m MockServiceKeeper) UpdateRequestContext(ctx sdk.Context,
	requestContextID cmn.HexBytes,
	providers []sdk.AccAddress,
	serviceFeeCap sdk.Coins,
	timeout int64,
	repeatedFreq uint64,
	repeatedTotal int64,
	consumer sdk.AccAddress) sdk.Error {
	return nil
}

func (m MockServiceKeeper) StartRequestContext(ctx sdk.Context, requestContextID cmn.HexBytes, consumer sdk.AccAddress) sdk.Error {
	reqCtx := m.cxtMap[string(requestContextID)]
	callback := m.callbackMap[reqCtx.ModuleName]
	for i := int64(reqCtx.BatchCounter + 1); i <= reqCtx.RepeatedTotal; i++ {
		reqCtx.BatchCounter = uint64(i)
		reqCtx.State = exported.RUNNING
		m.cxtMap[string(requestContextID)] = reqCtx
		ctx = ctx.WithBlockHeader(abci.Header{
			ChainID: ctx.BlockHeader().ChainID,
			Height:  ctx.BlockHeight() + 1,
			Time:    ctx.BlockTime().Add(2 * time.Minute),
		})
		callback(ctx, requestContextID, responses, nil)
	}
	return nil
}

func (m MockServiceKeeper) PauseRequestContext(ctx sdk.Context, requestContextID cmn.HexBytes, consumer sdk.AccAddress) sdk.Error {
	reqCtx := m.cxtMap[string(requestContextID)]
	reqCtx.State = exported.PAUSED
	m.cxtMap[string(requestContextID)] = reqCtx
	return nil
}
