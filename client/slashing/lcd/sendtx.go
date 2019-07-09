package lcd

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/app/v1/slashing"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

// Unrevoke TX body
type UnjailBody struct {
	BaseTx utils.BaseTx `json:"base_tx"`
}

func unrevokeRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		txCtx := utils.NewTxContextFromCLI().WithCodec(cliCtx.Codec)

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
		if !baseReq.ValidateBasic(w) {
			return
		}

		msg := slashing.NewMsgUnjail(validatorAddr)

		utils.WriteGenerateStdTxResponse(w, txCtx, []sdk.Msg{msg})
	}
}
