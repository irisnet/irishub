package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	serviceutils "github.com/irisnet/irismod/modules/service/client/utils"
	"github.com/irisnet/irismod/modules/service/types"
)

func registerQueryRoutes(cliCtx client.Context, r *mux.Router) {
	// query a service definition
	r.HandleFunc(fmt.Sprintf("/%s/definitions/{%s}", types.ModuleName, RestServiceName), queryDefinitionHandlerFn(cliCtx)).Methods("GET")
	// query a service binding
	r.HandleFunc(fmt.Sprintf("/%s/bindings/{%s}/{%s}", types.ModuleName, RestServiceName, RestProvider), queryBindingHandlerFn(cliCtx)).Methods("GET")
	// query all bindings of a service definition with an optional owner
	r.HandleFunc(fmt.Sprintf("/%s/bindings/{%s}", types.ModuleName, RestServiceName), queryBindingsHandlerFn(cliCtx)).Methods("GET")
	// query the withdrawal address of an owner
	r.HandleFunc(fmt.Sprintf("/%s/owners/{%s}/withdraw-address", types.ModuleName, RestOwner), queryWithdrawAddrHandlerFn(cliCtx)).Methods("GET")
	// query a request by ID
	r.HandleFunc(fmt.Sprintf("/%s/requests/{%s}", types.ModuleName, RestRequestID), queryRequestHandlerFn(cliCtx)).Methods("GET")
	// query active requests by the service binding or request context ID
	r.HandleFunc(fmt.Sprintf("/%s/requests/{%s}/{%s}", types.ModuleName, RestArg1, RestArg2), queryRequestsHandlerFn(cliCtx)).Methods("GET")
	// query a response
	r.HandleFunc(fmt.Sprintf("/%s/responses/{%s}", types.ModuleName, RestRequestID), queryResponseHandlerFn(cliCtx)).Methods("GET")
	// query a request context
	r.HandleFunc(fmt.Sprintf("/%s/contexts/{%s}", types.ModuleName, RestRequestContextID), queryRequestContextHandlerFn(cliCtx)).Methods("GET")
	// query active responses by the request context ID and batch counter
	r.HandleFunc(fmt.Sprintf("/%s/responses/{%s}/{%s}", types.ModuleName, RestRequestContextID, RestBatchCounter), queryResponsesHandlerFn(cliCtx)).Methods("GET")
	// query the earned fees of a provider
	r.HandleFunc(fmt.Sprintf("/%s/fees/{%s}", types.ModuleName, RestProvider), queryEarnedFeesHandlerFn(cliCtx)).Methods("GET")
	// query the system schema by the schema name
	r.HandleFunc(fmt.Sprintf("/%s/schemas/{%s}", types.ModuleName, RestSchemaName), querySchemaHandlerFn(cliCtx)).Methods("GET")
	// query the current service parameter values
	r.HandleFunc(fmt.Sprintf("/%s/params", types.ModuleName), queryParamsHandlerFn(cliCtx)).Methods("GET")
}

func queryDefinitionHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[RestServiceName]

		if err := types.ValidateServiceName(serviceName); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QueryDefinitionParams{
			ServiceName: serviceName,
		}

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDefinition)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryBindingHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[RestServiceName]
		providerStr := vars[RestProvider]

		if err := types.ValidateServiceName(serviceName); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QueryBindingParams{
			ServiceName: serviceName,
			Provider:    provider,
		}

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QueryBinding)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryBindingsHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[RestServiceName]
		ownerStr := r.FormValue("owner")

		if err := types.ValidateServiceName(serviceName); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var err error
		var owner sdk.AccAddress

		if len(ownerStr) > 0 {
			owner, err = sdk.AccAddressFromBech32(ownerStr)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QueryBindingsParams{
			ServiceName: serviceName,
			Owner:       owner,
		}

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QueryBindings)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryWithdrawAddrHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ownerStr := vars[RestOwner]

		owner, err := sdk.AccAddressFromBech32(ownerStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QueryWithdrawAddressParams{
			Owner: owner,
		}

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QueryWithdrawAddress)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryRequestHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestIDStr := vars[RestRequestID]

		requestID, err := types.ConvertRequestID(requestIDStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QueryRequestParams{
			RequestID: requestID,
		}

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QueryRequest)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var request types.Request
		_ = cliCtx.LegacyAmino.UnmarshalJSON(res, &request)
		if request.Empty() {
			request, err = serviceutils.QueryRequestByTxQuery(cliCtx, types.RouterKey, params.RequestID)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		if request.Empty() {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("unknown request: %s", params.RequestID))
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryRequestsHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		arg1 := vars[RestArg1]
		arg2 := vars[RestArg2]

		queryByBinding := true

		provider, err := sdk.AccAddressFromBech32(arg2)
		if err != nil {
			queryByBinding = false
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		var res []types.Request
		var height int64

		if queryByBinding {
			res, height, err = serviceutils.QueryRequestsByBinding(cliCtx, types.RouterKey, arg1, provider)
		} else {
			res, height, err = serviceutils.QueryRequestsByReqCtx(cliCtx, types.RouterKey, arg1, arg2)
		}

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryResponseHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestIDStr := vars[RestRequestID]

		requestID, err := types.ConvertRequestID(requestIDStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QueryResponseParams{
			RequestID: requestID,
		}

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QueryResponse)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var response types.Response
		_ = cliCtx.LegacyAmino.UnmarshalJSON(res, &response)
		if response.Empty() {
			response, err = serviceutils.QueryResponseByTxQuery(cliCtx, types.RouterKey, params.RequestID)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		if response.Empty() {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("unknown request: %s", params.RequestID))
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, response)
	}
}

func queryRequestContextHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestContextIDStr := vars[RestRequestContextID]

		requestContextID, err := hex.DecodeString(requestContextIDStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QueryRequestContextRequest{
			RequestContextId: tmbytes.HexBytes(requestContextID).String(),
		}

		requestContext, err := serviceutils.QueryRequestContext(cliCtx, types.RouterKey, params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, requestContext)
	}
}

func queryResponsesHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestContextIDStr := vars[RestRequestContextID]
		batchCounterStr := vars[RestBatchCounter]

		requestContextID, err := hex.DecodeString(requestContextIDStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		batchCounter, err := strconv.ParseUint(batchCounterStr, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QueryResponsesParams{
			RequestContextID: requestContextID,
			BatchCounter:     batchCounter,
		}

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QueryResponses)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryEarnedFeesHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		providerStr := vars[RestProvider]

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QueryEarnedFeesParams{
			Provider: provider,
		}

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QueryEarnedFees)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func querySchemaHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QuerySchemaParams{
			SchemaName: vars[RestSchemaName],
		}

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QuerySchema)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var schema serviceutils.SchemaType
		if err := cliCtx.LegacyAmino.UnmarshalJSON(res, &schema); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, schema)
	}
}

func queryParamsHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParameters)
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
