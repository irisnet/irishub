package utils

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"

	"github.com/irisnet/irismod/modules/service/types"
)

// QueryRequestContext queries a single request context
func QueryRequestContext(
	cliCtx client.Context, queryRoute string, params types.QueryRequestContextRequest,
) (
	requestContext types.RequestContext, err error,
) {
	queryClient := types.NewQueryClient(cliCtx)
	res, err := queryClient.RequestContext(context.Background(), &params)
	if err != nil {
		return requestContext, err
	}

	requestContext = *res.RequestContext
	if requestContext.Empty() {
		if requestContext, err = QueryRequestContextByTxQuery(cliCtx, queryRoute, params); err != nil {
			return requestContext, err
		}
	}

	if requestContext.Empty() {
		return requestContext, fmt.Errorf("unknown request context: %s", hex.EncodeToString(params.RequestContextId))
	}

	return requestContext, nil
}

// QueryRequestContextByTxQuery will query for a single request context via a direct txs tags query.
func QueryRequestContextByTxQuery(cliCtx client.Context, queryRoute string, params types.QueryRequestContextRequest) (
	requestContext types.RequestContext, err error) {
	txHash, _, err := types.SplitRequestContextID(params.RequestContextId)
	if err != nil {
		return requestContext, err
	}

	// NOTE: QueryTx is used to facilitate the txs query which does not currently
	txInfo, err := authclient.QueryTx(cliCtx, txHash.String())
	if err != nil {
		return requestContext, err
	}

	var msgIndex int
	var found bool
I:
	for i, log := range txInfo.Logs {
		for _, event := range log.Events {
			if event.Type == sdk.EventTypeMessage {
				for _, attribute := range event.Attributes {
					if attribute.Key == types.AttributeKeyRequestContextID &&
						attribute.Value == params.RequestContextId.String() {
						msgIndex = i
						found = true
						break I
					}
				}
			}
		}
	}

	if !found {
		return requestContext, fmt.Errorf("unknown request context: %s", hex.EncodeToString(params.RequestContextId))
	}

	if len(txInfo.GetTx().GetMsgs()) > msgIndex {
		if msg := txInfo.GetTx().GetMsgs()[msgIndex]; msg.Type() == types.TypeMsgCallService {
			requestMsg := msg.(*types.MsgCallService)
			return types.NewRequestContext(
				requestMsg.ServiceName, requestMsg.Providers, requestMsg.Consumer,
				requestMsg.Input, requestMsg.ServiceFeeCap, requestMsg.Timeout,
				requestMsg.SuperMode, requestMsg.Repeated, requestMsg.RepeatedFrequency,
				requestMsg.RepeatedTotal, uint64(requestMsg.RepeatedTotal),
				0, 0, 0, types.BATCHCOMPLETED, types.COMPLETED, 0, "",
			), nil
		}
	}

	return requestContext, nil
}

// QueryRequestByTxQuery will query for a single request via a direct txs tags query.
func QueryRequestByTxQuery(
	cliCtx client.Context, queryRoute string, requestID tmbytes.HexBytes,
) (
	request types.Request, err error,
) {

	contextID, _, requestHeight, batchRequestIndex, err := types.SplitRequestID(requestID)
	if err != nil {
		return request, err
	}

	// query request context
	requestContext, err := QueryRequestContext(
		cliCtx,
		queryRoute,
		types.QueryRequestContextRequest{RequestContextId: contextID},
	)
	if err != nil {
		return request, err
	}

	// query batch request by requestHeight
	node, err := cliCtx.GetNode()
	if err != nil {
		return request, err
	}

	blockResult, err := node.BlockResults(context.Background(), &requestHeight)
	if err != nil {
		return request, err
	}

	for _, event := range blockResult.EndBlockEvents {
		if event.Type == types.EventTypeNewBatchRequest {
			var found bool
			var requests []types.CompactRequest
			var requestsBz []byte
			for _, attribute := range event.Attributes {
				if string(attribute.Key) == types.AttributeKeyRequests {
					requestsBz = attribute.GetValue()
				}
				if string(attribute.Key) == types.AttributeKeyRequestContextID &&
					string(attribute.GetValue()) == contextID.String() {
					found = true
				}
			}
			if found {
				if err := json.Unmarshal(requestsBz, &requests); err != nil {
					return request, err
				}

				if len(requests) > int(batchRequestIndex) {
					compactRequest := requests[batchRequestIndex]
					return types.NewRequest(
						requestID,
						requestContext.ServiceName,
						compactRequest.Provider,
						requestContext.Consumer,
						requestContext.Input,
						compactRequest.ServiceFee,
						requestContext.SuperMode,
						compactRequest.RequestHeight,
						compactRequest.ExpirationHeight,
						compactRequest.RequestContextId,
						compactRequest.RequestContextBatchCounter,
					), nil
				}
			}
		}
	}

	return request, nil
}

// QueryResponseByTxQuery will query for a single request via a direct txs tags query.
func QueryResponseByTxQuery(
	cliCtx client.Context, queryRoute string, requestID tmbytes.HexBytes,
) (
	response types.Response, err error,
) {

	events := []string{
		fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, types.TypeMsgRespondService),
		fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, types.AttributeKeyRequestID, []byte(fmt.Sprintf("%d", requestID))),
	}

	// NOTE: SearchTxs is used to facilitate the txs query which does not currently
	// support configurable pagination.
	result, err := authclient.QueryTxsByEvents(cliCtx, events, 1, 1, "")
	if err != nil {
		return response, err
	}

	if len(result.Txs) == 0 {
		return response, fmt.Errorf("unknown response: %s", requestID)
	}

	contextID, batchCounter, _, _, err := types.SplitRequestID(requestID)
	if err != nil {
		return response, err
	}

	// query request context
	requestContext, err := QueryRequestContext(
		cliCtx,
		queryRoute,
		types.QueryRequestContextRequest{RequestContextId: contextID},
	)
	if err != nil {
		return response, err
	}

	for _, msg := range result.Txs[0].GetTx().GetMsgs() {
		if msg.Type() == types.TypeMsgRespondService {
			responseMsg := msg.(*types.MsgRespondService)
			if responseMsg.RequestId.String() != requestID.String() {
				continue
			}
			return types.NewResponse(
				responseMsg.Provider, requestContext.Consumer,
				responseMsg.Result, responseMsg.Output,
				contextID, batchCounter,
			), nil
		}
	}

	return response, nil
}

// QueryRequestsByBinding queries active requests by the service binding
func QueryRequestsByBinding(
	cliCtx client.Context, queryRoute string, serviceName string, provider sdk.AccAddress,
) (
	[]types.Request, int64, error,
) {
	params := types.QueryRequestsParams{
		ServiceName: serviceName,
		Provider:    provider,
	}

	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return nil, 0, err
	}

	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryRequests)
	res, height, err := cliCtx.QueryWithData(route, bz)
	if err != nil {
		return nil, 0, err
	}

	var requests []types.Request
	if err := cliCtx.LegacyAmino.UnmarshalJSON(res, &requests); err != nil {
		return nil, 0, err
	}

	return requests, height, nil
}

// QueryRequestsByReqCtx queries active requests by the request context ID
func QueryRequestsByReqCtx(
	cliCtx client.Context, queryRoute, reqCtxIDStr, batchCounterStr string,
) (
	[]types.Request, int64, error,
) {
	requestContextID, err := hex.DecodeString(reqCtxIDStr)
	if err != nil {
		return nil, 0, err
	}

	batchCounter, err := strconv.ParseUint(batchCounterStr, 10, 64)
	if err != nil {
		return nil, 0, err
	}

	params := types.QueryRequestsByReqCtxParams{
		RequestContextID: requestContextID,
		BatchCounter:     batchCounter,
	}

	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return nil, 0, err
	}

	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryRequestsByReqCtx)
	res, height, err := cliCtx.QueryWithData(route, bz)
	if err != nil {
		return nil, 0, err
	}

	var requests []types.Request
	if err := cliCtx.LegacyAmino.UnmarshalJSON(res, &requests); err != nil {
		return nil, 0, err
	}

	return requests, height, nil
}
