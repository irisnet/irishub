package cmd

// DONTCOVER

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	evmtypes "github.com/evmos/ethermint/x/evm/types"
	tokentypes "github.com/irisnet/irismod/modules/token/types"

	"github.com/spf13/cobra"

	tmconfig "github.com/cometbft/cometbft/config"
	tmos "github.com/cometbft/cometbft/libs/os"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	"github.com/cometbft/cometbft/types"
	tmtime "github.com/cometbft/cometbft/types/time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/server"
	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	servercfg "github.com/evmos/ethermint/server/config"

	guardiantypes "github.com/irisnet/irishub/v2/modules/guardian/types"
	iristypes "github.com/irisnet/irishub/v2/types"
	randomtypes "github.com/irisnet/irismod/modules/random/types"
	servicetypes "github.com/irisnet/irismod/modules/service/types"
	tokentypesv1 "github.com/irisnet/irismod/modules/token/types/v1"
)

var (
	flagNodeDirPrefix     = "node-dir-prefix"
	flagNumValidators     = "v"
	flagOutputDir         = "output-dir"
	flagNodeDaemonHome    = "node-daemon-home"
	flagNodeCLIHome       = "node-cli-home"
	flagStartingIPAddress = "starting-ip-address"
)

const nativeIrisMinUnit = "uiris"

var PowerReduction = sdk.NewIntFromUint64(1000000000000000000)

// get cmd to initialize all files for tendermint testnet and application
func testnetCmd(
	mbm module.BasicManager,
	genBalIterator banktypes.GenesisBalancesIterator,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "testnet",
		Short: "Initialize files for a simapp testnet",
		Long: `testnet will create "v" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).
Note, strict routability for addresses is turned off in the config file.
Example:
	iris testnet --v 4 --output-dir ./output --starting-ip-address 192.168.10.2
	`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			outputDir, _ := cmd.Flags().GetString(flagOutputDir)
			keyringBackend, _ := cmd.Flags().GetString(flags.FlagKeyringBackend)
			chainID, _ := cmd.Flags().GetString(flags.FlagChainID)
			minGasPrices, _ := cmd.Flags().GetString(server.FlagMinGasPrices)
			nodeDirPrefix, _ := cmd.Flags().GetString(flagNodeDirPrefix)
			nodeDaemonHome, _ := cmd.Flags().GetString(flagNodeDaemonHome)
			nodeCLIHome, _ := cmd.Flags().GetString(flagNodeCLIHome)
			startingIPAddress, _ := cmd.Flags().GetString(flagStartingIPAddress)
			numValidators, _ := cmd.Flags().GetInt(flagNumValidators)
			algo, _ := cmd.Flags().GetString(flags.FlagKeyType)

			return InitTestnet(
				clientCtx,
				cmd,
				config,
				mbm,
				genBalIterator,
				outputDir,
				chainID,
				minGasPrices,
				nodeDirPrefix,
				nodeDaemonHome,
				nodeCLIHome,
				startingIPAddress,
				keyringBackend,
				algo,
				numValidators,
			)
		},
	}

	cmd.Flags().Int(flagNumValidators, 4, "Number of validators to initialize the testnet with")
	cmd.Flags().
		StringP(flagOutputDir, "o", "./mytestnet", "Directory to store initialization data for the testnet")
	cmd.Flags().
		String(flagNodeDirPrefix, "node", "Prefix the directory name for each node with (node results in node0, node1, ...)")
	cmd.Flags().
		String(flagNodeDaemonHome, "iris", "Home directory of the node's daemon configuration")
	cmd.Flags().String(flagNodeCLIHome, "iriscli", "Home directory of the node's cli configuration")
	cmd.Flags().
		String(flagStartingIPAddress, "192.168.0.1", "Starting IP address (192.168.0.1 results in persistent peers list ID0@192.168.0.1:46656, ID1@192.168.0.2:46656, ...)")
	cmd.Flags().
		String(flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().
		String(server.FlagMinGasPrices, fmt.Sprintf("0.000006%s", sdk.DefaultBondDenom), "Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum (e.g. 0.01photino,0.001stake)")
	cmd.Flags().
		String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	cmd.Flags().
		String(flags.FlagKeyType, string(hd.Secp256k1Type), "Key signing algorithm to generate keys for")

	return cmd
}

const nodeDirPerm = 0755

// InitTestnet the testnet
func InitTestnet(
	clientCtx client.Context,
	cmd *cobra.Command,
	nodeConfig *tmconfig.Config,
	mbm module.BasicManager,
	genBalIterator banktypes.GenesisBalancesIterator,
	outputDir,
	chainID,
	minGasPrices,
	nodeDirPrefix,
	nodeDaemonHome,
	nodeCLIHome,
	startingIPAddress,
	keyringBackend,
	algoStr string,
	numValidators int,
) error {
	if chainID == "" {
		chainID = "chain-" + tmrand.Str(6)
	}

	nodeIDs := make([]string, numValidators)
	valPubKeys := make([]cryptotypes.PubKey, numValidators)

	simappConfig := servercfg.DefaultConfig()
	simappConfig.MinGasPrices = minGasPrices
	simappConfig.API.Enable = true
	simappConfig.Telemetry.Enabled = true
	simappConfig.Telemetry.PrometheusRetentionTime = 60
	simappConfig.Telemetry.EnableHostnameLabel = false
	simappConfig.Telemetry.GlobalLabels = [][]string{{"chain_id", chainID}}

	var (
		genAccounts []authtypes.GenesisAccount
		genBalances []banktypes.Balance
		genFiles    []string
	)

	inBuf := bufio.NewReader(cmd.InOrStdin())
	// generate private keys, node IDs, and initial transactions
	for i := 0; i < numValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", nodeDirPrefix, i)
		nodeDir := filepath.Join(outputDir, nodeDirName, nodeDaemonHome)
		clientDir := filepath.Join(outputDir, nodeDirName, nodeCLIHome)
		gentxsDir := filepath.Join(outputDir, "gentxs")

		nodeConfig.SetRoot(nodeDir)
		nodeConfig.RPC.ListenAddress = "tcp://0.0.0.0:26657"

		if err := os.MkdirAll(filepath.Join(nodeDir, "config"), nodeDirPerm); err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		if err := os.MkdirAll(clientDir, nodeDirPerm); err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		nodeConfig.Moniker = nodeDirName

		ip, err := getIP(i, startingIPAddress)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		nodeIDs[i], valPubKeys[i], err = genutil.InitializeNodeValidatorFiles(nodeConfig)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		memo := fmt.Sprintf("%s@%s:26656", nodeIDs[i], ip)
		genFiles = append(genFiles, nodeConfig.GenesisFile())

		kb, err := keyring.New(
			sdk.KeyringServiceName(),
			keyringBackend,
			nodeDir,
			inBuf,
			clientCtx.Codec,
		)
		if err != nil {
			return err
		}

		keyringAlgos, _ := kb.SupportedAlgorithms()
		algo, err := keyring.NewSigningAlgoFromString(algoStr, keyringAlgos)
		if err != nil {
			return err
		}

		addr, secret, err := testutil.GenerateSaveCoinKey(kb, nodeDirName, "", true, algo)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		info := map[string]string{"secret": secret}

		cliPrint, err := json.Marshal(info)
		if err != nil {
			return err
		}

		// save private key seed words
		if err := writeFile(fmt.Sprintf("%v.json", "key_seed"), clientDir, cliPrint); err != nil {
			return err
		}

		accTokens := sdk.TokensFromConsensusPower(1000, sdk.DefaultPowerReduction)
		accStakingTokens := sdk.TokensFromConsensusPower(500, sdk.DefaultPowerReduction)
		accEvmTokens := sdk.TokensFromConsensusPower(5000, PowerReduction)
		accIrisTokens := sdk.TokensFromConsensusPower(5000, sdk.DefaultPowerReduction)

		coins := sdk.Coins{
			sdk.NewCoin(fmt.Sprintf("%stoken", nodeDirName), accTokens),
			sdk.NewCoin(sdk.DefaultBondDenom, accStakingTokens),
			sdk.NewCoin(iristypes.EvmToken.MinUnit, accEvmTokens),
			sdk.NewCoin(nativeIrisMinUnit, accIrisTokens),
		}

		genBalances = append(
			genBalances,
			banktypes.Balance{Address: addr.String(), Coins: coins.Sort()},
		)
		genAccounts = append(genAccounts, authtypes.NewBaseAccount(addr, nil, 0, 0))

		valTokens := sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)
		createValMsg, err := stakingtypes.NewMsgCreateValidator(
			sdk.ValAddress(addr),
			valPubKeys[i],
			sdk.NewCoin(sdk.DefaultBondDenom, valTokens),
			stakingtypes.NewDescription(nodeDirName, "", "", "", ""),
			stakingtypes.NewCommissionRates(sdk.OneDec(), sdk.OneDec(), sdk.OneDec()),
			sdk.OneInt(),
		)
		if err != nil {
			return err
		}

		txBuilder := clientCtx.TxConfig.NewTxBuilder()
		if err := txBuilder.SetMsgs(createValMsg); err != nil {
			return err
		}

		txBuilder.SetMemo(memo)

		txFactory := tx.Factory{}
		txFactory = txFactory.
			WithChainID(chainID).
			WithMemo(memo).
			WithKeybase(kb).
			WithTxConfig(clientCtx.TxConfig)

		if err := tx.Sign(txFactory, nodeDirName, txBuilder, true); err != nil {
			return err
		}

		txBz, err := clientCtx.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
		if err != nil {
			return err
		}

		if err := writeFile(fmt.Sprintf("%v.json", nodeDirName), gentxsDir, txBz); err != nil {
			return err
		}

		customAppTemplate, _ := servercfg.AppConfig(iristypes.NativeToken.MinUnit)
		srvconfig.SetConfigTemplate(customAppTemplate)

		srvconfig.WriteConfigFile(filepath.Join(nodeDir, "config/app.toml"), simappConfig)
	}

	if err := initGenFiles(clientCtx, mbm, chainID, genAccounts, genBalances, genFiles, numValidators); err != nil {
		return err
	}

	err := collectGenFiles(
		clientCtx, nodeConfig, chainID, nodeIDs, valPubKeys, numValidators,
		outputDir, nodeDirPrefix, nodeDaemonHome, genBalIterator,
	)
	if err != nil {
		return err
	}

	cmd.PrintErrf("Successfully initialized %d node directories\n", numValidators)
	return nil
}

func initGenFiles(
	clientCtx client.Context, mbm module.BasicManager, chainID string,
	genAccounts []authtypes.GenesisAccount, genBalances []banktypes.Balance,
	genFiles []string, numValidators int,
) error {
	appGenState := mbm.DefaultGenesis(clientCtx.Codec)

	// add the profiler and trustees in the genesis state
	var guardianGenState guardiantypes.GenesisState
	clientCtx.Codec.MustUnmarshalJSON(appGenState[guardiantypes.ModuleName], &guardianGenState)

	for _, account := range genAccounts {
		guardian := guardiantypes.NewSuper(
			"genesis", guardiantypes.Genesis,
			account.GetAddress(), account.GetAddress(),
		)
		guardianGenState.Supers = append(guardianGenState.Supers, guardian)
	}
	appGenState[guardiantypes.ModuleName] = clientCtx.Codec.MustMarshalJSON(&guardianGenState)

	// set the accounts in the genesis state
	var authGenState authtypes.GenesisState
	clientCtx.Codec.MustUnmarshalJSON(appGenState[authtypes.ModuleName], &authGenState)

	accounts, err := authtypes.PackAccounts(genAccounts)
	if err != nil {
		return err
	}

	authGenState.Accounts = accounts
	appGenState[authtypes.ModuleName] = clientCtx.Codec.MustMarshalJSON(&authGenState)

	// set the balances in the genesis state
	var bankGenState banktypes.GenesisState
	clientCtx.Codec.MustUnmarshalJSON(appGenState[banktypes.ModuleName], &bankGenState)
	bankGenState.Balances = genBalances
	appGenState[banktypes.ModuleName] = clientCtx.Codec.MustMarshalJSON(&bankGenState)

	// set the point token in the genesis state
	var tokenGenState tokentypesv1.GenesisState
	clientCtx.Codec.MustUnmarshalJSON(appGenState[tokentypes.ModuleName], &tokenGenState)
	tokenGenState.Tokens = append(tokenGenState.Tokens, iristypes.EvmToken)
	appGenState[tokentypes.ModuleName] = clientCtx.Codec.MustMarshalJSON(&tokenGenState)

	//set system service in the genesis state
	var serviceGenState servicetypes.GenesisState
	clientCtx.Codec.MustUnmarshalJSON(appGenState[servicetypes.ModuleName], &serviceGenState)
	serviceGenState.Definitions = append(
		serviceGenState.Definitions,
		servicetypes.GenOraclePriceSvcDefinition(),
	)
	serviceGenState.Bindings = append(
		serviceGenState.Bindings,
		servicetypes.GenOraclePriceSvcBinding(iristypes.NativeToken.MinUnit),
	)
	serviceGenState.Definitions = append(
		serviceGenState.Definitions,
		randomtypes.GetSvcDefinition(),
	)
	appGenState[servicetypes.ModuleName] = clientCtx.Codec.MustMarshalJSON(&serviceGenState)

	// set the evm fee token denom in the genesis state
	var evmGenState evmtypes.GenesisState
	clientCtx.Codec.MustUnmarshalJSON(appGenState[evmtypes.ModuleName], &evmGenState)
	evmGenState.Params.EvmDenom = iristypes.EvmToken.MinUnit
	appGenState[evmtypes.ModuleName] = clientCtx.Codec.MustMarshalJSON(&evmGenState)

	appGenStateJSON, err := json.MarshalIndent(appGenState, "", "  ")
	if err != nil {
		return err
	}

	genDoc := types.GenesisDoc{
		ChainID:    chainID,
		AppState:   appGenStateJSON,
		Validators: nil,
	}

	// generate empty genesis files for each validator and save
	for i := 0; i < numValidators; i++ {
		if err := genDoc.SaveAs(genFiles[i]); err != nil {
			return err
		}
	}
	return nil
}

func collectGenFiles(
	clientCtx client.Context,
	nodeConfig *tmconfig.Config,
	chainID string,
	nodeIDs []string,
	valPubKeys []cryptotypes.PubKey,
	numValidators int,
	outputDir, nodeDirPrefix, nodeDaemonHome string,
	genBalIterator banktypes.GenesisBalancesIterator,
) error {

	var appState json.RawMessage
	genTime := tmtime.Now()

	for i := 0; i < numValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", nodeDirPrefix, i)
		nodeDir := filepath.Join(outputDir, nodeDirName, nodeDaemonHome)
		gentxsDir := filepath.Join(outputDir, "gentxs")
		nodeConfig.Moniker = nodeDirName

		nodeConfig.SetRoot(nodeDir)

		nodeID, valPubKey := nodeIDs[i], valPubKeys[i]
		initCfg := genutiltypes.NewInitConfig(chainID, gentxsDir, nodeID, valPubKey)

		genDoc, err := types.GenesisDocFromFile(nodeConfig.GenesisFile())
		if err != nil {
			return err
		}

		nodeAppState, err := genutil.GenAppStateFromConfig(
			clientCtx.Codec,
			clientCtx.TxConfig,
			nodeConfig,
			initCfg,
			*genDoc,
			genBalIterator,
			genutiltypes.DefaultMessageValidator,
		)
		if err != nil {
			return err
		}

		if appState == nil {
			// set the canonical application state (they should not differ)
			appState = nodeAppState
		}
		genDoc.ConsensusParams.Block.MaxGas = 20000000
		genDoc.ChainID = chainID
		genDoc.GenesisTime = genTime
		genDoc.AppState = appState

		genFile := nodeConfig.GenesisFile()
		// overwrite each validator's genesis file to have a canonical genesis time
		if err := genutil.ExportGenesisFile(genDoc, genFile); err != nil {
			return err
		}
	}

	return nil
}

func getIP(i int, startingIPAddr string) (ip string, err error) {
	if len(startingIPAddr) == 0 {
		ip, err = server.ExternalIP()
		if err != nil {
			return "", err
		}
		return ip, nil
	}
	return calculateIP(startingIPAddr, i)
}

func calculateIP(ip string, i int) (string, error) {
	ipv4 := net.ParseIP(ip).To4()
	if ipv4 == nil {
		return "", fmt.Errorf("%v: non ipv4 address", ip)
	}

	for j := 0; j < i; j++ {
		ipv4[3]++
	}

	return ipv4.String(), nil
}

func writeFile(name string, dir string, contents []byte) error {
	writePath := filepath.Join(dir)
	file := filepath.Join(writePath, name)

	if err := tmos.EnsureDir(writePath, 0755); err != nil {
		return err
	}

	if err := ioutil.WriteFile(file, contents, 0644); err != nil {
		return err
	}

	return nil
}
