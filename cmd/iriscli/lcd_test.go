package main

import (
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/tools/prometheus"
	"github.com/spf13/viper"
	"testing"
)

func TestRestServer(t *testing.T) {
	cdc := app.MakeCodec()
	comm := ServeCommand(cdc)
	viper.Set("chain-id", "fuxi")
	viper.Set("node", "tcp://localhost:26657")
	viper.Set("laddr", "tcp://localhost:1317")

	comm.ExecuteC()
}

func TestMetricsCmd(t *testing.T) {
	cdc := app.MakeCodec()
	comm := prometheus.MonitorCommand("stake", cdc)
	viper.Set("node", "tcp://0.0.0.0:46657")
	viper.Set("chain-id", "fuxi")
	viper.Set("commands", "iris start ;htop")
	/*
		viper.Set("paths","/etc ;/home/")
		viper.Set("recursively", false)
	*/
	viper.Set("paths", "/Users/zhangzhiqiang/.iris")
	viper.Set("disks", "/")
	viper.Set("recursively", true)
	comm.ExecuteC()
}
