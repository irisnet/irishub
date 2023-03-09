package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"

	etherminttypes "github.com/evmos/ethermint/x/evm/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	"github.com/irisnet/irishub/types"
	tokentypes "github.com/irisnet/irismod/modules/token/types"
)

func genesisMigrateCmd(appCodec codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: "genesis-migrate",
		Short: `The Migrate command only applies to modify the GENESIS data exported by version of IRISHub V1.4.1. On this basis, t
		he following needs to be modified:
			1. consensusparams.block.maxgas is adjusted from -1 to 20000000.
			2. Added EVM and Feemarket related initialization data.
			3. The token module adds the token definition required by EVM.
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
			genFile := config.GenesisFile()
			return migrateGenesis(appCodec, genFile)
		},
	}
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
		var tokenGenState tokentypes.GenesisState
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
