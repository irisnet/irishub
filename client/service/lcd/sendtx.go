package lcd

import (
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

	// set a new withdrawal address for a service binding
	r.HandleFunc(
		fmt.Sprintf("/service/bindings/{%s}/{%s}/withdraw-address", ServiceName, Provider),
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

	// // Add a request for a service binding
	// r.HandleFunc(
	// 	fmt.Sprintf("/service/requests"),
	// 	requestAddHandlerFn(cdc, cliCtx),
	// ).Methods("POST")

	// // Add a response for a service request
	// r.HandleFunc(
	// 	fmt.Sprintf("/service/responses"),
	// 	responseAddHandlerFn(cdc, cliCtx),
	// ).Methods("POST")

	// // refund fees from return fees
	// r.HandleFunc(
	// 	fmt.Sprintf("/service/fees/{%s}/refund", Consumer),
	// 	FeesRefundHandlerFn(cdc, cliCtx),
	// ).Methods("POST")

	// // withdraw fees from incoming fees
	// r.HandleFunc(
	// 	fmt.Sprintf("/service/fees/{%s}/withdraw", Provider),
	// 	FeesWithdrawHandlerFn(cdc, cliCtx),
	// ).Methods("POST")
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
	BaseTx          utils.BaseTx `json:"base_tx"` // basic tx info
	ServiceName     string       `json:"service_name"`
	Provider        string       `json:"provider"`
	Deposit         string       `json:"deposit"`
	Pricing         string       `json:"pricing"`
	WithdrawAddress string       `json:"withdraw_address"`
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

type serviceRequest struct {
	ServiceName string `json:"service_name"`
	BindChainID string `json:"bind_chain_id"`
	DefChainID  string `json:"def_chain_id"`
	MethodID    int16  `json:"method_id"`
	Provider    string `json:"provider"`
	Consumer    string `json:"consumer"`
	ServiceFee  string `json:"service_fee"`
	Data        string `json:"data"`
	Profiling   bool   `json:"profiling"`
}

type serviceRequestWithBasic struct {
	BaseTx   utils.BaseTx     `json:"base_tx"` // basic tx info
	Requests []serviceRequest `json:"requests"`
}

type serviceResponse struct {
	BaseTx     utils.BaseTx `json:"base_tx"` // basic tx info
	ReqChainID string       `json:"req_chain_id"`
	RequestID  string       `json:"request_id"`
	Data       string       `json:"data"`
	Provider   string       `json:"provider"`
	ErrorMsg   string       `json:"error_msg"`
}

type basicReq struct {
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

		withdrawAddr, err := sdk.AccAddressFromBech32(req.WithdrawAddress)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := service.NewMsgBindService(req.ServiceName, provider, deposit, req.Pricing, withdrawAddr)
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

// HTTP request handler to set a new withdrawal address for a service binding.
func setWithdrawAddrHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[ServiceName]
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

		msg := service.NewMsgSetWithdrawAddress(serviceName, provider, withdrawAddr)
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

// func requestAddHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var req serviceRequestWithBasic
// 		err := utils.ReadPostBody(w, r, cdc, &req)
// 		if err != nil {
// 			return
// 		}

// 		baseReq := req.BaseTx.Sanitize()
// 		if !baseReq.ValidateBasic(w) {
// 			return
// 		}

// 		var msgs []sdk.Msg
// 		for _, request := range req.Requests {
// 			consumerStr := request.Consumer
// 			consumer, err := sdk.AccAddressFromBech32(consumerStr)
// 			if err != nil {
// 				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 				return
// 			}

// 			providerStr := request.Provider
// 			provider, err := sdk.AccAddressFromBech32(providerStr)
// 			if err != nil {
// 				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 				return
// 			}

// 			inputString := request.Data
// 			input, err := hex.DecodeString(inputString)
// 			if err != nil {
// 				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 				return
// 			}

// 			serviceFeeStr := request.ServiceFee
// 			serviceFee, err := cliCtx.ParseCoins(serviceFeeStr)
// 			if err != nil {
// 				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 				return
// 			}

// 			msg := service.NewMsgSvcRequest(request.DefChainID, request.ServiceName, request.BindChainID, baseReq.ChainID, consumer, provider, request.MethodID, input, serviceFee, request.Profiling)
// 			err = msg.ValidateBasic()
// 			if err != nil {
// 				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 				return
// 			}
// 			msgs = append(msgs, msg)
// 		}

// 		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

// 		utils.WriteGenerateStdTxResponse(w, txCtx, msgs)
// 	}
// }

// func responseAddHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var req serviceResponse
// 		err := utils.ReadPostBody(w, r, cdc, &req)
// 		if err != nil {
// 			return
// 		}

// 		baseReq := req.BaseTx.Sanitize()
// 		if !baseReq.ValidateBasic(w) {
// 			return
// 		}

// 		providerStr := req.Provider
// 		provider, err := sdk.AccAddressFromBech32(providerStr)
// 		if err != nil {
// 			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		outputString := req.Data
// 		output, err := hex.DecodeString(outputString)
// 		if err != nil {
// 			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		errMsgString := req.ErrorMsg
// 		errMsg, err := hex.DecodeString(errMsgString)
// 		if err != nil {
// 			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		msg := service.NewMsgSvcResponse(req.ReqChainID, req.RequestID, provider, output, errMsg)
// 		err = msg.ValidateBasic()
// 		if err != nil {
// 			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

// 		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
// 	}
// }

// func FeesRefundHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)
// 		bechConsumerAddr := vars[Consumer]

// 		consumerAddr, err := sdk.AccAddressFromBech32(bechConsumerAddr)
// 		if err != nil {
// 			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		var req basicReq
// 		err = utils.ReadPostBody(w, r, cdc, &req)
// 		if err != nil {
// 			return
// 		}

// 		baseReq := req.BaseTx.Sanitize()
// 		if !baseReq.ValidateBasic(w) {
// 			return
// 		}

// 		msg := service.NewMsgSvcRefundFees(consumerAddr)
// 		err = msg.ValidateBasic()
// 		if err != nil {
// 			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

// 		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
// 	}
// }

// func FeesWithdrawHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)
// 		bechProviderAddr := vars[Provider]

// 		providerAddr, err := sdk.AccAddressFromBech32(bechProviderAddr)
// 		if err != nil {
// 			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		var req basicReq
// 		err = utils.ReadPostBody(w, r, cdc, &req)
// 		if err != nil {
// 			return
// 		}

// 		baseReq := req.BaseTx.Sanitize()
// 		if !baseReq.ValidateBasic(w) {
// 			return
// 		}

// 		msg := service.NewMsgSvcWithdrawFees(providerAddr)
// 		err = msg.ValidateBasic()
// 		if err != nil {
// 			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

// 		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
// 	}
// }
