package main

import (
	_ "github.com/irisnet/irishub/client/lcd/statik"
	"github.com/spf13/cobra"
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/client/lcd"
)

// rootCmd is the entry point for this binary
var (
	rootCmd = &cobra.Command{
		Use:   "irislcd",
		Short: "irishub lite server interface",
	}
)

func main() {
	cobra.EnableCommandSorting = false
	cdc := app.MakeCodec()

	rootCmd.AddCommand(
		lcd.ServeLCDCommand(cdc),
	)
}

