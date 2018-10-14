package lcd

import (
	"fmt"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/context"
	recordClient "github.com/irisnet/irishub/client/record"
	"github.com/irisnet/irishub/client/record/cli"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/modules/record"
	"github.com/spf13/viper"
)

// nolint: gocyclo
func queryRecordsWithParameterFn(cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		hashHexStr := r.URL.Query().Get(RestRecordHash)
		if len(hashHexStr) == 0 {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Tx hash '%s' can't be empty", hashHexStr))
			return
		}

		accountAddress := r.URL.Query().Get(RestAccountAddress)
		addr, err := sdk.AccAddressFromBech32(accountAddress)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Account address '%s' is wrong", accountAddress))
			return
		}

		trustNode := viper.GetBool(client.FlagTrustNode)
		ipfsHash, err := cli.GetDataHash(cdc, cliCtx, hashHexStr, trustNode)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("get data hash '%s' failed", hashHexStr))
			return
		}

		res, err := cliCtx.QueryStore(record.KeyRecord(addr, ipfsHash), storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusNotFound, "no exist yet and dataHash has not been set")
			return
		}

		if len(res) == 0 || err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Data hash [%s] is not existe", hashHexStr))
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

func queryRecordHandlerFn(cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hashHexStr := vars[RestRecordHash]
		if len(hashHexStr) == 0 {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Tx hash '%s' can't be empty", hashHexStr))
			return
		}

		accountAddress := vars[RestAccountAddress]
		addr, err := sdk.AccAddressFromBech32(accountAddress)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Account address '%s' is wrong", accountAddress))
			return
		}

		trustNode := viper.GetBool(client.FlagTrustNode)
		ipfsHash, err := cli.GetDataHash(cdc, cliCtx, hashHexStr, trustNode)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("get data hash '%s' failed", hashHexStr))
			return
		}

		res, err := cliCtx.QueryStore(record.KeyRecord(addr, ipfsHash), storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusNotFound, "no exist yet and dataHash has not been set")
			return
		}

		if len(res) == 0 || err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Data hash [%s] is not existe", hashHexStr))
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
