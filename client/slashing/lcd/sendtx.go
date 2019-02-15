package lcd

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/slashing"
	sdk "github.com/irisnet/irishub/types"
)

// Unrevoke TX body
type UnjailBody struct {
	BaseTx utils.BaseTx `json:"base_tx"`
}

func unrevokeRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx = utils.InitReqCliCtx(cliCtx, r)
		vars := mux.Vars(r)

		validatorAddr, err := sdk.ValAddressFromBech32(vars["validatorAddr"])
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var m UnjailBody
		err = utils.ReadPostBody(w, r, cdc, &m)
		if err != nil {
			return
		}

		baseReq := m.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w, cliCtx) {
			return
		}

		msg := slashing.NewMsgUnjail(validatorAddr)

		utils.SendOrReturnUnsignedTx(w, cliCtx, m.BaseTx, []sdk.Msg{msg})
	}
}
