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

	// query all active requests of a service binding
	r.HandleFunc(
		fmt.Sprintf("/service/requests/{%s}/{%s}", ServiceName, Provider),
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
		qeuryRequestContextHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// query requests by the request context ID and batch counter
	r.HandleFunc(
		fmt.Sprintf("/service/requests/{%s}/{%s}", RequestContextID, BatchCounter),
		queryRequestsByReqCtxHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// query responses by the request context ID and batch counter
	r.HandleFunc(
		fmt.Sprintf("/service/responses/{%s}/{%s}", RequestContextID, BatchCounter),
		queryResponsesHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// query the earned fees
	r.HandleFunc(
		fmt.Sprintf("/service/fees/{%s}", Provider),
		queryEarnedFeesHandlerFn(cliCtx, cdc),
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

func queryRequestsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[ServiceName]
		providerStr := vars[Provider]

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := service.QueryRequestsParams{
			ServiceName: serviceName,
			Provider:    provider,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryRequests)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cliCtx.Codec, res, cliCtx.Indent)
	}
}

func queryResponseHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestID := vars[RequestID]

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

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func qeuryRequestContextHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
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

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryRequestContext)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func queryRequestsByReqCtxHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
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

		params := service.QueryRequestsByReqCtxParams{
			RequestContextID: requestContextID,
			BatchCounter:     batchCounter,
		}

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryRequestsByReqCtx)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
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

		params := service.QueryFeesParams{
			Address: provider,
		}

		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", protocol.ServiceRoute, service.QueryFees)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
