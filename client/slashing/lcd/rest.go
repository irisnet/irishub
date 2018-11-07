package lcd

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
)

// RegisterRoutes registers staking-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/slashing/validators/{validatorPubKey}/signing_info",
		signingInfoHandlerFn(cliCtx, "slashing", cdc)).Methods("GET")
	r.HandleFunc("/slashing/validators/{validatorAddr}/unjail",
		unrevokeRequestHandlerFn(cdc, cliCtx)).Methods("POST")
}
