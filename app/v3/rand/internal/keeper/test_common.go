package keeper

import (
	"time"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/irisnet/irishub/app/v3/rand/internal/types"
	"github.com/irisnet/irishub/app/v3/service"
	"github.com/irisnet/irishub/app/v3/service/exported"
	sdk "github.com/irisnet/irishub/types"
)

var (
	_ types.ServiceKeeper = MockServiceKeeper{}

	responses = []string{
		// TODO
		`{}`,
		`{}`,
		`{}`,
		`{}`,
	}

	mockReqCtxID = []byte("mockRequest")
)

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
	requestContextID []byte) (exported.RequestContext, bool) {
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
	moduleName string) ([]byte, sdk.Error) {

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

func (m MockServiceKeeper) UpdateRequestContext(ctx sdk.Context,
	requestContextID []byte,
	providers []sdk.AccAddress,
	serviceFeeCap sdk.Coins,
	timeout int64,
	repeatedFreq uint64,
	repeatedTotal int64) sdk.Error {
	return nil
}

func (m MockServiceKeeper) StartRequestContext(ctx sdk.Context, requestContextID []byte) sdk.Error {
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
		callback(ctx, requestContextID, responses)
	}
	return nil
}

func (m MockServiceKeeper) PauseRequestContext(ctx sdk.Context, requestContextID []byte) sdk.Error {
	reqCtx := m.cxtMap[string(requestContextID)]
	reqCtx.State = exported.PAUSED
	m.cxtMap[string(requestContextID)] = reqCtx
	return nil
}

func (m MockServiceKeeper) GetParamSet(ctx sdk.Context) service.Params {
	return service.DefaultParams()
}
