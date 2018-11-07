package app

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/stake"
	banksim "github.com/irisnet/irishub/simulation/bank"
	govsim "github.com/irisnet/irishub/simulation/gov"
	"github.com/irisnet/irishub/simulation/mock/simulation"
	slashingsim "github.com/irisnet/irishub/simulation/slashing"
	stakesim "github.com/irisnet/irishub/simulation/stake"
	"os"
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
	flag.BoolVar(&enabled, "SimulationEnabled", false, "Enable the simulation")
	flag.BoolVar(&verbose, "SimulationVerbose", false, "Verbose log output")
	flag.BoolVar(&commit, "SimulationCommit", false, "Have the simulation commit")
}

func appStateFn(r *rand.Rand, accs []simulation.Account) json.RawMessage {
	var genesisAccounts []GenesisAccount

	amt := int64(10000)

	// Randomly generate some genesis accounts
	for _, acc := range accs {
		coins := sdk.Coins{sdk.Coin{"iris-atto", sdk.NewInt(amt)}}
		genesisAccounts = append(genesisAccounts, GenesisAccount{
			Address: acc.Address,
			Coins:   coins,
		})
	}

	// Default genesis state
	//govGenesis := gov.DefaultGenesisState()
	stakeGenesis := stake.DefaultGenesisState()
	slashingGenesis := slashing.DefaultGenesisState()
	var validators []stake.Validator
	var delegations []stake.Delegation

	// XXX Try different numbers of initially bonded validators
	numInitiallyBonded := int64(50)
	valAddrs := make([]sdk.ValAddress, numInitiallyBonded)
	for i := 0; i < int(numInitiallyBonded); i++ {
		valAddr := sdk.ValAddress(accs[i].Address)
		valAddrs[i] = valAddr

		validator := stake.NewValidator(valAddr, accs[i].PubKey, stake.Description{})
		validator.Tokens = sdk.NewDec(amt)
		validator.DelegatorShares = sdk.NewDec(amt)
		delegation := stake.Delegation{accs[i].Address, valAddr, sdk.NewDec(amt), 0}
		validators = append(validators, validator)
		delegations = append(delegations, delegation)
	}
	stakeGenesis.Pool.LooseTokens = sdk.NewDec(amt*250 + (numInitiallyBonded * amt))
	stakeGenesis.Validators = validators
	stakeGenesis.Bonds = delegations
	mintGenesis := mint.DefaultGenesisState()

	genesis := GenesisState{
		Accounts:     genesisAccounts,
		StakeData:    stakeGenesis,
		MintData:     mintGenesis,
		DistrData:    distr.DefaultGenesisWithValidators(valAddrs),
		SlashingData: slashingGenesis,
		// TODO: can't use govGenesis as GovData
		//GovData:      govGenesis,
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
	return []simulation.Invariant{
		banksim.NonnegativeBalanceInvariant(app.accountMapper),
		govsim.AllInvariants(),
		stakesim.AllInvariants(app.bankKeeper, app.stakeKeeper,
			app.feeCollectionKeeper, app.distrKeeper, app.accountMapper),
		slashingsim.AllInvariants(),
	}
}

// Profile with:
// /usr/local/go/bin/go test -benchmem -run=^$ github.com/cosmos/cosmos-sdk/cmd/gaia/app -bench ^BenchmarkFullGaiaSimulation$ -SimulationCommit=true -cpuprofile cpu.out
func BenchmarkFullIrisSimulation(b *testing.B) {
	// Setup Gaia application
	var logger log.Logger
	logger = log.NewNopLogger()
	var db dbm.DB
	dir := os.TempDir()
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

	// Setup Gaia application
	var logger log.Logger
	if verbose {
		logger = log.TestingLogger()
	} else {
		logger = log.NewNopLogger()
	}
	db := dbm.NewMemDB()
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
		fmt.Println("Database Size", db.Stats()["database.size"])
	}
	require.Nil(t, err)
}

// TODO: Make another test for the fuzzer itself, which just has noOp txs
// and doesn't depend on gaia
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
