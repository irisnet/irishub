package lcd

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/service"
	sdk "github.com/irisnet/irishub/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {

	// get a single definition info
	r.HandleFunc(
		fmt.Sprintf("/service/definitions/{%s}/{%s}", DefChainId, ServiceName),
		definitionGetHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// get a single binding info
	r.HandleFunc(
		fmt.Sprintf("/service/bindings/{%s}/{%s}/{%s}/{%s}", DefChainId, ServiceName, BindChainId, Provider),
		bindingHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// get all bindings of a definition
	r.HandleFunc(
		fmt.Sprintf("/service/bindings/{%s}/{%s}", DefChainId, ServiceName),
		bindingsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// get all active requests of a binding
	r.HandleFunc(
		fmt.Sprintf("/service/requests/{%s}/{%s}/{%s}/{%s}", DefChainId, ServiceName, BindChainId, Provider),
		requestsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// get a single response
	r.HandleFunc(
		fmt.Sprintf("/service/responses/{%s}/{%s}", ReqChainId, ReqId),
		responseGetHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// get return fee and incoming fee of a account
	r.HandleFunc(
		fmt.Sprintf("/service/fees/{%s}", Address),
		feesHandlerFn(cliCtx, cdc),
	).Methods("GET")
}

func definitionGetHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		defChainId := vars[DefChainId]
		serviceName := vars[ServiceName]

		params := service.QueryServiceParams{
			DefChainID:  defChainId,
			ServiceName: serviceName,
		}

		bz, err := cdc.MarshalJSON(params)
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

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func bindingHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		defChainId := vars[DefChainId]
		serviceName := vars[ServiceName]
		bindChainId := vars[BindChainId]
		bechProviderAddr := vars[Provider]

		provider, err := sdk.AccAddressFromBech32(bechProviderAddr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := service.QueryBindingParams{
			DefChainID:  defChainId,
			ServiceName: serviceName,
			BindChainId: bindChainId,
			Provider:    provider,
		}

		bz, err := cdc.MarshalJSON(params)
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
		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func bindingsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		defChainId := vars[DefChainId]
		serviceName := vars[ServiceName]

		params := service.QueryServiceParams{
			DefChainID:  defChainId,
			ServiceName: serviceName,
		}

		bz, err := cdc.MarshalJSON(params)
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

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func requestsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		defChainId := vars[DefChainId]
		serviceName := vars[ServiceName]
		bindChainId := vars[BindChainId]
		bechProviderAddr := vars[Provider]

		provider, err := sdk.AccAddressFromBech32(bechProviderAddr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := service.QueryBindingParams{
			DefChainID:  defChainId,
			ServiceName: serviceName,
			BindChainId: bindChainId,
			Provider:    provider,
		}

		bz, err := cdc.MarshalJSON(params)
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

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func responseGetHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		reqChainId := vars[ReqChainId]
		reqId := vars[ReqId]

		params := service.QueryResponseParams{
			ReqChainId: reqChainId,
			RequestId:  reqId,
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

func feesHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bechAddress := vars[Address]

		address, err := sdk.AccAddressFromBech32(bechAddress)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := service.QueryFeesParams{
			Address: address,
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
