package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/irisnet/irishub/modules/oracle/types"
)

func registerTxRoutes(cliCtx client.Context, r *mux.Router) {
	// define a feed
	r.HandleFunc("/oracle/feeds", createFeedHandlerFn(cliCtx)).Methods("POST")
	// edit a feed
	r.HandleFunc(fmt.Sprintf("/oracle/feeds/{%s}", FeedName), editFeedHandlerFn(cliCtx)).Methods("PUT")
	// start a feed
	r.HandleFunc(fmt.Sprintf("/oracle/feeds/{%s}/start", FeedName), startFeedHandlerFn(cliCtx)).Methods("POST")
	// pause a feed
	r.HandleFunc(fmt.Sprintf("/oracle/feeds/{%s}/pause", FeedName), pauseFeedHandlerFn(cliCtx)).Methods("POST")
}

type createFeedReq struct {
	BaseReq           rest.BaseReq `json:"base_req" yaml:"base_req"`
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
	ResponseThreshold uint32       `json:"response_threshold"`
}

type editFeedReq struct {
	BaseReq           rest.BaseReq `json:"base_req" yaml:"base_req"`
	FeedName          string       `json:"feed_name"`
	Description       string       `json:"description"`
	LatestHistory     uint64       `json:"latest_history"`
	Creator           string       `json:"creator"`
	Providers         []string     `json:"providers"`
	Timeout           int64        `json:"timeout"`
	ServiceFeeCap     string       `json:"service_fee_cap"`
	RepeatedFrequency uint64       `json:"repeated_frequency"`
	ResponseThreshold uint32       `json:"response_threshold"`
}

type startFeedReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Creator string       `json:"creator"`
}

type pauseFeedReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Creator string       `json:"creator"`
}

func createFeedHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createFeedReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		creator, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var providers []sdk.AccAddress
		for _, addr := range req.Providers {
			provider, err := sdk.AccAddressFromBech32(addr)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			providers = append(providers, provider)
		}

		serviceFeeCap, err := sdk.ParseCoins(req.ServiceFeeCap)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := &types.MsgCreateFeed{
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
			ResponseThreshold: req.ResponseThreshold,
			AggregateFunc:     req.AggregateFunc,
			ValueJsonPath:     req.ValueJsonPath,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func editFeedHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req editFeedReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		creator, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var providers []sdk.AccAddress
		for _, addr := range req.Providers {
			provider, err := sdk.AccAddressFromBech32(addr)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			providers = append(providers, provider)
		}

		serviceFeeCap, err := sdk.ParseCoins(req.ServiceFeeCap)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := &types.MsgEditFeed{
			FeedName:          req.FeedName,
			LatestHistory:     req.LatestHistory,
			Description:       req.Description,
			Providers:         providers,
			Timeout:           req.Timeout,
			ServiceFeeCap:     serviceFeeCap,
			RepeatedFrequency: req.RepeatedFrequency,
			ResponseThreshold: req.ResponseThreshold,
			Creator:           creator,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func startFeedHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req startFeedReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		creator, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		vars := mux.Vars(r)
		feedName := vars[FeedName]

		msg := &types.MsgStartFeed{
			FeedName: feedName,
			Creator:  creator,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func pauseFeedHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pauseFeedReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		creator, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		vars := mux.Vars(r)
		feedName := vars[FeedName]

		msg := &types.MsgStartFeed{
			FeedName: feedName,
			Creator:  creator,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}
