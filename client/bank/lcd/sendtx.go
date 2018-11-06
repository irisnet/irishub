package lcd

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/bank"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"net/http"
)

type sendBody struct {
	Amount string         `json:"amount"`
	Sender string         `json:"sender"`
	BaseTx context.BaseTx `json:"base_tx"`
}

// SendRequestHandlerFn - http request handler to send coins to a address
// nolint: gocyclo
func SendRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Init context and read request parameters
		cliCtx = utils.InitReqCliCtx(cliCtx, r)
		vars := mux.Vars(r)
		bech32addr := vars["address"]
		to, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		var m sendBody
		err = utils.ReadPostBody(w, r, cdc, &m)
		if err != nil {
			return
		}
		baseReq := m.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w, cliCtx) {
			return
		}
		// Build message
		amount, err := cliCtx.ParseCoins(m.Amount)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		sender, err := sdk.AccAddressFromBech32(m.Sender)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't decode delegator. Error: %s", err.Error())))
			return
		}
		msg := bank.BuildMsg(sender, to, amount)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		// Broadcast or return unsigned transaction
		utils.SendOrReturnUnsignedTx(w, cliCtx, m.BaseTx, []sdk.Msg{msg})
	}
}
