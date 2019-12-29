package token

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/irisnet/irishub/modules/asset/01-token/client/cli"
	"github.com/irisnet/irishub/modules/asset/01-token/client/rest"
)

// Name returns the token module name
func Name() string {
	return SubModuleName
}

// RegisterRESTRoutes registers the REST routes for the token module
func RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router, queryRoute string) {
	rest.RegisterRoutes(ctx, rtr, fmt.Sprintf("%s/%s", queryRoute, SubModuleName))
}

// GetTxCmd returns the root tx command for the token module
func GetTxCmd(cdc *codec.Codec, storeKey string) *cobra.Command {
	return cli.GetTxCmd(fmt.Sprintf("%s/%s", storeKey, SubModuleName), cdc)
}

// GetQueryCmd returns no root query command for the token module
func GetQueryCmd(cdc *codec.Codec, queryRoute string) *cobra.Command {
	return cli.GetQueryCmd(fmt.Sprintf("%s/%s", queryRoute, SubModuleName), cdc)
}
