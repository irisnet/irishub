package prometheus

import (
	"github.com/spf13/cobra"
	"log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	cmn "github.com/tendermint/tmlibs/common"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/viper"
	"strings"
	"github.com/irisnet/irishub/tools"
)


func MonitorCommand(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "monitor",
		Short: "irishub monitor",
		RunE: func(cmd *cobra.Command, args []string) error {
			//TODO
			csMetrics,p2pMetrics,memMetrics, sysMetrics:= DefaultMetricsProvider()
			ctx := tools.NewContext()

			//监控共识参数
			csMetrics.Monitor(ctx,cdc,storeName)
			//监控p2p参数
			p2pMetrics.Monitor(ctx)
			//监控mempool参数
			memMetrics.Monitor(ctx)

			paths := viper.GetString("paths")
			commands := viper.GetString("commands")
			disks := viper.GetString("disks")

			for _, command := range strings.Split(commands, ";"){
				if strings.TrimSpace(command) != ""{
					sysMetrics.AddProcess(strings.TrimSpace(command))
				}
			}

			for _, disk_path := range strings.Split(disks, ";"){
				if strings.TrimSpace(disk_path) != ""{
					sysMetrics.AddDisk(strings.TrimSpace(disk_path))
				}
			}

			for _, path := range strings.Split(paths, ";"){
				if strings.TrimSpace(path) != ""{
					sysMetrics.AddPath(strings.TrimSpace(path))
				}
			}

			recursively := viper.GetBool("recursively")
			sysMetrics.SetRecursively(recursively)

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
	cmd.Flags().StringP("paths", "p", "", `path to config and data files/directories, separated by semicolons ';'.
cannot use ~ and environment variables. eg: --paths="/etc;/home;
size of files in sub-directories is excluded. to compute the size recursively, you can use --recursively=true"`)
	cmd.Flags().BoolP("recursively", "r", false, `specify whether the files in sub-directories is included, 
excluded by default. If there are many files & sub-directory in given directories, this program may be very slow!`)
	return cmd
}
