package rpc

import (
	"fmt"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/irisnet/irishub/modules/legacy/types"
)

// Register REST endpoints.
func RegisterRoutes(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/%s/blocks/latest", types.ModuleName), LatestBlockRequestHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/blocks/{height}", types.ModuleName), BlockRequestHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/block-results/latest", types.ModuleName), LatestBlockResultRequestHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/block-results/{height}", types.ModuleName), BlockResultRequestHandlerFn(clientCtx)).Methods("GET")
}
