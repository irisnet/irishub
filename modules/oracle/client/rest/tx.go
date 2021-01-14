package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/irisnet/irismod/modules/oracle/types"
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

		if _, err := sdk.AccAddressFromBech32(req.Creator); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		for _, addr := range req.Providers {
			if _, err := sdk.AccAddressFromBech32(addr); err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		serviceFeeCap, err := sdk.ParseCoinsNormalized(req.ServiceFeeCap)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := &types.MsgCreateFeed{
			FeedName:          req.FeedName,
			LatestHistory:     req.LatestHistory,
			Description:       req.Description,
			Creator:           req.Creator,
			ServiceName:       req.ServiceName,
			Providers:         req.Providers,
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

		if _, err := sdk.AccAddressFromBech32(req.Creator); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		for _, addr := range req.Providers {
			if _, err := sdk.AccAddressFromBech32(addr); err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		serviceFeeCap, err := sdk.ParseCoinsNormalized(req.ServiceFeeCap)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		vars := mux.Vars(r)
		msg := &types.MsgEditFeed{
			FeedName:          vars[FeedName],
			LatestHistory:     req.LatestHistory,
			Description:       req.Description,
			Providers:         req.Providers,
			Timeout:           req.Timeout,
			ServiceFeeCap:     serviceFeeCap,
			RepeatedFrequency: req.RepeatedFrequency,
			ResponseThreshold: req.ResponseThreshold,
			Creator:           req.Creator,
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

		if _, err := sdk.AccAddressFromBech32(req.Creator); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		vars := mux.Vars(r)
		msg := &types.MsgStartFeed{
			FeedName: vars[FeedName],
			Creator:  req.Creator,
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

		if _, err := sdk.AccAddressFromBech32(req.Creator); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		vars := mux.Vars(r)
		msg := &types.MsgStartFeed{
			FeedName: vars[FeedName],
			Creator:  req.Creator,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}
