package rpc

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	log "github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/lite"
	lclient "github.com/tendermint/tendermint/lite/client"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

const (
	remoteNodeRPC = "remote-node-rpc"
	genesisFile   = "genesis-file"
	exportHeight  = "export-height"
)

func VerifyExportState(cdc *codec.Codec) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "start",
		Short: "start to verify exported genesis state",
		RunE: func(cmd *cobra.Command, args []string) error {
			exportHeight := viper.GetInt64(exportHeight)
			if exportHeight <= 0 {
				return fmt.Errorf("missing export height")
			}

			genesisFile := viper.GetString(genesisFile)
			if _, err := os.Stat(genesisFile); os.IsNotExist(err) {
				return fmt.Errorf("invalid genesis file path")
			}

			remoteURL := viper.GetString(remoteNodeRPC)
			if len(remoteURL) == 0 {
				return fmt.Errorf("missing remote-node-rpc")
			}
			remoteRPC := rpcclient.NewHTTP(remoteURL, "/websocket")

			viper.Set(client.FlagTrustNode, true)
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout)

			cliCtx = cliCtx.WithClient(remoteRPC)
			_ = cliCtx
			return nil
		},
	}
	cmd.Flags().String(remoteNodeRPC, "", "Remote node rpc")
	cmd.Flags().String(genesisFile, "", "Exported genesis file path")
	cmd.Flags().Int64(exportHeight, 0, "Exported height")
	cmd.Flags().String(client.FlagChainID, "", "Chain ID of tendermint node")
	cmd.Flags().String(client.FlagNode, "tcp://localhost:26657", "Trusted node rpc")
	cmd.MarkFlagFilename(genesisFile)
	return cmd
}

func NewVerifier(chainID, rootDir string, client lclient.SignStatusClient, logger log.Logger, cacheSize int) (*lite.DynamicVerifier, error) {

	logger = logger.With("module", "lite/proxy")
	logger.Info("lite/proxy/NewVerifier()...", "chainID", chainID, "rootDir", rootDir, "client", client)

	memProvider := lite.NewDBProvider("trusted.mem", dbm.NewMemDB()).SetLimit(cacheSize)
	lvlProvider := lite.NewDBProvider("trusted.lvl", dbm.NewDB("trust-base", dbm.LevelDBBackend, rootDir))
	trust := lite.NewMultiProvider(
		memProvider,
		lvlProvider,
	)
	source := lclient.NewProvider(chainID, client)
	cert := lite.NewDynamicVerifier(chainID, trust, source)
	cert.SetLogger(logger) // Sets logger recursively.

	// TODO: Make this more secure, e.g. make it interactive in the console?
	_, err := trust.LatestFullCommit(chainID, 1, 1<<63-1)
	if err != nil {
		logger.Info("lite/proxy/NewVerifier found no trusted full commit, initializing from source from height 1...")
		fc, err := source.LatestFullCommit(chainID, 1, 1)
		if err != nil {
			return nil, cmn.ErrorWrap(err, "fetching source full commit @ height 1")
		}
		err = trust.SaveFullCommit(fc)
		if err != nil {
			return nil, cmn.ErrorWrap(err, "saving full commit to trusted")
		}
	}

	return cert, nil
}
