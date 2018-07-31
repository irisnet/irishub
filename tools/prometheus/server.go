package prometheus

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/app"
  //"github.com/irisnet/irishub/tools"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	cmn "github.com/tendermint/tendermint/libs/common"
	"log"
	"net/http"
  "github.com/spf13/viper"
	"fmt"
)

func MonitorCommand(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "monitor",
		Short: "irishub monitor",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := app.NewContext().WithCodeC(cdc)
      //ctx := tools.NewContext(storeName, cdc)
			monitor := DefaultMonitor(ctx)
			monitor.Start()

			port := viper.GetInt("port")
			srv := &http.Server{
				Addr:    fmt.Sprintf(":%d", port),
				Handler: promhttp.Handler(),
			}
			go func() {
				if err := srv.ListenAndServe(); err != http.ErrServerClosed {
					log.Println("got ", err)
				}
			}()

			cmn.TrapSignal(func() {
				srv.Close()
			})

			return nil
		},
	}
	cmd.Flags().Int("port",  36660, "port to connect to")
	cmd.Flags().StringP("node", "n", "tcp://localhost:46657", "Node to connect to")
	cmd.Flags().String("chain-id", "fuxi", "Chain ID of tendermint node")
	cmd.Flags().StringP("address", "a", "", `hex address of the validator that you want to 
monitor`)

	cmd.Flags().String("irishome", app.DefaultNodeHome, fmt.Sprintf("your iris home, %s by default",
		app.DefaultNodeHome))

	cmd.Flags().StringP("commands", "c", "iris start", `the processes you want to monitor that started 
by these commands, separated by semicolons ';'. 
eg: --commands="command 0;command 1;command 2", --commands=iris by default`)
	cmd.Flags().StringP("disks", "d", "/", `mounted paths of storage devices, separated by semicolons ';'. 
eg: --disks="/;/mnt1;/mnt2"`)
	cmd.Flags().StringP("paths", "p", app.DefaultNodeHome, `path to config and data files/directories, separated by semicolons ';'.
cannot use ~ and environment variables. eg: --paths="/etc;/home;
size of files in sub-directories is excluded. to compute the size recursively, you can use --recursively=true`)
	cmd.Flags().BoolP("recursively", "r", false, `specify whether the files in sub-directories is included, 
excluded by default. If there are many files & sub-directory in given directories, this program may be very slow!`)
	return cmd
}
/*
TODO:
1.  修改连接不上远程node时程序崩溃的问题
2.  使用prometheus的{}中的选项
3.  监控round

1. 将数据轮询间隔做成可配置的
2. 将数据轮询间隔做成可配置的

如果有出块奖励的话：
	1.    监控上一个块的出块人
*/
