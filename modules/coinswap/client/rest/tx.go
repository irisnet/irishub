package rest

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/irisnet/irishub/modules/coinswap/internal/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/coinswap/liquidities/{id}/deposit", addLiquidityHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/coinswap/liquidities/{id}/withdraw", removeLiquidityHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/coinswap/liquidities/buy", swapOrderHandlerFn(cliCtx, true)).Methods("POST")
	r.HandleFunc("/coinswap/liquidities/sell", swapOrderHandlerFn(cliCtx, false)).Methods("POST")
}

func addLiquidityHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		uniDenom, err := types.GetUniDenom(id)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tokenDenom, err := types.GetCoinMinDenomFromUniDenom(uniDenom)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req AddLiquidityReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		senderAddress, e := sdk.AccAddressFromBech32(req.Sender)
		if e != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, e.Error())
			return
		}

		duration, e := time.ParseDuration(req.Deadline)
		if e != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, e.Error())
			return
		}

		deadline := time.Now().Add(duration)

		maxToken, ok := sdk.NewIntFromString(req.MaxToken)
		if !ok || !maxToken.IsPositive() {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid max token amount: "+req.MaxToken)
			return
		}

		exactIrisAmt, ok := sdk.NewIntFromString(req.ExactIrisAmt)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid exact iris amount: "+req.ExactIrisAmt)
			return
		}

		minLiquidity, ok := sdk.NewIntFromString(req.MinLiquidity)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid min liquidity amount: "+req.MinLiquidity)
			return
		}

		msg := types.NewMsgAddLiquidity(sdk.NewCoin(tokenDenom, maxToken), exactIrisAmt, minLiquidity, deadline.Unix(), senderAddress)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseTx, []sdk.Msg{msg})
	}
}

func removeLiquidityHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		uniDenom, err := types.GetUniDenom(id)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req RemoveLiquidityReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		senderAddress, e := sdk.AccAddressFromBech32(req.Sender)
		if e != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, e.Error())
			return
		}

		duration, e := time.ParseDuration(req.Deadline)
		if e != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, e.Error())
			return
		}

		deadline := time.Now().Add(duration)

		minToken, ok := sdk.NewIntFromString(req.MinToken)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid min token amount: "+req.MinToken)
			return
		}

		minIris, ok := sdk.NewIntFromString(req.MinIrisAmt)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid min iris amount: "+req.MinIrisAmt)
			return
		}

		liquidityAmt, ok := sdk.NewIntFromString(req.WithdrawLiquidity)
		if !ok || !liquidityAmt.IsPositive() {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid liquidity amount: "+req.WithdrawLiquidity)
			return
		}

		msg := types.NewMsgRemoveLiquidity(minToken, sdk.NewCoin(uniDenom, liquidityAmt), minIris, deadline.Unix(), senderAddress)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseTx, []sdk.Msg{msg})
	}
}

func swapOrderHandlerFn(cliCtx context.CLIContext, isBuyOrder bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SwapOrderReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		senderAddress, err := sdk.AccAddressFromBech32(req.Input.Address)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var recipientAddress sdk.AccAddress
		if len(req.Output.Address) > 0 {
			recipientAddress, err = sdk.AccAddressFromBech32(req.Output.Address)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		duration, err := time.ParseDuration(req.Deadline)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		input := types.Input{Address: senderAddress, Coin: req.Input.Coin}
		output := types.Output{Address: recipientAddress, Coin: req.Output.Coin}
		deadline := time.Now().Add(duration)

		msg := types.NewMsgSwapOrder(input, output, deadline.Unix(), isBuyOrder)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseTx, []sdk.Msg{msg})
	}
}
