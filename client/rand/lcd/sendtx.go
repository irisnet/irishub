package lcd

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/app/v1/rand"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// request a rand
	r.HandleFunc(
		"/rand/rands",
		requestRandHandlerFn(cdc, cliCtx),
	).Methods("POST")
}

type requestRandReq struct {
	BaseTx        utils.BaseTx   `json:"base_tx"`        // base tx
	Consumer      sdk.AccAddress `json:"consumer"`       // request address
	BlockInterval uint64         `json:"block_interval"` // block interval
}

func requestRandHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req requestRandReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgRequestRand message
		msg := rand.NewMsgRequestRand(req.Consumer, req.BlockInterval)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txCtx := utils.BuildReqTxCtx(cliCtx, baseReq, w)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}
