package cmd

import (
	"encoding/base64"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	"github.com/irisnet/irishub/v3/app"
	"github.com/spf13/cobra"
)

// ICACommands returns all of ica commands
func ICACommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ica",
		Short: "Manage your application's keys",
		Long:  ``,
	}
	cdcconf := app.RegisterEncodingConfig()

	cmd.AddCommand(build(cdcconf.Marshaler))
	return cmd
}

func build(cdc codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "build",
		Args:  cobra.ExactArgs(1),
		Short: "Build an ICA transaction",
		RunE: func(cmd *cobra.Command, args []string) error {
			var msg sdk.Msg
			if err := cdc.UnmarshalInterfaceJSON([]byte(args[0]), &msg); err != nil {
				return err
			}
			anyMsg, err := codectypes.NewAnyWithValue(msg)
			if err != nil {
				return err
			}

			tx := &icatypes.CosmosTx{
				Messages: []*codectypes.Any{anyMsg},
			}

			bz, err := cdc.Marshal(tx)
			if err != nil {
				return err
			}
			fmt.Println(base64.StdEncoding.EncodeToString(bz))
			return nil
		},
	}
}
