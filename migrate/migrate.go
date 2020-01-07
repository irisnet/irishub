package migrate

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/genutil"
)

const (
	flagGenesisTime = "genesis-time"
	flagChainID     = "chain-id"
)

// MigrateGenesisCmd returns a command to execute genesis state migration.
// nolint: funlen
func MigrateGenesisCmd(_ *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate [target-version] [genesis-file]",
		Short: "Migrate genesis to a specified target version",
		Long: fmt.Sprintf(`Migrate the source genesis into the target version and print to STDOUT.

Example:
$ %s migrate /path/to/genesis.json --chain-id=cosmoshub-3 --genesis-time=2019-04-22T17:00:00Z
`, version.ServerName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			importGenesis := args[0]

			genDoc, err := types.GenesisDocFromFile(importGenesis)
			if err != nil {
				return errors.Wrapf(err, "failed to read genesis document from file %s", importGenesis)
			}

			var initialState genutil.AppMap
			if err := cdc.UnmarshalJSON(genDoc.AppState, &initialState); err != nil {
				return errors.Wrap(err, "failed to JSON unmarshal initial genesis state")
			}

			//newGenState := Migrate(initialState)
			newGenState := initialState

			genDoc.AppState, err = cdc.MarshalJSON(newGenState)
			if err != nil {
				return errors.Wrap(err, "failed to JSON marshal migrated genesis state")
			}

			genesisTime := cmd.Flag(flagGenesisTime).Value.String()
			if genesisTime != "" {
				var t time.Time

				err := t.UnmarshalText([]byte(genesisTime))
				if err != nil {
					return errors.Wrap(err, "failed to unmarshal genesis time")
				}

				genDoc.GenesisTime = t
			}

			chainID := cmd.Flag(flagChainID).Value.String()
			if chainID != "" {
				genDoc.ChainID = chainID
			}

			bz, err := cdc.MarshalJSONIndent(genDoc, "", "  ")
			if err != nil {
				return errors.Wrap(err, "failed to marshal genesis doc")
			}

			sortedBz, err := sdk.SortJSON(bz)
			if err != nil {
				return errors.Wrap(err, "failed to sort JSON genesis doc")
			}

			fmt.Println(string(sortedBz))
			return nil
		},
	}

	cmd.Flags().String(flagGenesisTime, "", "override genesis_time with this flag")
	cmd.Flags().String(flagChainID, "", "override chain_id with this flag")

	return cmd
}

func Migrate(appState genutil.AppMap) genutil.AppMap {
	//v034Codec := codec.New()
	//codec.RegisterCrypto(v034Codec)
	//v034gov.RegisterCodec(v034Codec)
	//
	//v036Codec := codec.New()
	//codec.RegisterCrypto(v036Codec)
	//v036gov.RegisterCodec(v036Codec)
	//
	//// migrate genesis accounts state
	//if appState[v034genAccounts.ModuleName] != nil {
	//	var genAccs v034genAccounts.GenesisState
	//	v034Codec.MustUnmarshalJSON(appState[v034genAccounts.ModuleName], &genAccs)
	//
	//	var authGenState v034auth.GenesisState
	//	v034Codec.MustUnmarshalJSON(appState[v034auth.ModuleName], &authGenState)
	//
	//	var govGenState v034gov.GenesisState
	//	v034Codec.MustUnmarshalJSON(appState[v034gov.ModuleName], &govGenState)
	//
	//	var distrGenState v034distr.GenesisState
	//	v034Codec.MustUnmarshalJSON(appState[v034distr.ModuleName], &distrGenState)
	//
	//	var stakingGenState v034staking.GenesisState
	//	v034Codec.MustUnmarshalJSON(appState[v034staking.ModuleName], &stakingGenState)
	//
	//	delete(appState, v034genAccounts.ModuleName) // delete old key in case the name changed
	//	appState[v036genAccounts.ModuleName] = v036Codec.MustMarshalJSON(
	//		v036genAccounts.Migrate(
	//			genAccs, authGenState.CollectedFees, distrGenState.FeePool.CommunityPool, govGenState.Deposits,
	//			stakingGenState.Validators, stakingGenState.UnbondingDelegations, distrGenState.OutstandingRewards,
	//			stakingGenState.Params.BondDenom, v036distr.ModuleName, v036gov.ModuleName,
	//		),
	//	)
	//}
	//
	//// migrate auth state
	//if appState[v034auth.ModuleName] != nil {
	//	var authGenState v034auth.GenesisState
	//	v034Codec.MustUnmarshalJSON(appState[v034auth.ModuleName], &authGenState)
	//
	//	delete(appState, v034auth.ModuleName) // delete old key in case the name changed
	//	appState[v036auth.ModuleName] = v036Codec.MustMarshalJSON(v036auth.Migrate(authGenState))
	//}
	//
	//// migrate gov state
	//if appState[v034gov.ModuleName] != nil {
	//	var govGenState v034gov.GenesisState
	//	v034Codec.MustUnmarshalJSON(appState[v034gov.ModuleName], &govGenState)
	//
	//	delete(appState, v034gov.ModuleName) // delete old key in case the name changed
	//	appState[v036gov.ModuleName] = v036Codec.MustMarshalJSON(v036gov.Migrate(govGenState))
	//}
	//
	//// migrate distribution state
	//if appState[v034distr.ModuleName] != nil {
	//	var slashingGenState v034distr.GenesisState
	//	v034Codec.MustUnmarshalJSON(appState[v034distr.ModuleName], &slashingGenState)
	//
	//	delete(appState, v034distr.ModuleName) // delete old key in case the name changed
	//	appState[v036distr.ModuleName] = v036Codec.MustMarshalJSON(v036distr.Migrate(slashingGenState))
	//}
	//
	//// migrate staking state
	//if appState[v034staking.ModuleName] != nil {
	//	var stakingGenState v034staking.GenesisState
	//	v034Codec.MustUnmarshalJSON(appState[v034staking.ModuleName], &stakingGenState)
	//
	//	delete(appState, v034staking.ModuleName) // delete old key in case the name changed
	//	appState[v036staking.ModuleName] = v036Codec.MustMarshalJSON(v036staking.Migrate(stakingGenState))
	//}
	//
	//// migrate supply state
	//appState[v036supply.ModuleName] = v036Codec.MustMarshalJSON(v036supply.EmptyGenesisState())
	//
	//return appState
}
