package cmd

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"

	ethermintconfig "github.com/evmos/ethermint/server/config"
	etherminttypes "github.com/evmos/ethermint/x/evm/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	"github.com/irisnet/irishub/types"
	tokentypes "github.com/irisnet/irismod/modules/token/types"
	tokenv1 "github.com/irisnet/irismod/modules/token/types/v1"
)

func migrateCmd(appCodec codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: "migrate",
		Short: `The Migrate command only applies to modify the genesis data exported by version of IRISHub V1.4.1. On this basis, 
the following needs to be modified:
	1. consensusparams.block.maxgas is adjusted from -1 to 20000000.
	2. added EVM and Feemarket related initialization data.
	3. the token module adds the token definition required by EVM.
	4. add the default EVM configuration in app.toml, you can also manually modify the default configuration in app.toml.
When IRISHub successfully upgrade to version 2.0, the command will be deleted.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)
			if err := migrateAppConfig(serverCtx); err != nil {
				return err
			}

			genFile := config.GenesisFile()
			return migrateGenesis(appCodec, genFile)
		},
	}
}

func migrateAppConfig(serverCtx *server.Context) error {
	appConf, err := ethermintconfig.GetConfig(serverCtx.Viper)
	if err != nil {
		return err
	}

	customTemplate, _ := initAppConfig()

	appConf.EVM = *ethermintconfig.DefaultEVMConfig()
	appConf.JSONRPC = *ethermintconfig.DefaultJSONRPCConfig()
	appConf.TLS = *ethermintconfig.DefaultTLSConfig()

	rootDir := serverCtx.Viper.GetString(flags.FlagHome)
	configPath := filepath.Join(rootDir, "config")
	appCfgFilePath := filepath.Join(configPath, "app.toml")

	config.SetConfigTemplate(customTemplate)
	config.WriteConfigFile(appCfgFilePath, appConf)
	return nil
}

func migrateGenesis(appCodec codec.Codec, genFile string) error {
	appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
	if err != nil {
		return fmt.Errorf("failed to unmarshal genesis state: %w", err)
	}
	genDoc.ConsensusParams.Block.MaxGas = 20000000
	migrateAppState(appCodec, genDoc.InitialHeight, appState)

	appStateJSON, err := json.Marshal(appState)
	if err != nil {
		return fmt.Errorf("failed to marshal application genesis state: %w", err)
	}

	genDoc.AppState = appStateJSON
	return genutil.ExportGenesisFile(genDoc, genFile)
}

func migrateAppState(appCodec codec.Codec, initialHeight int64, appState map[string]json.RawMessage) {
	// add evm genesis
	if _, ok := appState[etherminttypes.ModuleName]; !ok {
		evmGenState := etherminttypes.GenesisState{
			Accounts: []etherminttypes.GenesisAccount{},
			Params: etherminttypes.Params{
				EvmDenom:            types.EvmToken.MinUnit,
				EnableCreate:        true,
				EnableCall:          true,
				ChainConfig:         etherminttypes.DefaultChainConfig(),
				ExtraEIPs:           nil,
				AllowUnprotectedTxs: etherminttypes.DefaultAllowUnprotectedTxs,
			},
		}
		appState[etherminttypes.ModuleName] = appCodec.MustMarshalJSON(&evmGenState)
	}

	// add feemarket genesis
	if _, ok := appState[feemarkettypes.ModuleName]; !ok {
		evmGenState := feemarkettypes.GenesisState{
			Params: feemarkettypes.Params{
				NoBaseFee:                false,
				BaseFeeChangeDenominator: 8,                        // TODO
				ElasticityMultiplier:     2,                        // TODO
				EnableHeight:             initialHeight,            // TODO
				BaseFee:                  math.NewInt(1000000000),  // TODO
				MinGasPrice:              sdk.ZeroDec(),            // TODO
				MinGasMultiplier:         sdk.NewDecWithPrec(5, 1), // TODO
			},
			BlockGas: 0,
		}
		appState[feemarkettypes.ModuleName] = appCodec.MustMarshalJSON(&evmGenState)
	}

	// add token genesis
	{
		var tokenGenState tokenv1.GenesisState
		appCodec.MustUnmarshalJSON(appState[tokentypes.ModuleName], &tokenGenState)

		evmTokenExist := false
		for _, token := range tokenGenState.Tokens {
			if token.MinUnit == types.EvmToken.MinUnit {
				evmTokenExist = true
				break
			}
		}
		if !evmTokenExist {
			tokenGenState.Tokens = append(tokenGenState.Tokens, types.EvmToken)
		}
		appState[tokentypes.ModuleName] = appCodec.MustMarshalJSON(&tokenGenState)
	}
}
