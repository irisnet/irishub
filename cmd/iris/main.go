package main

import (
	"os"

	"github.com/irisnet/irishub/cmd/iris/cmd"
	_ "github.com/irisnet/irishub/lite/statik"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
