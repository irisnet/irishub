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

	"github.com/irisnet/irishub/modules/service"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/service/definitions", definitionPostHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/service/bindings", bindingAddHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/service/bindings/{%s}/{%s}/{%s}", RestDefChainID, RestServiceName, RestProvider), bindingUpdateHandlerFn(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/service/bindings/{%s}/{%s}/{%s}/disable", RestDefChainID, RestServiceName, RestProvider), bindingDisableHandlerFn(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/service/bindings/{%s}/{%s}/{%s}/enable", RestDefChainID, RestServiceName, RestProvider), bindingEnableHandlerFn(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/service/bindings/{%s}/{%s}/{%s}/deposit/refund", RestDefChainID, RestServiceName, RestProvider), bindingRefundHandlerFn(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/service/requests"), requestAddHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/service/responses"), responseAddHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/service/fees/{%s}/refund", RestConsumer), FeesRefundHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/service/fees/{%s}/withdraw", RestProvider), FeesWithdrawHandlerFn(cliCtx)).Methods("POST")
}

func definitionPostHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ServiceDefinitionReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		author, err := sdk.AccAddressFromBech32(req.AuthorAddr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := service.NewMsgSvcDef(req.ServiceName, req.BaseReq.ChainID, req.ServiceDescription, req.Tags, author, req.AuthorDescription, req.IdlContent)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

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

		bindingType, err := service.BindingTypeFromString(req.BindingType)
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

		msg := service.NewMsgSvcBind(req.DefChainId, req.ServiceName, req.BaseReq.ChainID, provider, bindingType, deposit, prices, req.Level)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

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

		var bindingType service.BindingType
		if req.BindingType != "" {
			bindingType, err = service.BindingTypeFromString(req.BindingType)
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

		msg := service.NewMsgSvcBindingUpdate(DefChainID, serviceName, req.BaseReq.ChainID, provider, bindingType, deposit, prices, req.Level)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

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

		msg := service.NewMsgSvcDisable(DefChainID, serviceName, req.BaseReq.ChainID, provider)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

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

		msg := service.NewMsgSvcEnable(DefChainID, serviceName, req.BaseReq.ChainID, provider, deposit)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

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

		msg := service.NewMsgSvcRefundDeposit(DefChainID, serviceName, req.BaseReq.ChainID, provider)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

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

			msg := service.NewMsgSvcRequest(request.DefChainId, request.ServiceName, request.BindChainId, req.BaseReq.ChainID, consumer, provider, request.MethodId, input, serviceFee, request.Profiling)
			if err := msg.ValidateBasic(); err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			msgs = append(msgs, msg)
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, msgs)
	}
}

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

		msg := service.NewMsgSvcResponse(req.ReqChainId, req.RequestId, provider, output, errMsg)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func FeesRefundHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
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

		msg := service.NewMsgSvcRefundFees(consumer)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func FeesWithdrawHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
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

		msg := service.NewMsgSvcWithdrawFees(provider)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
