package app

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/irisnet/irishub/modules/gov"
	banksim "github.com/irisnet/irishub/simulation/bank"
	govsim "github.com/irisnet/irishub/simulation/gov"
	"github.com/irisnet/irishub/simulation/mock/simulation"
	slashingsim "github.com/irisnet/irishub/simulation/slashing"
	stakesim "github.com/irisnet/irishub/simulation/stake"
)

var (
	seed      int64
	numBlocks int
	blockSize int
	enabled   bool
	verbose   bool
	commit    bool
)

func init() {
	flag.Int64Var(&seed, "SimulationSeed", 42, "Simulation random seed")
	flag.IntVar(&numBlocks, "SimulationNumBlocks", 500, "Number of blocks")
	flag.IntVar(&blockSize, "SimulationBlockSize", 200, "Operations per block")
	flag.BoolVar(&enabled, "SimulationEnabled", true, "Enable the simulation")
	flag.BoolVar(&verbose, "SimulationVerbose", false, "Verbose log output")
	flag.BoolVar(&commit, "SimulationCommit", false, "Have the simulation commit")
}

func appStateFn(r *rand.Rand, accs []simulation.Account) json.RawMessage {
	stakeGenesis := stake.DefaultGenesisState()
	fmt.Printf("Selected randomly generated staking parameters: %+v\n", stakeGenesis)

	var genesisAccounts []GenesisAccount

	amount := sdk.NewIntWithDecimal(100, 18)
	stakeAmount := sdk.NewIntWithDecimal(1, 2)
	numInitiallyBonded := int64(r.Intn(250))
	//numInitiallyBonded := int64(4)
	numAccs := int64(len(accs))
	if numInitiallyBonded > numAccs {
		numInitiallyBonded = numAccs
	}
	fmt.Printf("Selected randomly generated parameters for simulated genesis: {amount of iris-atto per account: %v, initially bonded validators: %v}\n", amount, numInitiallyBonded)

	// Randomly generate some genesis accounts
	for _, acc := range accs {
		coins := sdk.Coins{
			{
				Denom:  "iris-atto",
				Amount: amount,
			},
			{
				Denom:  stakeGenesis.Params.BondDenom,
				Amount: stakeAmount,
			},
		}
		genesisAccounts = append(genesisAccounts, GenesisAccount{
			Address: acc.Address,
			Coins:   coins,
		})
	}

	// Random genesis states
	govGenesis := gov.DefaultGenesisState()
	fmt.Printf("Selected randomly generated governance parameters: %+v\n", govGenesis)
	slashingGenesis := slashing.DefaultGenesisState()
	fmt.Printf("Selected randomly generated slashing parameters: %+v\n", slashingGenesis)
	mintGenesis := mint.DefaultGenesisState()
	fmt.Printf("Selected randomly generated minting parameters: %v\n", mintGenesis)
	var (
		validators  []stake.Validator
		delegations []stake.Delegation
	)

	decAmt := sdk.NewDecFromInt(sdk.NewIntWithDecimal(100, 18))
	valAddrs := make([]sdk.ValAddress, numInitiallyBonded)
	for i := 0; i < int(numInitiallyBonded); i++ {
		valAddr := sdk.ValAddress(accs[i].Address)
		valAddrs[i] = valAddr

		validator := stake.NewValidator(valAddr, accs[i].PubKey, stake.Description{})
		validator.Tokens = decAmt
		validator.DelegatorShares = decAmt
		delegation := stake.Delegation{
			DelegatorAddr: accs[i].Address,
			ValidatorAddr: valAddr,
			Shares:        decAmt,
			Height:        0,
		}
		validators = append(validators, validator)
		delegations = append(delegations, delegation)
	}
	stakeGenesis.Pool.LooseTokens = sdk.NewDecFromInt(sdk.NewIntWithDecimal(100, 30))
	stakeGenesis.Validators = validators
	stakeGenesis.Bonds = delegations

	genesis := GenesisState{
		Accounts:     genesisAccounts,
		StakeData:    stakeGenesis,
		MintData:     mintGenesis,
		DistrData:    distr.DefaultGenesisWithValidators(valAddrs),
		SlashingData: slashingGenesis,
		GovData:      govGenesis,
	}

	// Marshal genesis
	appState, err := MakeCodec().MarshalJSON(genesis)
	if err != nil {
		panic(err)
	}

	return appState
}

func testAndRunTxs(app *IrisApp) []simulation.WeightedOperation {
	return []simulation.WeightedOperation{
		{100, banksim.SingleInputSendMsg(app.accountMapper, app.bankKeeper)},
		{5, govsim.SimulateSubmittingVotingAndSlashingForProposal(app.govKeeper, app.stakeKeeper)},
		{100, govsim.SimulateMsgDeposit(app.govKeeper, app.stakeKeeper)},
		{100, stakesim.SimulateMsgCreateValidator(app.accountMapper, app.stakeKeeper)},
		{5, stakesim.SimulateMsgEditValidator(app.stakeKeeper)},
		{100, stakesim.SimulateMsgDelegate(app.accountMapper, app.stakeKeeper)},
		{100, stakesim.SimulateMsgBeginUnbonding(app.accountMapper, app.stakeKeeper)},
		{100, stakesim.SimulateMsgBeginRedelegate(app.accountMapper, app.stakeKeeper)},
		{100, slashingsim.SimulateMsgUnjail(app.slashingKeeper)},
	}
}

func invariants(app *IrisApp) []simulation.Invariant {
	return []simulation.Invariant{}
}

func BenchmarkFullIrisSimulation(b *testing.B) {
	// Setup Iris application
	var logger log.Logger
	logger = log.NewNopLogger()
	var db dbm.DB
	dir, _ := ioutil.TempDir("", "goleveldb-iris-sim")
	db, _ = dbm.NewGoLevelDB("Simulation", dir)
	defer func() {
		db.Close()
		os.RemoveAll(dir)
	}()
	app := NewIrisApp(logger, db, nil)

	// Run randomized simulation
	// TODO parameterize numbers, save for a later PR
	err := simulation.SimulateFromSeed(
		b, app.BaseApp, appStateFn, seed,
		testAndRunTxs(app),
		[]simulation.RandSetup{},
		invariants(app), // these shouldn't get ran
		numBlocks,
		blockSize,
		commit,
	)
	if err != nil {
		fmt.Println(err)
		b.Fail()
	}
	if commit {
		fmt.Println("GoLevelDB Stats")
		fmt.Println(db.Stats()["leveldb.stats"])
		fmt.Println("GoLevelDB cached block size", db.Stats()["leveldb.cachedblock"])
	}
}

func TestFullIrisSimulation(t *testing.T) {
	if !enabled {
		t.Skip("Skipping Iris simulation")
	}

	// Setup Iris application
	var logger log.Logger
	if verbose {
		logger = log.TestingLogger()
	} else {
		logger = log.NewNopLogger()
	}
	var db dbm.DB
	dir, _ := ioutil.TempDir("", "goleveldb-iris-sim")
	db, _ = dbm.NewGoLevelDB("Simulation", dir)
	defer func() {
		db.Close()
		os.RemoveAll(dir)
	}()
	app := NewIrisApp(logger, db, nil)
	require.Equal(t, "IrisApp", app.Name())

	// Run randomized simulation
	err := simulation.SimulateFromSeed(
		t, app.BaseApp, appStateFn, seed,
		testAndRunTxs(app),
		[]simulation.RandSetup{},
		invariants(app),
		numBlocks,
		blockSize,
		commit,
	)
	if commit {
		// for memdb:
		// fmt.Println("Database Size", db.Stats()["database.size"])
		fmt.Println("GoLevelDB Stats")
		fmt.Println(db.Stats()["leveldb.stats"])
		fmt.Println("GoLevelDB cached block size", db.Stats()["leveldb.cachedblock"])
	}
	require.Nil(t, err)
}

func TestIrisImportExport(t *testing.T) {
	if !enabled {
		t.Skip("Skipping Iris import/export simulation")
	}

	// Setup Iris application
	var logger log.Logger
	if verbose {
		logger = log.TestingLogger()
	} else {
		logger = log.NewNopLogger()
	}
	var db dbm.DB
	dir, _ := ioutil.TempDir("", "goleveldb-iris-sim")
	db, _ = dbm.NewGoLevelDB("Simulation", dir)
	defer func() {
		db.Close()
		os.RemoveAll(dir)
	}()
	app := NewIrisApp(logger, db, nil)
	require.Equal(t, "IrisApp", app.Name())

	// Run randomized simulation
	err := simulation.SimulateFromSeed(
		t, app.BaseApp, appStateFn, seed,
		testAndRunTxs(app),
		[]simulation.RandSetup{},
		invariants(app),
		numBlocks,
		blockSize,
		commit,
	)
	if commit {
		// for memdb:
		// fmt.Println("Database Size", db.Stats()["database.size"])
		fmt.Println("GoLevelDB Stats")
		fmt.Println(db.Stats()["leveldb.stats"])
		fmt.Println("GoLevelDB cached block size", db.Stats()["leveldb.cachedblock"])
	}
	require.Nil(t, err)

	fmt.Printf("Exporting genesis...\n")

	appState, _, err := app.ExportAppStateAndValidators()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Importing genesis...\n")

	newDir, _ := ioutil.TempDir("", "goleveldb-iris-sim-2")
	newDB, _ := dbm.NewGoLevelDB("Simulation-2", dir)
	defer func() {
		newDB.Close()
		os.RemoveAll(newDir)
	}()
	newApp := NewIrisApp(log.NewNopLogger(), newDB, nil)
	require.Equal(t, "IrisApp", newApp.Name())
	request := abci.RequestInitChain{
		AppStateBytes: appState,
	}
	newApp.InitChain(request)
	newApp.Commit()

	fmt.Printf("Comparing stores...\n")
	ctxA := app.NewContext(true, abci.Header{})
	ctxB := newApp.NewContext(true, abci.Header{})
	type StoreKeysPrefixes struct {
		A        sdk.StoreKey
		B        sdk.StoreKey
		Prefixes [][]byte
	}
	storeKeysPrefixes := []StoreKeysPrefixes{
		{app.keyMain, newApp.keyMain, [][]byte{}},
		{app.keyAccount, newApp.keyAccount, [][]byte{}},
		{app.keyStake, newApp.keyStake, [][]byte{stake.UnbondingQueueKey, stake.RedelegationQueueKey, stake.ValidatorQueueKey}}, // ordering may change but it doesn't matter
		{app.keySlashing, newApp.keySlashing, [][]byte{}},
		{app.keyMint, newApp.keyMint, [][]byte{}},
		{app.keyDistr, newApp.keyDistr, [][]byte{}},
		{app.keyFeeCollection, newApp.keyFeeCollection, [][]byte{}},
		{app.keyParams, newApp.keyParams, [][]byte{}},
		{app.keyGov, newApp.keyGov, [][]byte{}},
	}
	for _, storeKeysPrefix := range storeKeysPrefixes {
		storeKeyA := storeKeysPrefix.A
		storeKeyB := storeKeysPrefix.B
		prefixes := storeKeysPrefix.Prefixes
		storeA := ctxA.KVStore(storeKeyA)
		storeB := ctxB.KVStore(storeKeyB)
		kvA, kvB, count, equal := sdk.DiffKVStores(storeA, storeB, prefixes)
		fmt.Printf("Compared %d key/value pairs between %s and %s\n", count, storeKeyA, storeKeyB)
		require.True(t, equal, "unequal stores: %s / %s:\nstore A %s (%X) => %s (%X)\nstore B %s (%X) => %s (%X)",
			storeKeyA, storeKeyB, kvA.Key, kvA.Key, kvA.Value, kvA.Value, kvB.Key, kvB.Key, kvB.Value, kvB.Value)
	}

}

// TODO: Make another test for the fuzzer itself, which just has noOp txs
// and doesn't depend on iris
func TestAppStateDeterminism(t *testing.T) {
	if !enabled {
		t.Skip("Skipping Iris simulation")
	}

	numSeeds := 3
	numTimesToRunPerSeed := 5
	appHashList := make([]json.RawMessage, numTimesToRunPerSeed)

	for i := 0; i < numSeeds; i++ {
		seed := rand.Int63()
		for j := 0; j < numTimesToRunPerSeed; j++ {
			logger := log.NewNopLogger()
			db := dbm.NewMemDB()
			app := NewIrisApp(logger, db, nil)

			// Run randomized simulation
			simulation.SimulateFromSeed(
				t, app.BaseApp, appStateFn, seed,
				testAndRunTxs(app),
				[]simulation.RandSetup{},
				[]simulation.Invariant{},
				50,
				100,
				true,
			)
			//app.Commit()
			appHash := app.LastCommitID().Hash
			appHashList[j] = appHash
		}
		for k := 1; k < numTimesToRunPerSeed; k++ {
			require.Equal(t, appHashList[0], appHashList[k], "appHash list: %v", appHashList)
		}
	}
}
