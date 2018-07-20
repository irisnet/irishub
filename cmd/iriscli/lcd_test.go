package main

import (
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/tools/prometheus"
	"github.com/spf13/viper"
	"testing"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"encoding/json"
	"fmt"
)

func TestRestServer(t *testing.T) {
	cdc := app.MakeCodec()
	comm := ServeCommand(cdc)
	viper.Set("chain-id", "fuxi-1001")
	viper.Set("node", "tcp://localhost:26657")
	viper.Set("laddr", "tcp://localhost:1317")

	comm.ExecuteC()
}

func TestMetricsCmd(t *testing.T) {
	cdc := app.MakeCodec()
	comm := prometheus.MonitorCommand(cdc)
	viper.Set("node", "tcp://0.0.0.0:26657")
	comm.ExecuteC()
}

func TestCoins(t *testing.T) {
	someCoins := sdk.Coins{sdk.NewCoin("atom", 123)}
	bz,_ := json.Marshal(someCoins)
	fmt.Println(string(bz))
}
