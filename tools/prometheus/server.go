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
	"strings"
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


			paths := viper.GetString("paths")
			commands := viper.GetString("commands")

			for _, command := range strings.Split(commands, ";"){
				sysMertrics.AddProcess(strings.TrimSpace(command))
			}

			for _, path := range strings.Split(paths, ";"){
				sysMertrics.AddDirectory(strings.TrimSpace(path))
			}

			sysMertrics.Monitor()

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
	cmd.Flags().StringP("commands", "c", "iris", `the processes you want to monitor that started 
by these commands, separated by semicolons ';'. 
eg: --commands="command 0;command 1;command 2", --commands=iris by default`)
	cmd.Flags().StringP("paths", "p", "", `directories for config and data, separated by semicolons ';'. 
eg: --paths="/;/etc/;/root"`)
	return cmd
}
