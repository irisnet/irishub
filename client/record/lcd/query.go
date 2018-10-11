package lcd

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/client/context"
	recordClient "github.com/irisnet/irishub/client/record"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/modules/record"
	cmn "github.com/tendermint/tendermint/libs/common"
)

// nolint: gocyclo
func queryRecordsWithParameterFn(cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		dataHash := r.URL.Query().Get(RecordHash)
		if len(dataHash) == 0 {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("'%s' can't be empty", dataHash))
			return
		}

		var tmpkey = cmn.HexBytes{}
		res, err := cliCtx.QueryStore(tmpkey /*record.KeyProposal(hashHexStr)*/, storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusNotFound, "no exist yet and dataHash has not been set")
			return
		}

		if len(res) == 0 || err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Data hash [%s] is not existe", dataHash))
			return
		}

		var submitFile record.MsgSubmitFile
		cdc.MustUnmarshalBinary(res, &submitFile)

		recordResponse, err := recordClient.ConvertRecordToRecordOutput(cliCtx, submitFile)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		output, err := wire.MarshalJSONIndent(cdc, recordResponse)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Write(output)
	}
}
