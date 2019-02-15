package lcd

import (
	"net/http"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/modules/service"
	"github.com/irisnet/irishub/client/utils"
	cmn "github.com/irisnet/irishub/client/service"
	sdk "github.com/irisnet/irishub/types"
)

const storeName = "service"

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {

	// get a single definition info
	r.HandleFunc(
		fmt.Sprintf("/service/definition/{%s}/{%s}", DefChainId, ServiceName),
		definitionGetHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// get a single binding info
	r.HandleFunc(
		fmt.Sprintf("/service/binding/{%s}/{%s}/{%s}/{%s}", DefChainId, ServiceName, BindChainId, Provider),
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
		fmt.Sprintf("/service/response/{%s}/{%s}", ReqChainId, ReqId),
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

		res, err := cliCtx.QueryStore(service.GetServiceDefinitionKey(defChainId, serviceName), storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if len(res) == 0 {
			utils.WriteErrorResponse(w, http.StatusNoContent, "")
			return
		}

		var svcDef service.SvcDef
		cdc.MustUnmarshalBinaryLengthPrefixed(res, &svcDef)

		res1, err := cliCtx.QuerySubspace(service.GetMethodsSubspaceKey(defChainId, serviceName), storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var methods []service.MethodProperty
		for _, re := range res1 {
			var method service.MethodProperty
			cdc.MustUnmarshalBinaryLengthPrefixed(re.Value, &method)
			methods = append(methods, method)
		}

		output, err := codec.MarshalJSONIndent(cdc, cmn.DefOutput{Definition: svcDef, Methods: methods})
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.PostProcessResponse(w, cdc, output, cliCtx.Indent)
	}
}

func bindingHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		defChainId := vars[DefChainId]
		serviceName := vars[ServiceName]
		bindChainId := vars[BindChainId]
		bechProviderAddr := vars[Provider]

		providerAddr, err := sdk.AccAddressFromBech32(bechProviderAddr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		res, err := cliCtx.QueryStore(service.GetServiceBindingKey(defChainId, serviceName, bindChainId, providerAddr), storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if len(res) == 0 {
			utils.WriteErrorResponse(w, http.StatusNoContent, "")
			return
		}
		var svcBinding service.SvcBinding
		cdc.MustUnmarshalBinaryLengthPrefixed(res, &svcBinding)
		output, err := codec.MarshalJSONIndent(cdc, svcBinding)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.PostProcessResponse(w, cdc, output, cliCtx.Indent)
	}
}

func bindingsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		defChainId := vars[DefChainId]
		serviceName := vars[ServiceName]

		res, err := cliCtx.QuerySubspace(service.GetBindingsSubspaceKey(defChainId, serviceName), storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if len(res) == 0 {
			utils.WriteErrorResponse(w, http.StatusNoContent, "")
			return
		}
		var bindings []service.SvcBinding
		for _, re := range res {
			var binding service.SvcBinding
			cdc.MustUnmarshalBinaryLengthPrefixed(re.Value, &binding)
			bindings = append(bindings, binding)
		}

		output, err := codec.MarshalJSONIndent(cdc, bindings)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.PostProcessResponse(w, cdc, output, cliCtx.Indent)
	}
}

func requestsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		defChainId := vars[DefChainId]
		serviceName := vars[ServiceName]
		bindChainId := vars[BindChainId]
		bechProviderAddr := vars[Provider]

		providerAddr, err := sdk.AccAddressFromBech32(bechProviderAddr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QuerySubspace(service.GetSubActiveRequestKey(defChainId, serviceName, bindChainId, providerAddr), storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if len(res) == 0 {
			utils.WriteErrorResponse(w, http.StatusNoContent, "")
			return
		}
		var reqs []service.SvcRequest
		for _, re := range res {
			var req service.SvcRequest
			cdc.MustUnmarshalBinaryLengthPrefixed(re.Value, &req)
			reqs = append(reqs, req)
		}

		output, err := codec.MarshalJSONIndent(cdc, reqs)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.PostProcessResponse(w, cdc, output, cliCtx.Indent)
	}
}

func responseGetHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		reqChainId := vars[ReqChainId]
		reqId := vars[ReqId]

		eHeight, rHeight, counter, err := service.ConvertRequestID(reqId)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryStore(service.GetResponseKey(reqChainId, eHeight, rHeight, counter), storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if len(res) == 0 {
			utils.WriteErrorResponse(w, http.StatusNoContent, "")
			return
		}
		var resp service.SvcResponse
		cdc.MustUnmarshalBinaryLengthPrefixed(res, &resp)
		output, err := codec.MarshalJSONIndent(cdc, resp)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.PostProcessResponse(w, cdc, output, cliCtx.Indent)
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

		res, err := cliCtx.QueryStore(service.GetReturnedFeeKey(address), storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		var returnedFee service.ReturnedFee
		if len(res) > 0 {
			cdc.MustUnmarshalBinaryLengthPrefixed(res, &returnedFee)
		}

		res1, err := cliCtx.QueryStore(service.GetIncomingFeeKey(address), storeName)
		var incomingFee service.IncomingFee
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if len(res1) > 0 {
			cdc.MustUnmarshalBinaryLengthPrefixed(res1, &incomingFee)
		}

		output, err := codec.MarshalJSONIndent(cdc, cmn.FeesOutput{ReturnedFee: returnedFee.Coins, IncomingFee: incomingFee.Coins})
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.PostProcessResponse(w, cdc, output, cliCtx.Indent)
	}
}
