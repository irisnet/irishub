package lite

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client"
	assethandler "github.com/irisnet/irishub/client/asset/lcd"
	bankhandler "github.com/irisnet/irishub/client/bank/lcd"
	coinswaphandler "github.com/irisnet/irishub/client/coinswap/lcd"
	"github.com/irisnet/irishub/client/context"
	distributionhandler "github.com/irisnet/irishub/client/distribution/lcd"
	govhandler "github.com/irisnet/irishub/client/gov/lcd"
	htlchandler "github.com/irisnet/irishub/client/htlc/lcd"
	paramshandler "github.com/irisnet/irishub/client/params/lcd"
	randhandler "github.com/irisnet/irishub/client/rand/lcd"
	servicehandler "github.com/irisnet/irishub/client/service/lcd"
	slashinghandler "github.com/irisnet/irishub/client/slashing/lcd"
	stakehandler "github.com/irisnet/irishub/client/stake/lcd"
	rpchandler "github.com/irisnet/irishub/client/tendermint/rpc"
	ttxhandler "github.com/irisnet/irishub/client/tendermint/tx"
	txhandler "github.com/irisnet/irishub/client/tx/lcd"
	"github.com/irisnet/irishub/codec"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	tmserver "github.com/tendermint/tendermint/rpc/lib/server"
)

// ServeLCDStartCommand will start irislcd node, which provides rest APIs with swagger-ui
func ServeLCDStartCommand(cdc *codec.Codec) *cobra.Command {
	flagListenAddr := "laddr"
	flagCORS := "cors"
	flagMaxOpenConnections := "max-open"

	cmd := &cobra.Command{
		Use:     "start",
		Short:   "Start IRISLCD (IRISHUB light-client daemon), a local REST server with swagger-ui: http://localhost:1317/swagger-ui/",
		Example: "irislcd start --chain-id=<chain-id> --trust-node --node=tcp://localhost:26657",
		RunE: func(cmd *cobra.Command, args []string) error {
			listenAddr := viper.GetString(flagListenAddr)
			router := createHandler(cdc)

			statikFS, err := fs.New()
			if err != nil {
				panic(err)
			}
			staticServer := http.FileServer(statikFS)
			router.PathPrefix("/swagger-ui/").Handler(http.StripPrefix("/swagger-ui/", staticServer))

			logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "irislcd")
			maxOpen := viper.GetInt(flagMaxOpenConnections)

			listener, err := tmserver.Listen(
				listenAddr,
				tmserver.Config{MaxOpenConnections: maxOpen},
			)
			if err != nil {
				return err
			}

			logger.Info("Starting IRISLCD service...")

			err = tmserver.StartHTTPServer(listener, router, logger)
			if err != nil {
				return err
			}

			logger.Info("IRISLCD server started")

			// wait forever and cleanup
			cmn.TrapSignal(func() {
				err := listener.Close()
				logger.Error("error closing listener", "err", err)
			})

			return nil
		},
	}

	cmd.Flags().String(flagListenAddr, "tcp://localhost:1317", "The address for the server to listen on")
	cmd.Flags().String(flagCORS, "", "Set the domains that can make CORS requests (* for all)")
	cmd.Flags().String(client.FlagChainID, "", "Chain ID of tendermint node")
	cmd.Flags().String(client.FlagNode, "tcp://localhost:26657", "Address of the node to connect to")
	cmd.Flags().Int(flagMaxOpenConnections, 1000, "The number of maximum open connections")
	cmd.Flags().Bool(client.FlagTrustNode, false, "Don't verify proofs for responses")
	cmd.Flags().Bool(client.FlagIndentResponse, true, "Add indent to JSON response")

	return cmd
}

func createHandler(cdc *codec.Codec) *mux.Router {
	r := mux.NewRouter()

	cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout)

	r.HandleFunc("/version", CLIVersionRequestHandler).Methods("GET")
	r.HandleFunc("/node-version", NodeVersionRequestHandler(cliCtx)).Methods("GET")

	assethandler.RegisterRoutes(cliCtx, r, cdc)
	randhandler.RegisterRoutes(cliCtx, r, cdc)
	bankhandler.RegisterRoutes(cliCtx, r, cdc)
	txhandler.RegisterRoutes(cliCtx, r, cdc)
	distributionhandler.RegisterRoutes(cliCtx, r, cdc)
	slashinghandler.RegisterRoutes(cliCtx, r, cdc)
	stakehandler.RegisterRoutes(cliCtx, r, cdc)
	govhandler.RegisterRoutes(cliCtx, r, cdc)
	servicehandler.RegisterRoutes(cliCtx, r, cdc)
	paramshandler.RegisterRoutes(cliCtx, r, cdc)
	coinswaphandler.RegisterRoutes(cliCtx, r, cdc)
	htlchandler.RegisterRoutes(cliCtx, r, cdc)
	// tendermint apis
	rpchandler.RegisterRoutes(cliCtx, r, cdc)
	ttxhandler.RegisterRoutes(cliCtx, r, cdc)
	return r
}
