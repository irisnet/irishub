package simapp

// import (
// 	"encoding/hex"
// 	"encoding/json"
// 	"fmt"
// 	"math/rand"
// 	"os"
// 	"runtime/debug"
// 	"strings"
// 	"testing"

// 	dbm "github.com/cometbft/cometbft-db"
// 	abci "github.com/cometbft/cometbft/abci/types"
// 	"github.com/cometbft/cometbft/libs/log"
// 	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
// 	"github.com/stretchr/testify/require"

// 	"github.com/cosmos/cosmos-sdk/baseapp"
// 	"github.com/cosmos/cosmos-sdk/client/flags"
// 	"github.com/cosmos/cosmos-sdk/server"
// 	"github.com/cosmos/cosmos-sdk/store"
// 	storetypes "github.com/cosmos/cosmos-sdk/store/types"
// 	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
// 	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
// 	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
// 	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
// 	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
// 	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
// 	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
// 	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
// 	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
// 	"github.com/cosmos/cosmos-sdk/x/simulation"
// 	simcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"
// 	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
// 	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

// 	coinswaptypes "github.com/irisnet/irismod/modules/coinswap/types"
// 	htlctypes "github.com/irisnet/irismod/modules/htlc/types"
// 	mttypes "github.com/irisnet/irismod/modules/mt/types"
// 	nfttypes "github.com/irisnet/irismod/modules/nft/types"
// 	oracletypes "github.com/irisnet/irismod/modules/oracle/types"
// 	randomtypes "github.com/irisnet/irismod/modules/random/types"
// 	servicetypes "github.com/irisnet/irismod/modules/service/types"
// 	tokentypes "github.com/irisnet/irismod/modules/token/types"
// 	"github.com/irisnet/irismod/simapp/helpers"
// )

// // SimAppChainID hardcoded chainID for simulation
// const SimAppChainID = "simulation-app"

// // Get flags every time the simulator is run
// func init() {
// 	simcli.GetSimulatorFlags()
// }

// type StoreKeysPrefixes struct {
// 	A        storetypes.StoreKey
// 	B        storetypes.StoreKey
// 	Prefixes [][]byte
// }

// // fauxMerkleModeOpt returns a BaseApp option to use a dbStoreAdapter instead of
// // an IAVLStore for faster simulation speed.
// func fauxMerkleModeOpt(bapp *baseapp.BaseApp) {
// 	bapp.SetFauxMerkleMode()
// }

// // interBlockCacheOpt returns a BaseApp option function that sets the persistent
// // inter-block write-through cache.
// func interBlockCacheOpt() func(*baseapp.BaseApp) {
// 	return baseapp.SetInterBlockCache(store.NewCommitKVStoreCacheManager())
// }

// func TestFullAppSimulation(t *testing.T) {
// 	config := simcli.NewConfigFromFlags()
// 	config.ChainID = SimAppChainID

// 	db, dir, logger, skip, err := simtestutil.SetupSimulation(
// 		config,
// 		"leveldb-app-sim",
// 		"Simulation",
// 		simcli.FlagVerboseValue,
// 		simcli.FlagEnabledValue,
// 	)
// 	if skip {
// 		t.Skip("skipping application simulation")
// 	}
// 	require.NoError(t, err, "simulation setup failed")

// 	defer func() {
// 		require.NoError(t, db.Close())
// 		require.NoError(t, os.RemoveAll(dir))
// 	}()

// 	appOptions := make(simtestutil.AppOptionsMap, 0)
// 	appOptions[flags.FlagHome] = DefaultNodeHome
// 	appOptions[server.FlagInvCheckPeriod] = simcli.FlagPeriodValue

// 	app := NewSimApp(
// 		logger,
// 		db,
// 		nil,
// 		true,
// 		appOptions,
// 		fauxMerkleModeOpt,
// 		baseapp.SetChainID(config.ChainID),
// 	)
// 	require.Equal(t, "SimApp", app.Name())

// 	// run randomized simulation
// 	_, simParams, simErr := simulation.SimulateFromSeed(
// 		t,
// 		os.Stdout,
// 		app.BaseApp,
// 		AppStateFn(app.AppCodec(), app.SimulationManager()),
// 		simtypes.RandomAccounts, // Replace with own random account function if using keys other than secp256k1
// 		simtestutil.SimulationOperations(app, app.AppCodec(), config),
// 		BlockedAddresses(),
// 		config,
// 		app.AppCodec(),
// 	)

// 	// export state and simParams before the simulation error is checked
// 	err = simtestutil.CheckExportSimulation(app, config, simParams)
// 	require.NoError(t, err)
// 	require.NoError(t, simErr)

// 	if config.Commit {
// 		simtestutil.PrintStats(db)
// 	}
// }

// func TestAppImportExport(t *testing.T) {
// 	config := simcli.NewConfigFromFlags()
// 	config.ChainID = SimAppChainID

// 	db, dir, logger, skip, err := simtestutil.SetupSimulation(
// 		config,
// 		"leveldb-app-sim",
// 		"Simulation",
// 		simcli.FlagVerboseValue,
// 		simcli.FlagEnabledValue,
// 	)
// 	if skip {
// 		t.Skip("skipping application simulation")
// 	}
// 	require.NoError(t, err, "simulation setup failed")

// 	defer func() {
// 		require.NoError(t, db.Close())
// 		require.NoError(t, os.RemoveAll(dir))
// 	}()

// 	appOptions := make(simtestutil.AppOptionsMap, 0)
// 	appOptions[flags.FlagHome] = DefaultNodeHome
// 	appOptions[server.FlagInvCheckPeriod] = simcli.FlagPeriodValue

// 	app := NewSimApp(
// 		logger,
// 		db,
// 		nil,
// 		true,
// 		appOptions,
// 		fauxMerkleModeOpt,
// 		baseapp.SetChainID(config.ChainID),
// 	)
// 	require.Equal(t, "SimApp", app.Name())

// 	// Run randomized simulation
// 	_, simParams, simErr := simulation.SimulateFromSeed(
// 		t,
// 		os.Stdout,
// 		app.BaseApp,
// 		AppStateFn(app.AppCodec(), app.SimulationManager()),
// 		simtypes.RandomAccounts, // Replace with own random account function if using keys other than secp256k1
// 		simtestutil.SimulationOperations(app, app.AppCodec(), config),
// 		BlockedAddresses(),
// 		config,
// 		app.AppCodec(),
// 	)

// 	// export state and simParams before the simulation error is checked
// 	err = simtestutil.CheckExportSimulation(app, config, simParams)
// 	require.NoError(t, err)
// 	require.NoError(t, simErr)

// 	if config.Commit {
// 		simtestutil.PrintStats(db)
// 	}

// 	fmt.Printf("exporting genesis...\n")

// 	exported, err := app.ExportAppStateAndValidators(false, []string{}, []string{})
// 	require.NoError(t, err)

// 	fmt.Printf("importing genesis...\n")

// 	newDB, newDir, _, _, err := simtestutil.SetupSimulation(
// 		config,
// 		"leveldb-app-sim-2",
// 		"Simulation-2",
// 		simcli.FlagVerboseValue,
// 		simcli.FlagEnabledValue,
// 	)
// 	require.NoError(t, err, "simulation setup failed")

// 	defer func() {
// 		require.NoError(t, newDB.Close())
// 		require.NoError(t, os.RemoveAll(newDir))
// 	}()

// 	newApp := NewSimApp(
// 		log.NewNopLogger(),
// 		newDB,
// 		nil,
// 		true,
// 		appOptions,
// 		fauxMerkleModeOpt,
// 		baseapp.SetChainID(config.ChainID),
// 	)
// 	require.Equal(t, "SimApp", newApp.Name())

// 	var genesisState GenesisState
// 	err = json.Unmarshal(exported.AppState, &genesisState)
// 	require.NoError(t, err)

// 	defer func() {
// 		if r := recover(); r != nil {
// 			err := fmt.Sprintf("%v", r)
// 			if !strings.Contains(err, "validator set is empty after InitGenesis") {
// 				panic(r)
// 			}
// 			logger.Info("Skipping simulation as all validators have been unbonded")
// 			logger.Info("err", err, "stacktrace", string(debug.Stack()))
// 		}
// 	}()

// 	ctxA := app.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()})
// 	ctxB := newApp.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()})
// 	newApp.ModuleManager.InitGenesis(ctxB, app.AppCodec(), genesisState)
// 	newApp.StoreConsensusParams(ctxB, exported.ConsensusParams)

// 	fmt.Printf("comparing stores...\n")

// 	storeKeysPrefixes := []StoreKeysPrefixes{
// 		{app.GetKey(authtypes.StoreKey), newApp.GetKey(authtypes.StoreKey), [][]byte{}},
// 		{app.GetKey(stakingtypes.StoreKey), newApp.GetKey(stakingtypes.StoreKey),
// 			[][]byte{
// 				stakingtypes.UnbondingQueueKey, stakingtypes.RedelegationQueueKey,
// 				stakingtypes.ValidatorQueueKey, stakingtypes.HistoricalInfoKey,
// 			}}, // ordering may change but it doesn't matter
// 		{app.GetKey(slashingtypes.StoreKey), newApp.GetKey(slashingtypes.StoreKey), [][]byte{}},
// 		{app.GetKey(minttypes.StoreKey), newApp.GetKey(minttypes.StoreKey), [][]byte{}},
// 		{app.GetKey(distrtypes.StoreKey), newApp.GetKey(distrtypes.StoreKey), [][]byte{}},
// 		{
// 			app.GetKey(banktypes.StoreKey),
// 			newApp.GetKey(banktypes.StoreKey),
// 			[][]byte{banktypes.BalancesPrefix},
// 		},
// 		{app.GetKey(paramtypes.StoreKey), newApp.GetKey(paramtypes.StoreKey), [][]byte{}},
// 		{app.GetKey(govtypes.StoreKey), newApp.GetKey(govtypes.StoreKey), [][]byte{}},
// 		{app.GetKey(evidencetypes.StoreKey), newApp.GetKey(evidencetypes.StoreKey), [][]byte{}},
// 		{app.GetKey(capabilitytypes.StoreKey), newApp.GetKey(capabilitytypes.StoreKey), [][]byte{}},

// 		// check irismod module
// 		{app.GetKey(tokentypes.StoreKey), newApp.GetKey(tokentypes.StoreKey), [][]byte{}},
// 		{app.GetKey(oracletypes.StoreKey), newApp.GetKey(oracletypes.StoreKey), [][]byte{}},
// 		//mt.Supply is InitSupply, can be not equal to TotalSupply
// 		{app.GetKey(mttypes.StoreKey), newApp.GetKey(mttypes.StoreKey), [][]byte{mttypes.PrefixMT}},
// 		{app.GetKey(nfttypes.StoreKey), newApp.GetKey(nfttypes.StoreKey), [][]byte{{0x05}}},
// 		{
// 			app.GetKey(servicetypes.StoreKey),
// 			newApp.GetKey(servicetypes.StoreKey),
// 			[][]byte{servicetypes.InternalCounterKey},
// 		},
// 		{
// 			app.GetKey(randomtypes.StoreKey),
// 			newApp.GetKey(randomtypes.StoreKey),
// 			[][]byte{randomtypes.RandomKey},
// 		},
// 		//{app.keys[recordtypes.StoreKey), newApp.keys[recordtypes.StoreKey), [][]byte{recordtypes.IntraTxCounterKey}},
// 		{app.GetKey(htlctypes.StoreKey), newApp.GetKey(htlctypes.StoreKey), [][]byte{}},
// 		{app.GetKey(coinswaptypes.StoreKey), newApp.GetKey(coinswaptypes.StoreKey), [][]byte{}},
// 	}

// 	for _, skp := range storeKeysPrefixes {
// 		storeA := ctxA.KVStore(skp.A)
// 		storeB := ctxB.KVStore(skp.B)

// 		failedKVAs, failedKVBs := sdk.DiffKVStores(storeA, storeB, skp.Prefixes)
// 		require.Equal(t, len(failedKVAs), len(failedKVBs), "unequal sets of key-values to compare")

// 		fmt.Printf(
// 			"compared %d different key/value pairs between %s and %s\n",
// 			len(failedKVAs),
// 			skp.A,
// 			skp.B,
// 		)
// 		for _, kv := range failedKVAs {
// 			fmt.Printf("storeKey: %s,\n failedKVBs: %v,\n ", skp.A.Name(), kv.Key)
// 		}
// 		require.Equal(
// 			t,
// 			0,
// 			len(failedKVAs),
// 			simtestutil.GetSimulationLog(
// 				skp.A.Name(),
// 				app.SimulationManager().StoreDecoders,
// 				failedKVAs,
// 				failedKVBs,
// 			),
// 		)
// 	}
// }

// func TestAppSimulationAfterImport(t *testing.T) {
// 	config := simcli.NewConfigFromFlags()
// 	config.ChainID = SimAppChainID

// 	db, dir, logger, skip, err := simtestutil.SetupSimulation(
// 		config,
// 		"leveldb-app-sim",
// 		"Simulation",
// 		simcli.FlagVerboseValue,
// 		simcli.FlagEnabledValue,
// 	)
// 	if skip {
// 		t.Skip("skipping application simulation")
// 	}
// 	require.NoError(t, err, "simulation setup failed")

// 	defer func() {
// 		require.NoError(t, db.Close())
// 		require.NoError(t, os.RemoveAll(dir))
// 	}()

// 	appOptions := make(simtestutil.AppOptionsMap, 0)
// 	appOptions[flags.FlagHome] = DefaultNodeHome
// 	appOptions[server.FlagInvCheckPeriod] = simcli.FlagPeriodValue

// 	app := NewSimApp(
// 		logger,
// 		db,
// 		nil,
// 		true,
// 		appOptions,
// 		fauxMerkleModeOpt,
// 		baseapp.SetChainID(config.ChainID),
// 	)
// 	require.Equal(t, "SimApp", app.Name())

// 	// Run randomized simulation
// 	stopEarly, simParams, simErr := simulation.SimulateFromSeed(
// 		t,
// 		os.Stdout,
// 		app.BaseApp,
// 		AppStateFn(app.AppCodec(), app.SimulationManager()),
// 		simtypes.RandomAccounts, // Replace with own random account function if using keys other than secp256k1
// 		simtestutil.SimulationOperations(app, app.AppCodec(), config),
// 		BlockedAddresses(),
// 		config,
// 		app.AppCodec(),
// 	)

// 	// export state and simParams before the simulation error is checked
// 	err = simtestutil.CheckExportSimulation(app, config, simParams)
// 	require.NoError(t, err)
// 	require.NoError(t, simErr)

// 	if config.Commit {
// 		simtestutil.PrintStats(db)
// 	}

// 	if stopEarly {
// 		fmt.Println("can't export or import a zero-validator genesis, exiting test...")
// 		return
// 	}

// 	fmt.Printf("exporting genesis...\n")

// 	exported, err := app.ExportAppStateAndValidators(true, []string{}, []string{})
// 	require.NoError(t, err)

// 	fmt.Printf("importing genesis...\n")

// 	newDB, newDir, _, _, err := simtestutil.SetupSimulation(
// 		config,
// 		"leveldb-app-sim-2",
// 		"Simulation-2",
// 		simcli.FlagVerboseValue,
// 		simcli.FlagEnabledValue,
// 	)
// 	require.NoError(t, err, "simulation setup failed")

// 	defer func() {
// 		require.NoError(t, newDB.Close())
// 		require.NoError(t, os.RemoveAll(newDir))
// 	}()

// 	newApp := NewSimApp(
// 		log.NewNopLogger(),
// 		newDB,
// 		nil,
// 		true,
// 		appOptions,
// 		fauxMerkleModeOpt,
// 		baseapp.SetChainID(config.ChainID),
// 	)
// 	require.Equal(t, "SimApp", newApp.Name())

// 	newApp.InitChain(abci.RequestInitChain{
// 		AppStateBytes: exported.AppState,
// 	})
// 	_, _, err = simulation.SimulateFromSeed(
// 		t,
// 		os.Stdout,
// 		newApp.BaseApp,
// 		AppStateFn(app.AppCodec(), app.SimulationManager()),
// 		simtypes.RandomAccounts, // Replace with own random account function if using keys other than secp256k1
// 		simtestutil.SimulationOperations(newApp, newApp.AppCodec(), config),
// 		BlockedAddresses(),
// 		config,
// 		app.AppCodec(),
// 	)
// 	require.NoError(t, err)
// }

// // TODO: Make another test for the fuzzer itself, which just has noOp txs
// // and doesn't depend on the application.
// func TestAppStateDeterminism(t *testing.T) {
// 	if !simcli.FlagEnabledValue {
// 		t.Skip("skipping application simulation")
// 	}

// 	config := simcli.NewConfigFromFlags()
// 	config.InitialBlockHeight = 1
// 	config.ExportParamsPath = ""
// 	config.OnOperation = false
// 	config.AllInvariants = false
// 	config.ChainID = helpers.SimAppChainID

// 	numSeeds := 3
// 	numTimesToRunPerSeed := 5
// 	appHashList := make([]json.RawMessage, numTimesToRunPerSeed)

// 	appOptions := make(simtestutil.AppOptionsMap, 0)
// 	appOptions[flags.FlagHome] = DefaultNodeHome
// 	appOptions[server.FlagInvCheckPeriod] = simcli.FlagPeriodValue

// 	for i := 0; i < numSeeds; i++ {
// 		config.Seed = rand.Int63()

// 		for j := 0; j < numTimesToRunPerSeed; j++ {
// 			var logger log.Logger
// 			if simcli.FlagVerboseValue {
// 				logger = log.TestingLogger()
// 			} else {
// 				logger = log.NewNopLogger()
// 			}

// 			db := dbm.NewMemDB()
// 			app := NewSimApp(
// 				logger,
// 				db,
// 				nil,
// 				true,
// 				appOptions,
// 				interBlockCacheOpt(),
// 				baseapp.SetChainID(config.ChainID),
// 			)

// 			fmt.Printf(
// 				"running non-determinism simulation; seed %d: %d/%d, attempt: %d/%d\n",
// 				config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
// 			)

// 			_, _, err := simulation.SimulateFromSeed(
// 				t,
// 				os.Stdout,
// 				app.BaseApp,
// 				AppStateFn(app.AppCodec(), app.SimulationManager()),
// 				simtypes.RandomAccounts, // Replace with own random account function if using keys other than secp256k1
// 				simtestutil.SimulationOperations(app, app.AppCodec(), config),
// 				BlockedAddresses(),
// 				config,
// 				app.AppCodec(),
// 			)
// 			require.NoError(t, err)

// 			if config.Commit {
// 				simtestutil.PrintStats(db)
// 			}

// 			appHash := app.LastCommitID().Hash
// 			appHashList[j] = appHash

// 			if j != 0 {
// 				require.Equal(
// 					t, hex.EncodeToString(appHashList[0]), hex.EncodeToString(appHashList[j]),
// 					"non-determinism in seed %d: %d/%d, attempt: %d/%d\n",
// 					config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
// 				)
// 			}
// 		}
// 	}
// }
