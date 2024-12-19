package app

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"runtime/debug"
	"strings"
	"testing"

	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	"github.com/stretchr/testify/require"
	coinswaptypes "mods.irisnet.org/modules/coinswap/types"
	htlctypes "mods.irisnet.org/modules/htlc/types"
	mttypes "mods.irisnet.org/modules/mt/types"
	nfttypes "mods.irisnet.org/modules/nft/types"
	oracletypes "mods.irisnet.org/modules/oracle/types"
	randomtypes "mods.irisnet.org/modules/random/types"
	servicetypes "mods.irisnet.org/modules/service/types"
	tokentypes "mods.irisnet.org/modules/token/types"

	"github.com/irisnet/irishub/v4/app/params"
	iristypes "github.com/irisnet/irishub/v4/types"

	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	storetypes "cosmossdk.io/store/types"
	evidencetypes "cosmossdk.io/x/evidence/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	simcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
)

// AppChainID hardcoded chainID for simulation
const AppChainID = "irishub-1"

// Get flags every time the simulator is run
func init() {
	simcli.GetSimulatorFlags()
}

type StoreKeysPrefixes struct {
	A        storetypes.StoreKey
	B        storetypes.StoreKey
	Prefixes [][]byte
}

// fauxMerkleModeOpt returns a BaseApp option to use a dbStoreAdapter instead of
// an IAVLStore for faster simulation speed.
func fauxMerkleModeOpt(bapp *baseapp.BaseApp) {
	bapp.SetFauxMerkleMode()
}

// interBlockCacheOpt returns a BaseApp option function that sets the persistent
// inter-block write-through cache.
func interBlockCacheOpt() func(*baseapp.BaseApp) {
	return baseapp.SetInterBlockCache(store.NewCommitKVStoreCacheManager())
}

func TestFullAppSimulation(t *testing.T) {
	config := simcli.NewConfigFromFlags()
	config.ChainID = AppChainID

	db, dir, logger, skip, err := simtestutil.SetupSimulation(
		config,
		"goleveldb-app-sim",
		"Simulation",
		simcli.FlagVerboseValue,
		simcli.FlagEnabledValue,
	)

	if skip {
		t.Skip("skipping benchmark application simulation")
	}
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		db.Close()
		require.NoError(t, os.RemoveAll(dir))
	}()

	encfg := RegisterEncodingConfig()

	app := createApp(logger, db, encfg, fauxMerkleModeOpt)
	require.Equal(t, "IrisApp", app.Name())

	// run randomized simulation
	_, simParams, simErr := simulation.SimulateFromSeed(
		t,
		os.Stdout,
		app.BaseApp,
		simtestutil.AppStateFn(app.AppCodec(), app.SimulationManager(), app.DefaultGenesis()),
		simtypes.RandomAccounts, // Replace with own random account function if using keys other than secp256k1
		simtestutil.SimulationOperations(app, app.AppCodec(), config),
		app.ModuleAccountAddrs(),
		config,
		app.AppCodec(),
	)

	// export state and simParams before the simulation error is checked
	err = simtestutil.CheckExportSimulation(app, config, simParams)
	require.NoError(t, err)
	require.NoError(t, simErr)

	if config.Commit {
		simtestutil.PrintStats(db)
	}
}

func TestAppImportExport(t *testing.T) {
	config := simcli.NewConfigFromFlags()
	config.ChainID = AppChainID

	sdk.DefaultBondDenom = iristypes.NativeToken.Symbol

	db, dir, logger, skip, err := simtestutil.SetupSimulation(
		config,
		"goleveldb-app-sim",
		"Simulation",
		simcli.FlagVerboseValue,
		simcli.FlagEnabledValue,
	)

	if skip {
		t.Skip("skipping benchmark application simulation")
	}
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		db.Close()
		require.NoError(t, os.RemoveAll(dir))
	}()

	encfg := RegisterEncodingConfig()
	app := createApp(logger, db, encfg, fauxMerkleModeOpt)
	require.Equal(t, "IrisApp", app.Name())

	// Run randomized simulation
	_, simParams, simErr := simulation.SimulateFromSeed(
		t,
		os.Stdout,
		app.BaseApp,
		simtestutil.AppStateFn(app.AppCodec(), app.SimulationManager(), app.DefaultGenesis()),
		simtypes.RandomAccounts, // Replace with own random account function if using keys other than secp256k1
		simtestutil.SimulationOperations(app, app.AppCodec(), config),
		app.ModuleAccountAddrs(),
		config,
		app.AppCodec(),
	)

	// export state and simParams before the simulation error is checked
	err = simtestutil.CheckExportSimulation(app, config, simParams)
	require.NoError(t, err)
	require.NoError(t, simErr)

	if config.Commit {
		simtestutil.PrintStats(db)
	}

	fmt.Printf("exporting genesis...\n")

	exported, err := app.ExportAppStateAndValidators(false, []string{}, []string{})
	require.NoError(t, err)

	fmt.Printf("importing genesis...\n")

	newDB, newDir, _, _, err := simtestutil.SetupSimulation(
		config,
		"goleveldb-app-sim-2",
		"Simulation-2",
		simcli.FlagVerboseValue,
		simcli.FlagEnabledValue,
	)
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		newDB.Close()
		require.NoError(t, os.RemoveAll(newDir))
	}()

	newApp := createApp(logger, newDB, encfg, fauxMerkleModeOpt)
	require.Equal(t, "IrisApp", newApp.Name())

	var genesisState iristypes.GenesisState
	err = json.Unmarshal(exported.AppState, &genesisState)
	require.NoError(t, err)

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("%v", r)
			if !strings.Contains(err, "validator set is empty after InitGenesis") {
				panic(r)
			}
			logger.Info("Skipping simulation as all validators have been unbonded")
			logger.Info("err", err, "stacktrace", string(debug.Stack()))
		}
	}()

	ctxA := app.NewContextLegacy(true, tmproto.Header{Height: app.LastBlockHeight()})
	ctxB := newApp.NewContextLegacy(true, tmproto.Header{Height: app.LastBlockHeight()})
	newApp.mm.InitGenesis(ctxB, app.AppCodec(), genesisState)
	newApp.StoreConsensusParams(ctxB, exported.ConsensusParams)

	fmt.Printf("comparing stores...\n")

	storeKeysPrefixes := []StoreKeysPrefixes{
		{app.AppKeepers.KvStoreKeys()[authtypes.StoreKey], newApp.AppKeepers.KvStoreKeys()[authtypes.StoreKey], [][]byte{}},
		{
			app.AppKeepers.KvStoreKeys()[stakingtypes.StoreKey], newApp.AppKeepers.KvStoreKeys()[stakingtypes.StoreKey],
			[][]byte{
				stakingtypes.UnbondingQueueKey, stakingtypes.RedelegationQueueKey, stakingtypes.ValidatorQueueKey,
				stakingtypes.HistoricalInfoKey,
			},
		}, // ordering may change but it doesn't matter
		{app.AppKeepers.KvStoreKeys()[slashingtypes.StoreKey], newApp.AppKeepers.KvStoreKeys()[slashingtypes.StoreKey], [][]byte{}},
		{app.AppKeepers.KvStoreKeys()[minttypes.StoreKey], newApp.AppKeepers.KvStoreKeys()[minttypes.StoreKey], [][]byte{}},
		{app.AppKeepers.KvStoreKeys()[distrtypes.StoreKey], newApp.AppKeepers.KvStoreKeys()[distrtypes.StoreKey], [][]byte{}},
		{
			app.AppKeepers.KvStoreKeys()[banktypes.StoreKey],
			newApp.AppKeepers.KvStoreKeys()[banktypes.StoreKey],
			[][]byte{banktypes.BalancesPrefix},
		},
		{app.AppKeepers.KvStoreKeys()[paramtypes.StoreKey], newApp.AppKeepers.KvStoreKeys()[paramtypes.StoreKey], [][]byte{}},
		{app.AppKeepers.KvStoreKeys()[govtypes.StoreKey], newApp.AppKeepers.KvStoreKeys()[govtypes.StoreKey], [][]byte{}},
		{app.AppKeepers.KvStoreKeys()[evidencetypes.StoreKey], newApp.AppKeepers.KvStoreKeys()[evidencetypes.StoreKey], [][]byte{}},
		{app.AppKeepers.KvStoreKeys()[capabilitytypes.StoreKey], newApp.AppKeepers.KvStoreKeys()[capabilitytypes.StoreKey], [][]byte{}},
		{app.AppKeepers.KvStoreKeys()[ibcexported.StoreKey], newApp.AppKeepers.KvStoreKeys()[ibcexported.StoreKey], [][]byte{}},
		{app.AppKeepers.KvStoreKeys()[ibctransfertypes.StoreKey], newApp.AppKeepers.KvStoreKeys()[ibctransfertypes.StoreKey], [][]byte{}},

		// check irismod module
		{app.AppKeepers.KvStoreKeys()[tokentypes.StoreKey], newApp.AppKeepers.KvStoreKeys()[tokentypes.StoreKey], [][]byte{}},
		{app.AppKeepers.KvStoreKeys()[oracletypes.StoreKey], newApp.AppKeepers.KvStoreKeys()[oracletypes.StoreKey], [][]byte{}},
		// mt.Supply is InitSupply, can be not equal to TotalSupply
		{app.AppKeepers.KvStoreKeys()[mttypes.StoreKey], newApp.AppKeepers.KvStoreKeys()[mttypes.StoreKey], [][]byte{mttypes.PrefixMT}},
		{app.AppKeepers.KvStoreKeys()[nfttypes.StoreKey], newApp.AppKeepers.KvStoreKeys()[nfttypes.StoreKey], [][]byte{{0x05}}},
		{
			app.AppKeepers.KvStoreKeys()[servicetypes.StoreKey],
			newApp.AppKeepers.KvStoreKeys()[servicetypes.StoreKey],
			[][]byte{servicetypes.InternalCounterKey},
		},
		{
			app.AppKeepers.KvStoreKeys()[randomtypes.StoreKey],
			newApp.AppKeepers.KvStoreKeys()[randomtypes.StoreKey],
			[][]byte{randomtypes.RandomKey},
		},
		{app.AppKeepers.KvStoreKeys()[htlctypes.StoreKey], newApp.AppKeepers.KvStoreKeys()[htlctypes.StoreKey], [][]byte{}},
		{app.AppKeepers.KvStoreKeys()[coinswaptypes.StoreKey], newApp.AppKeepers.KvStoreKeys()[coinswaptypes.StoreKey], [][]byte{}},
	}

	for _, skp := range storeKeysPrefixes {
		storeA := ctxA.KVStore(skp.A)
		storeB := ctxB.KVStore(skp.B)

		failedKVAs, failedKVBs := simtestutil.DiffKVStores(storeA, storeB, skp.Prefixes)
		require.Equal(t, len(failedKVAs), len(failedKVBs), "unequal sets of key-values to compare")

		fmt.Printf(
			"compared %d different key/value pairs between %s and %s\n",
			len(failedKVAs),
			skp.A,
			skp.B,
		)
		require.Equal(
			t,
			len(failedKVAs),
			0,
			simtestutil.GetSimulationLog(
				skp.A.Name(),
				app.SimulationManager().StoreDecoders,
				failedKVAs,
				failedKVBs,
			),
		)
	}
}

func TestAppSimulationAfterImport(t *testing.T) {
	config := simcli.NewConfigFromFlags()
	config.ChainID = AppChainID
	encfg := RegisterEncodingConfig()

	db, dir, logger, skip, err := simtestutil.SetupSimulation(
		config,
		"goleveldb-app-sim",
		"Simulation",
		simcli.FlagVerboseValue,
		simcli.FlagEnabledValue,
	)

	if skip {
		t.Skip("skipping benchmark application simulation")
	}
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		db.Close()
		require.NoError(t, os.RemoveAll(dir))
	}()

	app := createApp(logger, db, encfg, fauxMerkleModeOpt)
	require.Equal(t, "IrisApp", app.Name())

	// Run randomized simulation
	stopEarly, simParams, simErr := simulation.SimulateFromSeed(
		t,
		os.Stdout,
		app.BaseApp,
		simtestutil.AppStateFn(app.AppCodec(), app.SimulationManager(), app.DefaultGenesis()),
		simtypes.RandomAccounts, // Replace with own random account function if using keys other than secp256k1
		simtestutil.SimulationOperations(app, app.AppCodec(), config),
		app.ModuleAccountAddrs(),
		config,
		app.AppCodec(),
	)

	// export state and simParams before the simulation error is checked
	err = simtestutil.CheckExportSimulation(app, config, simParams)
	require.NoError(t, err)
	require.NoError(t, simErr)

	if config.Commit {
		simtestutil.PrintStats(db)
	}

	if stopEarly {
		fmt.Println("can't export or import a zero-validator genesis, exiting test...")
		return
	}

	fmt.Printf("exporting genesis...\n")

	exported, err := app.ExportAppStateAndValidators(true, []string{}, []string{})
	require.NoError(t, err)

	fmt.Printf("importing genesis...\n")

	newDB, newDir, _, _, err := simtestutil.SetupSimulation(
		config,
		"goleveldb-app-sim-2",
		"Simulation-2",
		simcli.FlagVerboseValue,
		simcli.FlagEnabledValue,
	)
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		newDB.Close()
		require.NoError(t, os.RemoveAll(newDir))
	}()

	newApp := createApp(logger, newDB, encfg, fauxMerkleModeOpt)
	require.Equal(t, "IrisApp", newApp.Name())

	_, err = newApp.InitChain(&abci.RequestInitChain{
		AppStateBytes: exported.AppState,
		ChainId:       AppChainID,
	})

	require.NoError(t, err)

	_, _, err = simulation.SimulateFromSeed(
		t,
		os.Stdout,
		newApp.BaseApp,
		simtestutil.AppStateFn(app.AppCodec(), app.SimulationManager(), app.DefaultGenesis()),
		simtypes.RandomAccounts, // Replace with own random account function if using keys other than secp256k1
		simtestutil.SimulationOperations(newApp, newApp.AppCodec(), config),
		app.ModuleAccountAddrs(),
		config,
		newApp.AppCodec(),
	)
	require.NoError(t, err)
}

// TODO: Make another test for the fuzzer itself, which just has noOp txs
// and doesn't depend on the application.
func TestAppStateDeterminism(t *testing.T) {
	if !simcli.FlagEnabledValue {
		t.Skip("skipping application simulation")
	}
	sdk.DefaultBondDenom = iristypes.NativeToken.Symbol

	config := simcli.NewConfigFromFlags()
	config.InitialBlockHeight = 1
	config.ExportParamsPath = ""
	config.OnOperation = false
	config.AllInvariants = false
	config.ChainID = AppChainID

	numSeeds := 3
	numTimesToRunPerSeed := 5
	appHashList := make([]json.RawMessage, numTimesToRunPerSeed)
	encfg := RegisterEncodingConfig()

	for i := 0; i < numSeeds; i++ {
		config.Seed = rand.Int63()

		for j := 0; j < numTimesToRunPerSeed; j++ {
			var logger log.Logger
			if simcli.FlagVerboseValue {
				logger = log.NewTestLogger(t)
			} else {
				logger = log.NewNopLogger()
			}

			db := dbm.NewMemDB()
			app := createApp(logger, db, encfg, interBlockCacheOpt())

			fmt.Printf(
				"running non-determinism simulation; seed %d: %d/%d, attempt: %d/%d\n",
				config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
			)

			_, _, err := simulation.SimulateFromSeed(
				t,
				os.Stdout,
				app.BaseApp,
				simtestutil.AppStateFn(
					app.AppCodec(),
					app.SimulationManager(),
					app.DefaultGenesis(),
				),
				simtypes.RandomAccounts, // Replace with own random account function if using keys other than secp256k1
				simtestutil.SimulationOperations(app, app.AppCodec(), config),
				app.ModuleAccountAddrs(),
				config,
				app.AppCodec(),
			)
			require.NoError(t, err)

			if config.Commit {
				simtestutil.PrintStats(db)
			}

			appHash := app.LastCommitID().Hash
			appHashList[j] = appHash

			if j != 0 {
				require.Equal(
					t,
					string(appHashList[0]),
					string(appHashList[j]),
					"non-determinism in seed %d: %d/%d, attempt: %d/%d\n",
					config.Seed,
					i+1,
					numSeeds,
					j+1,
					numTimesToRunPerSeed,
				)
			}
		}
	}
}

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(o string) interface{} {
	return nil
}

// SimTestAppOptions is a stub implementing AppOptions
type SimTestAppOptions struct {
	options map[string]interface{}
}

// Get implements AppOptions
func (o SimTestAppOptions) Get(key string) interface{} {
	return o.options[key]
}

func createApp(
	logger log.Logger,
	db dbm.DB,
	encodingConfig params.EncodingConfig,
	baseAppOptions ...func(*baseapp.BaseApp),
) *IrisApp {
	if baseAppOptions == nil {
		baseAppOptions = []func(*baseapp.BaseApp){}
	}
	baseAppOptions = append(baseAppOptions, baseapp.SetChainID(AppChainID))
	return NewIrisApp(
		logger,
		db,
		nil,
		true,
		encodingConfig,
		SimTestAppOptions{
			options: map[string]interface{}{params.SimulationTest: true},
		},
		baseAppOptions...,
	)
}
