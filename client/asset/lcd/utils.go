package lcd

import (
	"net/http"

	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/codec"
)

// contains checks if the a given query contains one of the tx types
func contains(stringSlice []string, txType string) bool {
	for _, word := range stringSlice {
		if word == txType {
			return true
		}
	}
	return false
}

// queryGateway queries a gateway of the given moniker from the specified endpoint
func queryGateway(cliCtx context.CLIContext, cdc *codec.Codec, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// queryGateways queries all gateways of an owner from the specified endpoint
func queryGateways(cliCtx context.CLIContext, cdc *codec.Codec, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}
