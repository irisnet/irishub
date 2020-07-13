package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/gorilla/mux"

	"github.com/irisnet/irishub/modules/random/types"
)

func registerTxRoutes(cliCtx client.Context, r *mux.Router) {
	// request rands
	r.HandleFunc("/rand/rands", requestRandomHandlerFn(cliCtx)).Methods("POST")
}

// HTTP request handler to request rand.
func requestRandomHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestRandomReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// create the MsgRequestRandom message
		msg := types.NewMsgRequestRandom(req.Consumer, req.BlockInterval)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		authclient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
