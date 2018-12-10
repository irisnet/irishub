package lcd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/irisnet/irishub/codec"
	crkeys "github.com/irisnet/irishub/crypto/keys"
	"github.com/irisnet/irishub/server"
	"github.com/irisnet/irishub/tests"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/auth"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	"github.com/irisnet/irishub/modules/stake"
	irisapp "github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/app/v0"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/keys"
	"github.com/irisnet/irishub/modules/gov"
	"github.com/irisnet/irishub/modules/upgrade"
	abci "github.com/tendermint/tendermint/abci/types"
	tmcfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	nm "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	pvm "github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	tmrpc "github.com/tendermint/tendermint/rpc/lib/server"
	tmtypes "github.com/tendermint/tendermint/types"
)

// makePathname creates a unique pathname for each test. It will panic if it
// cannot get the current working directory.
func makePathname() string {
	p, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	sep := string(filepath.Separator)
	return strings.Replace(p, sep, "_", -1)
}

// GetConfig returns a Tendermint config for the test cases.
func GetConfig() *tmcfg.Config {
	pathname := makePathname()
	config := tmcfg.ResetTestRoot(pathname)

	tmAddr, _, err := server.FreeTCPAddr()
	if err != nil {
		panic(err)
	}

	rcpAddr, _, err := server.FreeTCPAddr()
	if err != nil {
		panic(err)
	}

	grpcAddr, _, err := server.FreeTCPAddr()
	if err != nil {
		panic(err)
	}

	config.P2P.ListenAddress = tmAddr
	config.RPC.ListenAddress = rcpAddr
	config.RPC.GRPCListenAddress = grpcAddr

	return config
}

// GetKeyBase returns the LCD test keybase. It also requires that a directory
// could be made and a keybase could be fetched.
//
// NOTE: memDB cannot be used because the request is expecting to interact with
// the default location.
func GetKeyBase(t *testing.T) crkeys.Keybase {
	dir, err := ioutil.TempDir("", "lcd_test")
	require.NoError(t, err)

	viper.Set(cli.HomeFlag, dir)

	keybase, err := keys.GetKeyBase()
	require.NoError(t, err)

	return keybase
}

// CreateAddr adds an address to the key store and returns an address and seed.
// It also requires that the key could be created.
func CreateAddr(t *testing.T, name, password string, kb crkeys.Keybase) (sdk.AccAddress, string) {
	var (
		err  error
		info crkeys.Info
		seed string
	)

	info, seed, err = kb.CreateMnemonic(name, crkeys.English, password, crkeys.Secp256k1)
	require.NoError(t, err)

	return sdk.AccAddress(info.GetPubKey().Address()), seed
}

func extractPortFromAddress(listenAddress string) string {
	stringList := strings.Split(listenAddress, ":")
	length := len(stringList)
	if length != 3 {
		panic(fmt.Errorf("expected listen address: tcp://0.0.0.0:12345, got %s", listenAddress))
	}
	return stringList[2]
}

// InitializeTestLCD starts Tendermint and the LCD in process, listening on
// their respective sockets where nValidators is the total number of validators
// and initAddrs are the accounts to initialize with some steak tokens. It
// returns a cleanup function, a set of validator public keys, and a port.
func InitializeTestLCD(
	t *testing.T, nValidators int, initAddrs []sdk.AccAddress,
) (cleanup func(), valConsPubKeys []crypto.PubKey, valOperAddrs []sdk.ValAddress, port string) {

	if nValidators < 1 {
		panic("InitializeTestLCD must use at least one validator")
	}

	config := GetConfig()
	config.Consensus.TimeoutCommit = 100
	config.Consensus.SkipTimeoutCommit = false
	config.TxIndex.IndexAllTags = true

	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	logger = log.NewFilter(logger, log.AllowError())

	privValidatorFile := config.PrivValidatorFile()
	privVal := pvm.LoadOrGenFilePV(privValidatorFile)
	privVal.Reset()

	db := dbm.NewMemDB()
	app := irisapp.NewIrisApp(logger, db, nil)
	cdc = irisapp.MakeCodec()

	genesisFile := config.GenesisFile()
	genDoc, err := tmtypes.GenesisDocFromFile(genesisFile)
	require.Nil(t, err)
	genDoc.Validators = nil
	genDoc.SaveAs(genesisFile)
	genTxs := []json.RawMessage{}

	// append any additional (non-proposing) validators
	for i := 0; i < nValidators; i++ {
		operPrivKey := secp256k1.GenPrivKey()
		operAddr := operPrivKey.PubKey().Address()
		pubKey := privVal.PubKey
		delegation := 20
		if i > 0 {
			pubKey = ed25519.GenPrivKey().PubKey()
			delegation = 18
		}
		msg := stake.NewMsgCreateValidator(
			sdk.ValAddress(operAddr),
			pubKey,
			sdk.NewCoin("iris-atto", sdk.NewIntWithDecimal(1, delegation)),
			stake.Description{Moniker: fmt.Sprintf("validator-%d", i+1)},
			stake.NewCommissionMsg(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
		)
		stdSignMsg := client.StdSignMsg{
			ChainID: genDoc.ChainID,
			Msgs:    []sdk.Msg{msg},
		}
		sig, err := operPrivKey.Sign(stdSignMsg.Bytes())
		require.Nil(t, err)
		tx := auth.NewStdTx([]sdk.Msg{msg}, auth.StdFee{}, []auth.StdSignature{{Signature: sig, PubKey: operPrivKey.PubKey()}}, "")
		txBytes, err := cdc.MarshalJSON(tx)
		require.Nil(t, err)
		genTxs = append(genTxs, txBytes)
		valConsPubKeys = append(valConsPubKeys, pubKey)
		valOperAddrs = append(valOperAddrs, sdk.ValAddress(operAddr))
	}

	genesisState, err := v0.IrisAppGenState(cdc, *genDoc, genTxs)
	require.NoError(t, err)

	// add some tokens to init accounts
	for _, addr := range initAddrs {
		accAuth := auth.NewBaseAccountWithAddress(addr)
		accAuth.Coins = sdk.Coins{sdk.NewCoin("iris-atto", sdk.NewIntWithDecimal(1, 20))}
		acc := v0.NewGenesisFileAccount(&accAuth)
		genesisState.Accounts = append(genesisState.Accounts, acc)
		genesisState.StakeData.Pool.LooseTokens = genesisState.StakeData.Pool.LooseTokens.Add(sdk.NewDecFromInt(sdk.NewIntWithDecimal(1, 20)))
	}

	genesisState.GovData = gov.DefaultGenesisStateForLCDTest()
	genesisState.UpgradeData = upgrade.DefaultGenesisStateForTest()

	appState, err := codec.MarshalJSONIndent(cdc, genesisState)
	require.NoError(t, err)
	genDoc.AppState = appState

	listenAddr, port, err := server.FreeTCPAddr()
	require.NoError(t, err)

	// XXX: Need to set this so LCD knows the tendermint node address!
	viper.Set(client.FlagNode, config.RPC.ListenAddress)
	viper.Set(client.FlagChainID, genDoc.ChainID)
	// TODO Set to false once the upstream Tendermint proof verification issue is fixed.
	viper.Set(client.FlagTrustNode, true)
	dir, err := ioutil.TempDir("", "lcd_test")
	require.NoError(t, err)
	viper.Set(cli.HomeFlag, dir)

	node, err := startTM(config, logger, genDoc, privVal, app)
	require.NoError(t, err)

	tests.WaitForNextHeightTM(tests.ExtractPortFromAddress(config.RPC.ListenAddress))
	lcd, err := startLCD(logger, listenAddr, cdc)
	require.NoError(t, err)

	tests.WaitForLCDStart(port)
	tests.WaitForHeight(1, port)

	cleanup = func() {
		logger.Debug("cleaning up LCD initialization")
		node.Stop()
		node.Wait()
		lcd.Close()
	}

	return cleanup, valConsPubKeys, valOperAddrs, port
}

// startTM creates and starts an in-process Tendermint node with memDB and
// in-process ABCI application. It returns the new node or any error that
// occurred.
//
// TODO: Clean up the WAL dir or enable it to be not persistent!
func startTM(
	tmcfg *tmcfg.Config, logger log.Logger, genDoc *tmtypes.GenesisDoc,
	privVal tmtypes.PrivValidator, app abci.Application,
) (*nm.Node, error) {
	genDocProvider := func() (*tmtypes.GenesisDoc, error) { return genDoc, nil }
	dbProvider := func(*nm.DBContext) (dbm.DB, error) { return dbm.NewMemDB(), nil }
	nodeKey, err := p2p.LoadOrGenNodeKey(tmcfg.NodeKeyFile())
	if err != nil {
		return nil, err
	}
	node, err := nm.NewNode(
		tmcfg,
		privVal,
		nodeKey,
		proxy.NewLocalClientCreator(app),
		genDocProvider,
		dbProvider,
		nm.DefaultMetricsProvider(tmcfg.Instrumentation),
		logger.With("module", "node"),
	)
	if err != nil {
		return nil, err
	}

	err = node.Start()
	if err != nil {
		return nil, err
	}

	tests.WaitForRPC(tmcfg.RPC.ListenAddress)
	logger.Info("Tendermint running!")

	return node, err
}

// startLCD starts the LCD.
//
// NOTE: This causes the thread to block.
func startLCD(logger log.Logger, listenAddr string, cdc *codec.Codec) (net.Listener, error) {
	listener, err := tmrpc.Listen(listenAddr, tmrpc.Config{})
	if err != nil {
		return nil, err
	}

	err = tmrpc.StartHTTPServer(listener, createHandler(cdc), logger)
	if err != nil {
		return nil, err
	}
	return listener, nil
}

// Request makes a test LCD test request. It returns a response object and a
// stringified response body.
func Request(t *testing.T, port, method, path string, payload []byte) (*http.Response, string) {
	var (
		err error
		res *http.Response
	)
	url := fmt.Sprintf("http://localhost:%v%v", port, path)
	fmt.Println("REQUEST " + method + " " + url)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	require.Nil(t, err)

	res, err = http.DefaultClient.Do(req)
	require.Nil(t, err)

	output, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	require.Nil(t, err)

	return res, string(output)
}
