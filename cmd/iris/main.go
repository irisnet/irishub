package main

import (
	"os"

	"github.com/irisnet/irishub/cmd/iris/cmd"
	_ "github.com/irisnet/irishub/lite/statik"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := cmd.Execute(rootCmd); err != nil {
		os.Exit(1)
	}
}
