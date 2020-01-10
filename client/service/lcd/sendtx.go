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

	// Add a new service binding
	r.HandleFunc(
		"/service/bindings",
		bindingAddHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// Update a service binding
	r.HandleFunc(
		fmt.Sprintf("/service/bindings/{%s}/{%s}/{%s}", DefChainID, ServiceName, Provider),
		bindingUpdateHandlerFn(cdc, cliCtx),
	).Methods("PUT")

	// disable a service binding
	r.HandleFunc(
		fmt.Sprintf("/service/bindings/{%s}/{%s}/{%s}/disable", DefChainID, ServiceName, Provider),
		bindingDisableHandlerFn(cdc, cliCtx),
	).Methods("PUT")

	// enable a service binding
	r.HandleFunc(
		fmt.Sprintf("/service/bindings/{%s}/{%s}/{%s}/enable", DefChainID, ServiceName, Provider),
		bindingEnableHandlerFn(cdc, cliCtx),
	).Methods("PUT")

	// refund deposit from a service binding
	r.HandleFunc(
		fmt.Sprintf("/service/bindings/{%s}/{%s}/{%s}/deposit/refund", DefChainID, ServiceName, Provider),
		bindingRefundHandlerFn(cdc, cliCtx),
	).Methods("PUT")

	// Add a request for a service binding
	r.HandleFunc(
		fmt.Sprintf("/service/requests"),
		requestAddHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// Add a response for a service request
	r.HandleFunc(
		fmt.Sprintf("/service/responses"),
		responseAddHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// refund fees from return fees
	r.HandleFunc(
		fmt.Sprintf("/service/fees/{%s}/refund", Consumer),
		FeesRefundHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// withdraw fees from incoming fees
	r.HandleFunc(
		fmt.Sprintf("/service/fees/{%s}/withdraw", Provider),
		FeesWithdrawHandlerFn(cdc, cliCtx),
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

type binding struct {
	BaseTx      utils.BaseTx  `json:"base_tx"` // basic tx info
	ServiceName string        `json:"service_name"`
	DefChainID  string        `json:"def_chain_id"`
	BindingType string        `json:"binding_type"`
	Deposit     string        `json:"deposit"`
	Prices      []string      `json:"prices"`
	Level       service.Level `json:"level"`
	Provider    string        `json:"provider"`
}

type bindingUpdate struct {
	BaseTx      utils.BaseTx  `json:"base_tx"` // basic tx info
	BindingType string        `json:"binding_type"`
	Deposit     string        `json:"deposit"`
	Prices      []string      `json:"prices"`
	Level       service.Level `json:"level"`
}

type bindingEnable struct {
	BaseTx  utils.BaseTx `json:"base_tx"` // basic tx info
	Deposit string       `json:"deposit"`
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

// HTTP request handler to define service.
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

func bindingAddHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req binding
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		providerAddr, err := sdk.AccAddressFromBech32(req.Provider)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		bindingType, err := service.BindingTypeFromString(req.BindingType)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		deposit, err := cliCtx.ParseCoins(req.Deposit)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var prices []sdk.Coin
		for _, ip := range req.Prices {
			price, err := cliCtx.ParseCoin(ip)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			prices = append(prices, price)
		}

		msg := service.NewMsgSvcBind(req.DefChainID, req.ServiceName, baseReq.ChainID, providerAddr, bindingType, deposit, prices, req.Level)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

func bindingUpdateHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		DefChainId := vars[DefChainID]
		serviceName := vars[ServiceName]
		bechProviderAddr := vars[Provider]

		providerAddr, err := sdk.AccAddressFromBech32(bechProviderAddr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req bindingUpdate
		err = utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		var bindingType service.BindingType
		if req.BindingType != "" {
			bindingType, err = service.BindingTypeFromString(req.BindingType)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		var deposit sdk.Coins
		if req.Deposit != "" {
			deposit, err = cliCtx.ParseCoins(req.Deposit)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		var prices []sdk.Coin
		for _, ip := range req.Prices {
			price, err := cliCtx.ParseCoin(ip)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			prices = append(prices, price)
		}

		msg := service.NewMsgSvcBindingUpdate(DefChainId, serviceName, baseReq.ChainID, providerAddr, bindingType, deposit, prices, req.Level)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

func bindingDisableHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		DefChainId := vars[DefChainID]
		serviceName := vars[ServiceName]
		bechProviderAddr := vars[Provider]

		providerAddr, err := sdk.AccAddressFromBech32(bechProviderAddr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req basicReq
		err = utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		msg := service.NewMsgSvcDisable(DefChainId, serviceName, baseReq.ChainID, providerAddr)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

func bindingEnableHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		DefChainId := vars[DefChainID]
		serviceName := vars[ServiceName]
		bechProviderAddr := vars[Provider]

		providerAddr, err := sdk.AccAddressFromBech32(bechProviderAddr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req bindingEnable
		err = utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		deposit, err := cliCtx.ParseCoins(req.Deposit)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := service.NewMsgSvcEnable(DefChainId, serviceName, baseReq.ChainID, providerAddr, deposit)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

func bindingRefundHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		DefChainId := vars[DefChainID]
		serviceName := vars[ServiceName]
		bechProviderAddr := vars[Provider]

		providerAddr, err := sdk.AccAddressFromBech32(bechProviderAddr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req basicReq
		err = utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		msg := service.NewMsgSvcRefundDeposit(DefChainId, serviceName, baseReq.ChainID, providerAddr)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

func requestAddHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req serviceRequestWithBasic
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		var msgs []sdk.Msg
		for _, request := range req.Requests {
			consumerStr := request.Consumer
			consumer, err := sdk.AccAddressFromBech32(consumerStr)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			providerStr := request.Provider
			provider, err := sdk.AccAddressFromBech32(providerStr)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			inputString := request.Data
			input, err := hex.DecodeString(inputString)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			serviceFeeStr := request.ServiceFee
			serviceFee, err := cliCtx.ParseCoins(serviceFeeStr)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			msg := service.NewMsgSvcRequest(request.DefChainID, request.ServiceName, request.BindChainID, baseReq.ChainID, consumer, provider, request.MethodID, input, serviceFee, request.Profiling)
			err = msg.ValidateBasic()
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			msgs = append(msgs, msg)
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, msgs)
	}
}

func responseAddHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req serviceResponse
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		providerStr := req.Provider
		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		outputString := req.Data
		output, err := hex.DecodeString(outputString)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		errMsgString := req.ErrorMsg
		errMsg, err := hex.DecodeString(errMsgString)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := service.NewMsgSvcResponse(req.ReqChainID, req.RequestID, provider, output, errMsg)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

func FeesRefundHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bechConsumerAddr := vars[Consumer]

		consumerAddr, err := sdk.AccAddressFromBech32(bechConsumerAddr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req basicReq
		err = utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		msg := service.NewMsgSvcRefundFees(consumerAddr)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

func FeesWithdrawHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bechProviderAddr := vars[Provider]

		providerAddr, err := sdk.AccAddressFromBech32(bechProviderAddr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req basicReq
		err = utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		msg := service.NewMsgSvcWithdrawFees(providerAddr)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}
