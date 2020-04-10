package keeper

import (
	"strings"

	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/irisnet/irishub/app/v3/service/internal/types"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

// NewQuerier creates a querier for the service module
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryDefinition:
			return queryDefinition(ctx, req, k)
		case types.QueryBinding:
			return queryBinding(ctx, req, k)
		case types.QueryBindings:
			return queryBindings(ctx, req, k)
		case types.QueryWithdrawAddress:
			return queryWithdrawAddress(ctx, req, k)
		case types.QueryRequest:
			return queryRequest(ctx, req, k)
		case types.QueryRequests:
			return queryRequests(ctx, req, k)
		case types.QueryResponse:
			return queryResponse(ctx, req, k)
		case types.QueryRequestContext:
			return queryRequestContext(ctx, req, k)
		case types.QueryRequestsByReqCtx:
			return queryRequestsByReqCtx(ctx, req, k)
		case types.QueryResponses:
			return queryResponses(ctx, req, k)
		case types.QueryEarnedFees:
			return queryEarnedFees(ctx, req, k)
		case types.QuerySchema:
			return querySchema(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown service query endpoint")
		}
	}
}

func queryDefinition(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryDefinitionParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	svcDef, found := k.GetServiceDefinition(ctx, params.ServiceName)
	if !found {
		return nil, types.ErrUnknownServiceDefinition(types.DefaultCodespace, params.ServiceName)
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, svcDef)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

func queryBinding(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryBindingParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	svcBinding, found := k.GetServiceBinding(ctx, params.ServiceName, params.Provider)
	if !found {
		return nil, types.ErrUnknownServiceBinding(types.DefaultCodespace)
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, svcBinding)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

func queryBindings(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryBindingsParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	iterator := k.ServiceBindingsIterator(ctx, params.ServiceName)
	defer iterator.Close()

	bindings := make([]types.ServiceBinding, 0)

	for ; iterator.Valid(); iterator.Next() {
		var binding types.ServiceBinding
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &binding)

		bindings = append(bindings, binding)
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, bindings)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

func queryWithdrawAddress(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryWithdrawAddressParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	withdrawAddr := k.GetWithdrawAddress(ctx, params.Provider)

	bz, err := codec.MarshalJSONIndent(k.cdc, withdrawAddr)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

func queryRequest(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryRequestParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	if len(params.RequestID) != types.RequestIDLen {
		return nil, types.ErrInvalidRequestID(types.DefaultCodespace, params.RequestID)
	}

	request, _ := k.GetRequest(ctx, params.RequestID)

	bz, err := codec.MarshalJSONIndent(k.cdc, request)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

func queryRequests(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryRequestsParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	iterator := k.ActiveRequestsIterator(ctx, params.ServiceName, params.Provider)
	defer iterator.Close()

	requests := make([]types.Request, 0)

	for ; iterator.Valid(); iterator.Next() {
		var requestID cmn.HexBytes
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &requestID)

		request, _ := k.GetRequest(ctx, requestID)
		requests = append(requests, request)
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, requests)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

func queryResponse(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryResponseParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	if len(params.RequestID) != types.RequestIDLen {
		return nil, types.ErrInvalidRequestID(types.DefaultCodespace, params.RequestID)
	}

	response, _ := k.GetResponse(ctx, params.RequestID)

	bz, err := codec.MarshalJSONIndent(k.cdc, response)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

func queryRequestContext(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryRequestContextParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	requestContext, _ := k.GetRequestContext(ctx, params.RequestContextID)
	bz, err := codec.MarshalJSONIndent(k.cdc, requestContext)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

func queryRequestsByReqCtx(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryRequestsByReqCtxParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	iterator := k.RequestsIteratorByReqCtx(ctx, params.RequestContextID, params.BatchCounter)
	defer iterator.Close()

	requests := make([]types.Request, 0)

	for ; iterator.Valid(); iterator.Next() {
		requestID := iterator.Key()[1:]
		request, _ := k.GetRequest(ctx, requestID)

		requests = append(requests, request)
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, requests)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

func queryResponses(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryResponsesParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	iterator := k.ResponsesIteratorByReqCtx(ctx, params.RequestContextID, params.BatchCounter)
	defer iterator.Close()

	responses := make([]types.Response, 0)

	for ; iterator.Valid(); iterator.Next() {
		var response types.Response
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &response)

		responses = append(responses, response)
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, responses)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

func queryEarnedFees(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryEarnedFeesParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	fees, found := k.GetEarnedFees(ctx, params.Provider)
	if !found {
		return nil, types.ErrNoEarnedFees(types.DefaultCodespace, params.Provider)
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, fees)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

func querySchema(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QuerySchemaParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	var schema string

	if strings.ToLower(params.SchemaName) == "pricing" {
		schema = types.PricingSchema
	} else if strings.ToLower(params.SchemaName) == "result" {
		schema = types.ResultSchema
	} else {
		return nil, types.ErrInvalidSchemaName(types.DefaultCodespace)
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, schema)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}
