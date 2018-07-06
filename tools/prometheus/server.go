package prometheus

import (
	"github.com/spf13/cobra"
	"log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	cmn "github.com/tendermint/tmlibs/common"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire" // XXX fix
	"github.com/irisnet/irishub/tools/prometheus/consensus"
)


func MonitorCommand(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "monitor",
		Short: "irishub monitor",
		RunE: func(cmd *cobra.Command, args []string) error {
			//TODO
			csMetrics,_,_ := DefaultMetricsProvider()
			ctx := context.NewCoreContextFromViper()

			//监控共识参数
			consensus.Monitor(ctx,*csMetrics,cdc,storeName)

			srv := &http.Server{
				Addr:    ":26660",
				Handler: promhttp.Handler(),
			}
			go func() {
				if err := srv.ListenAndServe(); err != http.ErrServerClosed {
					log.Println("got ", err)
				}
			}()

			cmn.TrapSignal(func() {
				ctx.Client.Stop()
				srv.Close()
			})

			return nil
		},
	}
	cmd.Flags().StringP("node", "n", "tcp://localhost:46657", "Node to connect to")
	cmd.Flags().String("chain-id", "", "Chain ID of tendermint node")
	return cmd
}
