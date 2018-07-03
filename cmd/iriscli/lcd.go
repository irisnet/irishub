package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tmlibs/log"

	tmserver "github.com/tendermint/tendermint/rpc/lib/server"
	cmn "github.com/tendermint/tmlibs/common"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/irisnet/irishub/version"
	"github.com/cosmos/cosmos-sdk/wire"
	auth "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	bank "github.com/cosmos/cosmos-sdk/x/bank/client/rest"
	ibc "github.com/cosmos/cosmos-sdk/x/ibc/client/rest"
	stake "github.com/cosmos/cosmos-sdk/x/stake/client/rest"
	"github.com/irisnet/irishub/app"
)

// ServeCommand will generate a long-running rest server
// (aka Light Client Daemon) that exposes functionality similar
// to the cli, but over rest
func ServeCommand(cdc *wire.Codec) *cobra.Command {
	flagListenAddr := "laddr"
	flagCORS := "cors"

	cmd := &cobra.Command{
		Use:   "rest-server",
		Short: "Start LCD (light-client daemon), a local REST server",
		RunE: func(cmd *cobra.Command, args []string) error {
			listenAddr := viper.GetString(flagListenAddr)
			handler := createHandler(cdc)
			logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).
				With("module", "rest-server")
			listener, err := tmserver.StartHTTPServer(listenAddr, handler, logger)
			if err != nil {
				return err
			}
			logger.Info("REST server started")

			// Wait forever and cleanup
			cmn.TrapSignal(func() {
				err := listener.Close()
				logger.Error("Error closing listener", "err", err)
			})
			return nil
		},
	}
	cmd.Flags().StringP(flagListenAddr, "a", "tcp://localhost:1317", "Address for server to listen on")
	cmd.Flags().String(flagCORS, "", "Set to domains that can make CORS requests (* for all)")
	cmd.Flags().StringP(client.FlagChainID, "c", "", "ID of chain we connect to")
	cmd.Flags().StringP(client.FlagNode, "n", "tcp://localhost:46657", "Node to connect to")
	return cmd
}

func createHandler(cdc *wire.Codec) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/version", version.RequestHandler).Methods("GET")

	kb, err := keys.GetKeyBase() //XXX
	if err != nil {
		panic(err)
	}

	ctx := app.NewContext()

	// TODO make more functional? aka r = keys.RegisterRoutes(r)
	keys.RegisterRoutes(r)
	rpc.RegisterRoutes(ctx.GetCosmosCtx(), r)
	tx.RegisterRoutes(ctx.GetCosmosCtx(), r, cdc)
	auth.RegisterRoutes(ctx.GetCosmosCtx(), r, cdc, "acc")
	bank.RegisterRoutes(ctx.GetCosmosCtx(), r, cdc, kb)
	ibc.RegisterRoutes(ctx.GetCosmosCtx(), r, cdc, kb)
	stake.RegisterRoutes(ctx.GetCosmosCtx(), r, cdc, kb)
	RegisterRoutes(ctx,r,cdc,kb)
	return r
}
