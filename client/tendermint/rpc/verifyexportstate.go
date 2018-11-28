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
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
)

const (
	remoteNodeRPC = "remote-node-rpc"
	genesisFile   = "genesis-file"
	exportHeight  = "export-height"

	accountStore = "acc"
)
// VerifyExportState create a command to verify exported state file
func VerifyExportState(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify",
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
			logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
			verifier, err := newVerifier(logger, chainID, home, trustRPC, remoteRPC, exportHeight)
			if err != nil {
				return err
			}

			// Disable verifier creation in NewCLIContext
			viper.Set(client.FlagTrustNode, true)
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout)
			// Set trust-node to false
			cliCtx = cliCtx.WithCertifier(verifier).WithTrustNode(false).WithClient(remoteRPC).WithHeight(exportHeight)

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

			err = verifyExportGenesisnState(cliCtx, genesisState)
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

func verifyExportGenesisnState(cliCtx context.CLIContext, genesisState app.GenesisState) error {
	var err error
	err = verifyAccountsState(cliCtx, genesisState.Accounts)
	if err != nil {
		return err
	}
	err = verifyAuthState(cliCtx, genesisState.AuthData)
	if err != nil {
		return err
	}
	err = verifyStakeState(cliCtx, genesisState.StakeData)
	if err != nil {
		return err
	}
	err = verifyMintState(cliCtx, genesisState.MintData)
	if err != nil {
		return err
	}
	err = verifyDistrState(cliCtx, genesisState.DistrData)
	if err != nil {
		return err
	}
	err = verifyGovState(cliCtx, genesisState.GovData)
	if err != nil {
		return err
	}
	err = verifySlashingState(cliCtx, genesisState.SlashingData)
	if err != nil {
		return err
	}
	return nil
}

func verifyAccountsState(cliCtx context.CLIContext, accountsState []app.GenesisAccount) error {
	decoder := authcmd.GetAccountDecoder(cliCtx.Codec)
	for _, acc := range accountsState {
		res, err := cliCtx.QueryStore(auth.AddressStoreKey(acc.Address), accountStore)
		if err != nil {
			return err
		}
		if len(res) == 0 {
			return fmt.Errorf("account %s doesn't exist", acc.Address.String())
		}
		account, err := decoder(res)
		if err != nil {
			return fmt.Errorf("account %s: failed to decode account info", acc.Address.String())
		}
		if !acc.Coins.IsEqual(account.GetCoins()) {
			return fmt.Errorf("account %s: token amount doesn't match, expect %s, got %s", acc.Address.String(), acc.Coins.String(), account.GetCoins().String())
		}
		if acc.AccountNumber !=  account.GetAccountNumber() {
			return fmt.Errorf("account %s: account number doesn't match, expect %d, got %d", acc.Address.String(), acc.AccountNumber, account.GetAccountNumber())
		}
		if acc.Sequence != account.GetSequence() {
			return fmt.Errorf("account %s: account sequence doesn't match, expect %d, got %d", acc.Address.String(), acc.Sequence, account.GetSequence())
		}
	}
	return nil
}

func verifyAuthState(cliCtx context.CLIContext, authState auth.GenesisState) error {
	return nil
}

func verifyStakeState(cliCtx context.CLIContext, stakeState stake.GenesisState) error {
	return nil
}

func verifyMintState(cliCtx context.CLIContext, mintState mint.GenesisState) error {
	return nil
}

func verifyDistrState(cliCtx context.CLIContext, distrState distr.GenesisState) error {
	return nil
}

func verifyGovState(cliCtx context.CLIContext, govState gov.GenesisState) error {
	return nil
}

func verifySlashingState(cliCtx context.CLIContext, slashingState slashing.GenesisState) error {
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
