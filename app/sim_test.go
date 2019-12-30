package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	banksim "github.com/cosmos/cosmos-sdk/x/bank/simulation"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrsim "github.com/cosmos/cosmos-sdk/x/distribution/simulation"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govsim "github.com/cosmos/cosmos-sdk/x/gov/simulation"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsim "github.com/cosmos/cosmos-sdk/x/params/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingsim "github.com/cosmos/cosmos-sdk/x/slashing/simulation"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingsim "github.com/cosmos/cosmos-sdk/x/staking/simulation"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/irisnet/irishub/modules/mint"
	randmodule "github.com/irisnet/irishub/modules/rand"
	randsim "github.com/irisnet/irishub/modules/rand/simulation"
)

func init() {
	simapp.GetSimulatorFlags()
}

func testAndRunTxs(app *IrisApp, config simulation.Config) []simulation.WeightedOperation {
	ap := make(simulation.AppParams)

	paramChanges := app.sm.GenerateParamChanges(config.Seed)

	if config.ParamsFile != "" {
		bz, err := ioutil.ReadFile(config.ParamsFile)
		if err != nil {
			panic(err)
		}

		app.cdc.MustUnmarshalJSON(bz, &ap)
	}

	// nolint: govet
	return []simulation.WeightedOperation{{
		Weight: func(_ *rand.Rand) int {
			var v int
			ap.GetOrGenerate(app.cdc, OpWeightMsgSend, &v, nil,
				func(_ *rand.Rand) {
					v = 100
				})
			return v
		}(nil),
		Op: banksim.SimulateMsgSend(app.accountKeeper, app.bankKeeper),
	}, {
		Weight: func(_ *rand.Rand) int {
			var v int
			ap.GetOrGenerate(app.cdc, OpWeightMsgMultiSend, &v, nil,
				func(_ *rand.Rand) {
					v = 40
				})
			return v
		}(nil),
		Op: banksim.SimulateMsgMultiSend(app.accountKeeper, app.bankKeeper),
	}, {
		Weight: func(_ *rand.Rand) int {
			var v int
			ap.GetOrGenerate(app.cdc, OpWeightMsgSetWithdrawAddress, &v, nil,
				func(_ *rand.Rand) {
					v = 50
				})
			return v
		}(nil),
		Op: distrsim.SimulateMsgSetWithdrawAddress(app.accountKeeper, app.distrKeeper),
	}, {
		Weight: func(_ *rand.Rand) int {
			var v int
			ap.GetOrGenerate(app.cdc, OpWeightMsgWithdrawDelegationReward, &v, nil,
				func(_ *rand.Rand) {
					v = 50
				})
			return v
		}(nil),
		Op: distrsim.SimulateMsgWithdrawDelegatorReward(app.accountKeeper, app.distrKeeper, app.stakingKeeper),
	}, {
		Weight: func(_ *rand.Rand) int {
			var v int
			ap.GetOrGenerate(app.cdc, OpWeightMsgWithdrawValidatorCommission, &v, nil,
				func(_ *rand.Rand) {
					v = 50
				})
			return v
		}(nil),
		Op: distrsim.SimulateMsgWithdrawValidatorCommission(app.accountKeeper, app.distrKeeper, app.stakingKeeper),
	}, {
		Weight: func(_ *rand.Rand) int {
			var v int
			ap.GetOrGenerate(app.cdc, OpWeightSubmitTextProposal, &v, nil,
				func(_ *rand.Rand) {
					v = 20
				})
			return v
		}(nil),
		Op: govsim.SimulateSubmitProposal(app.accountKeeper, app.govKeeper, govsim.SimulateTextProposalContent),
	}, {
		Weight: func(_ *rand.Rand) int {
			var v int
			ap.GetOrGenerate(app.cdc, OpWeightSubmitCommunitySpendProposal, &v, nil,
				func(_ *rand.Rand) {
					v = 20
				})
			return v
		}(nil),
		Op: govsim.SimulateSubmitProposal(app.accountKeeper, app.govKeeper, distrsim.SimulateCommunityPoolSpendProposalContent(app.distrKeeper)),
	}, {
		Weight: func(_ *rand.Rand) int {
			var v int
			ap.GetOrGenerate(app.cdc, OpWeightSubmitParamChangeProposal, &v, nil,
				func(_ *rand.Rand) {
					v = 20
				})
			return v
		}(nil),
		Op: govsim.SimulateSubmitProposal(app.accountKeeper, app.govKeeper, paramsim.SimulateParamChangeProposalContent(paramChanges)),
	}, {
		Weight: func(_ *rand.Rand) int {
			var v int
			ap.GetOrGenerate(app.cdc, OpWeightMsgDeposit, &v, nil,
				func(_ *rand.Rand) {
					v = 100
				})
			return v
		}(nil),
		Op: govsim.SimulateMsgDeposit(app.accountKeeper, app.govKeeper),
	}, {
		Weight: func(_ *rand.Rand) int {
			var v int
			ap.GetOrGenerate(app.cdc, OpWeightMsgVote, &v, nil,
				func(_ *rand.Rand) {
					v = 100
				})
			return v
		}(nil),
		Op: govsim.SimulateMsgVote(app.accountKeeper, app.govKeeper),
	}, {
		Weight: func(_ *rand.Rand) int {
			var v int
			ap.GetOrGenerate(app.cdc, OpWeightMsgCreateValidator, &v, nil,
				func(_ *rand.Rand) {
					v = 100
				})
			return v
		}(nil),
		Op: stakingsim.SimulateMsgCreateValidator(app.accountKeeper, app.stakingKeeper),
	}, {
		Weight: func(_ *rand.Rand) int {
			var v int
			ap.GetOrGenerate(app.cdc, OpWeightMsgEditValidator, &v, nil,
				func(_ *rand.Rand) {
					v = 20
				})
			return v
		}(nil),
		Op: stakingsim.SimulateMsgEditValidator(app.accountKeeper, app.stakingKeeper),
	}, {
		Weight: func(_ *rand.Rand) int {
			var v int
			ap.GetOrGenerate(app.cdc, OpWeightMsgDelegate, &v, nil,
				func(_ *rand.Rand) {
					v = 100
				})
			return v
		}(nil),
		Op: stakingsim.SimulateMsgDelegate(app.accountKeeper, app.stakingKeeper),
	}, {
		Weight: func(_ *rand.Rand) int {
			var v int
			ap.GetOrGenerate(app.cdc, OpWeightMsgUndelegate, &v, nil,
				func(_ *rand.Rand) {
					v = 100
				})
			return v
		}(nil),
		Op: stakingsim.SimulateMsgUndelegate(app.accountKeeper, app.stakingKeeper),
	}, {
		Weight: func(_ *rand.Rand) int {
			var v int
			ap.GetOrGenerate(app.cdc, OpWeightMsgBeginRedelegate, &v, nil,
				func(_ *rand.Rand) {
					v = 100
				})
			return v
		}(nil),
		Op: stakingsim.SimulateMsgBeginRedelegate(app.accountKeeper, app.stakingKeeper),
	}, {
		Weight: func(_ *rand.Rand) int {
			var v int
			ap.GetOrGenerate(app.cdc, OpWeightMsgUnjail, &v, nil,
				func(_ *rand.Rand) {
					v = 100
				})
			return v
		}(nil),
		Op: slashingsim.SimulateMsgUnjail(app.accountKeeper, app.slashingKeeper, app.stakingKeeper),
	}, {
		Weight: func(_ *rand.Rand) int {
			var v int
			ap.GetOrGenerate(app.cdc, OpWeightMsgRequestRand, &v, nil,
				func(_ *rand.Rand) {
					v = 20
				})
			return v
		}(nil),
		Op: randsim.SimulateMsgRequestRand(app.accountKeeper, app.randKeeper),
	}}
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
	if !simapp.FlagEnabledValue {
		t.Skip("skipping application simulation")
	}

	var logger log.Logger
	config := simapp.NewConfigFromFlags()

	if simapp.FlagVerboseValue {
		logger = log.TestingLogger()
	} else {
		logger = log.NewNopLogger()
	}

	var db dbm.DB
	dir, _ := ioutil.TempDir("", "goleveldb-app-sim")
	db, _ = sdk.NewLevelDB("Simulation", dir)

	defer func() {
		db.Close()
		_ = os.RemoveAll(dir)
	}()

	app := NewIrisApp(logger, db, nil, true, simapp.FlagPeriodValue, fauxMerkleModeOpt)
	require.Equal(t, "IrisApp", app.Name())

	// Run randomized simulation
	_, simParams, simErr := simulation.SimulateFromSeed(
		t, os.Stdout, app.BaseApp, simapp.AppStateFn(app.Codec(), app.sm),
		testAndRunTxs(app, config), app.ModuleAccountAddrs(), config,
	)

	// export state and params before the simulation error is checked
	if config.ExportStatePath != "" {
		err := ExportStateToJSON(app, config.ExportStatePath)
		require.NoError(t, err)
	}

	if config.ExportParamsPath != "" {
		err := simapp.CheckExportSimulation(app, config, simParams)
		require.NoError(t, err)
	}

	require.NoError(t, simErr)

	if config.Commit {
		// for memdb:
		// fmt.Println("Database Size", db.Stats()["database.size"])
		fmt.Println("\nGoLevelDB Stats")
		fmt.Println(db.Stats()["leveldb.stats"])
		fmt.Println("GoLevelDB cached block size", db.Stats()["leveldb.cachedblock"])
	}
}

func TestAppImportExport(t *testing.T) {
	if !simapp.FlagEnabledValue {
		t.Skip("skipping application import/export simulation")
	}

	var logger log.Logger
	config := simapp.NewConfigFromFlags()

	if simapp.FlagVerboseValue {
		logger = log.TestingLogger()
	} else {
		logger = log.NewNopLogger()
	}

	var db dbm.DB
	dir, _ := ioutil.TempDir("", "goleveldb-app-sim")
	db, _ = sdk.NewLevelDB("Simulation", dir)

	defer func() {
		db.Close()
		_ = os.RemoveAll(dir)
	}()

	app := NewIrisApp(logger, db, nil, true, simapp.FlagPeriodValue, fauxMerkleModeOpt)
	require.Equal(t, "SimApp", app.Name())

	// Run randomized simulation
	_, simParams, simErr := simulation.SimulateFromSeed(
		t, os.Stdout, app.BaseApp, simapp.AppStateFn(app.Codec(), app.sm),
		testAndRunTxs(app, config), app.ModuleAccountAddrs(), config,
	)

	// export state and simParams before the simulation error is checked
	err := simapp.CheckExportSimulation(app, config, simParams)
	require.NoError(t, err)
	require.NoError(t, simErr)

	if config.Commit {
		// for memdb:
		// fmt.Println("Database Size", db.Stats()["database.size"])
		fmt.Println("\nGoLevelDB Stats")
		fmt.Println(db.Stats()["leveldb.stats"])
		fmt.Println("GoLevelDB cached block size", db.Stats()["leveldb.cachedblock"])
	}

	fmt.Printf("exporting genesis...\n")

	appState, _, err := app.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err)
	fmt.Printf("importing genesis...\n")

	newDir, _ := ioutil.TempDir("", "goleveldb-app-sim-2")
	newDB, _ := sdk.NewLevelDB("Simulation-2", dir)

	defer func() {
		newDB.Close()
		_ = os.RemoveAll(newDir)
	}()

	newApp := NewIrisApp(log.NewNopLogger(), newDB, nil, true, simapp.FlagPeriodValue, fauxMerkleModeOpt)
	require.Equal(t, "SimApp", newApp.Name())

	var genesisState simapp.GenesisState
	err = app.cdc.UnmarshalJSON(appState, &genesisState)
	require.NoError(t, err)

	ctxB := newApp.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
	newApp.mm.InitGenesis(ctxB, genesisState)

	fmt.Printf("comparing stores...\n")
	ctxA := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})

	type StoreKeysPrefixes struct {
		A        sdk.StoreKey
		B        sdk.StoreKey
		Prefixes [][]byte
	}

	storeKeysPrefixes := []StoreKeysPrefixes{
		{app.keys[baseapp.MainStoreKey], newApp.keys[baseapp.MainStoreKey], [][]byte{}},
		{app.keys[auth.StoreKey], newApp.keys[auth.StoreKey], [][]byte{}},
		{app.keys[staking.StoreKey], newApp.keys[staking.StoreKey],
			[][]byte{
				staking.UnbondingQueueKey, staking.RedelegationQueueKey, staking.ValidatorQueueKey,
			}}, // ordering may change but it doesn't matter
		{app.keys[slashing.StoreKey], newApp.keys[slashing.StoreKey], [][]byte{}},
		{app.keys[mint.StoreKey], newApp.keys[mint.StoreKey], [][]byte{}},
		{app.keys[distr.StoreKey], newApp.keys[distr.StoreKey], [][]byte{}},
		{app.keys[supply.StoreKey], newApp.keys[supply.StoreKey], [][]byte{}},
		{app.keys[params.StoreKey], newApp.keys[params.StoreKey], [][]byte{}},
		{app.keys[gov.StoreKey], newApp.keys[gov.StoreKey], [][]byte{}},
		{app.keys[randmodule.StoreKey], newApp.keys[randmodule.StoreKey], [][]byte{}},
	}

	for _, storeKeysPrefix := range storeKeysPrefixes {
		storeKeyA := storeKeysPrefix.A
		storeKeyB := storeKeysPrefix.B
		prefixes := storeKeysPrefix.Prefixes

		storeA := ctxA.KVStore(storeKeyA)
		storeB := ctxB.KVStore(storeKeyB)

		failedKVAs, failedKVBs := sdk.DiffKVStores(storeA, storeB, prefixes)
		require.Equal(t, len(failedKVAs), len(failedKVBs), "unequal sets of key-values to compare")

		fmt.Printf("compared %d key/value pairs between %s and %s\n", len(failedKVAs), storeKeyA, storeKeyB)
		require.Len(t, failedKVAs, 0, simapp.GetSimulationLog(storeKeyA.Name(), app.sm.StoreDecoders, app.cdc, failedKVAs, failedKVBs))
	}
}

func TestAppSimulationAfterImport(t *testing.T) {
	if !simapp.FlagEnabledValue {
		t.Skip("skipping application simulation after import")
	}

	var logger log.Logger
	config := simapp.NewConfigFromFlags()

	if simapp.FlagVerboseValue {
		logger = log.TestingLogger()
	} else {
		logger = log.NewNopLogger()
	}

	dir, _ := ioutil.TempDir("", "goleveldb-app-sim")
	db, _ := sdk.NewLevelDB("Simulation", dir)

	defer func() {
		db.Close()
		_ = os.RemoveAll(dir)
	}()

	app := NewIrisApp(logger, db, nil, true, simapp.FlagPeriodValue, fauxMerkleModeOpt)
	require.Equal(t, "IrisApp", app.Name())

	// Run randomized simulation
	// Run randomized simulation
	stopEarly, simParams, simErr := simulation.SimulateFromSeed(
		t, os.Stdout, app.BaseApp, simapp.AppStateFn(app.Codec(), app.sm),
		testAndRunTxs(app, config), app.ModuleAccountAddrs(), config,
	)

	// export state and params before the simulation error is checked
	err := simapp.CheckExportSimulation(app, config, simParams)
	require.NoError(t, err)
	require.NoError(t, simErr)

	if config.Commit {
		// for memdb:
		// fmt.Println("Database Size", db.Stats()["database.size"])
		fmt.Println("\nGoLevelDB Stats")
		fmt.Println(db.Stats()["leveldb.stats"])
		fmt.Println("GoLevelDB cached block size", db.Stats()["leveldb.cachedblock"])
	}

	if stopEarly {
		// we can't export or import a zero-validator genesis
		fmt.Printf("We can't export or import a zero-validator genesis, exiting test...\n")
		return
	}

	fmt.Printf("Exporting genesis...\n")

	appState, _, err := app.ExportAppStateAndValidators(true, []string{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Importing genesis...\n")

	newDir, _ := ioutil.TempDir("", "goleveldb-app-sim-2")
	newDB, _ := sdk.NewLevelDB("Simulation-2", dir)

	defer func() {
		newDB.Close()
		_ = os.RemoveAll(newDir)
	}()

	newApp := NewIrisApp(log.NewNopLogger(), newDB, nil, true, 0, fauxMerkleModeOpt)
	require.Equal(t, "IrisApp", newApp.Name())

	newApp.InitChain(abci.RequestInitChain{
		AppStateBytes: appState,
	})

	// Run randomized simulation on imported app
	_, _, err = simulation.SimulateFromSeed(
		t, os.Stdout, newApp.BaseApp, simapp.AppStateFn(app.Codec(), app.sm),
		testAndRunTxs(newApp, config), newApp.ModuleAccountAddrs(), config,
	)

	require.NoError(t, err)
}

// TODO: Make another test for the fuzzer itself, which just has noOp txs
// and doesn't depend on the application.
func TestAppStateDeterminism(t *testing.T) {
	if !simapp.FlagEnabledValue {
		t.Skip("skipping application simulation")
	}

	config := simapp.NewConfigFromFlags()
	config.InitialBlockHeight = 1
	config.ExportParamsPath = ""
	config.OnOperation = false
	config.AllInvariants = false

	numSeeds := 3
	numTimesToRunPerSeed := 5
	appHashList := make([]json.RawMessage, numTimesToRunPerSeed)

	for i := 0; i < numSeeds; i++ {
		config.Seed = rand.Int63()

		for j := 0; j < numTimesToRunPerSeed; j++ {
			logger := log.NewNopLogger()
			db := dbm.NewMemDB()
			app := NewIrisApp(logger, db, nil, true, simapp.FlagPeriodValue, interBlockCacheOpt())

			fmt.Printf(
				"running non-determinism simulation; seed %d: %d/%d, attempt: %d/%d\n",
				config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
			)

			_, _, err := simulation.SimulateFromSeed(
				t, os.Stdout, app.BaseApp, simapp.AppStateFn(app.Codec(), app.sm),
				testAndRunTxs(app, config), app.ModuleAccountAddrs(), config,
			)
			require.NoError(t, err)

			appHash := app.LastCommitID().Hash
			appHashList[j] = appHash

			if j != 0 {
				require.Equal(
					t, appHashList[0], appHashList[j],
					"non-determinism in seed %d: %d/%d, attempt: %d/%d\n", config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
				)
			}
		}
	}
}
