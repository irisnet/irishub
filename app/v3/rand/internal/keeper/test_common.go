package keeper

import (
	"time"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v3/rand/internal/types"
	"github.com/irisnet/irishub/app/v3/service"
	"github.com/irisnet/irishub/app/v3/service/exported"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

var (
	emptyByte         = []byte{0x00}
	serviceBindingKey = []byte{0x02}
)

var (
	_ types.ServiceKeeper = MockServiceKeeper{}

	responses = []string{
		`{"seed":"3132333435363738393031323334353637383930313233343536373839303132"}`,
	}

	mockReqCtxID = []byte("mockRequest")
)

type MockServiceKeeper struct {
	storeKey    sdk.StoreKey
	cdc         *codec.Codec
	cxtMap      map[string]exported.RequestContext
	callbackMap map[string]exported.ResponseCallback
}

func NewMockServiceKeeper() MockServiceKeeper {
	storeKey := protocol.KeyService
	cdc := codec.New()
	service.RegisterCodec(cdc)
	cxtMap := make(map[string]exported.RequestContext)
	callbackMap := make(map[string]exported.ResponseCallback)
	return MockServiceKeeper{
		storeKey:    storeKey,
		cdc:         cdc,
		cxtMap:      cxtMap,
		callbackMap: callbackMap,
	}
}

func (m MockServiceKeeper) RegisterResponseCallback(
	moduleName string,
	respCallback exported.ResponseCallback,
) sdk.Error {
	m.callbackMap[moduleName] = respCallback
	return nil
}

func (m MockServiceKeeper) GetRequestContext(
	ctx sdk.Context,
	requestContextID []byte,
) (exported.RequestContext, bool) {
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
	respThreshold uint16,
	moduleName string,
) ([]byte, sdk.Error) {
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
	requestContextID []byte,
	providers []sdk.AccAddress,
	serviceFeeCap sdk.Coins,
	timeout int64,
	repeatedFreq uint64,
	repeatedTotal int64,
	consumer sdk.AccAddress,
) sdk.Error {
	return nil
}

func (m MockServiceKeeper) StartRequestContext(
	ctx sdk.Context,
	requestContextID []byte,
	consumer sdk.AccAddress,
) sdk.Error {
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

func (m MockServiceKeeper) PauseRequestContext(
	ctx sdk.Context,
	requestContextID []byte,
	consumer sdk.AccAddress,
) sdk.Error {
	reqCtx := m.cxtMap[string(requestContextID)]
	reqCtx.State = exported.PAUSED
	m.cxtMap[string(requestContextID)] = reqCtx
	return nil
}

func (m MockServiceKeeper) GetParamSet(ctx sdk.Context) service.Params {
	return service.DefaultParams()
}

func (m MockServiceKeeper) SetServiceBinding(ctx sdk.Context, svcBinding service.ServiceBinding) {
	store := ctx.KVStore(m.storeKey)

	bz := m.cdc.MustMarshalBinaryLengthPrefixed(svcBinding)
	store.Set(GetServiceBindingKey(svcBinding.ServiceName, svcBinding.Provider), bz)
}

func (m MockServiceKeeper) ServiceBindingsIterator(ctx sdk.Context, serviceName string) sdk.Iterator {
	store := ctx.KVStore(m.storeKey)
	return sdk.KVStorePrefixIterator(store, GetBindingsSubspace(serviceName))
}

func GetServiceBindingKey(serviceName string, provider sdk.AccAddress) []byte {
	return append(serviceBindingKey, getStringsKey([]string{serviceName, provider.String()})...)
}

// GetBindingsSubspace returns the key for retrieving all bindings of the specified service
func GetBindingsSubspace(serviceName string) []byte {
	return append(append(serviceBindingKey, []byte(serviceName)...), emptyByte...)
}

func getStringsKey(ss []string) (result []byte) {
	for _, s := range ss {
		result = append(append(result, []byte(s)...), emptyByte...)
	}

	if len(result) > 0 {
		return result[0 : len(result)-1]
	}

	return
}
