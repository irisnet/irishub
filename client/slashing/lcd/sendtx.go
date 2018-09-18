package lcd

import (
	"bytes"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"net/http"
)

// Unrevoke TX body
type UnrevokeBody struct {
	BaseTx        context.BaseTx `json:"base_tx"`
	ValidatorAddr string         `json:"validator_addr"`
}

func unrevokeRequestHandlerFn(cdc *wire.Codec, kb keys.Keybase, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m UnrevokeBody
		err := utils.ReadPostBody(w, r, cdc, &m)
		if err != nil {
			return
		}
		cliCtx = utils.InitRequestClictx(cliCtx, r, m.BaseTx.LocalAccountName, m.ValidatorAddr)
		txCtx, err := context.NewTxContextFromBaseTx(cliCtx, cdc, m.BaseTx)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		validatorAddr, err := sdk.AccAddressFromBech32(m.ValidatorAddr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't decode validator. Error: %s", err.Error()))
			return
		}

		if !cliCtx.GenerateOnly {
			fromAddress, err := cliCtx.GetFromAddress()
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			if !bytes.Equal(fromAddress, validatorAddr) {
				utils.WriteErrorResponse(w, http.StatusUnauthorized, "Must use own validator address")
				return
			}
		}

		msg := slashing.NewMsgUnrevoke(validatorAddr)

		utils.SendOrReturnUnsignedTx(w, cliCtx, txCtx, m.BaseTx, []sdk.Msg{msg})
	}
}
