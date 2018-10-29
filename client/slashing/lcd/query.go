package lcd

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"net/http"
)

// http request handler to query signing info
func signingInfoHandlerFn(cliCtx context.CLIContext, storeName string, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		pk, err := sdk.GetValPubKeyBech32(vars["validator_pub"])
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		key := slashing.GetValidatorSigningInfoKey(sdk.ValAddress(pk.Address()))

		res, err := cliCtx.QueryStore(key, storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("couldn't query signing info. Error: %s", err.Error()))
			return
		}
		if len(res) ==0 {
			utils.WriteErrorResponse(w, http.StatusBadRequest,
				fmt.Sprintf("the signing information of this validator %s is empty, please make sure its existence", vars["validator_pub"]))
			return
		}

		var signingInfo slashing.ValidatorSigningInfo

		err = cdc.UnmarshalBinary(res, &signingInfo)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("couldn't decode signing info. Error: %s", err.Error()))
			return
		}

		output, err := cdc.MarshalJSONIndent(signingInfo, "", "  ")
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Write(output)
	}
}
