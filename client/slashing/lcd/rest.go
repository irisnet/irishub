package lcd

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
)

// RegisterRoutes registers staking-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	r.HandleFunc("/slashing/signing_info/{validator}",
		signingInfoHandlerFn(cliCtx, "slashing", cdc)).Methods("GET")
	r.HandleFunc("/slashing/unrevoke",
		unrevokeRequestHandlerFn(cdc, kb, cliCtx)).Methods("POST")
}
