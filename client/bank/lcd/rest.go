package lcd

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/wire"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	r.HandleFunc(fmt.Sprintf("/accounts/{address}/send/{%s}/{%s}", utils.BroadcastSync, utils.GenerateOnly),
		SendRequestHandlerFn(cdc, kb, cliCtx)).Methods("POST")
	r.HandleFunc("/accounts/{address}",
		QueryAccountRequestHandlerFn("acc", cdc, authcmd.GetAccountDecoder(cdc), cliCtx)).Methods("GET")
	r.HandleFunc("/accounts/{coin-type}",
		QueryCoinTypeRequestHandlerFn(cdc, cliCtx)).Methods("GET")
}
