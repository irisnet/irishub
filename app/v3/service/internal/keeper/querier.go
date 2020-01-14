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

	svcBinding, found := k.GetServiceBinding(ctx, params.DefChainID, params.ServiceName, params.BindChainID, params.Provider)
	if !found {
		return nil, types.ErrSvcBindingNotExists(types.DefaultCodespace)
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

	iterator := k.ServiceBindingsIterator(ctx, params.DefChainID, params.ServiceName)
	defer iterator.Close()

	var bindings []types.SvcBinding
	for ; iterator.Valid(); iterator.Next() {
		var binding types.SvcBinding
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
	var params types.QueryBindingParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	iterator := k.ActiveBindRequestsIterator(ctx, params.DefChainID, params.ServiceName, params.BindChainID, params.Provider)
	defer iterator.Close()

	var requests []types.SvcRequest
	for ; iterator.Valid(); iterator.Next() {
		var request types.SvcRequest
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &request)
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

	eHeight, rHeight, counter, err := types.ConvertRequestID(params.RequestID)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(err.Error())
	}

	response, found := k.GetResponse(ctx, params.ReqChainID, eHeight, rHeight, counter)
	if !found {
		return nil, types.ErrNoResponseFound(types.DefaultCodespace, params.RequestID)
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

	var feesOutput types.FeesOutput

	if returnFee, found := k.GetReturnFee(ctx, params.Address); found {
		feesOutput.ReturnedFee = returnFee.Coins
	}

	if incomingFee, found := k.GetIncomingFee(ctx, params.Address); found {
		feesOutput.IncomingFee = incomingFee.Coins
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, feesOutput)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}

	return bz, nil
}
