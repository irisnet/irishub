package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

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
		case types.QueryRequests:
			return queryRequests(ctx, req, k)
		case types.QueryResponse:
			return queryResponse(ctx, req, k)
		case types.QueryFees:
			return queryFees(ctx, req, k)
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
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
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
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	iterator := k.ServiceBindingsIterator(ctx, params.ServiceName)
	defer iterator.Close()

	var bindings []types.ServiceBinding
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

func queryRequests(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryRequestsParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	iterator := k.ActiveRequestsIterator(ctx, params.ServiceName, params.Provider)
	defer iterator.Close()

	var requests []types.Request
	for ; iterator.Valid(); iterator.Next() {
		var requestID []byte
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
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	requestID, err := types.ConvertRequestID(params.RequestID)
	if err != nil {
		return nil, types.ErrInvalidRequestID(types.DefaultCodespace, params.RequestID)
	}

	response, found := k.GetResponse(ctx, requestID)
	if !found {
		return nil, types.ErrInvalidRequestID(types.DefaultCodespace, params.RequestID)
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, response)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}

func queryFees(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryFeesParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	fees, found := k.GetEarnedFees(ctx, params.Address)
	if !found {
		return nil, types.ErrNoEarnedFees(types.DefaultCodespace, params.Address)
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, fees)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}
