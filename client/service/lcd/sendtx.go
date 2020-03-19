package lcd

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/irisnet/irishub/app/v3/service"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// define a service
	r.HandleFunc(
		"/service/definitions",
		defineServiceHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// bind a service
	r.HandleFunc(
		"/service/bindings",
		bindServiceHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// update a service binding
	r.HandleFunc(
		fmt.Sprintf("/service/bindings/{%s}/{%s}", ServiceName, Provider),
		updateServiceBindingHandlerFn(cdc, cliCtx),
	).Methods("PUT")

	// set a new withdrawal address for a provider
	r.HandleFunc(
		fmt.Sprintf("/service/providers/{%s}/withdraw-address", Provider),
		setWithdrawAddrHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// disable a service binding
	r.HandleFunc(
		fmt.Sprintf("/service/bindings/{%s}/{%s}/disable", ServiceName, Provider),
		disableServiceHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// enable a service binding
	r.HandleFunc(
		fmt.Sprintf("/service/bindings/{%s}/{%s}/enable", ServiceName, Provider),
		enableServiceHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// refund deposit from a service binding
	r.HandleFunc(
		fmt.Sprintf("/service/bindings/{%s}/{%s}/refund-deposit", ServiceName, Provider),
		refundServiceDepositHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// call a service
	r.HandleFunc(
		"/service/requests",
		requestServiceHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// respond to a service request
	r.HandleFunc(
		"/service/responses",
		respondServiceHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// pause a request context
	r.HandleFunc(
		fmt.Sprintf("/service/contexts/{%s}/pause", RequestContextID),
		pauseRequestContextHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// start a paused request context
	r.HandleFunc(
		fmt.Sprintf("/service/contexts/{%s}/start", RequestContextID),
		startRequestContextHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// kill a request context
	r.HandleFunc(
		fmt.Sprintf("/service/contexts/{%s}/kill", RequestContextID),
		killRequestContextHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// update a request context
	r.HandleFunc(
		fmt.Sprintf("/service/contexts/{%s}", RequestContextID),
		updateRequestContextHandlerFn(cdc, cliCtx),
	).Methods("PUT")

	// withdraw the earned fees of a provider
	r.HandleFunc(
		fmt.Sprintf("/service/fees/{%s}/withdraw", Provider),
		withdrawEarnedFeesHandlerFn(cdc, cliCtx),
	).Methods("POST")
}

type defineServiceReq struct {
	BaseTx            utils.BaseTx `json:"base_tx"` // basic tx info
	Name              string       `json:"name"`
	Description       string       `json:"description"`
	Tags              []string     `json:"tags"`
	Author            string       `json:"author"`
	AuthorDescription string       `json:"author_description"`
	Schemas           string       `json:"schemas"`
}

type bindServiceReq struct {
	BaseTx      utils.BaseTx `json:"base_tx"` // basic tx info
	ServiceName string       `json:"service_name"`
	Provider    string       `json:"provider"`
	Deposit     string       `json:"deposit"`
	Pricing     string       `json:"pricing"`
}

type updateServiceBindingReq struct {
	BaseTx  utils.BaseTx `json:"base_tx"` // basic tx info
	Deposit string       `json:"deposit"`
	Pricing string       `json:"pricing"`
}

type setWithdrawAddrReq struct {
	BaseTx          utils.BaseTx `json:"base_tx"` // basic tx info
	WithdrawAddress string       `json:"withdraw_address"`
}

type disableServiceReq struct {
	BaseTx utils.BaseTx `json:"base_tx"` // basic tx info`
}

type enableServiceReq struct {
	BaseTx  utils.BaseTx `json:"base_tx"` // basic tx info
	Deposit string       `json:"deposit"`
}

type refundServiceDepositReq struct {
	BaseTx utils.BaseTx `json:"base_tx"` // basic tx info`
}

type requestServiceReq struct {
	BaseTx            utils.BaseTx `json:"base_tx"` // basic tx info
	ServiceName       string       `json:"service_name"`
	Providers         []string     `json:"providers"`
	Consumer          string       `json:"consumer"`
	Input             string       `json:"input"`
	ServiceFeeCap     string       `json:"service_fee_cap"`
	Timeout           int64        `json:"timeout"`
	SuperMode         bool         `json:"super_mode"`
	Repeated          bool         `json:"repeated"`
	RepeatedFrequency uint64       `json:"repeated_frequency"`
	RepeatedTotal     int64        `json:"repeated_total"`
}

type respondServiceReq struct {
	BaseTx    utils.BaseTx `json:"base_tx"` // basic tx info
	RequestID string       `json:"request_id"`
	Provider  string       `json:"provider"`
	Result    string       `json:"result"`
	Output    string       `json:"output"`
}

type pauseRequestContextReq struct {
	BaseTx   utils.BaseTx `json:"base_tx"` // basic tx info
	Consumer string       `json:"consumer"`
}

type startRequestContextReq struct {
	BaseTx   utils.BaseTx `json:"base_tx"` // basic tx info
	Consumer string       `json:"consumer"`
}

type killRequestContextReq struct {
	BaseTx   utils.BaseTx `json:"base_tx"` // basic tx info
	Consumer string       `json:"consumer"`
}

type updateRequestContextReq struct {
	BaseTx            utils.BaseTx `json:"base_tx"` // basic tx info
	Providers         []string     `json:"providers"`
	ServiceFeeCap     string       `json:"service_fee_cap"`
	Timeout           int64        `json:"timeout"`
	RepeatedFrequency uint64       `json:"repeated_frequency"`
	RepeatedTotal     int64        `json:"repeated_total"`
	Consumer          string       `json:"consumer"`
}

type withdrawEarnedFeesReq struct {
	BaseTx utils.BaseTx `json:"base_tx"` // basic tx info
}

// HTTP request handler to define a service.
func defineServiceHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req defineServiceReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		author, err := sdk.AccAddressFromBech32(req.Author)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := service.NewMsgDefineService(req.Name, req.Description, req.Tags, author, req.AuthorDescription, req.Schemas)
		if err := msg.ValidateBasic(); err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

// HTTP request handler to bind a service.
func bindServiceHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req bindServiceReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		provider, err := sdk.AccAddressFromBech32(req.Provider)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		deposit, err := cliCtx.ParseCoins(req.Deposit)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := service.NewMsgBindService(req.ServiceName, provider, deposit, req.Pricing)
		if err := msg.ValidateBasic(); err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

// HTTP request handler to update a service binding.
func updateServiceBindingHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[ServiceName]
		providerStr := vars[Provider]

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req updateServiceBindingReq
		err = utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		var deposit sdk.Coins
		if req.Deposit != "" {
			deposit, err = cliCtx.ParseCoins(req.Deposit)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		msg := service.NewMsgUpdateServiceBinding(serviceName, provider, deposit, req.Pricing)
		if err := msg.ValidateBasic(); err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

// HTTP request handler to set a withdrawal address for a provider.
func setWithdrawAddrHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		providerStr := vars[Provider]

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req setWithdrawAddrReq
		err = utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		withdrawAddr, err := sdk.AccAddressFromBech32(req.WithdrawAddress)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := service.NewMsgSetWithdrawAddress(provider, withdrawAddr)
		if err := msg.ValidateBasic(); err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

// HTTP request handler to disable a service.
func disableServiceHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[ServiceName]
		providerStr := vars[Provider]

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req disableServiceReq
		err = utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		msg := service.NewMsgDisableService(serviceName, provider)
		if err := msg.ValidateBasic(); err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

// HTTP request handler to enable a service.
func enableServiceHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[ServiceName]
		providerStr := vars[Provider]

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req enableServiceReq
		err = utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		var deposit sdk.Coins
		if len(req.Deposit) != 0 {
			deposit, err = cliCtx.ParseCoins(req.Deposit)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		msg := service.NewMsgEnableService(serviceName, provider, deposit)
		if err := msg.ValidateBasic(); err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

// HTTP request handler to refund deposit from a service binding
func refundServiceDepositHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[ServiceName]
		providerStr := vars[Provider]

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req refundServiceDepositReq
		err = utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		msg := service.NewMsgRefundServiceDeposit(serviceName, provider)
		if err = msg.ValidateBasic(); err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

func requestServiceHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req requestServiceReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		consumer, err := sdk.AccAddressFromBech32(req.Consumer)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		serviceFeeCap, err := cliCtx.ParseCoins(req.ServiceFeeCap)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var providers []sdk.AccAddress
		for _, p := range req.Providers {
			provider, err := sdk.AccAddressFromBech32(p)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			providers = append(providers, provider)
		}

		msg := service.NewMsgRequestService(
			req.ServiceName, providers, consumer, req.Input, serviceFeeCap,
			req.Timeout, req.SuperMode, req.Repeated, req.RepeatedFrequency, req.RepeatedTotal,
		)
		if err = msg.ValidateBasic(); err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

func respondServiceHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req respondServiceReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		provider, err := sdk.AccAddressFromBech32(req.Provider)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := service.NewMsgRespondService(req.RequestID, provider, req.Result, req.Output)
		if err = msg.ValidateBasic(); err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

func pauseRequestContextHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestContextIDStr := vars[RequestContextID]

		requestContextID, err := hex.DecodeString(requestContextIDStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req pauseRequestContextReq
		err = utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		consumer, err := sdk.AccAddressFromBech32(req.Consumer)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := service.NewMsgPauseRequestContext(requestContextID, consumer)
		if err = msg.ValidateBasic(); err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

func startRequestContextHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestContextIDStr := vars[RequestContextID]

		requestContextID, err := hex.DecodeString(requestContextIDStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req startRequestContextReq
		err = utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		consumer, err := sdk.AccAddressFromBech32(req.Consumer)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := service.NewMsgStartRequestContext(requestContextID, consumer)
		if err = msg.ValidateBasic(); err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

func killRequestContextHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestContextIDStr := vars[RequestContextID]

		requestContextID, err := hex.DecodeString(requestContextIDStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req killRequestContextReq
		err = utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		consumer, err := sdk.AccAddressFromBech32(req.Consumer)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := service.NewMsgKillRequestContext(requestContextID, consumer)
		if err = msg.ValidateBasic(); err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

func updateRequestContextHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestContextIDStr := vars[RequestContextID]

		requestContextID, err := hex.DecodeString(requestContextIDStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req updateRequestContextReq
		err = utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		consumer, err := sdk.AccAddressFromBech32(req.Consumer)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var serviceFeeCap sdk.Coins

		if len(req.ServiceFeeCap) != 0 {
			serviceFeeCap, err = cliCtx.ParseCoins(req.ServiceFeeCap)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		var providers []sdk.AccAddress
		for _, p := range req.Providers {
			provider, err := sdk.AccAddressFromBech32(p)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			providers = append(providers, provider)
		}

		msg := service.NewMsgUpdateRequestContext(
			requestContextID, providers, serviceFeeCap, req.Timeout,
			req.RepeatedFrequency, req.RepeatedTotal, consumer,
		)
		if err = msg.ValidateBasic(); err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

func withdrawEarnedFeesHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		providerStr := vars[Provider]

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req withdrawEarnedFeesReq
		err = utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		msg := service.NewMsgWithdrawEarnedFees(provider)
		if err = msg.ValidateBasic(); err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}
