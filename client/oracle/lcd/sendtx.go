package lcd

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/app/v3/oracle"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// define a feed
	r.HandleFunc(
		"/oracle/feeds",
		createFeedHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// edit a feed
	r.HandleFunc(
		fmt.Sprintf("/oracle/feeds/{%s}", FeedName),
		editFeedHandlerFn(cdc, cliCtx),
	).Methods("PUT")

	// start a feed
	r.HandleFunc(
		fmt.Sprintf("/oracle/feeds/{%s}/start", FeedName),
		startFeedHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// pause a feed
	r.HandleFunc(
		fmt.Sprintf("/oracle/feeds/{%s}/pause", FeedName),
		pauseFeedHandlerFn(cdc, cliCtx),
	).Methods("POST")
}

type createFeedReq struct {
	BaseTx            utils.BaseTx `json:"base_tx"` // basic tx info
	FeedName          string       `json:"feed_name"`
	AggregateFunc     string       `json:"aggregate_func"`
	ValueJsonPath     string       `json:"value_json_path"`
	LatestHistory     uint64       `json:"latest_history"`
	Description       string       `json:"description"`
	Creator           string       `json:"creator"`
	ServiceName       string       `json:"service_name"`
	Providers         []string     `json:"providers"`
	Input             string       `json:"input"`
	Timeout           int64        `json:"timeout"`
	ServiceFeeCap     string       `json:"service_fee_cap"`
	RepeatedFrequency uint64       `json:"repeated_frequency"`
	RepeatedTotal     int64        `json:"repeated_total"`
	ResponseThreshold uint16       `json:"response_threshold"`
}

type editFeedReq struct {
	BaseTx            utils.BaseTx `json:"base_tx"` // basic tx info
	FeedName          string       `json:"feed_name"`
	Description       string       `json:"description"`
	LatestHistory     uint64       `json:"latest_history"`
	Creator           string       `json:"creator"`
	Providers         []string     `json:"providers"`
	Timeout           int64        `json:"timeout"`
	ServiceFeeCap     string       `json:"service_fee_cap"`
	RepeatedFrequency uint64       `json:"repeated_frequency"`
	RepeatedTotal     int64        `json:"repeated_total"`
	ResponseThreshold uint16       `json:"response_threshold"`
}

type startFeedReq struct {
	BaseTx  utils.BaseTx `json:"base_tx"` // basic tx info
	Creator string       `json:"creator"`
}

type pauseFeedReq struct {
	BaseTx  utils.BaseTx `json:"base_tx"` // basic tx info
	Creator string       `json:"creator"`
}

func createFeedHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createFeedReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		creator, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var providers []sdk.AccAddress
		for _, addr := range req.Providers {
			provider, err := sdk.AccAddressFromBech32(addr)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			providers = append(providers, provider)
		}

		serviceFeeCap, err := sdk.ParseCoins(req.ServiceFeeCap)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := oracle.MsgCreateFeed{
			FeedName:          req.FeedName,
			LatestHistory:     req.LatestHistory,
			Description:       req.Description,
			Creator:           creator,
			ServiceName:       req.ServiceName,
			Providers:         providers,
			Input:             req.Input,
			Timeout:           req.Timeout,
			ServiceFeeCap:     serviceFeeCap,
			RepeatedFrequency: req.RepeatedFrequency,
			RepeatedTotal:     req.RepeatedTotal,
			ResponseThreshold: req.ResponseThreshold,
			AggregateFunc:     req.AggregateFunc,
			ValueJsonPath:     req.ValueJsonPath,
		}
		if err := msg.ValidateBasic(); err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)
		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

func editFeedHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req editFeedReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		creator, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var providers []sdk.AccAddress
		for _, addr := range req.Providers {
			provider, err := sdk.AccAddressFromBech32(addr)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			providers = append(providers, provider)
		}

		serviceFeeCap, err := sdk.ParseCoins(req.ServiceFeeCap)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := oracle.MsgEditFeed{
			FeedName:          req.FeedName,
			LatestHistory:     req.LatestHistory,
			Description:       req.Description,
			Providers:         providers,
			Timeout:           req.Timeout,
			ServiceFeeCap:     serviceFeeCap,
			RepeatedFrequency: req.RepeatedFrequency,
			RepeatedTotal:     req.RepeatedTotal,
			ResponseThreshold: req.ResponseThreshold,
			Creator:           creator,
		}
		if err := msg.ValidateBasic(); err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)
		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

func startFeedHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req startFeedReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		creator, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		vars := mux.Vars(r)
		feedName := vars[FeedName]

		msg := oracle.MsgStartFeed{
			FeedName: feedName,
			Creator:  creator,
		}
		if err := msg.ValidateBasic(); err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)
		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

func pauseFeedHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pauseFeedReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		creator, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		vars := mux.Vars(r)
		feedName := vars[FeedName]

		msg := oracle.MsgStartFeed{
			FeedName: feedName,
			Creator:  creator,
		}
		if err := msg.ValidateBasic(); err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)
		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}
