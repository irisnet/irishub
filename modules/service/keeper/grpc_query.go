package keeper

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gogotypes "github.com/gogo/protobuf/types"

	"github.com/irisnet/irismod/modules/service/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Definition(c context.Context, req *types.QueryDefinitionRequest) (*types.QueryDefinitionResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	definition, found := k.GetServiceDefinition(ctx, req.ServiceName)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrUnknownServiceDefinition, req.ServiceName)
	}

	return &types.QueryDefinitionResponse{ServiceDefinition: &definition}, nil
}

func (k Keeper) Binding(c context.Context, req *types.QueryBindingRequest) (*types.QueryBindingResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	binding, found := k.GetServiceBinding(ctx, req.ServiceName, req.Provider)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrUnknownServiceBinding, "service: %s, provider: %s", req.ServiceName, req.Provider.String())
	}

	return &types.QueryBindingResponse{ServiceBinding: &binding}, nil
}

func (k Keeper) Bindings(c context.Context, req *types.QueryBindingsRequest) (*types.QueryBindingsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	bindings := make([]*types.ServiceBinding, 0)
	if req.Owner.Empty() {
		iterator := k.ServiceBindingsIterator(ctx, req.ServiceName)
		defer iterator.Close()

		for ; iterator.Valid(); iterator.Next() {
			var binding types.ServiceBinding
			k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &binding)

			bindings = append(bindings, &binding)
		}
	} else {
		bindings = k.GetOwnerServiceBindings(ctx, req.Owner, req.ServiceName)
	}

	return &types.QueryBindingsResponse{ServiceBindings: bindings}, nil
}

func (k Keeper) WithdrawAddress(c context.Context, req *types.QueryWithdrawAddressRequest) (*types.QueryWithdrawAddressResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	withdrawAddr := k.GetWithdrawAddress(ctx, req.Owner)

	return &types.QueryWithdrawAddressResponse{WithdrawAddress: withdrawAddr}, nil
}

func (k Keeper) RequestContext(c context.Context, req *types.QueryRequestContextRequest) (*types.QueryRequestContextResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	requestContext, _ := k.GetRequestContext(ctx, req.RequestContextId)

	return &types.QueryRequestContextResponse{RequestContext: &requestContext}, nil
}

func (k Keeper) Request(c context.Context, req *types.QueryRequestRequest) (*types.QueryRequestResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	if len(req.RequestId) != types.RequestIDLen {
		return nil, sdkerrors.Wrapf(
			types.ErrInvalidRequestID,
			"invalid length, expected: %d, got: %d",
			types.RequestIDLen, len(req.RequestId),
		)
	}

	request, _ := k.GetRequest(ctx, req.RequestId)

	return &types.QueryRequestResponse{Request: &request}, nil
}

func (k Keeper) Requests(c context.Context, req *types.QueryRequestsRequest) (*types.QueryRequestsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	iterator := k.ActiveRequestsIterator(ctx, req.ServiceName, req.Provider)
	defer iterator.Close()

	requests := make([]*types.Request, 0)

	for ; iterator.Valid(); iterator.Next() {
		var requestID gogotypes.BytesValue

		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &requestID)

		request, _ := k.GetRequest(ctx, requestID.Value)
		requests = append(requests, &request)
	}

	return &types.QueryRequestsResponse{Requests: requests}, nil
}

func (k Keeper) RequestsByReqCtx(c context.Context, req *types.QueryRequestsByReqCtxRequest) (*types.QueryRequestsByReqCtxResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	iterator := k.RequestsIteratorByReqCtx(ctx, req.RequestContextId, req.BatchCounter)
	defer iterator.Close()

	requests := make([]*types.Request, 0)
	for ; iterator.Valid(); iterator.Next() {
		requestID := iterator.Key()[1:]
		request, _ := k.GetRequest(ctx, requestID)

		requests = append(requests, &request)
	}

	return &types.QueryRequestsByReqCtxResponse{Requests: requests}, nil
}

func (k Keeper) Response(c context.Context, req *types.QueryResponseRequest) (*types.QueryResponseResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	if len(req.RequestId) != types.RequestIDLen {
		return nil, sdkerrors.Wrapf(
			types.ErrInvalidRequestID,
			"invalid length, expected: %d, got: %d",
			types.RequestIDLen, len(req.RequestId),
		)
	}

	response, _ := k.GetResponse(ctx, req.RequestId)

	return &types.QueryResponseResponse{Response: &response}, nil
}

func (k Keeper) Responses(c context.Context, req *types.QueryResponsesRequest) (*types.QueryResponsesResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	iterator := k.ResponsesIteratorByReqCtx(ctx, req.RequestContextId, req.BatchCounter)
	defer iterator.Close()

	responses := make([]*types.Response, 0)
	for ; iterator.Valid(); iterator.Next() {
		var response types.Response
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &response)

		responses = append(responses, &response)
	}

	return &types.QueryResponsesResponse{Responses: responses}, nil
}

func (k Keeper) EarnedFees(c context.Context, req *types.QueryEarnedFeesRequest) (*types.QueryEarnedFeesResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	fees, found := k.GetEarnedFees(ctx, req.Provider)
	if !found {
		return nil, sdkerrors.Wrapf(
			types.ErrNoEarnedFees, "no earned fees for %s", req.Provider.String(),
		)
	}

	return &types.QueryEarnedFeesResponse{Fees: fees}, nil
}

func (k Keeper) Schema(c context.Context, req *types.QuerySchemaRequest) (*types.QuerySchemaResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	var schemaName = strings.ToLower(req.SchemaName)
	var schema string
	switch schemaName {
	case "pricing":
		schema = types.PricingSchema
	case "result":
		schema = types.ResultSchema
	default:
		return nil, sdkerrors.Wrap(types.ErrInvalidSchemaName, schema)
	}

	return &types.QuerySchemaResponse{Schema: schema}, nil
}

func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}
