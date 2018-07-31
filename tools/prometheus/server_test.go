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
	comm.ExecuteC()
}
