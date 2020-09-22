package keeper

import (
	"strings"

	gogotypes "github.com/gogo/protobuf/types"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/service/types"
)

// NewQuerier creates a new service Querier instance
func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryDefinition:
			return queryServiceDefinition(ctx, path[1:], req, k, legacyQuerierCdc)

		case types.QueryBinding:
			return queryBinding(ctx, req, k, legacyQuerierCdc)

		case types.QueryBindings:
			return queryBindings(ctx, req, k, legacyQuerierCdc)

		case types.QueryWithdrawAddress:
			return queryWithdrawAddress(ctx, req, k, legacyQuerierCdc)

		case types.QueryRequest:
			return queryRequest(ctx, req, k, legacyQuerierCdc)

		case types.QueryRequests:
			return queryRequests(ctx, req, k, legacyQuerierCdc)

		case types.QueryResponse:
			return queryResponse(ctx, req, k, legacyQuerierCdc)

		case types.QueryRequestContext:
			return queryRequestContext(ctx, req, k, legacyQuerierCdc)

		case types.QueryRequestsByReqCtx:
			return queryRequestsByReqCtx(ctx, req, k, legacyQuerierCdc)

		case types.QueryResponses:
			return queryResponses(ctx, req, k, legacyQuerierCdc)

		case types.QueryEarnedFees:
			return queryEarnedFees(ctx, req, k, legacyQuerierCdc)

		case types.QuerySchema:
			return querySchema(ctx, req, k, legacyQuerierCdc)

		case types.QueryParameters:
			return queryParams(ctx, k, legacyQuerierCdc)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query path: %s", types.ModuleName, path[0])
		}
	}
}

func queryServiceDefinition(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryDefinitionParams
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	definition, found := k.GetServiceDefinition(ctx, params.ServiceName)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrUnknownServiceDefinition, params.ServiceName)
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, definition)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryBinding(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryBindingParams
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	svcBinding, found := k.GetServiceBinding(ctx, params.ServiceName, params.Provider)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrUnknownServiceBinding, "")
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, svcBinding)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryBindings(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryBindingsParams
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	bindings := make([]*types.ServiceBinding, 0)

	if params.Owner.Empty() {
		iterator := k.ServiceBindingsIterator(ctx, params.ServiceName)
		defer iterator.Close()

		for ; iterator.Valid(); iterator.Next() {
			var binding types.ServiceBinding
			k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &binding)

			bindings = append(bindings, &binding)
		}
	} else {
		bindings = k.GetOwnerServiceBindings(ctx, params.Owner, params.ServiceName)
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, bindings)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryWithdrawAddress(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryWithdrawAddressParams
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	withdrawAddr := k.GetWithdrawAddress(ctx, params.Owner)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, withdrawAddr)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryRequest(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryRequestParams
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	if len(params.RequestID) != types.RequestIDLen {
		return nil, sdkerrors.Wrapf(
			types.ErrInvalidRequestID,
			"invalid length, expected: %d, got: %d",
			types.RequestIDLen, len(params.RequestID),
		)
	}

	request, _ := k.GetRequest(ctx, params.RequestID)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, request)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryRequests(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryRequestsParams
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	iterator := k.ActiveRequestsIterator(ctx, params.ServiceName, params.Provider)
	defer iterator.Close()

	requests := make([]types.Request, 0)

	for ; iterator.Valid(); iterator.Next() {
		var requestID gogotypes.BytesValue

		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &requestID)

		request, _ := k.GetRequest(ctx, requestID.Value)
		requests = append(requests, request)
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, requests)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryResponse(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryResponseParams
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	if len(params.RequestID) != types.RequestIDLen {
		return nil, sdkerrors.Wrapf(types.ErrInvalidRequestID, "invalid length, expected: %d, got: %d",
			types.RequestIDLen, len(params.RequestID))
	}

	response, _ := k.GetResponse(ctx, params.RequestID)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, response)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryRequestContext(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryRequestContextParams
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	requestContext, _ := k.GetRequestContext(ctx, params.RequestContextID)
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, requestContext)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryRequestsByReqCtx(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryRequestsByReqCtxParams
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	iterator := k.RequestsIteratorByReqCtx(ctx, params.RequestContextID, params.BatchCounter)
	defer iterator.Close()

	requests := make([]types.Request, 0)

	for ; iterator.Valid(); iterator.Next() {
		requestID := iterator.Key()[1:]
		request, _ := k.GetRequest(ctx, requestID)

		requests = append(requests, request)
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, requests)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryResponses(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryResponsesParams
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	iterator := k.ResponsesIteratorByReqCtx(ctx, params.RequestContextID, params.BatchCounter)
	defer iterator.Close()

	responses := make([]types.Response, 0)

	for ; iterator.Valid(); iterator.Next() {
		var response types.Response
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &response)

		responses = append(responses, response)
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, responses)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryEarnedFees(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryEarnedFeesParams
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	fees, found := k.GetEarnedFees(ctx, params.Provider)
	if !found {
		return nil, sdkerrors.Wrapf(
			types.ErrNoEarnedFees, "no earned fees for %s", params.Provider.String(),
		)
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, fees)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func querySchema(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QuerySchemaParams
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var schemaName = strings.ToLower(params.SchemaName)
	var schema string

	switch schemaName {
	case "pricing":
		schema = types.PricingSchema
	case "result":
		schema = types.ResultSchema
	default:
		return nil, sdkerrors.Wrap(types.ErrInvalidSchemaName, schema)
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, schema)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryParams(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	params := k.GetParams(ctx)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
