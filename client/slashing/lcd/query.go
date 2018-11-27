package lcd

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/slashing"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"net/http"
)

// http request handler to query signing info
func signingInfoHandlerFn(cliCtx context.CLIContext, storeName string, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		pk, err := sdk.GetValPubKeyBech32(vars["validatorPubKey"])
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		key := slashing.GetValidatorSigningInfoKey(sdk.ConsAddress(pk.Address()))

		res, err := cliCtx.QueryStore(key, storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("couldn't query signing info. Error: %s", err.Error()))
			return
		}
		if len(res) == 0 {
			utils.WriteErrorResponse(w, http.StatusNoContent, "")
			return
		}

		var signingInfo slashing.ValidatorSigningInfo

		err = cdc.UnmarshalBinaryLengthPrefixed(res, &signingInfo)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("couldn't decode signing info. Error: %s", err.Error()))
			return
		}

		utils.PostProcessResponse(w, cliCtx.Codec, signingInfo, cliCtx.Indent)
	}
}
