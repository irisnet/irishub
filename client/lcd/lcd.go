package lcd

import (
	"os"

	"net/http"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/context"
	govhandler "github.com/irisnet/irishub/client/gov/lcd"
	recordhandle "github.com/irisnet/irishub/client/record/lcd"
	rpchandler "github.com/irisnet/irishub/client/tendermint/rpc"
	txhandler "github.com/irisnet/irishub/client/tendermint/tx"
	stakehandler "github.com/irisnet/irishub/client/stake/lcd"
	"github.com/irisnet/irishub/client/keys"
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
		Use:   "start",
		Short: "Start IRISLCD (IRISHUB light-client daemon), a local REST server with swagger-ui: http://localhost:1317/swagger-ui/",
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

			listener, err := tmserver.StartHTTPServer(
				listenAddr, router, logger,
				tmserver.Config{MaxOpenConnections: maxOpen},
			)
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
	cmd.Flags().String(client.FlagChainID, "", "The chain ID to connect to")
	cmd.Flags().String(client.FlagNode, "tcp://localhost:26657", "Address of the node to connect to")
	cmd.Flags().Int(flagMaxOpenConnections, 1000, "The number of maximum open connections")
	cmd.Flags().Bool(client.FlagTrustNode, false, "Don't verify proofs for responses")

	return cmd
}

func createHandler(cdc *codec.Codec) *mux.Router {
	r := mux.NewRouter()

	kb, err := keys.GetKeyBase() //XXX
	if err != nil {
		panic(err)
	}

	cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout)

	r.HandleFunc("/version", CLIVersionRequestHandler).Methods("GET")
	r.HandleFunc("/node_version", NodeVersionRequestHandler(cliCtx)).Methods("GET")

	//keyshandler.RegisterRoutes(r)
	//bankhandler.RegisterRoutes(cliCtx, r, cdc)
	//slashinghandler.RegisterRoutes(cliCtx, r, cdc)
	stakehandler.RegisterRoutes(cliCtx, r, cdc, kb)
	govhandler.RegisterRoutes(cliCtx, r, cdc)
	recordhandle.RegisterRoutes(cliCtx, r, cdc)
	// tendermint apis
	rpchandler.RegisterRoutes(cliCtx, r, cdc)
	txhandler.RegisterRoutes(cliCtx, r, cdc)
	return r
}
