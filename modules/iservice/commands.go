package iservice

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/iris-hub/modules/iservice/bind"
	"github.com/irisnet/iris-hub/modules/iservice/def"
	"github.com/spf13/cobra"
)

func Commands(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "iservice",
		Short: "service for irisnet",
		//Long: ``,
	}
	cmd.AddCommand(
		client.PostCommands(def.SvcDefTxCmd(cdc))...,
	)
	cmd.AddCommand(
		client.PostCommands(bind.SvcBindTxCmd(cdc))...,
	)

	return cmd
}
