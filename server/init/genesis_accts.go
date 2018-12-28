package init

import (
	"fmt"
	"os"

	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/server"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/irisnet/irishub/app/v0"
)

// AddGenesisAccountCmd returns add-genesis-account cobra Command
func AddGenesisAccountCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-account [address] [coin][,[coin]]",
		Short: "Add genesis account to genesis.json",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))
			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(utils.GetAccountDecoder(cdc))

			coins, err := cliCtx.ParseCoins(args[1])
			if err != nil {
				return err
			}
			coins.Sort()
			genFile := config.GenesisFile()
			if !common.FileExists(genFile) {
				return fmt.Errorf("%s does not exist, run `iris init` first", genFile)
			}
			genDoc, err := loadGenesisDoc(cdc, genFile)
			if err != nil {
				return err
			}
			var genesisState v0.GenesisFileState
			if err = cdc.UnmarshalJSON(genDoc.AppState, &genesisState); err != nil {
				return err
			}
			acc := auth.NewBaseAccountWithAddress(addr)
			acc.Coins = coins
			genesisState.Accounts = append(genesisState.Accounts, v0.NewGenesisFileAccount(&acc))
			appStateJSON, err := cdc.MarshalJSON(genesisState)
			if err != nil {
				return err
			}
			return ExportGenesisFile(genFile, genDoc.ChainID, nil, appStateJSON)
		},
	}
	cmd.Flags().String(cli.HomeFlag, app.DefaultNodeHome, "node's home directory")
	return cmd
}
