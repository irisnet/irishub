package lcd

import (
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/app/v2/coinswap"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"net/http"
	"time"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(
		"/coinswap/liquidities/{id}/deposit",
		addLiquidityHandlerFn(cdc, cliCtx),
	).Methods("POST")

	r.HandleFunc(
		"/coinswap/liquidities/{id}/withdraw",
		removeLiquidityHandlerFn(cdc, cliCtx),
	).Methods("POST")

	r.HandleFunc(
		"/coinswap/liquidities/buy",
		swapOrderHandlerFn(cdc, cliCtx, true),
	).Methods("POST")

	r.HandleFunc(
		"/coinswap/liquidities/sell",
		swapOrderHandlerFn(cdc, cliCtx, false),
	).Methods("POST")
}

type addLiquidityReq struct {
	BaseTx       utils.BaseTx `json:"base_tx"`
	Id           string       `json:"id"`             // the unique liquidity id
	MaxToken     string       `json:"max_token"`      // token to be deposited as liquidity with an upper bound for its amount
	ExactIrisAmt string       `json:"exact_iris_amt"` // exact amount of iris-atto being add to the liquidity pool
	MinLiquidity string       `json:"min_liquidity"`  // lower bound UNI sender is willing to accept for deposited coins
	Deadline     string       `json:"deadline"`       // deadline duration, e.g. 10m
	Sender       string       `json:"sender"`
}

type removeLiquidityReq struct {
	BaseTx            utils.BaseTx `json:"base_tx"`
	Id                string       `json:"id"`                 // the unique liquidity id
	MinToken          string       `json:"min_token"`          // coin to be withdrawn with a lower bound for its amount
	WithdrawLiquidity string       `json:"withdraw_liquidity"` // amount of UNI to be burned to withdraw liquidity from a reserve pool
	MinIrisAmt        string       `json:"min_iris_amt"`       // minimum amount of the native asset the sender is willing to accept
	Deadline          string       `json:"deadline"`           // deadline duration, e.g. 10m
	Sender            string       `json:"sender"`
}

type input struct {
	Address string   `json:"address"`
	Coin    sdk.Coin `json:"coin"`
}

type output struct {
	Address string   `json:"address"`
	Coin    sdk.Coin `json:"coin"`
}

type swapOrderReq struct {
	BaseTx   utils.BaseTx `json:"base_tx"`
	Input    input        `json:"input"`    // the amount the sender is trading
	Output   output       `json:"output"`   // the amount the sender is receiving
	Deadline string       `json:"deadline"` // deadline for the transaction to still be considered valid
}

func addLiquidityHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		uniDenom, err := coinswap.GetUniDenom(id)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tokenDenom, err := coinswap.GetCoinMinDenomFromUniDenom(uniDenom)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req addLiquidityReq
		err1 := utils.ReadPostBody(w, r, cdc, &req)
		if err1 != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err1.Error())
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		senderAddress, err1 := sdk.AccAddressFromBech32(req.Sender)
		if err1 != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err1.Error())
			return
		}

		duration, err1 := time.ParseDuration(req.Deadline)
		if err1 != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err1.Error())
			return
		}

		deadline := time.Now().Add(duration)

		maxToken, ok := sdk.NewIntFromString(req.MaxToken)
		if !ok || !maxToken.IsPositive() {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid max token amount: "+req.MaxToken)
			return
		}

		exactIrisAmt, ok := sdk.NewIntFromString(req.ExactIrisAmt)
		if !ok {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid exact iris amount: "+req.ExactIrisAmt)
			return
		}

		minLiquidity, ok := sdk.NewIntFromString(req.MinLiquidity)
		if !ok {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid min liquidity amount: "+req.MinLiquidity)
			return
		}

		msg := coinswap.NewMsgAddLiquidity(sdk.NewCoin(tokenDenom, maxToken), exactIrisAmt, minLiquidity, deadline.Unix(), senderAddress)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

func removeLiquidityHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		uniDenom, err := coinswap.GetUniDenom(id)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req removeLiquidityReq
		err1 := utils.ReadPostBody(w, r, cdc, &req)
		if err1 != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err1.Error())
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		senderAddress, err1 := sdk.AccAddressFromBech32(req.Sender)
		if err1 != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err1.Error())
			return
		}

		duration, err1 := time.ParseDuration(req.Deadline)
		if err1 != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err1.Error())
			return
		}

		deadline := time.Now().Add(duration)

		minToken, ok := sdk.NewIntFromString(req.MinToken)
		if !ok {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid min token amount: "+req.MinToken)
			return
		}

		minIris, ok := sdk.NewIntFromString(req.MinIrisAmt)
		if !ok {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid min iris amount: "+req.MinIrisAmt)
			return
		}

		liquidityAmt, ok := sdk.NewIntFromString(req.WithdrawLiquidity)
		if !ok || !liquidityAmt.IsPositive() {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid liquidity amount: "+req.WithdrawLiquidity)
			return
		}

		msg := coinswap.NewMsgRemoveLiquidity(minToken, sdk.NewCoin(uniDenom, liquidityAmt), minIris, deadline.Unix(), senderAddress)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}

func swapOrderHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext, isBuyOrder bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req swapOrderReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		senderAddress, err := sdk.AccAddressFromBech32(req.Input.Address)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var recipientAddress sdk.AccAddress
		if len(req.Output.Address) > 0 {
			recipientAddress, err = sdk.AccAddressFromBech32(req.Output.Address)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		duration, err := time.ParseDuration(req.Deadline)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		input := coinswap.Input{Address: senderAddress, Coin: req.Input.Coin}
		output := coinswap.Output{Address: recipientAddress, Coin: req.Output.Coin}
		deadline := time.Now().Add(duration)

		msg := coinswap.NewMsgSwapOrder(input, output, deadline.Unix(), isBuyOrder)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}
