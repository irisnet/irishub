package init

import (
	"encoding/json"
	v2 "github.com/irisnet/irishub/app/v2"
	"path/filepath"

	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/types"
)

type initConfig struct {
	ChainID   string
	GenTxsDir string
	Name      string
	NodeID    string
	ValPubKey crypto.PubKey
}

// nolint
func CollectGenTxsCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collect-gentxs",
		Short: "Collect genesis txs and output a genesis.json file",
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))
			name := viper.GetString(client.FlagName)

			nodeID, valPubKey, err := LoadNodeValidatorFiles(config)
			if err != nil {
				return err
			}

			genDoc, err := loadGenesisDoc(cdc, config.GenesisFile())
			if err != nil {
				return err
			}

			toPrint := printInfo{
				Moniker: config.Moniker,
				ChainID: genDoc.ChainID,
				NodeID:  nodeID,
			}

			initCfg := initConfig{
				ChainID:   genDoc.ChainID,
				GenTxsDir: filepath.Join(config.RootDir, "config", "gentx"),
				Name:      name,
				NodeID:    nodeID,
				ValPubKey: valPubKey,
			}

			appMessage, err := genAppStateFromConfig(cdc, config, initCfg, genDoc)
			if err != nil {
				return err
			}

			toPrint.AppMessage = appMessage

			// print out some key information
			return displayInfo(cdc, toPrint)
		},
	}

	cmd.Flags().String(cli.HomeFlag, app.DefaultNodeHome, "node's home directory")
	return cmd
}

func genAppStateFromConfig(
	cdc *codec.Codec, config *cfg.Config, initCfg initConfig, genDoc types.GenesisDoc,
) (appState json.RawMessage, err error) {

	genFile := config.GenesisFile()
	var (
		appGenTxs       []auth.StdTx
		persistentPeers string
		genTxs          []json.RawMessage
		jsonRawTx       json.RawMessage
	)

	// process genesis transactions, else create default genesis.json
	appGenTxs, persistentPeers, err = v2.CollectStdTxs(
		cdc, config.Moniker, initCfg.GenTxsDir, genDoc,
	)
	if err != nil {
		return
	}

	genTxs = make([]json.RawMessage, len(appGenTxs))
	config.P2P.PersistentPeers = persistentPeers

	for i, stdTx := range appGenTxs {
		jsonRawTx, err = cdc.MarshalJSON(stdTx)
		if err != nil {
			return
		}
		genTxs[i] = jsonRawTx
	}

	cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)

	appState, err = v2.IrisAppGenStateJSON(cdc, genDoc, genTxs)
	if err != nil {
		return
	}

	err = ExportGenesisFile(genFile, initCfg.ChainID, nil, appState)
	return
}
