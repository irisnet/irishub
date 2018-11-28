package rpc

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/modules/gov"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/lite"
	lclient "github.com/tendermint/tendermint/lite/client"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tmlibs/cli"
)

const (
	remoteNodeRPC = "remote-node-rpc"
	genesisFile   = "genesis-file"
	exportHeight  = "export-height"
)
// VerifyExportState create a command to verify exported state file
func VerifyExportState(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start to verify exported genesis state",
		RunE: func(cmd *cobra.Command, args []string) error {
			chainID := viper.GetString(client.FlagChainID)
			if len(chainID) == 0 {
				return fmt.Errorf("missing chain-id")
			}

			home := viper.GetString(cli.HomeFlag)
			if len(home) == 0 {
				return fmt.Errorf("missing home")
			}

			exportHeight := viper.GetInt64(exportHeight)
			if exportHeight <= 0 {
				return fmt.Errorf("missing export height")
			}

			genesisFile := viper.GetString(genesisFile)
			if _, err := os.Stat(genesisFile); os.IsNotExist(err) {
				return fmt.Errorf("invalid genesis file path")
			}

			trustURL := viper.GetString(client.FlagNode)
			if len(trustURL) == 0 {
				return fmt.Errorf("missing trusted node rpc")
			}
			trustRPC := rpcclient.NewHTTP(trustURL, "/websocket")

			remoteURL := viper.GetString(remoteNodeRPC)
			if len(remoteURL) == 0 {
				return fmt.Errorf("missing remote node rpc")
			}
			remoteRPC := rpcclient.NewHTTP(remoteURL, "/websocket")

			verifier, err := newVerifier(log.NewNopLogger(), chainID, home, trustRPC, remoteRPC, exportHeight)
			if err != nil {
				return err
			}

			viper.Set(client.FlagTrustNode, true)
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout)
			cliCtx = cliCtx.WithCertifier(verifier).WithTrustNode(false).WithClient(remoteRPC)

			genContents, err := ioutil.ReadFile(genesisFile)
			if err != nil {
				return err
			}
			var genDoc types.GenesisDoc
			if err := cdc.UnmarshalJSON(genContents, &genDoc); err != nil {
				return err
			}
			var genesisFileState app.GenesisFileState
			err = cdc.UnmarshalJSON(genDoc.AppState, &genesisFileState)
			if err != nil {
				return err
			}
			genesisState := app.ConvertToGenesisState(genesisFileState)

			err = verifyExportGenesisnState(cliCtx, genesisState, exportHeight)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String(remoteNodeRPC, "", "Remote node rpc")
	cmd.Flags().String(genesisFile, "", "Exported genesis file path")
	cmd.Flags().Int64(exportHeight, 0, "Exported height")
	cmd.Flags().String(client.FlagChainID, "", "Chain ID of tendermint node")
	cmd.Flags().String(client.FlagNode, "tcp://localhost:26657", "Trusted node rpc")
	cmd.MarkFlagFilename(genesisFile)
	cmd.MarkFlagRequired(remoteNodeRPC)
	cmd.MarkFlagRequired(genesisFile)
	cmd.MarkFlagRequired(exportHeight)
	cmd.MarkFlagRequired(client.FlagChainID)
	return cmd
}

func verifyExportGenesisnState(cliCtx context.CLIContext, genesisState app.GenesisState, exportHeight int64) error {
	var err error
	err = verifyAccountsState(cliCtx, genesisState.Accounts, exportHeight)
	if err != nil {
		return err
	}
	err = verifyAuthState(cliCtx, genesisState.AuthData, exportHeight)
	if err != nil {
		return err
	}
	err = verifyStakeState(cliCtx, genesisState.StakeData, exportHeight)
	if err != nil {
		return err
	}
	err = verifyMintState(cliCtx, genesisState.MintData, exportHeight)
	if err != nil {
		return err
	}
	err = verifyDistrState(cliCtx, genesisState.DistrData, exportHeight)
	if err != nil {
		return err
	}
	err = verifyGovState(cliCtx, genesisState.GovData, exportHeight)
	if err != nil {
		return err
	}
	err = verifySlashingState(cliCtx, genesisState.SlashingData, exportHeight)
	if err != nil {
		return err
	}
	return nil
}

func verifyAccountsState(cliCtx context.CLIContext, accountsState []app.GenesisAccount, exportHeight int64) error {
	return nil
}

func verifyAuthState(cliCtx context.CLIContext, AuthState auth.GenesisState, exportHeight int64) error {
	return nil
}

func verifyStakeState(cliCtx context.CLIContext, stakeState stake.GenesisState, exportHeight int64) error {
	return nil
}

func verifyMintState(cliCtx context.CLIContext, mintState mint.GenesisState, exportHeight int64) error {
	return nil
}

func verifyDistrState(cliCtx context.CLIContext, distrState distr.GenesisState, exportHeight int64) error {
	return nil
}

func verifyGovState(cliCtx context.CLIContext, govState gov.GenesisState, exportHeight int64) error {
	return nil
}

func verifySlashingState(cliCtx context.CLIContext, slashingState slashing.GenesisState, exportHeight int64) error {
	return nil
}

func newVerifier(logger log.Logger, chainID, rootDir string, trustClient lclient.SignStatusClient, remoteClient lclient.SignStatusClient, exportHeight int64) (*lite.DynamicVerifier, error) {

	memProvider := lite.NewDBProvider("trusted.mem", dbm.NewMemDB()).SetLimit(10)
	lvlProvider := lite.NewDBProvider("trusted.lvl", dbm.NewDB("trust-base", dbm.LevelDBBackend, rootDir))
	trust := lite.NewMultiProvider(
		memProvider,
		lvlProvider,
	)
	trustSource := lclient.NewProvider(chainID, trustClient)
	remoteSource := lclient.NewProvider(chainID, remoteClient)
	verifier := lite.NewDynamicVerifier(chainID, trust, remoteSource)
	verifier.SetLogger(logger) // Sets logger recursively.

	fc, err := trust.LatestFullCommit(chainID, 1, exportHeight)
	if fc.Height() < exportHeight || err != nil {
		fc, err := trustSource.LatestFullCommit(chainID, 1, exportHeight)
		if err != nil {
			return nil, cmn.ErrorWrap(err, "fetching source full commit @ height 1")
		}
		err = trust.SaveFullCommit(fc)
		if err != nil {
			return nil, cmn.ErrorWrap(err, "saving full commit to trusted")
		}
	}

	return verifier, nil
}
