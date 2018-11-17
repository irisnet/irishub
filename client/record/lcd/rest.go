package lcd

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {

	r.HandleFunc("/record/records", postRecordHandlerFn(cdc, cliCtx)).Methods("POST")

	r.HandleFunc(fmt.Sprintf("/record/records/{%s}", RestRecordID), queryRecordHandlerFn(cdc, cliCtx)).Methods("GET")

	r.HandleFunc("/record/records", queryRecordsWithParameterFn(cdc, cliCtx)).Methods("GET")

}
