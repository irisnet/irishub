package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/types/module"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cobra"

	"github.com/irisnet/irishub/v5/app"
)

// Keep the SDK command's flags and aliases, with application-wide validation
// because legacy LSM conversion must inspect more than the staking module.
func validateGenesisCmd(basics module.BasicManager) *cobra.Command {
	cmd := genutilcli.ValidateGenesisCmd(basics)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		path := server.GetServerContextFromCmd(cmd).Config.GenesisFile()
		if len(args) != 0 {
			path = args[0]
		}
		genesis, err := genutiltypes.AppGenesisFromFile(path)
		if err != nil {
			return err
		}
		if err := genesis.ValidateAndComplete(); err != nil {
			return err
		}
		var state map[string]json.RawMessage
		if err := json.Unmarshal(genesis.AppState, &state); err != nil {
			return err
		}
		clientCtx := client.GetClientContextFromCmd(cmd)
		if err := app.ValidateGenesis(clientCtx.Codec, clientCtx.TxConfig, basics, state); err != nil {
			return fmt.Errorf("error validating genesis file %s: %w", path, err)
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "File at %s is a valid genesis file\n", path)
		return err
	}
	return cmd
}
