package simapp

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	tmtypes "github.com/cometbft/cometbft/types"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	pruningtypes "github.com/cosmos/cosmos-sdk/store/pruning/types"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	minttypes "github.com/irisnet/irishub/v2/modules/mint/types"
)

func setup(withGenesis bool, invCheckPeriod uint) (*SimApp, GenesisState) {
	db := dbm.NewMemDB()
	encCdc := MakeTestEncodingConfig()
	app := NewSimApp(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		DefaultNodeHome,
		invCheckPeriod,
		encCdc,
		EmptyAppOptions{},
	)
	if withGenesis {
		return app, NewDefaultGenesisState(encCdc.Codec)
	}
	return app, GenesisState{}
}

// Setup initializes a new SimApp. A Nop logger is set in SimApp.
func Setup(t *testing.T, _ bool) *SimApp {
	t.Helper()

	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)

	// create validator set with single validator
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	acc := authtypes.NewBaseAccount(
		senderPrivKey.PubKey().Address().Bytes(),
		senderPrivKey.PubKey(),
		0,
		0,
	)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000000000000))),
	}

	app := SetupWithGenesisValSet(t, valSet, []authtypes.GenesisAccount{acc}, balance)

	return app
}

func NewConfig() network.Config {
	cfg := network.DefaultConfig(NewTestNetworkFixture)
	encCfg := MakeTestEncodingConfig() // redundant
	cfg.Codec = encCfg.Codec
	cfg.TxConfig = encCfg.TxConfig
	cfg.LegacyAmino = encCfg.Amino
	cfg.InterfaceRegistry = encCfg.InterfaceRegistry
	cfg.AppConstructor = func(val network.ValidatorI) servertypes.Application {
		return NewSimApp(
			val.GetCtx().Logger,
			dbm.NewMemDB(),
			nil,
			true,
			nil,
			DefaultNodeHome,
			0,
			encCfg,
			EmptyAppOptions{},
			bam.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
			bam.SetChainID(cfg.ChainID),
		)
	}
	cfg.GenesisState = NewDefaultGenesisState(cfg.Codec)
	return cfg
}

func SimAppConstructor(val network.Validator) servertypes.Application {
	return NewSimApp(
		val.Ctx.Logger, dbm.NewMemDB(), nil, true, make(map[int64]bool),
		val.Ctx.Config.RootDir, 0, MakeTestEncodingConfig(), EmptyAppOptions{},
		bam.SetMinGasPrices(val.AppConfig.MinGasPrices),
	)
}

// SetupWithGenesisValSet initializes a new SimApp with a validator set and genesis accounts
// that also act as delegators. For simplicity, each validator is bonded with a delegation
// of one consensus engine unit (10^6) in the default token of the simapp from first genesis
// account. A Nop logger is set in SimApp.
func SetupWithGenesisValSet(
	t *testing.T,
	valSet *tmtypes.ValidatorSet,
	genAccs []authtypes.GenesisAccount,
	balances ...banktypes.Balance,
) *SimApp {
	t.Helper()

	app, genesisState := setup(true, 5)
	genesisState, err := GenesisStateWithValSet(
		app.AppCodec(),
		genesisState,
		valSet,
		genAccs,
		balances...)
	require.NoError(t, err)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)

	// init chain will set the validator set and initialize the genesis accounts
	app.InitChain(
		abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: simtestutil.DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

	// commit genesis changes
	app.Commit()
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{
		Height:             app.LastBlockHeight() + 1,
		AppHash:            app.LastCommitID().Hash,
		ValidatorsHash:     valSet.Hash(),
		NextValidatorsHash: valSet.Hash(),
	}})

	return app
}

// GenesisStateWithValSet returns a new genesis state with the validator set
func GenesisStateWithValSet(codec codec.Codec, genesisState map[string]json.RawMessage,
	valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount,
	balances ...banktypes.Balance,
) (map[string]json.RawMessage, error) {
	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = codec.MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.DefaultPowerReduction

	for _, val := range valSet.Validators {
		pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		if err != nil {
			return nil, fmt.Errorf("failed to convert pubkey: %w", err)
		}

		pkAny, err := codectypes.NewAnyWithValue(pk)
		if err != nil {
			return nil, fmt.Errorf("failed to create new any: %w", err)
		}

		validator := stakingtypes.Validator{
			OperatorAddress: sdk.ValAddress(val.Address).String(),
			ConsensusPubkey: pkAny,
			Jailed:          false,
			Status:          stakingtypes.Bonded,
			Tokens:          bondAmt,
			DelegatorShares: sdk.OneDec(),
			Description:     stakingtypes.Description{},
			UnbondingHeight: int64(0),
			UnbondingTime:   time.Unix(0, 0).UTC(),
			Commission: stakingtypes.NewCommission(
				sdk.ZeroDec(),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
			),
			MinSelfDelegation: sdk.ZeroInt(),
		}
		validators = append(validators, validator)
		delegations = append(
			delegations,
			stakingtypes.NewDelegation(genAccs[0].GetAddress(), val.Address.Bytes(), sdk.OneDec()),
		)

	}
	// set validators and delegations
	stakingGenesis := stakingtypes.NewGenesisState(
		stakingtypes.DefaultParams(),
		validators,
		delegations,
	)
	genesisState[stakingtypes.ModuleName] = codec.MustMarshalJSON(stakingGenesis)

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens to total supply
		totalSupply = totalSupply.Add(b.Coins...)
	}

	for range delegations {
		// add delegated tokens to total supply
		totalSupply = totalSupply.Add(sdk.NewCoin(sdk.DefaultBondDenom, bondAmt))
	}

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, bondAmt)},
	})

	// update total supply
	bankGenesis := banktypes.NewGenesisState(
		banktypes.DefaultGenesisState().Params,
		balances,
		totalSupply,
		[]banktypes.Metadata{},
		nil,
	)
	genesisState[banktypes.ModuleName] = codec.MustMarshalJSON(bankGenesis)

	return genesisState, nil
}

type GenerateAccountStrategy func(int) []sdk.AccAddress

// createRandomAccounts is a strategy used by addTestAddrs() in order to generated addresses in random order.
func createRandomAccounts(accNum int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, accNum)
	for i := 0; i < accNum; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	return testAddrs
}

// createIncrementalAccounts is a strategy used by addTestAddrs() in order to generated addresses in ascending order.
func createIncrementalAccounts(accNum int) []sdk.AccAddress {
	var addresses []sdk.AccAddress
	var buffer bytes.Buffer

	// start at 100 so we can make up to 999 test addresses with valid test addresses
	for i := 100; i < (accNum + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") //base address string

		buffer.WriteString(numString) //adding on final two digits to make addresses unique
		res, _ := sdk.AccAddressFromHexUnsafe(buffer.String())
		bech := res.String()
		addr, _ := TestAddr(buffer.String(), bech)

		addresses = append(addresses, addr)
		buffer.Reset()
	}

	return addresses
}

// AddTestAddrsFromPubKeys adds the addresses into the SimApp providing only the public keys.
func AddTestAddrsFromPubKeys(
	app *SimApp,
	ctx sdk.Context,
	pubKeys []cryptotypes.PubKey,
	accAmt sdk.Int,
) {
	initCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), accAmt))
	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, pubKey := range pubKeys {
		saveAccount(app, ctx, sdk.AccAddress(pubKey.Address()), initCoins)
	}
}

// AddTestAddrs constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrs(app *SimApp, ctx sdk.Context, accNum int, accAmt sdk.Int) []sdk.AccAddress {
	return addTestAddrs(app, ctx, accNum, accAmt, createRandomAccounts)
}

// AddTestAddrsIncremental constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrsIncremental(
	app *SimApp,
	ctx sdk.Context,
	accNum int,
	accAmt sdk.Int,
) []sdk.AccAddress {
	return addTestAddrs(app, ctx, accNum, accAmt, createIncrementalAccounts)
}

func addTestAddrs(
	app *SimApp,
	ctx sdk.Context,
	accNum int,
	accAmt sdk.Int,
	strategy GenerateAccountStrategy,
) []sdk.AccAddress {
	testAddrs := strategy(accNum)

	initCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), accAmt))

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, addr := range testAddrs {
		saveAccount(app, ctx, addr, initCoins)
	}

	return testAddrs
}

// saveAccount saves the provided account into the simapp with balance based on initCoins.
func saveAccount(app *SimApp, ctx sdk.Context, addr sdk.AccAddress, initCoins sdk.Coins) {
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	app.AccountKeeper.SetAccount(ctx, acc)
	if err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins); err != nil {
		panic(err)
	}
	if err := app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, initCoins); err != nil {
		panic(err)
	}
}

// ConvertAddrsToValAddrs converts the provided addresses to ValAddress.
func ConvertAddrsToValAddrs(addrs []sdk.AccAddress) []sdk.ValAddress {
	valAddrs := make([]sdk.ValAddress, len(addrs))

	for i, addr := range addrs {
		valAddrs[i] = sdk.ValAddress(addr)
	}

	return valAddrs
}

func TestAddr(addr string, bech string) (sdk.AccAddress, error) {
	res, err := sdk.AccAddressFromHexUnsafe(addr)
	if err != nil {
		return nil, err
	}
	bechexpected := res.String()
	if bech != bechexpected {
		return nil, fmt.Errorf("bech encoding doesn't match reference")
	}

	bechres, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(bechres, res) {
		return nil, err
	}

	return res, nil
}

// CreateTestPubKeys returns a total of numPubKeys public keys in ascending order.
func CreateTestPubKeys(numPubKeys int) []cryptotypes.PubKey {
	var publicKeys []cryptotypes.PubKey
	var buffer bytes.Buffer

	// start at 10 to avoid changing 1 to 01, 2 to 02, etc
	for i := 100; i < (numPubKeys + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString(
			"0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF",
		) // base pubkey string
		buffer.WriteString(
			numString,
		) // adding on final two digits to make pubkeys unique
		publicKeys = append(publicKeys, NewPubKeyFromHex(buffer.String()))
		buffer.Reset()
	}

	return publicKeys
}

// NewPubKeyFromHex returns a PubKey from a hex string.
func NewPubKeyFromHex(pk string) (res cryptotypes.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	if len(pkBytes) != ed25519.PubKeySize {
		panic(errors.Wrap(errors.ErrInvalidPubKey, "invalid pubkey size"))
	}
	return &ed25519.PubKey{Key: pkBytes}
}

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(o string) interface{} {
	return nil
}

// NewTestNetworkFixture returns a new simapp AppConstructor for network simulation tests
func NewTestNetworkFixture() network.TestFixture {
	dir, err := os.MkdirTemp("", "simapp")
	if err != nil {
		panic(fmt.Sprintf("failed creating temporary directory: %v", err))
	}
	defer os.RemoveAll(dir)

	encodingConfig := MakeTestEncodingConfig()

	app := NewSimApp(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		nil,
		DefaultNodeHome,
		0,
		encodingConfig, // redundant
		simtestutil.NewAppOptionsWithFlagHome(dir),
	)

	appCtr := func(val network.ValidatorI) servertypes.Application {
		return NewSimApp(
			val.GetCtx().Logger,
			dbm.NewMemDB(),
			nil,
			true,
			nil,
			DefaultNodeHome,
			0,
			encodingConfig, // redundant
			simtestutil.NewAppOptionsWithFlagHome(val.GetCtx().Config.RootDir),
			bam.SetPruning(pruningtypes.NewPruningOptionsFromString(val.GetAppConfig().Pruning)),
			bam.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
			bam.SetChainID(val.GetCtx().Viper.GetString(flags.FlagChainID)),
		)
	}

	return network.TestFixture{
		AppConstructor: appCtr,
		GenesisState:   app.DefaultGenesis(),
		EncodingConfig: testutil.TestEncodingConfig{
			InterfaceRegistry: app.InterfaceRegistry(),
			Codec:             app.AppCodec(),
			TxConfig:          app.TxConfig(),
			Amino:             app.LegacyAmino(),
		},
	}
}

func ExecTxCmdWithResult(
	t *testing.T,
	network *network.Network,
	clientCtx client.Context,
	cmd *cobra.Command,
	extraArgs []string,
) *coretypes.ResultTx {
	buf, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, extraArgs)
	require.NoError(t, err, "ExecTestCLICmd failed")
	require.NoError(t, network.WaitForNextBlock(), "network.WaitForNextBlock failed")

	respType := proto.Message(&sdk.TxResponse{})
	require.NoError(t, clientCtx.Codec.UnmarshalJSON(buf.Bytes(), respType), buf.String())

	txResp := respType.(*sdk.TxResponse)
	require.Equal(t, uint32(0), txResp.Code)
	return WaitForTx(t, network, clientCtx, txResp.TxHash)
}

func WaitForTx(
	t *testing.T,
	network *network.Network,
	clientCtx client.Context,
	txHash string,
) *coretypes.ResultTx {
	var (
		result *coretypes.ResultTx
		err    error
		tryCnt = 3
	)

	txHashBz, err := hex.DecodeString(txHash)
	require.NoError(t, err, "hex.DecodeString failed")

reTry:
	result, err = clientCtx.Client.Tx(context.Background(), txHashBz, false)
	if err != nil && strings.Contains(err.Error(), "not found") && tryCnt > 0 {
		require.NoError(t, network.WaitForNextBlock(), "network.WaitForNextBlock failed")
		tryCnt--
		goto reTry
	}

	require.NoError(t, err, "query tx failed")
	return result
}
