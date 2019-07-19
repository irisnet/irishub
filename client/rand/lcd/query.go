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
	// Get rands with optional consumer
	r.HandleFunc(
		"/rand/rands",
		queryRandsHandlerFn(cliCtx, cdc),
	).Methods("GET")
	// Get the pending rand request from queue
	r.HandleFunc(
		"/rand/queue",
		queryQueueHandlerFn(cliCtx, cdc),
	).Methods("GET")
}

// queryRandHandlerFn performs rand query by the request id
func queryRandHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryRand(cliCtx, cdc, "custom/rand/rand/")
}

// queryRandsHandlerFn performs rands query by an optional consumer
func queryRandsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryRands(cliCtx, cdc, "custom/rand/rands")
}

// queryQueueHandlerFn performs rand request queue query by an optional heigth
func queryQueueHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryQueue(cliCtx, cdc, "custom/rand/queue")
}
