package cli

import (
	"bufio"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/irisnet/irishub/modules/rand"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	randTxCmd := &cobra.Command{
		Use:                        rand.ModuleName,
		Short:                      "Rand transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	randTxCmd.AddCommand(client.PostCommands(
		GetCmdRequestRand(cdc),
	)...)

	return randTxCmd
}

// GetCmdRequestRand implements the request-rand command
func GetCmdRequestRand(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "request-rand",
		Short:   "Request a random number",
		Example: "iriscli rand request-rand --block-interval=100",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			consumer, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := rand.NewMsgRequestRand(consumer, uint64(viper.GetInt64(FlagBlockInterval)))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsRequestRand)

	return cmd
}
