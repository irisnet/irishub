package lcd

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/codec"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// Get rand by the request id
	r.HandleFunc(
		"/rand/rands/{request-id}",
		queryRandHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get the pending rand requests from queue
	r.HandleFunc(
		"/rand/queue",
		queryQueueHandlerFn(cliCtx, cdc),
	).Methods("GET")
}

// queryRandHandlerFn performs rand query by the request id
func queryRandHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryRand(cliCtx, cdc, "custom/rand/rand/")
}

// queryQueueHandlerFn performs rand request queue query by an optional heigth
func queryQueueHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryQueue(cliCtx, cdc, "custom/rand/queue")
}
