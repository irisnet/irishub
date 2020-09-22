package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
)

// Rest variable names
// nolint
const (
	RestServiceName      = "service-name"
	RestRequestID        = "request-id"
	RestOwner            = "owner"
	RestProvider         = "provider"
	RestConsumer         = "consumer"
	RestRequestContextID = "request-context-id"
	RestBatchCounter     = "batch-counter"
	RestArg1             = "arg1"
	RestArg2             = "arg2"
	RestSchemaName       = "schema-name"
)

// RegisterHandlers defines routes that get registered by the main application
func RegisterHandlers(cliCtx client.Context, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}
