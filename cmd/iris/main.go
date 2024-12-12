package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	_ "github.com/irisnet/irishub/v4/client/lite/statik"
	"github.com/irisnet/irishub/v4/cmd/iris/cmd"
	"github.com/irisnet/irishub/v4/types"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, "", types.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
