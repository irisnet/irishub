package lcd

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
	recordClient "github.com/irisnet/irishub/client/record"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/modules/record"
)

// nolint: gocyclo
func queryRecordsWithParameterFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		recordID := r.URL.Query().Get(RestRecordID)
		if len(recordID) == 0 {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Record ID '%s' can't be empty", recordID))
			return
		}

		res, err := cliCtx.QueryStore([]byte(recordID), storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusNotFound, "no exist yet and record ID has not been set")
			return
		}

		if len(res) == 0 || err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Record ID [%s] is not existe", recordID))
			return
		}

		var submitFile record.MsgSubmitRecord
		cdc.MustUnmarshalBinary(res, &submitFile)

		recordResponse, err := recordClient.ConvertRecordToRecordOutput(cliCtx, submitFile)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		output, err := codec.MarshalJSONIndent(cdc, recordResponse)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Write(output)
	}
}

func queryRecordHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		recordID := vars[RestRecordID]
		if len(recordID) == 0 {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Record ID '%s' can't be empty", recordID))
			return
		}

		res, err := cliCtx.QueryStore([]byte(recordID), storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusNotFound, "no exist yet and record ID has not been set")
			return
		}

		if len(res) == 0 || err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Record ID [%s] is not existe", recordID))
			return
		}

		var submitFile record.MsgSubmitRecord
		cdc.MustUnmarshalBinary(res, &submitFile)

		recordResponse, err := recordClient.ConvertRecordToRecordOutput(cliCtx, submitFile)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		output, err := codec.MarshalJSONIndent(cdc, recordResponse)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Write(output)
	}
}
