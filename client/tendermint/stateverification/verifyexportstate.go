package stateverification

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/context"
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

	accountStore = "acc"
	authStore    = "fee"
)

// VerifyExportState create a command to verify exported state file
func VerifyExportState(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "start to verify exported genesis state",
		RunE: func(cmd *cobra.Command, args []string) error {
			exportHeight := viper.GetInt64(exportHeight)
			if exportHeight <= 0 {
				return fmt.Errorf("export height shoule be greater than 0")
			}
			chainID := viper.GetString(client.FlagChainID)
			home := viper.GetString(cli.HomeFlag)
			genesisFile := viper.GetString(genesisFile)
			trustURL := viper.GetString(client.FlagNode)
			remoteURL := viper.GetString(remoteNodeRPC)

			trustRPC := rpcclient.NewHTTP(trustURL, "/websocket")
			remoteRPC := rpcclient.NewHTTP(remoteURL, "/websocket")
			logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
			verifier, err := newVerifier(logger, chainID, home, trustRPC, remoteRPC, exportHeight)
			if err != nil {
				return fmt.Errorf("encountered error in creating verifier: %s", err.Error())
			}

			// Disable verifier creation in NewCLIContext
			viper.Set(client.FlagTrustNode, true)
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout)
			// Set trust-node to false
			cliCtx = cliCtx.WithCertifier(verifier).WithTrustNode(false).WithClient(remoteRPC).WithHeight(exportHeight)

			// Parser exported genesis file
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

			// Verify genesis state
			err = verifyExportGenesisnState(logger, cliCtx, genesisState)
			if err != nil {
				return err
			}
			logger.Info("Verification passed")
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
	if fc.SignedHeader.Header != nil && fc.Height() < exportHeight || err != nil {
		newFc, err := trustSource.LatestFullCommit(chainID, 1, exportHeight)
		if err != nil {
			return nil, cmn.ErrorWrap(err, "fetching source full commit @ height 1")
		}
		err = trust.SaveFullCommit(newFc)
		if err != nil {
			return nil, cmn.ErrorWrap(err, "saving full commit to trusted")
		}
	}

	return verifier, nil
}

func verifyExportGenesisnState(logger log.Logger, cliCtx context.CLIContext, genesisState app.GenesisState) error {
	var err error
	err = verifyAccountsState(logger, cliCtx, genesisState.Accounts)
	if err != nil {
		return err
	}
	err = verifyAuthState(logger, cliCtx, genesisState.AuthData)
	if err != nil {
		return err
	}
	err = verifyStakeState(logger, cliCtx, genesisState.StakeData)
	if err != nil {
		return err
	}
	err = verifyMintState(logger, cliCtx, genesisState.MintData)
	if err != nil {
		return err
	}
	err = verifyDistrState(logger, cliCtx, genesisState.DistrData)
	if err != nil {
		return err
	}
	err = verifyGovState(logger, cliCtx, genesisState.GovData)
	if err != nil {
		return err
	}
	err = verifySlashingState(logger, cliCtx, genesisState.SlashingData)
	if err != nil {
		return err
	}
	return nil
}
