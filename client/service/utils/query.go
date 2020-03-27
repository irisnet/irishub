package utils

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v3/service"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/tendermint/tx"
	sdk "github.com/irisnet/irishub/types"
)

// QueryRequestContext queries a single request context
func QueryRequestContext(cliCtx context.CLIContext, params service.QueryRequestContextParams) (
	requestContext service.RequestContext, err error) {
	bz, err := cliCtx.Codec.MarshalJSON(params)
	if err != nil {
		return requestContext, err
	}

	route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryRequestContext)
	res, err := cliCtx.QueryWithData(route, bz)
	if err != nil {
		return requestContext, err
	}

	_ = cliCtx.Codec.UnmarshalJSON(res, &requestContext)
	if requestContext.Empty() {
		requestContext, err = QueryRequestContextByTxQuery(cliCtx, params)
		if err != nil {
			return requestContext, err
		}
	}

	if requestContext.Empty() {
		return requestContext, fmt.Errorf("unknown request context: %s", hex.EncodeToString(params.RequestContextID))
	}
	return requestContext, nil
}

// QueryRequestContextByTxQuery will query for a single request context via a direct txs tags query.
func QueryRequestContextByTxQuery(cliCtx context.CLIContext, params service.QueryRequestContextParams) (
	requestContext service.RequestContext, err error) {
	txHash, msgIndex, err := service.SplitRequestContextID(params.RequestContextID)
	if err != nil {
		return requestContext, err
	}

	// NOTE: QueryTx is used to facilitate the txs query which does not currently
	txInfo, err := tx.QueryTx(cliCtx, txHash)
	if err != nil {
		return requestContext, err
	}

	if int64(len(txInfo.Tx.GetMsgs())) > msgIndex {
		msg := txInfo.Tx.GetMsgs()[msgIndex]
		if msg.Type() == service.TypeMsgRequestService {
			requestMsg := msg.(service.MsgRequestService)
			requestContext := service.NewRequestContext(
				requestMsg.ServiceName, requestMsg.Providers,
				requestMsg.Consumer, requestMsg.Input, requestMsg.ServiceFeeCap,
				requestMsg.Timeout, requestMsg.SuperMode, requestMsg.Repeated,
				requestMsg.RepeatedFrequency, requestMsg.RepeatedTotal,
				uint64(requestMsg.RepeatedTotal), 0, 0,
				service.BATCHCOMPLETED, service.COMPLETED, 0, "",
			)
			return requestContext, nil
		}
	}

	return requestContext, nil
}

// QueryRequestByTxQuery will query for a single request via a direct txs tags query.
func QueryRequestByTxQuery(cliCtx context.CLIContext, params service.QueryRequestParams) (
	request service.Request, err error) {
	requestID := params.RequestID
	if err != nil {
		return request, nil
	}

	contextID, _, requestHeight, batchRequestIndex, err := service.SplitRequestID(requestID)
	if err != nil {
		return request, err
	}

	// query request context
	requestContext, err := QueryRequestContext(cliCtx, service.QueryRequestContextParams{
		RequestContextID: contextID,
	})

	if err != nil {
		return request, err
	}

	// query batch request by requestHeight
	node, err := cliCtx.GetNode()
	if err != nil {
		return request, err
	}

	blockResult, err := node.BlockResults(&requestHeight)
	if err != nil {
		return request, err
	}

	for _, tag := range blockResult.Results.EndBlock.Tags {
		if string(tag.Key) == sdk.ActionTag(service.ActionNewBatchRequest, contextID.String()) {
			var requests []service.CompactRequest
			err := json.Unmarshal(tag.GetValue(), &requests)
			if err != nil {
				return request, err
			}
			if len(requests) > int(batchRequestIndex) {
				compactRequest := requests[batchRequestIndex]
				request = service.NewRequest(
					requestID,
					requestContext.ServiceName,
					compactRequest.Provider,
					requestContext.Consumer,
					requestContext.Input,
					compactRequest.ServiceFee,
					requestContext.SuperMode,
					compactRequest.RequestHeight,
					compactRequest.RequestHeight+requestContext.Timeout,
					compactRequest.RequestContextID,
					compactRequest.RequestContextBatchCounter,
				)
				return request, nil
			}
		}
	}

	return request, nil
}

// QueryResponseByTxQuery will query for a single request via a direct txs tags query.
func QueryResponseByTxQuery(cliCtx context.CLIContext, params service.QueryResponseParams) (
	response service.Response, err error) {
	var tmTags []string
	tmTags = append(tmTags, fmt.Sprintf("%s='%s'", sdk.TagAction, service.TypeMsgRespondService))
	tmTags = append(tmTags, fmt.Sprintf("%s='%s'", service.TagRequestID, params.RequestID))

	result, err := tx.SearchTxs(cliCtx, cliCtx.Codec, tmTags, 1, 1)
	if err != nil {
		return response, err
	}

	if len(result.Txs) == 0 {
		return response, fmt.Errorf("unknown response: %s", params.RequestID)
	}

	requestID := params.RequestID

	contextID, batchCounter, _, _, err := service.SplitRequestID(requestID)
	if err != nil {
		return response, err
	}

	// query request context
	requestContext, err := QueryRequestContext(cliCtx, service.QueryRequestContextParams{
		RequestContextID: contextID,
	})

	if err != nil {
		return response, err
	}

	for _, msg := range result.Txs[0].Tx.GetMsgs() {
		if msg.Type() == service.TypeMsgRespondService {
			responseMsg := msg.(service.MsgRespondService)
			if responseMsg.RequestID.String() != params.RequestID.String() {
				continue
			}
			response := service.NewResponse(
				responseMsg.Provider, requestContext.Consumer,
				responseMsg.Result, responseMsg.Output,
				contextID, batchCounter,
			)
			return response, nil
		}
	}

	return response, nil
}

// QueryRequestsByBinding queries active requests by the service binding
func QueryRequestsByBinding(cliCtx context.CLIContext, serviceName string, provider sdk.AccAddress) (service.Requests, error) {
	params := service.QueryRequestsParams{
		ServiceName: serviceName,
		Provider:    provider,
	}

	bz, err := cliCtx.Codec.MarshalJSON(params)
	if err != nil {
		return nil, err
	}

	route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryRequests)
	res, err := cliCtx.QueryWithData(route, bz)
	if err != nil {
		return nil, err
	}

	var requests service.Requests
	if err := cliCtx.Codec.UnmarshalJSON(res, &requests); err != nil {
		return nil, err
	}

	return requests, nil
}

// QueryRequestsByReqCtx queries active requests by the request context ID
func QueryRequestsByReqCtx(cliCtx context.CLIContext, reqCtxIDStr, batchCounterStr string) (service.Requests, error) {
	requestContextID, err := hex.DecodeString(reqCtxIDStr)
	if err != nil {
		return nil, err
	}

	batchCounter, err := strconv.ParseUint(batchCounterStr, 10, 64)
	if err != nil {
		return nil, err
	}

	params := service.QueryRequestsByReqCtxParams{
		RequestContextID: requestContextID,
		BatchCounter:     batchCounter,
	}

	bz, err := cliCtx.Codec.MarshalJSON(params)
	if err != nil {
		return nil, err
	}

	route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryRequestsByReqCtx)
	res, err := cliCtx.QueryWithData(route, bz)
	if err != nil {
		return nil, err
	}

	var requests service.Requests
	if err := cliCtx.Codec.UnmarshalJSON(res, &requests); err != nil {
		return nil, err
	}

	return requests, nil
}
