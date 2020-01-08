package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/irisnet/irishub/modules/service/internal/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// define service
	r.HandleFunc("/service/definitions", defineServiceHandlerFn(cliCtx)).Methods("POST")
	// create service bind
	r.HandleFunc("/service/bindings", bindingAddHandlerFn(cliCtx)).Methods("POST")
	// update service bind
	r.HandleFunc(fmt.Sprintf("/service/bindings/{%s}/{%s}/{%s}", RestDefChainID, RestServiceName, RestProvider), bindingUpdateHandlerFn(cliCtx)).Methods("PUT")
	// disable service binding
	r.HandleFunc(fmt.Sprintf("/service/bindings/{%s}/{%s}/{%s}/disable", RestDefChainID, RestServiceName, RestProvider), bindingDisableHandlerFn(cliCtx)).Methods("PUT")
	// enable service binding
	r.HandleFunc(fmt.Sprintf("/service/bindings/{%s}/{%s}/{%s}/enable", RestDefChainID, RestServiceName, RestProvider), bindingEnableHandlerFn(cliCtx)).Methods("PUT")
	// refund all deposit from a service binding
	r.HandleFunc(fmt.Sprintf("/service/bindings/{%s}/{%s}/{%s}/deposit/refund", RestDefChainID, RestServiceName, RestProvider), bindingRefundHandlerFn(cliCtx)).Methods("PUT")
	// call a service method
	r.HandleFunc(fmt.Sprintf("/service/requests"), requestAddHandlerFn(cliCtx)).Methods("POST")
	// respond a service method invocation
	r.HandleFunc(fmt.Sprintf("/service/responses"), responseAddHandlerFn(cliCtx)).Methods("POST")
	// refund all fees from service call timeout
	r.HandleFunc(fmt.Sprintf("/service/fees/{%s}/refund", RestConsumer), feesRefundHandlerFn(cliCtx)).Methods("POST")
	// withdraw all fees from service call reward
	r.HandleFunc(fmt.Sprintf("/service/fees/{%s}/withdraw", RestProvider), feesWithdrawHandlerFn(cliCtx)).Methods("POST")
}

// HTTP request handler to define service.
func defineServiceHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DefineServiceReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		author, err := sdk.AccAddressFromBech32(req.Author)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgDefineService(req.Name, req.Description, req.Tags, author, req.AuthorDescription, req.Schemas)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// HTTP request handler to create service bind.
func bindingAddHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ServiceBindingReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		provider, err := sdk.AccAddressFromBech32(req.Provider)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		bindingType, err := types.BindingTypeFromString(req.BindingType)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		deposit, err := sdk.ParseCoins(req.Deposit)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var prices []sdk.Coin
		for _, p := range req.Prices {
			price, err := sdk.ParseCoin(p)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			prices = append(prices, price)
		}

		msg := types.NewMsgSvcBind(req.DefChainID, req.ServiceName, req.BaseReq.ChainID, provider, bindingType, deposit, prices, req.Level)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// HTTP request handler to update service bind.
func bindingUpdateHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		DefChainID := vars[RestDefChainID]
		serviceName := vars[RestServiceName]
		providerStr := vars[RestProvider]

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req ServiceBindingUpdateReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		var bindingType types.BindingType
		if req.BindingType != "" {
			bindingType, err = types.BindingTypeFromString(req.BindingType)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		var deposit sdk.Coins
		if req.Deposit != "" {
			deposit, err = sdk.ParseCoins(req.Deposit)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		var prices []sdk.Coin
		for _, p := range req.Prices {
			price, err := sdk.ParseCoin(p)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			prices = append(prices, price)
		}

		msg := types.NewMsgSvcBindingUpdate(DefChainID, serviceName, req.BaseReq.ChainID, provider, bindingType, deposit, prices, req.Level)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// HTTP request handler to disable service binding.
func bindingDisableHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		DefChainID := vars[RestDefChainID]
		serviceName := vars[RestServiceName]
		providerStr := vars[RestProvider]

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req BasicReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgSvcDisable(DefChainID, serviceName, req.BaseReq.ChainID, provider)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// HTTP request handler to enable service binding.
func bindingEnableHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		DefChainID := vars[RestDefChainID]
		serviceName := vars[RestServiceName]
		providerStr := vars[RestProvider]

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req ServiceBindingEnableReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		deposit, err := sdk.ParseCoins(req.Deposit)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgSvcEnable(DefChainID, serviceName, req.BaseReq.ChainID, provider, deposit)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// HTTP request handler to refund all deposit.
func bindingRefundHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		DefChainID := vars[RestDefChainID]
		serviceName := vars[RestServiceName]
		providerStr := vars[RestProvider]

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req BasicReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgSvcRefundDeposit(DefChainID, serviceName, req.BaseReq.ChainID, provider)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// HTTP request handler to call service method.
func requestAddHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ServiceRequestReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		var msgs []sdk.Msg
		for _, request := range req.Requests {
			consumer, err := sdk.AccAddressFromBech32(request.Consumer)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			provider, err := sdk.AccAddressFromBech32(request.Provider)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			input, err := hex.DecodeString(request.Data)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			serviceFee, err := sdk.ParseCoins(request.ServiceFee)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			msg := types.NewMsgSvcRequest(request.DefChainID, request.ServiceName, request.BindChainID, req.BaseReq.ChainID, consumer, provider, request.MethodID, input, serviceFee, request.Profiling)
			if err := msg.ValidateBasic(); err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			msgs = append(msgs, msg)
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, msgs)
	}
}

// HTTP request handler to respond service method invocation.
func responseAddHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ServiceResponseReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		provider, err := sdk.AccAddressFromBech32(req.Provider)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		output, err := hex.DecodeString(req.Data)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		errMsg, err := hex.DecodeString(req.ErrorMsg)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgSvcResponse(req.ReqChainID, req.RequestID, provider, output, errMsg)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// HTTP request handler to refund all fees.
func feesRefundHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		consumerStr := vars[RestConsumer]

		consumer, err := sdk.AccAddressFromBech32(consumerStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req BasicReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgSvcRefundFees(consumer)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// HTTP request handler to withdraw all fees.
func feesWithdrawHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		providerStr := vars[RestProvider]

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req BasicReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgSvcWithdrawFees(provider)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
