package lcd

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec) {

	r.HandleFunc("/record/records", queryRecordsWithParameterFn(cdc, cliCtx)).Methods("GET")

}
