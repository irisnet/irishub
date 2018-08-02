package prometheus

import (
	"testing"
	"github.com/spf13/viper"
	"github.com/irisnet/irishub/app"
)

func TestMetricsCmd(t *testing.T) {
	cdc := app.MakeCodec()
	comm := MonitorCommand(cdc)
	viper.Set("node", "tcp://0.0.0.0:26657")
	viper.Set("address", "92F517F241560761B5D201C0E103C1C3D5C277DB")
	viper.Set("home",  app.DefaultNodeHome)
	viper.Set("chain-id", "fuxi-test")
	viper.Set("recursively", true)
	viper.Set("port", 36660)
	comm.ExecuteC()
}
