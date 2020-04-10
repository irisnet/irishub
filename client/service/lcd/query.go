package lcd

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v3/service"
	"github.com/irisnet/irishub/client/context"
	serviceutils "github.com/irisnet/irishub/client/service/utils"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// query a service definition
	r.HandleFunc(
		fmt.Sprintf("/service/definitions/{%s}", ServiceName),
		queryDefinitionHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// query a service binding
	r.HandleFunc(
		fmt.Sprintf("/service/bindings/{%s}/{%s}", ServiceName, Provider),
		queryBindingHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// query all bindings of a service definition
	r.HandleFunc(
		fmt.Sprintf("/service/bindings/{%s}", ServiceName),
		queryBindingsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// query the withdrawal address of a provider
	r.HandleFunc(
		fmt.Sprintf("/service/providers/{%s}/withdraw-address", Provider),
		queryWithdrawAddrHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// query a request by ID
	r.HandleFunc(
		fmt.Sprintf("/service/requests/{%s}", RequestID),
		queryRequestHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// query active requests by the service binding or request context ID
	r.HandleFunc(
		fmt.Sprintf("/service/requests/{%s}/{%s}", Arg1, Arg2),
		queryRequestsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// query a response
	r.HandleFunc(
		fmt.Sprintf("/service/responses/{%s}", RequestID),
		queryResponseHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// query a request context
	r.HandleFunc(
		fmt.Sprintf("/service/contexts/{%s}", RequestContextID),
		queryRequestContextHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// query active responses by the request context ID and batch counter
	r.HandleFunc(
		fmt.Sprintf("/service/responses/{%s}/{%s}", RequestContextID, BatchCounter),
		queryResponsesHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// query the earned fees of a provider
	r.HandleFunc(
		fmt.Sprintf("/service/fees/{%s}", Provider),
		queryEarnedFeesHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// query the system schema by the schema name
	r.HandleFunc(
		fmt.Sprintf("/service/schemas/{%s}", SchemaName),
		querySchemaHandlerFn(cliCtx, cdc),
	).Methods("GET")
}

func queryDefinitionHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[ServiceName]

		params := service.QueryDefinitionParams{
			ServiceName: serviceName,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryDefinition)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cliCtx.Codec, res, cliCtx.Indent)
	}
}

func queryBindingHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[ServiceName]
		providerStr := vars[Provider]

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := service.QueryBindingParams{
			ServiceName: serviceName,
			Provider:    provider,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryBinding)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cliCtx.Codec, res, cliCtx.Indent)
	}
}

func queryBindingsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[ServiceName]

		params := service.QueryBindingsParams{
			ServiceName: serviceName,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryBindings)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cliCtx.Codec, res, cliCtx.Indent)
	}
}

func queryWithdrawAddrHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		providerStr := vars[Provider]

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := service.QueryWithdrawAddressParams{
			Provider: provider,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryWithdrawAddress)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cliCtx.Codec, res, cliCtx.Indent)
	}
}

func queryRequestHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestIDStr := vars[RequestID]

		requestID, err := service.ConvertRequestID(requestIDStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		params := service.QueryRequestParams{
			RequestID: requestID,
		}

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryRequest)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var request service.Request
		_ = cdc.UnmarshalJSON(res, &request)
		if request.Empty() {
			request, err = serviceutils.QueryRequestByTxQuery(cliCtx, params)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		if request.Empty() {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("unknown request: %s", params.RequestID))
			return
		}

		utils.PostProcessResponse(w, cdc, request, cliCtx.Indent)
	}
}

func queryRequestsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		arg1 := vars[Arg1]
		arg2 := vars[Arg2]

		queryByBinding := true

		provider, err := sdk.AccAddressFromBech32(arg2)
		if err != nil {
			queryByBinding = false
		}

		var requests service.Requests

		if queryByBinding {
			requests, err = serviceutils.QueryRequestsByBinding(cliCtx, arg1, provider)
		} else {
			requests, err = serviceutils.QueryRequestsByReqCtx(cliCtx, arg1, arg2)
		}

		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, requests, cliCtx.Indent)
	}
}

func queryResponseHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestIDStr := vars[RequestID]

		requestID, err := service.ConvertRequestID(requestIDStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := service.QueryResponseParams{
			RequestID: requestID,
		}

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryResponse)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var response service.Response
		_ = cdc.UnmarshalJSON(res, &response)
		if response.Empty() {
			response, err = serviceutils.QueryResponseByTxQuery(cliCtx, params)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		if response.Empty() {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("unknown request: %s", params.RequestID))
			return
		}

		utils.PostProcessResponse(w, cdc, response, cliCtx.Indent)
	}
}

func queryRequestContextHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestContextIDStr := vars[RequestContextID]

		requestContextID, err := hex.DecodeString(requestContextIDStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := service.QueryRequestContextParams{
			RequestContextID: requestContextID,
		}

		requestContext, err := serviceutils.QueryRequestContext(cliCtx, params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, requestContext, cliCtx.Indent)
	}
}

func queryResponsesHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestContextIDStr := vars[RequestContextID]
		batchCounterStr := vars[BatchCounter]

		requestContextID, err := hex.DecodeString(requestContextIDStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		batchCounter, err := strconv.ParseUint(batchCounterStr, 10, 64)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := service.QueryResponsesParams{
			RequestContextID: requestContextID,
			BatchCounter:     batchCounter,
		}

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryResponses)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryEarnedFeesHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		providerStr := vars[Provider]

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := service.QueryEarnedFeesParams{
			Provider: provider,
		}

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryEarnedFees)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func querySchemaHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		params := service.QuerySchemaParams{
			SchemaName: vars[SchemaName],
		}

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QuerySchema)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var schema serviceutils.SchemaType
		if err := cdc.UnmarshalJSON(res, &schema); err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, schema, cliCtx.Indent)
	}
}
