package prometheus

import (
	"github.com/spf13/cobra"
	"log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	cmn "github.com/tendermint/tmlibs/common"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/tools"
	"github.com/spf13/viper"
)


func MonitorCommand(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "monitor",
		Short: "irishub monitor",
		RunE: func(cmd *cobra.Command, args []string) error {
			//TODO
			csMetrics,p2pMetrics,memMetrics, sysMertrics:= DefaultMetricsProvider()
			ctx := tools.NewContext()

			//监控共识参数
			csMetrics.Monitor(ctx,cdc,storeName)
			//监控p2p参数
			p2pMetrics.Monitor(ctx)
			//监控mempool参数
			memMetrics.Monitor(ctx)

			//monitor system info, first parameter is the command of the process to be monitor
			// and the second parameter is the directory that you want to get total size of its' files
			path := viper.GetString("home")
			sysMertrics.Monitor("iris", path)

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
	cmd.Flags().String("chain-id", "fuxi", "Chain ID of tendermint node")
	cmd.Flags().StringP("home", "", "", "directory for config and data")
	return cmd
}
