package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/irisnet/irismod/modules/service/types"
)

func registerTxRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc("/service/definitions", defineServiceHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/service/bindings", bindServiceHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/service/bindings/{%s}/{%s}", RestServiceName, RestProvider), updateServiceBindingHandlerFn(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/service/owners/{%s}/withdraw-address", RestOwner), setWithdrawAddrHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/service/bindings/{%s}/{%s}/disable", RestServiceName, RestProvider), disableServiceBindingHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/service/bindings/{%s}/{%s}/enable", RestServiceName, RestProvider), enableServiceBindingHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/service/bindings/{%s}/{%s}/refund-deposit", RestServiceName, RestProvider), refundServiceDepositHandlerFn(cliCtx)).Methods("POST")
	// initiate a service call
	r.HandleFunc("/service/contexts", requestServiceHandlerFn(cliCtx)).Methods("POST")
	// respond to a service request
	r.HandleFunc("/service/responses", respondServiceHandlerFn(cliCtx)).Methods("POST")
	// pause a request context
	r.HandleFunc(fmt.Sprintf("/service/contexts/{%s}/pause", RestRequestContextID), pauseRequestContextHandlerFn(cliCtx)).Methods("POST")
	// start a paused request context
	r.HandleFunc(fmt.Sprintf("/service/contexts/{%s}/start", RestRequestContextID), startRequestContextHandlerFn(cliCtx)).Methods("POST")
	// kill a request context
	r.HandleFunc(fmt.Sprintf("/service/contexts/{%s}/kill", RestRequestContextID), killRequestContextHandlerFn(cliCtx)).Methods("POST")
	// update a request context
	r.HandleFunc(fmt.Sprintf("/service/contexts/{%s}", RestRequestContextID), updateRequestContextHandlerFn(cliCtx)).Methods("PUT")
	// withdraw the earned fees of a provider
	r.HandleFunc(fmt.Sprintf("/service/fees/{%s}/withdraw", RestProvider), withdrawEarnedFeesHandlerFn(cliCtx)).Methods("POST")
}

func defineServiceHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DefineServiceReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		if _, err := sdk.AccAddressFromBech32(req.Author); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgDefineService(req.Name, req.Description, req.Tags, req.Author, req.AuthorDescription, req.Schemas)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func bindServiceHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BindServiceReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		if _, err := sdk.AccAddressFromBech32(req.Owner); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		provider := req.Owner
		if len(req.Provider) > 0 {
			if _, err := sdk.AccAddressFromBech32(req.Provider); err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			provider = req.Provider
		}

		deposit, err := sdk.ParseCoins(req.Deposit)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgBindService(req.ServiceName, provider, deposit, req.Pricing, req.QoS, req.Options, req.Owner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)

	}
}

func updateServiceBindingHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[RestServiceName]
		provider := vars[RestProvider]

		if _, err := sdk.AccAddressFromBech32(provider); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req UpdateServiceBindingReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		owner := provider
		if len(req.Owner) > 0 {
			if _, err := sdk.AccAddressFromBech32(req.Owner); err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			owner = req.Owner
		}

		var deposit sdk.Coins
		var err error
		if req.Deposit != "" {
			deposit, err = sdk.ParseCoins(req.Deposit)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		msg := types.NewMsgUpdateServiceBinding(serviceName, provider, deposit, req.Pricing, req.QoS, req.Options, owner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)

	}
}

func setWithdrawAddrHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		owner := vars[RestOwner]

		if _, err := sdk.AccAddressFromBech32(owner); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req SetWithdrawAddrReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		if _, err := sdk.AccAddressFromBech32(req.WithdrawAddress); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgSetWithdrawAddress(owner, req.WithdrawAddress)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)

	}
}

func disableServiceBindingHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[RestServiceName]
		provider := vars[RestProvider]

		if _, err := sdk.AccAddressFromBech32(provider); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req DisableServiceBindingReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		owner := provider
		if len(req.Owner) > 0 {
			if _, err := sdk.AccAddressFromBech32(req.Owner); err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			owner = req.Owner
		}

		msg := types.NewMsgDisableServiceBinding(serviceName, provider, owner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func enableServiceBindingHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[RestServiceName]
		provider := vars[RestProvider]

		if _, err := sdk.AccAddressFromBech32(provider); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req EnableServiceBindingReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		owner := provider
		if len(req.Owner) > 0 {
			if _, err := sdk.AccAddressFromBech32(req.Owner); err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			owner = req.Owner
		}

		var deposit sdk.Coins
		var err error
		if len(req.Deposit) != 0 {
			deposit, err = sdk.ParseCoins(req.Deposit)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		msg := types.NewMsgEnableServiceBinding(serviceName, provider, deposit, owner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func refundServiceDepositHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[RestServiceName]
		provider := vars[RestProvider]

		if _, err := sdk.AccAddressFromBech32(provider); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req RefundServiceDepositReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		owner := provider
		if len(req.Owner) > 0 {
			if _, err := sdk.AccAddressFromBech32(req.Owner); err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			owner = req.Owner
		}

		msg := types.NewMsgRefundServiceDeposit(serviceName, provider, owner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func requestServiceHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req callServiceReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		if _, err := sdk.AccAddressFromBech32(req.Consumer); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		serviceFeeCap, err := sdk.ParseCoins(req.ServiceFeeCap)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		for _, p := range req.Providers {
			if _, err := sdk.AccAddressFromBech32(p); err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		msg := types.NewMsgCallService(
			req.ServiceName, req.Providers, req.Consumer, req.Input, serviceFeeCap,
			req.Timeout, req.SuperMode, req.Repeated, req.RepeatedFrequency, req.RepeatedTotal,
		)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func respondServiceHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req respondServiceReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		if _, err := types.ConvertRequestID(req.RequestID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if _, err := sdk.AccAddressFromBech32(req.Provider); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgRespondService(req.RequestID, req.Provider, req.Result, req.Output)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func pauseRequestContextHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestContextID := vars[RestRequestContextID]

		if _, err := hex.DecodeString(requestContextID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req pauseRequestContextReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		if _, err := sdk.AccAddressFromBech32(req.Consumer); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgPauseRequestContext(requestContextID, req.Consumer)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func startRequestContextHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestContextID := vars[RestRequestContextID]

		if _, err := hex.DecodeString(requestContextID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req startRequestContextReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		if _, err := sdk.AccAddressFromBech32(req.Consumer); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgStartRequestContext(requestContextID, req.Consumer)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func killRequestContextHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestContextID := vars[RestRequestContextID]

		if _, err := hex.DecodeString(requestContextID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req killRequestContextReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		if _, err := sdk.AccAddressFromBech32(req.Consumer); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgKillRequestContext(requestContextID, req.Consumer)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func updateRequestContextHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestContextID := vars[RestRequestContextID]

		if _, err := hex.DecodeString(requestContextID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req updateRequestContextReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		if _, err := sdk.AccAddressFromBech32(req.Consumer); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var serviceFeeCap sdk.Coins
		var err error
		if len(req.ServiceFeeCap) != 0 {
			serviceFeeCap, err = sdk.ParseCoins(req.ServiceFeeCap)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		for _, p := range req.Providers {
			if _, err := sdk.AccAddressFromBech32(p); err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		msg := types.NewMsgUpdateRequestContext(
			requestContextID, req.Providers, serviceFeeCap, req.Timeout,
			req.RepeatedFrequency, req.RepeatedTotal, req.Consumer,
		)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func withdrawEarnedFeesHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		provider := vars[RestProvider]

		if _, err := sdk.AccAddressFromBech32(provider); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req withdrawEarnedFeesReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		owner := provider
		if len(req.Owner) > 0 {
			if _, err := sdk.AccAddressFromBech32(req.Owner); err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			owner = req.Owner
		}

		msg := types.NewMsgWithdrawEarnedFees(owner, provider)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}
