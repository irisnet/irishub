package lcd

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/app/v1/asset"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// create a gateway
	r.HandleFunc(
		"/asset/gateways",
		createGatewayHandlerFn(cdc, cliCtx),
	).Methods("POST")

	// edit a gateway
	r.HandleFunc(
		"/asset/gateways/{moniker}",
		editGatewayHandlerFn(cdc, cliCtx),
	).Methods("PUT")
}

type createGatewayReq struct {
	BaseTx   utils.BaseTx   `json:"base_tx"`
	Owner    sdk.AccAddress `json:"owner"`    //  Owner of the gateway
	Moniker  string         `json:"moniker"`  //  Name of the gateway
	Identity string         `json:"identity"` //  Identity of the gateway
	Details  string         `json:"details"`  //  Description of the gateway
	Website  string         `json:"website"`  //  Website of the gateway
	Fee      sdk.Coin       `json:"fee"`      //  Creation fee of the gateway
}

type editGatewayReq struct {
	BaseTx   utils.BaseTx   `json:"base_tx"`
	Owner    sdk.AccAddress `json:"owner"`    //  Owner of the gateway
	Identity *string        `json:"identity"` //  Identity of the gateway
	Details  *string        `json:"details"`  //  Description of the gateway
	Website  *string        `json:"website"`  //  Website of the gateway
}

func createGatewayHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx = utils.InitReqCliCtx(cliCtx, r)

		var req createGatewayReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w, cliCtx) {
			return
		}

		// create the MsgCreateGateway message
		msg := asset.NewMsgCreateGateway(req.Owner, req.Moniker, req.Identity, req.Details, req.Website, req.Fee)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.SendOrReturnUnsignedTx(w, cliCtx, req.BaseTx, []sdk.Msg{msg})
	}
}

func editGatewayHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx = utils.InitReqCliCtx(cliCtx, r)

		vars := mux.Vars(r)
		moniker := vars["moniker"]

		var req editGatewayReq
		err := utils.ReadPostBody(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w, cliCtx) {
			return
		}

		// create the MsgEditGateway message
		msg := asset.NewMsgEditGateway(req.Owner, moniker, req.Identity, req.Details, req.Website)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.SendOrReturnUnsignedTx(w, cliCtx, req.BaseTx, []sdk.Msg{msg})
	}
}
