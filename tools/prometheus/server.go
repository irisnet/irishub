package prometheus

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/tools"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	cmn "github.com/tendermint/tmlibs/common"
	"log"
	"net/http"
	"github.com/irisnet/irishub/app"
)

func MonitorCommand(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "monitor",
		Short: "irishub monitor",
		RunE: func(cmd *cobra.Command, args []string) error {

			csMetrics, p2pMetrics, memMetrics, sysMetrics := DefaultMetricsProvider()
			ctx := tools.NewContext()

			//监控共识指标
			csMetrics.Monitor(ctx, cdc, storeName)
			//监控p2p指标
			p2pMetrics.Monitor(ctx)
			//监控mempool指标
			memMetrics.Monitor(ctx)
			//监控系统指标
			sysMetrics.Monitor()

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
	cmd.Flags().StringP("commands", "c", "iris start", `the processes you want to monitor that started 
by these commands, separated by semicolons ';'. 
eg: --commands="command 0;command 1;command 2", --commands=iris by default`)
	cmd.Flags().StringP("disks", "d", "/", `mounted paths of storage devices, separated by semicolons ';'. 
eg: --disks="/;/mnt1;/mnt2"`)
	cmd.Flags().StringP("paths", "p", app.DefaultNodeHome, `path to config and data files/directories, separated by semicolons ';'.
cannot use ~ and environment variables. eg: --paths="/etc;/home;
size of files in sub-directories is excluded. to compute the size recursively, you can use --recursively=true"`)
	cmd.Flags().BoolP("recursively", "r", false, `specify whether the files in sub-directories is included, 
excluded by default. If there are many files & sub-directory in given directories, this program may be very slow!`)
	return cmd
}
