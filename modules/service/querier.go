package service

import (
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryDefinition = "definition"
	QueryBinding    = "binding"
	QueryBindings   = "bindings"
	QueryRequests   = "requests"
	QueryResponse   = "response"
	QueryFees       = "fees"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryDefinition:
			return queryDefinition(ctx, req, k)
		case QueryBinding:
			return queryBinding(ctx, req, k)
		case QueryBindings:
			return queryBindings(ctx, req, k)
		case QueryRequests:
			return queryRequests(ctx, req, k)
		case QueryResponse:
			return queryResponse(ctx, req, k)
		case QueryFees:
			return queryFees(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown service query endpoint")
		}
	}
}

type QueryServiceParams struct {
	DefChainID  string
	ServiceName string
}

type DefinitionOutput struct {
	Definition SvcDef           `json:"definition"`
	Methods    []MethodProperty `json:"methods"`
}

func queryDefinition(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params QueryServiceParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}
	svcDef, found := k.GetServiceDefinition(ctx, params.DefChainID, params.ServiceName)
	if !found {
		return nil, ErrSvcDefNotExists(DefaultCodespace, params.DefChainID, params.ServiceName)
	}

	iterator := k.GetMethods(ctx, params.DefChainID, params.ServiceName)
	defer iterator.Close()
	var methods []MethodProperty
	for ; iterator.Valid(); iterator.Next() {
		var method MethodProperty
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &method)
		methods = append(methods, method)
	}

	definitionOutput := DefinitionOutput{Definition: svcDef, Methods: methods}

	bz, err := codec.MarshalJSONIndent(k.cdc, definitionOutput)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return bz, nil
}

type QueryBindingParams struct {
	DefChainID  string
	ServiceName string
	BindChainId string
	Provider    sdk.AccAddress
}

func queryBinding(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params QueryBindingParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}
	svcBinding, found := k.GetServiceBinding(ctx, params.DefChainID, params.ServiceName, params.BindChainId, params.Provider)
	if !found {
		return nil, ErrSvcBindingNotExists(DefaultCodespace)
	}
	bz, err := codec.MarshalJSONIndent(k.cdc, svcBinding)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return bz, nil
}

func queryBindings(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params QueryServiceParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	iterator := k.ServiceBindingsIterator(ctx, params.DefChainID, params.ServiceName)
	defer iterator.Close()
	var bindings []SvcBinding
	for ; iterator.Valid(); iterator.Next() {
		var binding SvcBinding
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
	var params QueryBindingParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	iterator := k.ActiveBindRequestsIterator(ctx, params.DefChainID, params.ServiceName, params.BindChainId, params.Provider)
	defer iterator.Close()
	var requests []SvcRequest
	for ; iterator.Valid(); iterator.Next() {
		var request SvcRequest
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &request)
		requests = append(requests, request)
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, requests)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return bz, nil
}

type QueryResponseParams struct {
	ReqChainId string
	RequestId  string
}

func queryResponse(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params QueryResponseParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	eHeight, rHeight, counter, err := ConvertRequestID(params.RequestId)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(err.Error())
	}
	response, found := k.GetResponse(ctx, params.ReqChainId, eHeight, rHeight, counter)
	if !found {
		return nil, ErrNoResponseFound(DefaultCodespace, params.RequestId)
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, response)
	if err != nil {
		return nil, sdk.MarshalResultErr(err)
	}
	return bz, nil
}

type QueryFeesParams struct {
	Address sdk.AccAddress
}

type FeesOutput struct {
	ReturnedFee sdk.Coins `json:"returned_fee"`
	IncomingFee sdk.Coins `json:"incoming_fee"`
}

func queryFees(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params QueryFeesParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ParseParamsErr(err)
	}

	var feesOutput FeesOutput

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
