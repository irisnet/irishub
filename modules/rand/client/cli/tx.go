package cli

import (
	"bufio"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/irisnet/irishub/modules/rand"
	"github.com/spf13/cobra"
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
		Short:   "Request a random number with an optional block interval",
		Example: "iriscli tx rand request-rand [block-interval]",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			var blockInterval uint64
			var err error

			if len(args) > 0 {
				blockInterval, err = strconv.ParseUint(args[0], 10, 64)
				if err != nil {
					return err
				}
			}

			consumer, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := rand.NewMsgRequestRand(consumer, blockInterval)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
