package lcd

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"net/http"
	"github.com/gorilla/mux"
)

// Unrevoke TX body
type UnjailBody struct {
	BaseTx        context.BaseTx `json:"base_tx"`
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
