package rpc

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
)

// Register REST endpoints.
func RegisterRoutes(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc("/v0.16.3/blocks/latest", LatestBlockRequestHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc("/v0.16.3/blocks/{height}", BlockRequestHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc("/v0.16.3/block-results/latest", LatestBlockResultRequestHandlerFn(clientCtx)).Methods("GET")
	r.HandleFunc("/v0.16.3/block-results/{height}", BlockResultRequestHandlerFn(clientCtx)).Methods("GET")
}
