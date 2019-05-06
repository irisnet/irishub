package lcd

import (
	"net/http"

	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/auth"
	sdk "github.com/irisnet/irishub/types"
)

type broadcastBody struct {
	Tx auth.StdTx `json:"tx"`
}

// BroadcastTxRequestHandlerFn returns the broadcast tx REST handler
func BroadcastTxRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx = utils.InitReqCliCtx(cliCtx, r)
		var m broadcastBody
		if err := utils.ReadPostBody(w, r, cliCtx.Codec, &m); err != nil {
			return
		}

		txBytes, err := cliCtx.Codec.MarshalBinaryLengthPrefixed(m.Tx)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if cliCtx.DryRun {
			rawRes, err := cliCtx.Query("/app/simulate", txBytes)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			var simulationResult sdk.Result
			if err := cdc.UnmarshalBinaryLengthPrefixed(rawRes, &simulationResult); err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			utils.WriteSimulationResponse(w, cliCtx, simulationResult.GasUsed, simulationResult)
			return
		}
		res, err := cliCtx.BroadcastTx(txBytes)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
