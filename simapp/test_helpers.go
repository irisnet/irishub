package simapp

//
//import (
//	"bytes"
//	"context"
//	"encoding/hex"
//	"encoding/json"
//	"fmt"
//	"math/rand"
//	"os"
//	"strconv"
//	"testing"
//	"time"
//
//	"cosmossdk.io/depinject"
//	errorsmod "cosmossdk.io/errors"
//	"cosmossdk.io/math"
//	dbm "github.com/cometbft/cometbft-db"
//	abci "github.com/cometbft/cometbft/abci/types"
//	"github.com/cometbft/cometbft/libs/log"
//	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
//	tmtypes "github.com/cometbft/cometbft/types"
//	bam "github.com/cosmos/cosmos-sdk/baseapp"
//	"github.com/cosmos/cosmos-sdk/client"
//	"github.com/cosmos/cosmos-sdk/client/flags"
//	"github.com/cosmos/cosmos-sdk/codec"
//	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
//	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
//	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
//	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
//	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
//	"github.com/cosmos/cosmos-sdk/runtime"
//	"github.com/cosmos/cosmos-sdk/server"
//	servertypes "github.com/cosmos/cosmos-sdk/server/types"
//	pruningtypes "github.com/cosmos/cosmos-sdk/store/pruning/types"
//	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
//	"github.com/cosmos/cosmos-sdk/testutil/mock"
//	"github.com/cosmos/cosmos-sdk/testutil/network"
//	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	"github.com/cosmos/cosmos-sdk/types/errors"
//	"github.com/cosmos/cosmos-sdk/types/module/testutil"
//	authcli "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
//	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
//	bankcli "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
//	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
//	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
//	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
//	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
//	"github.com/cosmos/gogoproto/proto"
//	"github.com/stretchr/testify/require"
//)
//
//// SetupOptions defines arguments that are passed into `Simapp` constructor.
//type SetupOptions struct {
//	Logger  log.Logger
//	DB      *dbm.MemDB
//	AppOpts servertypes.AppOptions
//}
//
//func setup(withGenesis bool, invCheckPeriod uint, depInjectOptions DepinjectOptions) (*SimApp, GenesisState) {
//	db := dbm.NewMemDB()
//
//	appOptions := make(simtestutil.AppOptionsMap, 0)
//	appOptions[flags.FlagHome] = DefaultNodeHome
//	appOptions[server.FlagInvCheckPeriod] = invCheckPeriod
//
//	app := NewSimApp(log.NewNopLogger(), db, nil, true, depInjectOptions, appOptions)
//	if withGenesis {
//		return app, app.DefaultGenesis()
//	}
//	return app, GenesisState{}
//}
//
//// Setup initializes a new SimApp. A Nop logger is set in SimApp.
//func Setup(t *testing.T, isCheckTx bool, depInjectOptions DepinjectOptions) *SimApp {
//	t.Helper()
//
//	privVal := mock.NewPV()
//	pubKey, err := privVal.GetPubKey()
//	require.NoError(t, err)
//
//	// create validator set with single validator
//	validator := tmtypes.NewValidator(pubKey, 1)
//	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})
//
//	// generate genesis account
//	senderPrivKey := secp256k1.GenPrivKey()
//	acc := authtypes.NewBaseAccount(
//		senderPrivKey.PubKey().Address().Bytes(),
//		senderPrivKey.PubKey(),
//		0,
//		0,
//	)
//	balance := banktypes.Balance{
//		Address: acc.GetAddress().String(),
//		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000000000000))),
//	}
//
//	app := SetupWithGenesisValSet(t, depInjectOptions, valSet, []authtypes.GenesisAccount{acc}, balance)
//
//	return app
//}
//
//func SetupWithGenesisStateFn(
//	t *testing.T,
//	depInjectOptions DepinjectOptions,
//	merge func(cdc codec.Codec, state GenesisState) GenesisState,
//) *SimApp {
//	t.Helper()
//	app, genesisState := setup(true, 5, depInjectOptions)
//
//	privVal := mock.NewPV()
//	pubKey, err := privVal.GetPubKey()
//	require.NoError(t, err)
//
//	// create validator set with single validator
//	validator := tmtypes.NewValidator(pubKey, 1)
//	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})
//
//	// generate genesis account
//	senderPrivKey := secp256k1.GenPrivKey()
//	acc := authtypes.NewBaseAccount(
//		senderPrivKey.PubKey().Address().Bytes(),
//		senderPrivKey.PubKey(),
//		0,
//		0,
//	)
//	balance := banktypes.Balance{
//		Address: acc.GetAddress().String(),
//		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000000000000))),
//	}
//	genesisState = genesisStateWithValSet(
//		t,
//		app,
//		genesisState,
//		valSet,
//		[]authtypes.GenesisAccount{acc},
//		balance,
//	)
//
//	if merge != nil {
//		genesisState = merge(app.appCodec, genesisState)
//	}
//	// init chain must be called to stop deliverState from being nil
//	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
//	if err != nil {
//		panic(err)
//	}
//
//	// Initialize the chain
//	app.InitChain(
//		abci.RequestInitChain{
//			Validators:      []abci.ValidatorUpdate{},
//			ConsensusParams: simtestutil.DefaultConsensusParams,
//			AppStateBytes:   stateBytes,
//		},
//	)
//	// commit genesis changes
//	app.Commit()
//	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{
//		Height:             app.LastBlockHeight() + 1,
//		AppHash:            app.LastCommitID().Hash,
//		ValidatorsHash:     valSet.Hash(),
//		NextValidatorsHash: valSet.Hash(),
//	}})
//	return app
//}
//
//func NewConfig(depInjectOptions DepinjectOptions) (network.Config, error) {
//	var (
//		appBuilder        *runtime.AppBuilder
//		txConfig          client.TxConfig
//		legacyAmino       *codec.LegacyAmino
//		cdc               codec.Codec
//		interfaceRegistry codectypes.InterfaceRegistry
//	)
//
//	providers := append(depInjectOptions.Providers, log.NewNopLogger())
//	if err := depinject.Inject(
//		depinject.Configs(
//			depInjectOptions.Config,
//			depinject.Supply(
//				providers...,
//			),
//		),
//		&appBuilder,
//		&txConfig,
//		&cdc,
//		&legacyAmino,
//		&interfaceRegistry,
//	); err != nil {
//		return network.Config{}, err
//	}
//
//	cfg := network.DefaultConfig(func() network.TestFixture {
//		return NewTestNetworkFixture(depInjectOptions)
//	})
//	cfg.Codec = cdc
//	cfg.TxConfig = txConfig
//	cfg.LegacyAmino = legacyAmino
//	cfg.InterfaceRegistry = interfaceRegistry
//	cfg.GenesisState = appBuilder.DefaultGenesis()
//	cfg.AppConstructor = func(val network.ValidatorI) servertypes.Application {
//		return NewSimApp(
//			val.GetCtx().Logger,
//			dbm.NewMemDB(),
//			nil,
//			true,
//			depInjectOptions,
//			EmptyAppOptions{},
//			bam.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
//			bam.SetChainID(cfg.ChainID),
//		)
//	}
//	return cfg, nil
//}
//
//// func SimAppConstructor(val network.ValidatorI) servertypes.Application {
//// 	return NewSimApp(
//// 		val.GetCtx().Logger, dbm.NewMemDB(), nil, true, EmptyAppOptions{},
//// 		bam.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
//// 	)
//// }
//
//func genesisStateWithValSet(t *testing.T,
//	app *SimApp, genesisState GenesisState,
//	valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount,
//	balances ...banktypes.Balance,
//) GenesisState {
//	t.Helper()
//	// set genesis accounts
//	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
//	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)
//
//	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
//	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))
//
//	bondAmt := sdk.DefaultPowerReduction
//
//	for _, val := range valSet.Validators {
//		pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
//		require.NoError(t, err)
//		pkAny, err := codectypes.NewAnyWithValue(pk)
//		require.NoError(t, err)
//		validator := stakingtypes.Validator{
//			OperatorAddress: sdk.ValAddress(val.Address).String(),
//			ConsensusPubkey: pkAny,
//			Jailed:          false,
//			Status:          stakingtypes.Bonded,
//			Tokens:          bondAmt,
//			DelegatorShares: sdk.OneDec(),
//			Description:     stakingtypes.Description{},
//			UnbondingHeight: int64(0),
//			UnbondingTime:   time.Unix(0, 0).UTC(),
//			Commission: stakingtypes.NewCommission(
//				sdk.ZeroDec(),
//				sdk.ZeroDec(),
//				sdk.ZeroDec(),
//			),
//			MinSelfDelegation: sdk.ZeroInt(),
//		}
//		validators = append(validators, validator)
//		delegations = append(
//			delegations,
//			stakingtypes.NewDelegation(genAccs[0].GetAddress(), val.Address.Bytes(), sdk.OneDec()),
//		)
//
//	}
//	// set validators and delegations
//	stakingGenesis := stakingtypes.NewGenesisState(
//		stakingtypes.DefaultParams(),
//		validators,
//		delegations,
//	)
//	genesisState[stakingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(stakingGenesis)
//
//	totalSupply := sdk.NewCoins()
//	for _, b := range balances {
//		// add genesis acc tokens to total supply
//		totalSupply = totalSupply.Add(b.Coins...)
//	}
//
//	for range delegations {
//		// add delegated tokens to total supply
//		totalSupply = totalSupply.Add(sdk.NewCoin(sdk.DefaultBondDenom, bondAmt))
//	}
//
//	// add bonded amount to bonded pool module account
//	balances = append(balances, banktypes.Balance{
//		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
//		Coins:   sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, bondAmt)},
//	})
//
//	// update total supply
//	bankGenesis := banktypes.NewGenesisState(
//		banktypes.DefaultGenesisState().Params,
//		balances,
//		totalSupply,
//		[]banktypes.Metadata{},
//		[]banktypes.SendEnabled{},
//	)
//	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)
//
//	return genesisState
//}
//
//// SetupWithGenesisValSet initializes a new SimApp with a validator set and genesis accounts
//// that also act as delegators. For simplicity, each validator is bonded with a delegation
//// of one consensus engine unit (10^6) in the default token of the simapp from first genesis
//// account. A Nop logger is set in SimApp.
//func SetupWithGenesisValSet(
//	t *testing.T,
//	depInjectOptions DepinjectOptions,
//	valSet *tmtypes.ValidatorSet,
//	genAccs []authtypes.GenesisAccount,
//	balances ...banktypes.Balance,
//) *SimApp {
//	t.Helper()
//
//	app, genesisState := setup(true, 5, depInjectOptions)
//	genesisState = genesisStateWithValSet(t, app, genesisState, valSet, genAccs, balances...)
//
//	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
//	require.NoError(t, err)
//
//	// init chain will set the validator set and initialize the genesis accounts
//	app.InitChain(
//		abci.RequestInitChain{
//			Validators:      []abci.ValidatorUpdate{},
//			ConsensusParams: simtestutil.DefaultConsensusParams,
//			AppStateBytes:   stateBytes,
//		},
//	)
//
//	// commit genesis changes
//	app.Commit()
//	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{
//		Height:             app.LastBlockHeight() + 1,
//		AppHash:            app.LastCommitID().Hash,
//		ValidatorsHash:     valSet.Hash(),
//		NextValidatorsHash: valSet.Hash(),
//	}})
//
//	return app
//}
//
//// SetupWithGenesisAccounts initializes a new SimApp with the provided genesis
//// accounts and possible balances.
//func SetupWithGenesisAccounts(
//	t *testing.T,
//	depInjectOptions DepinjectOptions,
//	genAccs []authtypes.GenesisAccount,
//	balances ...banktypes.Balance,
//) *SimApp {
//	t.Helper()
//
//	privVal := mock.NewPV()
//	pubKey, err := privVal.GetPubKey()
//	require.NoError(t, err)
//
//	// create validator set with single validator
//	validator := tmtypes.NewValidator(pubKey, 1)
//	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})
//
//	return SetupWithGenesisValSet(t, depInjectOptions, valSet, genAccs, balances...)
//}
//
//type GenerateAccountStrategy func(int) []sdk.AccAddress
//
//// createRandomAccounts is a strategy used by addTestAddrs() in order to generated addresses in random order.
//func createRandomAccounts(accNum int) []sdk.AccAddress {
//	testAddrs := make([]sdk.AccAddress, accNum)
//	for i := 0; i < accNum; i++ {
//		pk := ed25519.GenPrivKey().PubKey()
//		testAddrs[i] = sdk.AccAddress(pk.Address())
//	}
//
//	return testAddrs
//}
//
//// CreateTestAddrs creates test addresses
//func CreateTestAddrs(numAddrs int) []sdk.AccAddress {
//	var addresses []sdk.AccAddress
//	var buffer bytes.Buffer
//
//	// start at 100 so we can make up to 999 test addresses with valid test addresses
//	for i := 100; i < (numAddrs + 100); i++ {
//		numString := strconv.Itoa(i)
//		buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") // base address string
//
//		buffer.WriteString(numString) // adding on final two digits to make addresses unique
//		res, _ := sdk.AccAddressFromHexUnsafe(buffer.String())
//		bech := res.String()
//		addresses = append(addresses, testAddr(buffer.String(), bech))
//		buffer.Reset()
//	}
//
//	return addresses
//}
//
//// for incode address generation
//func testAddr(addr, bech string) sdk.AccAddress {
//	res, err := sdk.AccAddressFromHexUnsafe(addr)
//	if err != nil {
//		panic(err)
//	}
//	bechexpected := res.String()
//	if bech != bechexpected {
//		panic("Bech encoding doesn't match reference")
//	}
//
//	bechres, err := sdk.AccAddressFromBech32(bech)
//	if err != nil {
//		panic(err)
//	}
//	if !bytes.Equal(bechres, res) {
//		panic("Bech decode and hex decode don't match")
//	}
//
//	return res
//}
//
//// createIncrementalAccounts is a strategy used by addTestAddrs() in order to generated addresses in ascending order.
//func createIncrementalAccounts(accNum int) []sdk.AccAddress {
//	var addresses []sdk.AccAddress
//	var buffer bytes.Buffer
//
//	// start at 100 so we can make up to 999 test addresses with valid test addresses
//	for i := 100; i < (accNum + 100); i++ {
//		numString := strconv.Itoa(i)
//		buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") // base address string
//
//		buffer.WriteString(numString) // adding on final two digits to make addresses unique
//		res, _ := sdk.AccAddressFromHexUnsafe(buffer.String())
//		bech := res.String()
//		addr, _ := TestAddr(buffer.String(), bech)
//
//		addresses = append(addresses, addr)
//		buffer.Reset()
//	}
//
//	return addresses
//}
//
//// AddTestAddrsFromPubKeys adds the addresses into the SimApp providing only the public keys.
//func AddTestAddrsFromPubKeys(
//	app *SimApp,
//	ctx sdk.Context,
//	pubKeys []cryptotypes.PubKey,
//	accAmt math.Int,
//) {
//	initCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), accAmt))
//
//	for _, pk := range pubKeys {
//		initAccountWithCoins(app, ctx, sdk.AccAddress(pk.Address()), initCoins)
//	}
//}
//
//// AddTestAddrs constructs and returns accNum amount of accounts with an
//// initial balance of accAmt in random order
//func AddTestAddrs(app *SimApp, ctx sdk.Context, accNum int, accAmt math.Int) []sdk.AccAddress {
//	return addTestAddrs(app, ctx, accNum, accAmt, createRandomAccounts)
//}
//
//// AddTestAddrsIncremental constructs and returns accNum amount of accounts with an
//// initial balance of accAmt in random order
//func AddTestAddrsIncremental(
//	app *SimApp,
//	ctx sdk.Context,
//	accNum int,
//	accAmt math.Int,
//) []sdk.AccAddress {
//	return addTestAddrs(app, ctx, accNum, accAmt, createIncrementalAccounts)
//}
//
//func addTestAddrs(
//	app *SimApp,
//	ctx sdk.Context,
//	accNum int,
//	accAmt math.Int,
//	strategy GenerateAccountStrategy,
//) []sdk.AccAddress {
//	testAddrs := strategy(accNum)
//
//	initCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), accAmt))
//
//	for _, addr := range testAddrs {
//		initAccountWithCoins(app, ctx, addr, initCoins)
//	}
//
//	return testAddrs
//}
//
//func initAccountWithCoins(app *SimApp, ctx sdk.Context, addr sdk.AccAddress, coins sdk.Coins) {
//	err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
//	if err != nil {
//		panic(err)
//	}
//
//	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, coins)
//	if err != nil {
//		panic(err)
//	}
//}
//
//// ConvertAddrsToValAddrs converts the provided addresses to ValAddress.
//func ConvertAddrsToValAddrs(addrs []sdk.AccAddress) []sdk.ValAddress {
//	valAddrs := make([]sdk.ValAddress, len(addrs))
//
//	for i, addr := range addrs {
//		valAddrs[i] = sdk.ValAddress(addr)
//	}
//
//	return valAddrs
//}
//
//func TestAddr(addr, bech string) (sdk.AccAddress, error) {
//	res, err := sdk.AccAddressFromHexUnsafe(addr)
//	if err != nil {
//		return nil, err
//	}
//	bechexpected := res.String()
//	if bech != bechexpected {
//		return nil, fmt.Errorf("bech encoding doesn't match reference")
//	}
//
//	bechres, err := sdk.AccAddressFromBech32(bech)
//	if err != nil {
//		return nil, err
//	}
//	if !bytes.Equal(bechres, res) {
//		return nil, err
//	}
//
//	return res, nil
//}
//
//// CheckBalance checks the balance of an account.
//func CheckBalance(t *testing.T, app *SimApp, addr sdk.AccAddress, balances sdk.Coins) {
//	t.Helper()
//	ctxCheck := app.BaseApp.NewContext(true, tmproto.Header{})
//	require.True(t, balances.IsEqual(app.BankKeeper.GetAllBalances(ctxCheck, addr)))
//}
//
//// SignCheckDeliver checks a generated signed transaction and simulates a
//// block commitment with the given transaction. A test assertion is made using
//// the parameter 'expPass' against the result. A corresponding result is
//// returned.
//func SignCheckDeliver(
//	t *testing.T,
//	txCfg client.TxConfig,
//	app *bam.BaseApp,
//	header tmproto.Header,
//	msgs []sdk.Msg,
//	chainID string,
//	accNums, accSeqs []uint64,
//	expSimPass, expPass bool,
//	priv ...cryptotypes.PrivKey,
//) (sdk.GasInfo, *sdk.Result, error) {
//	t.Helper()
//	tx, err := simtestutil.GenSignedMockTx(
//		rand.New(rand.NewSource(time.Now().UnixNano())),
//		txCfg,
//		msgs,
//		sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)},
//		simtestutil.DefaultGenTxGas,
//		chainID,
//		accNums,
//		accSeqs,
//		priv...,
//	)
//	require.NoError(t, err)
//	txBytes, err := txCfg.TxEncoder()(tx)
//	require.Nil(t, err)
//
//	// Must simulate now as CheckTx doesn't run Msgs anymore
//	_, res, err := app.Simulate(txBytes)
//
//	if expSimPass {
//		require.NoError(t, err)
//		require.NotNil(t, res)
//	} else {
//		require.Error(t, err)
//		require.Nil(t, res)
//	}
//
//	// Simulate a sending a transaction and committing a block
//	app.BeginBlock(abci.RequestBeginBlock{Header: header})
//	gInfo, res, err := app.SimDeliver(txCfg.TxEncoder(), tx)
//
//	if expPass {
//		require.NoError(t, err)
//		require.NotNil(t, res)
//	} else {
//		require.Error(t, err)
//		require.Nil(t, res)
//	}
//
//	app.EndBlock(abci.RequestEndBlock{})
//	app.Commit()
//
//	return gInfo, res, err
//}
//
//// GenSequenceOfTxs generates a set of signed transactions of messages, such
//// that they differ only by having the sequence numbers incremented between
//// every transaction.
//func GenSequenceOfTxs(
//	txGen client.TxConfig,
//	msgs []sdk.Msg,
//	accNums []uint64,
//	initSeqNums []uint64,
//	numToGenerate int,
//	priv ...cryptotypes.PrivKey,
//) ([]sdk.Tx, error) {
//	txs := make([]sdk.Tx, numToGenerate)
//	var err error
//	for i := 0; i < numToGenerate; i++ {
//		txs[i], err = simtestutil.GenSignedMockTx(
//			rand.New(rand.NewSource(time.Now().UnixNano())),
//			txGen,
//			msgs,
//			sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)},
//			simtestutil.DefaultGenTxGas,
//			"",
//			accNums,
//			initSeqNums,
//			priv...,
//		)
//		if err != nil {
//			break
//		}
//		incrementAllSequenceNumbers(initSeqNums)
//	}
//
//	return txs, err
//}
//
//func incrementAllSequenceNumbers(initSeqNums []uint64) {
//	for i := 0; i < len(initSeqNums); i++ {
//		initSeqNums[i]++
//	}
//}
//
//// CreateTestPubKeys returns a total of numPubKeys public keys in ascending order.
//func CreateTestPubKeys(numPubKeys int) []cryptotypes.PubKey {
//	var publicKeys []cryptotypes.PubKey
//	var buffer bytes.Buffer
//
//	// start at 10 to avoid changing 1 to 01, 2 to 02, etc
//	for i := 100; i < (numPubKeys + 100); i++ {
//		numString := strconv.Itoa(i)
//		buffer.WriteString(
//			"0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF",
//		) // base pubkey string
//		buffer.WriteString(
//			numString,
//		) // adding on final two digits to make pubkeys unique
//		publicKeys = append(publicKeys, NewPubKeyFromHex(buffer.String()))
//		buffer.Reset()
//	}
//
//	return publicKeys
//}
//
//// NewPubKeyFromHex returns a PubKey from a hex string.
//func NewPubKeyFromHex(pk string) (res cryptotypes.PubKey) {
//	pkBytes, err := hex.DecodeString(pk)
//	if err != nil {
//		panic(err)
//	}
//	if len(pkBytes) != ed25519.PubKeySize {
//		panic(errorsmod.Wrap(errors.ErrInvalidPubKey, "invalid pubkey size"))
//	}
//	return &ed25519.PubKey{Key: pkBytes}
//}
//
//// EmptyAppOptions is a stub implementing AppOptions
//type EmptyAppOptions struct{}
//
//// Get implements AppOptions
//func (ao EmptyAppOptions) Get(o string) interface{} {
//	return nil
//}
//
//// FundAccount is a utility function that funds an account by minting and
//// sending the coins to the address. This should be used for testing purposes
//// only!
////
//// TODO: Instead of using the mint module account, which has the
//// permission of minting, create a "faucet" account. (@fdymylja)
//func FundAccount(
//	bankKeeper bankkeeper.Keeper,
//	ctx sdk.Context,
//	addr sdk.AccAddress,
//	amounts sdk.Coins,
//) error {
//	if err := bankKeeper.MintCoins(ctx, minttypes.ModuleName, amounts); err != nil {
//		return err
//	}
//
//	return bankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, amounts)
//}
//
//// FundModuleAccount is a utility function that funds a module account by
//// minting and sending the coins to the address. This should be used for testing
//// purposes only!
////
//// TODO: Instead of using the mint module account, which has the
//// permission of minting, create a "faucet" account. (@fdymylja)
//func FundModuleAccount(
//	bankKeeper bankkeeper.Keeper,
//	ctx sdk.Context,
//	recipientMod string,
//	amounts sdk.Coins,
//) error {
//	if err := bankKeeper.MintCoins(ctx, minttypes.ModuleName, amounts); err != nil {
//		return err
//	}
//
//	return bankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, recipientMod, amounts)
//}
//
//func QueryBalancesExec(
//	t *testing.T,
//	network Network,
//	clientCtx client.Context,
//	address string,
//	extraArgs ...string,
//) sdk.Coins {
//	t.Helper()
//	args := []string{
//		address,
//		fmt.Sprintf("--%s=json", "output"),
//	}
//	args = append(args, extraArgs...)
//
//	result := &banktypes.QueryAllBalancesResponse{}
//	network.ExecQueryCmd(t, clientCtx, bankcli.GetBalancesCmd(), args, result)
//	return result.Balances
//}
//
//func QueryBalanceExec(
//	t *testing.T,
//	network Network,
//	clientCtx client.Context,
//	address string,
//	denom string,
//	extraArgs ...string,
//) *sdk.Coin {
//	t.Helper()
//	args := []string{
//		address,
//		fmt.Sprintf("--%s=%s", bankcli.FlagDenom, denom),
//		fmt.Sprintf("--%s=json", "output"),
//	}
//	args = append(args, extraArgs...)
//
//	result := &sdk.Coin{}
//	network.ExecQueryCmd(t, clientCtx, bankcli.GetBalancesCmd(), args, result)
//	return result
//}
//
//func QueryAccountExec(
//	t *testing.T,
//	network Network,
//	clientCtx client.Context,
//	address string,
//	extraArgs ...string,
//) authtypes.AccountI {
//	t.Helper()
//	args := []string{
//		address,
//		fmt.Sprintf("--%s=json", "output"),
//	}
//	args = append(args, extraArgs...)
//	out, err := clitestutil.ExecTestCLICmd(clientCtx, authcli.GetAccountCmd(), args)
//	require.NoError(t, err, "QueryAccountExec  failed")
//
//	respType := proto.Message(&codectypes.Any{})
//	require.NoError(t, clientCtx.Codec.UnmarshalJSON(out.Bytes(), respType))
//
//	var account authtypes.AccountI
//	err = clientCtx.InterfaceRegistry.UnpackAny(respType.(*codectypes.Any), &account)
//	require.NoError(t, err, "UnpackAccount failed")
//
//	return account
//}
//
//func MsgSendExec(
//	t *testing.T,
//	network Network,
//	clientCtx client.Context,
//	from, to, amount fmt.Stringer,
//	extraArgs ...string,
//) *ResponseTx {
//	t.Helper()
//	args := []string{from.String(), to.String(), amount.String()}
//	args = append(args, extraArgs...)
//
//	return network.ExecTxCmdWithResult(t, clientCtx, bankcli.NewSendTxCmd(), args)
//}
//
//func QueryTx(t *testing.T, clientCtx client.Context, txHash string) abci.ResponseDeliverTx {
//	t.Helper()
//	txResult, _ := QueryTxWithHeight(t, clientCtx, txHash)
//	return txResult
//}
//
//func QueryTxWithHeight(
//	t *testing.T,
//	clientCtx client.Context,
//	txHash string,
//) (abci.ResponseDeliverTx, int64) {
//	t.Helper()
//	txHashBz, err := hex.DecodeString(txHash)
//	require.NoError(t, err, "query tx failed")
//
//	txResult, err := clientCtx.Client.Tx(context.Background(), txHashBz, false)
//	require.NoError(t, err, "query tx failed")
//	return txResult.TxResult, txResult.Height
//}
//
//// NewTestNetworkFixture returns a new simapp AppConstructor for network simulation tests
//func NewTestNetworkFixture(depInjectOptions DepinjectOptions) network.TestFixture {
//	dir, err := os.MkdirTemp("", "simapp")
//	if err != nil {
//		panic(fmt.Sprintf("failed creating temporary directory: %v", err))
//	}
//	defer os.RemoveAll(dir)
//
//	app := NewSimApp(
//		log.NewNopLogger(),
//		dbm.NewMemDB(),
//		nil,
//		true,
//		depInjectOptions,
//		simtestutil.NewAppOptionsWithFlagHome(dir),
//	)
//
//	appCtr := func(val network.ValidatorI) servertypes.Application {
//		return NewSimApp(
//			val.GetCtx().Logger, dbm.NewMemDB(), nil, true, depInjectOptions,
//			simtestutil.NewAppOptionsWithFlagHome(val.GetCtx().Config.RootDir),
//			bam.SetPruning(pruningtypes.NewPruningOptionsFromString(val.GetAppConfig().Pruning)),
//			bam.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
//			bam.SetChainID(val.GetCtx().Viper.GetString(flags.FlagChainID)),
//		)
//	}
//
//	return network.TestFixture{
//		AppConstructor: appCtr,
//		GenesisState:   app.DefaultGenesis(),
//		EncodingConfig: testutil.TestEncodingConfig{
//			InterfaceRegistry: app.InterfaceRegistry(),
//			Codec:             app.AppCodec(),
//			TxConfig:          app.TxConfig(),
//			Amino:             app.LegacyAmino(),
//		},
//	}
//}
