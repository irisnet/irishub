package keeper

import (
	"context"
	"encoding/hex"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

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

	provider, err := sdk.AccAddressFromBech32(req.Provider)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid provider address (%s)", err)
	}

	ctx := sdk.UnwrapSDKContext(c)

	binding, found := k.GetServiceBinding(ctx, req.ServiceName, provider)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrUnknownServiceBinding, "service: %s, provider: %s", req.ServiceName, req.Provider)
	}

	return &types.QueryBindingResponse{ServiceBinding: &binding}, nil
}

func (k Keeper) Bindings(c context.Context, req *types.QueryBindingsRequest) (*types.QueryBindingsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	bindings := make([]*types.ServiceBinding, 0)
	store := ctx.KVStore(k.storeKey)
	var pageRes *query.PageResponse
	var err error
	if len(req.Owner) == 0 {
		bindingStore := prefix.NewStore(store, types.GetBindingsSubspace(req.ServiceName))
		pageRes, err = query.Paginate(bindingStore, req.Pagination, func(key []byte, value []byte) error {
			var binding types.ServiceBinding
			k.cdc.MustUnmarshalBinaryBare(value, &binding)
			bindings = append(bindings, &binding)
			return nil
		})
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
		}
	} else {
		owner, err := sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
		}
		bindingStore := prefix.NewStore(store, types.GetOwnerBindingsSubspace(owner, req.ServiceName))
		pageRes, err = query.Paginate(bindingStore, req.Pagination, func(key []byte, value []byte) error {
			provider := sdk.AccAddress(key)

			if binding, found := k.GetServiceBinding(ctx, req.ServiceName, provider); found {
				bindings = append(bindings, &binding)
			}
			return nil
		})
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
		}
	}

	return &types.QueryBindingsResponse{ServiceBindings: bindings, Pagination: pageRes}, nil
}

func (k Keeper) WithdrawAddress(c context.Context, req *types.QueryWithdrawAddressRequest) (*types.QueryWithdrawAddressResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	owner, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	ctx := sdk.UnwrapSDKContext(c)

	withdrawAddr := k.GetWithdrawAddress(ctx, owner)

	return &types.QueryWithdrawAddressResponse{WithdrawAddress: withdrawAddr.String()}, nil
}

func (k Keeper) RequestContext(c context.Context, req *types.QueryRequestContextRequest) (*types.QueryRequestContextResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if len(req.RequestContextId) != types.ContextIDLen {
		return nil, sdkerrors.Wrapf(types.ErrInvalidRequestContextID, "length of the request context ID must be %d in bytes", types.ContextIDLen)
	}
	requestContextId, err := hex.DecodeString(req.RequestContextId)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequestContextID, "request context ID must be a hex encoded string")
	}

	ctx := sdk.UnwrapSDKContext(c)

	requestContext, _ := k.GetRequestContext(ctx, requestContextId)

	return &types.QueryRequestContextResponse{RequestContext: &requestContext}, nil
}

func (k Keeper) Request(c context.Context, req *types.QueryRequestRequest) (*types.QueryRequestResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if len(req.RequestId) != types.RequestIDLen {
		return nil, sdkerrors.Wrapf(
			types.ErrInvalidRequestID,
			"invalid length, expected: %d, got: %d",
			types.RequestIDLen, len(req.RequestId),
		)
	}

	requestId, err := hex.DecodeString(req.RequestId)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequestContextID, "request ID must be a hex encoded string")
	}

	ctx := sdk.UnwrapSDKContext(c)
	request, _ := k.GetRequest(ctx, requestId)

	return &types.QueryRequestResponse{Request: &request}, nil
}

func (k Keeper) Requests(c context.Context, req *types.QueryRequestsRequest) (*types.QueryRequestsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	provider, err := sdk.AccAddressFromBech32(req.Provider)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid provider address (%s)", err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	requests := make([]*types.Request, 0)
	store := ctx.KVStore(k.storeKey)
	requestStore := prefix.NewStore(store, types.GetActiveRequestSubspace(req.ServiceName, provider))
	pageRes, err := query.Paginate(requestStore, req.Pagination, func(key []byte, value []byte) error {
		var requestID gogotypes.BytesValue
		k.cdc.MustUnmarshalBinaryBare(value, &requestID)
		request, _ := k.GetRequest(ctx, requestID.Value)
		requests = append(requests, &request)
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryRequestsResponse{Requests: requests, Pagination: pageRes}, nil
}

func (k Keeper) RequestsByReqCtx(c context.Context, req *types.QueryRequestsByReqCtxRequest) (*types.QueryRequestsByReqCtxResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if len(req.RequestContextId) != types.ContextIDLen {
		return nil, sdkerrors.Wrapf(types.ErrInvalidRequestContextID, "length of the request context ID must be %d in bytes", types.ContextIDLen)
	}
	requestContextId, err := hex.DecodeString(req.RequestContextId)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequestContextID, "request context ID must be a hex encoded string")
	}

	ctx := sdk.UnwrapSDKContext(c)
	requests := make([]*types.Request, 0)
	store := ctx.KVStore(k.storeKey)
	requestStore := prefix.NewStore(store, types.GetRequestSubspaceByReqCtx(requestContextId, req.BatchCounter))
	pageRes, err := query.Paginate(requestStore, req.Pagination, func(key []byte, value []byte) error {
		requestID := append(append(requestContextId, sdk.Uint64ToBigEndian(req.BatchCounter)...),key...)
		request, _ := k.GetRequest(ctx, requestID)
		requests = append(requests, &request)
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryRequestsByReqCtxResponse{Requests: requests, Pagination: pageRes}, nil
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

	requestId, err := hex.DecodeString(req.RequestId)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequestContextID, "request ID must be a hex encoded string")
	}
	response, _ := k.GetResponse(ctx, requestId)

	return &types.QueryResponseResponse{Response: &response}, nil
}

func (k Keeper) Responses(c context.Context, req *types.QueryResponsesRequest) (*types.QueryResponsesResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if len(req.RequestContextId) != types.ContextIDLen {
		return nil, sdkerrors.Wrapf(types.ErrInvalidRequestContextID, "length of the request context ID must be %d in bytes", types.ContextIDLen)
	}
	requestContextId, err := hex.DecodeString(req.RequestContextId)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequestContextID, "request context ID must be a hex encoded string")
	}

	ctx := sdk.UnwrapSDKContext(c)
	responses := make([]*types.Response, 0)
	store := ctx.KVStore(k.storeKey)

	responseStore := prefix.NewStore(store, types.GetResponseSubspaceByReqCtx(requestContextId, req.BatchCounter))
	pageRes, err := query.Paginate(responseStore, req.Pagination, func(key []byte, value []byte) error {
		var response types.Response
		k.cdc.MustUnmarshalBinaryBare(value, &response)
		responses = append(responses, &response)
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryResponsesResponse{Responses: responses, Pagination: pageRes}, nil
}

func (k Keeper) EarnedFees(c context.Context, req *types.QueryEarnedFeesRequest) (*types.QueryEarnedFeesResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	provider, err := sdk.AccAddressFromBech32(req.Provider)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid provider address (%s)", err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	fees, found := k.GetEarnedFees(ctx, provider)
	if !found {
		return nil, sdkerrors.Wrapf(
			types.ErrNoEarnedFees, "no earned fees for %s", req.Provider,
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
