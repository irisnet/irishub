package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irismod/modules/oracle/keeper"
	"github.com/irisnet/irismod/modules/oracle/types"
	"github.com/irisnet/irismod/modules/service/exported"
	servicetypes "github.com/irisnet/irismod/modules/service/types"
	"github.com/irisnet/irismod/simapp"
)

var (
	testAddr1, _ = sdk.AccAddressFromHex(crypto.AddressHash([]byte("test1")).String())
	testAddr2, _ = sdk.AccAddressFromHex(crypto.AddressHash([]byte("test2")).String())

	addrs = []sdk.AccAddress{testAddr1, testAddr2}

	mockReqCtxID = []byte("mockRequest")
	responses    = []string{
		`{"header":{},"body":{"last":100,"high":100,"low":50}}`,
		`{"header":{},"body":{"last":100,"high":200,"low":50}}`,
		`{"header":{},"body":{"last":100,"high":300,"low":50}}`,
		`{"header":{},"body":{"last":100,"high":400,"low":50}}`,
	}
)

type KeeperTestSuite struct {
	suite.Suite

	cdc    *codec.LegacyAmino
	ctx    sdk.Context
	app    *simapp.SimApp
	keeper keeper.Keeper
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(false)

	suite.cdc = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(false, tmproto.Header{})
	suite.app = app

	serviceKeeper := NewMockServiceKeeper()
	suite.keeper = keeper.NewKeeper(app.AppCodec(), app.GetKey(types.StoreKey), app.GetSubspace(types.ModuleName), serviceKeeper)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestFeed() {
	msg := &types.MsgCreateFeed{
		FeedName:          "ethPrice",
		ServiceName:       "GetEthPrice",
		AggregateFunc:     "avg",
		ValueJsonPath:     "high",
		LatestHistory:     5,
		Providers:         []sdk.AccAddress{addrs[1]},
		Input:             "xxxx",
		Timeout:           10,
		ServiceFeeCap:     sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
		RepeatedFrequency: 1,
		ResponseThreshold: 1,
		Creator:           addrs[0],
	}

	// ================test CreateFeed start================
	err := suite.keeper.CreateFeed(suite.ctx, msg)
	suite.NoError(err)

	// check feed existed
	feed, existed := suite.keeper.GetFeed(suite.ctx, msg.FeedName)
	suite.True(existed)
	suite.EqualValues(
		types.Feed{
			FeedName:         msg.FeedName,
			AggregateFunc:    msg.AggregateFunc,
			ValueJsonPath:    msg.ValueJsonPath,
			LatestHistory:    msg.LatestHistory,
			RequestContextID: mockReqCtxID,
			Creator:          msg.Creator,
		},
		feed,
	)

	// check feed state
	var feeds []types.Feed
	suite.keeper.IteratorFeedsByState(suite.ctx, exported.PAUSED, func(feed types.Feed) {
		feeds = append(feeds, feed)
	})
	suite.Len(feeds, 1)
	suite.Equal(msg.FeedName, feeds[0].FeedName)
	// ================test CreateFeed end================

	// ================test StartFeed start================
	err = suite.keeper.StartFeed(suite.ctx, &types.MsgStartFeed{
		FeedName: msg.FeedName,
		Creator:  addrs[0],
	})

	// check feed result
	result := suite.keeper.GetFeedValues(suite.ctx, msg.FeedName)
	suite.NoError(err)
	suite.Equal("250.00000000", result[0].Data)

	// check feed state
	var feeds1 []types.Feed
	suite.keeper.IteratorFeedsByState(suite.ctx, exported.RUNNING, func(feed types.Feed) {
		feeds1 = append(feeds1, feed)
	})
	suite.Len(feeds1, 1)
	suite.Equal(msg.FeedName, feeds1[0].FeedName)

	// start again, will return error
	err = suite.keeper.StartFeed(suite.ctx, &types.MsgStartFeed{
		FeedName: msg.FeedName,
		Creator:  addrs[0],
	})
	suite.Error(err)
	// ================test StartFeed end================

	// ================test EditFeed start================
	latestHistory := uint64(1)
	err = suite.keeper.EditFeed(suite.ctx, &types.MsgEditFeed{
		FeedName:          msg.FeedName,
		LatestHistory:     latestHistory,
		Providers:         []sdk.AccAddress{addrs[0]},
		ServiceFeeCap:     sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
		RepeatedFrequency: 1,
		ResponseThreshold: 1,
		Creator:           addrs[0],
	})
	suite.NoError(err)

	//check feed existed
	feed, existed = suite.keeper.GetFeed(suite.ctx, msg.FeedName)
	suite.True(existed)
	suite.EqualValues(
		types.Feed{
			FeedName:         msg.FeedName,
			AggregateFunc:    msg.AggregateFunc,
			ValueJsonPath:    msg.ValueJsonPath,
			LatestHistory:    latestHistory,
			RequestContextID: feed.RequestContextID,
			Creator:          msg.Creator,
		},
		feed,
	)
	// ================test EditFeed end================

	// ================test PauseFeed start================
	err = suite.keeper.PauseFeed(suite.ctx, &types.MsgPauseFeed{
		FeedName: msg.FeedName,
		Creator:  addrs[0],
	})
	suite.NoError(err)

	reqCtx, existed := suite.keeper.GetRequestContext(suite.ctx, feed.RequestContextID)
	suite.True(existed)
	suite.Equal(exported.PAUSED, reqCtx.State)

	// pause again, will return error
	err = suite.keeper.PauseFeed(suite.ctx, &types.MsgPauseFeed{
		FeedName: msg.FeedName,
		Creator:  addrs[0],
	})
	suite.Error(err)

	// Start Feed again
	err = suite.keeper.StartFeed(suite.ctx, &types.MsgStartFeed{
		FeedName: msg.FeedName,
		Creator:  addrs[0],
	})

	// check feed result
	result = suite.keeper.GetFeedValues(suite.ctx, msg.FeedName)
	suite.NoError(err)
	suite.Len(result, int(latestHistory))
	suite.Equal("250.00000000", result[0].Data)

	// check feed state
	var feeds2 []types.Feed
	suite.keeper.IteratorFeedsByState(suite.ctx, exported.RUNNING, func(feed types.Feed) {
		feeds2 = append(feeds2, feed)
	})
	suite.Len(feeds2, 1)
	suite.Equal(msg.FeedName, feeds2[0].FeedName)
	// ================test PauseFeed end================
}

var _ types.ServiceKeeper = MockServiceKeeper{}

type MockServiceKeeper struct {
	cxtMap           map[string]exported.RequestContext
	callbackMap      map[string]exported.ResponseCallback
	stateCallbackMap map[string]exported.StateCallback
	moduleServiceMap map[string]*exported.ModuleService
}

func NewMockServiceKeeper() MockServiceKeeper {
	cxtMap := make(map[string]exported.RequestContext)
	callbackMap := make(map[string]exported.ResponseCallback)
	stateCallbackMap := make(map[string]exported.StateCallback)
	moduleServiceMap := make(map[string]*exported.ModuleService)
	return MockServiceKeeper{
		cxtMap:           cxtMap,
		callbackMap:      callbackMap,
		stateCallbackMap: stateCallbackMap,
		moduleServiceMap: moduleServiceMap,
	}
}

func (m MockServiceKeeper) RegisterStateCallback(moduleName string, stateCallback exported.StateCallback) error {
	m.stateCallbackMap[moduleName] = stateCallback
	return nil
}

func (m MockServiceKeeper) RegisterResponseCallback(moduleName string, respCallback exported.ResponseCallback) error {
	m.callbackMap[moduleName] = respCallback
	return nil
}

func (m MockServiceKeeper) RegisterModuleService(moduleName string, moduleService *exported.ModuleService) error {
	m.moduleServiceMap[moduleName] = moduleService
	return nil
}

func (m MockServiceKeeper) GetRequestContext(ctx sdk.Context, requestContextID tmbytes.HexBytes) (exported.RequestContext, bool) {
	reqCtx, ok := m.cxtMap[string(requestContextID)]
	return reqCtx, ok
}

func (m MockServiceKeeper) CreateRequestContext(
	ctx sdk.Context,
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
	respThreshold uint32,
	moduleName string,
) (tmbytes.HexBytes, error) {
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
	return mockReqCtxID, nil
}

func (m MockServiceKeeper) UpdateRequestContext(
	ctx sdk.Context,
	requestContextID tmbytes.HexBytes,
	providers []sdk.AccAddress,
	respThreshold uint32,
	serviceFeeCap sdk.Coins,
	timeout int64,
	repeatedFreq uint64,
	repeatedTotal int64,
	consumer sdk.AccAddress,
) error {
	return nil
}

func (m MockServiceKeeper) StartRequestContext(ctx sdk.Context, requestContextID tmbytes.HexBytes, consumer sdk.AccAddress) error {
	reqCtx := m.cxtMap[string(requestContextID)]
	callback := m.callbackMap[reqCtx.ModuleName]
	reqCtx.State = servicetypes.RUNNING
	callback(ctx, requestContextID, responses, nil)
	m.cxtMap[string(requestContextID)] = reqCtx
	return nil
}

func (m MockServiceKeeper) PauseRequestContext(ctx sdk.Context, requestContextID tmbytes.HexBytes, consumer sdk.AccAddress) error {
	reqCtx := m.cxtMap[string(requestContextID)]
	reqCtx.State = exported.PAUSED
	m.cxtMap[string(requestContextID)] = reqCtx
	return nil
}
