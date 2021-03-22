package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/irisnet/irismod/modules/coinswap/types"
)

func registerTxRoutes(cliCtx client.Context, r *mux.Router) {
	// add liquidity
	r.HandleFunc(fmt.Sprintf("/coinswap/liquidities/{%s}/deposit", RestPoolID), addLiquidityHandlerFn(cliCtx)).Methods("POST")
	// remove liquidity
	r.HandleFunc(fmt.Sprintf("/coinswap/liquidities/{%s}/withdraw", RestPoolID), removeLiquidityHandlerFn(cliCtx)).Methods("POST")
	// post a buy order
	r.HandleFunc("/coinswap/liquidities/buy", swapOrderHandlerFn(cliCtx, true)).Methods("POST")
	// post a sell order
	r.HandleFunc("/coinswap/liquidities/sell", swapOrderHandlerFn(cliCtx, false)).Methods("POST")
}

// HTTP request handler to add liquidity.
func addLiquidityHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars[RestPoolID]

		if err := sdk.ValidateDenom(denom); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req AddLiquidityReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		if _, err := sdk.AccAddressFromBech32(req.Sender); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		duration, e := time.ParseDuration(req.Deadline)
		if e != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, e.Error())
			return
		}

		status, e := cliCtx.Client.Status(context.Background())
		if e != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, e.Error())
			return
		}
		deadline := status.SyncInfo.LatestBlockTime.Add(duration)

		maxToken, ok := sdk.NewIntFromString(req.MaxToken)
		if !ok || !maxToken.IsPositive() {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid max token amount: "+req.MaxToken)
			return
		}

		exactStandardAmt, ok := sdk.NewIntFromString(req.ExactStandardAmt)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid exact standard token amount: "+req.ExactStandardAmt)
			return
		}

		minLiquidity, ok := sdk.NewIntFromString(req.MinLiquidity)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid min liquidity amount: "+req.MinLiquidity)
			return
		}

		msg := types.NewMsgAddLiquidity(sdk.NewCoin(denom, maxToken), exactStandardAmt, minLiquidity, deadline.Unix(), req.Sender)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

// HTTP request handler to remove liquidity.
func removeLiquidityHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars[RestPoolID]

		if err := sdk.ValidateDenom(denom); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req RemoveLiquidityReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		if _, err := sdk.AccAddressFromBech32(req.Sender); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		duration, e := time.ParseDuration(req.Deadline)
		if e != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, e.Error())
			return
		}

		status, e := cliCtx.Client.Status(context.Background())
		if e != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, e.Error())
			return
		}
		deadline := status.SyncInfo.LatestBlockTime.Add(duration)

		minToken, ok := sdk.NewIntFromString(req.MinToken)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid min token amount: "+req.MinToken)
			return
		}

		minStandard, ok := sdk.NewIntFromString(req.MinStandardAmt)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid min iris amount: "+req.MinStandardAmt)
			return
		}

		liquidityAmt, ok := sdk.NewIntFromString(req.WithdrawLiquidity)
		if !ok || !liquidityAmt.IsPositive() {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid liquidity amount: "+req.WithdrawLiquidity)
			return
		}

		uniDenom := types.GetUniDenomFromDenom(denom)

		msg := types.NewMsgRemoveLiquidity(
			minToken, sdk.NewCoin(uniDenom, liquidityAmt), minStandard, deadline.Unix(), req.Sender,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

// HTTP request handler to post order.
func swapOrderHandlerFn(cliCtx client.Context, isBuyOrder bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SwapOrderReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		senderAddress, err := sdk.AccAddressFromBech32(req.Input.Address)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		recipientAddress, err := sdk.AccAddressFromBech32(req.Output.Address)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		duration, err := time.ParseDuration(req.Deadline)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		input := types.Input{Address: senderAddress.String(), Coin: req.Input.Coin}
		output := types.Output{Address: recipientAddress.String(), Coin: req.Output.Coin}

		status, e := cliCtx.Client.Status(context.Background())
		if e != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, e.Error())
			return
		}
		deadline := status.SyncInfo.LatestBlockTime.Add(duration)

		msg := types.NewMsgSwapOrder(input, output, deadline.Unix(), isBuyOrder)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}
