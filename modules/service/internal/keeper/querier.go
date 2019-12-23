package keeper

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/service/internal/types"
)

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

func queryDefinition(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryDefinitionParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	svcDef, found := keeper.GetServiceDefinition(ctx, params.DefChainID, params.ServiceName)
	if !found {
		return nil, types.ErrSvcDefNotExists(types.DefaultCodespace, params.DefChainID, params.ServiceName)
	}

	iterator := keeper.GetMethods(ctx, params.DefChainID, params.ServiceName)
	defer iterator.Close()

	var methods []types.MethodProperty
	for ; iterator.Valid(); iterator.Next() {
		var method types.MethodProperty
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &method)
		methods = append(methods, method)
	}

	definitionOutput := types.DefinitionOutput{Definition: svcDef, Methods: methods}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, definitionOutput)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryBinding(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryBindingParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	svcBinding, found := keeper.GetServiceBinding(ctx, params.DefChainID, params.ServiceName, params.BindChainID, params.Provider)
	if !found {
		return nil, types.ErrSvcBindingNotExists(types.DefaultCodespace)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, svcBinding)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryBindings(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryBindingsParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	iterator := keeper.ServiceBindingsIterator(ctx, params.DefChainID, params.ServiceName)
	defer iterator.Close()

	var bindings []types.SvcBinding
	for ; iterator.Valid(); iterator.Next() {
		var binding types.SvcBinding
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &binding)
		bindings = append(bindings, binding)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, bindings)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryRequests(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryRequestsParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	iterator := keeper.ActiveBindRequestsIterator(ctx, params.DefChainID, params.ServiceName, params.BindChainID, params.Provider)
	defer iterator.Close()

	var requests []types.SvcRequest
	for ; iterator.Valid(); iterator.Next() {
		var request types.SvcRequest
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &request)
		requests = append(requests, request)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, requests)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryResponse(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryResponseParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	eHeight, rHeight, counter, err := types.ConvertRequestID(params.RequestID)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(err.Error())
	}

	response, found := keeper.GetResponse(ctx, params.ReqChainID, eHeight, rHeight, counter)
	if !found {
		return nil, types.ErrNoResponseFound(types.DefaultCodespace, params.RequestID)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, response)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryFees(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryFeesParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	var feesOutput types.FeesOutput

	if returnFee, found := keeper.GetReturnFee(ctx, params.Address); found {
		feesOutput.ReturnedFee = returnFee.Coins
	}

	if incomingFee, found := keeper.GetIncomingFee(ctx, params.Address); found {
		feesOutput.IncomingFee = incomingFee.Coins
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, feesOutput)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
