package prometheus

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/app"
	//"github.com/irisnet/irishub/tools"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cmn "github.com/tendermint/tendermint/libs/common"
	"log"
	"net/http"
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
	cmd.Flags().Int("port", 36660, "port to connect to")
	cmd.Flags().StringP("node", "n", "tcp://localhost:26657", "Node to connect to")
	cmd.Flags().StringP("chain-id", "c", "fuxi", "Chain ID of tendermint node")
	cmd.Flags().StringP("address", "a", "", `hex address of the validator that you want to 
monitor`)

	cmd.Flags().BoolP("recursively", "r", true, `specify whether the files in sub-directory is included, 
included by default. If there are many files & sub-directories in home directory, this program may be very slow!`)
	return cmd
}

/*
TODO:
1. 将数据轮询间隔做成可配置的
2. 修改连接不上远程node时程序崩溃的问题
3. 监控round
4. 监控上一个块的出块人
*/
