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
		"/swap/liquidities/{address}/deposit",
		addLiquidityHandlerFn(cdc, cliCtx),
	).Methods("POST")

	r.HandleFunc(
		"/swap/liquidities/{address}/withdraw",
		removeLiquidityHandlerFn(cdc, cliCtx),
	).Methods("POST")
}

type addLiquidityReq struct {
	BaseTx       utils.BaseTx `json:"base_tx"`
	MaxToken     sdk.Coin     `json:"max_token"`      // coin to be deposited as liquidity with an upper bound for its amount
	ExactIrisAmt uint64       `json:"exact_iris_amt"` // exact amount of iris-atto being add to the liquidity pool
	MinLiquidity uint64       `json:"min_liquidity"`  // lower bound UNI sender is willing to accept for deposited coins
	Deadline     string       `json:"deadline"`       // deadline duration, e.g. 10m
}

type removeLiquidityReq struct {
	BaseTx            utils.BaseTx `json:"base_tx"`
	MinToken          uint64       `json:"min_token"`          // coin to be withdrawn with a lower bound for its amount
	WithdrawLiquidity sdk.Coin     `json:"withdraw_liquidity"` // amount of UNI to be burned to withdraw liquidity from a reserve pool
	MinIrisAmt        uint64       `json:"min_iris_amt"`       // minimum amount of the native asset the sender is willing to accept
	Deadline          string       `json:"deadline"`           // deadline duration, e.g. 10m
}

func addLiquidityHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars["address"]
		senderAddress, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req addLiquidityReq
		err = utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		duration, err := time.ParseDuration(req.Deadline)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		deadline := time.Now().Add(duration)

		msg := coinswap.NewMsgAddLiquidity(req.MaxToken, sdk.NewIntFromUint64(req.ExactIrisAmt), sdk.NewIntFromUint64(req.MinLiquidity), deadline, senderAddress)
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
		address := vars["address"]
		senderAddress, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req removeLiquidityReq
		err = utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		duration, err := time.ParseDuration(req.Deadline)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		deadline := time.Now().Add(duration)

		msg := coinswap.NewMsgRemoveLiquidity(sdk.NewIntFromUint64(req.MinToken), req.WithdrawLiquidity, sdk.NewIntFromUint64(req.MinIrisAmt), deadline, senderAddress)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}
