package app

import (
	"encoding/json"
	"io"
	"os"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/irisnet/irishub/modules/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"errors"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/viper"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/node"
	sm "github.com/tendermint/tendermint/state"
)

const (
	appName    = "IrisApp"
	FlagReplay = "replay"
)

// default home directories for expected binaries
var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.iriscli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.iris")
)

// Extended ABCI application
type IrisApp struct {
	*bam.BaseApp
	cdc *wire.Codec

	// keys to access the substores
	keyMain          *sdk.KVStoreKey
	keyAccount       *sdk.KVStoreKey
	keyStake         *sdk.KVStoreKey
	keySlashing      *sdk.KVStoreKey
	keyGov           *sdk.KVStoreKey
	keyFeeCollection *sdk.KVStoreKey
	keyParams        *sdk.KVStoreKey

	// Manage getting and setting accounts
	accountMapper       auth.AccountMapper
	feeCollectionKeeper auth.FeeCollectionKeeper
	coinKeeper          bank.Keeper
	stakeKeeper         stake.Keeper
	slashingKeeper      slashing.Keeper
	govKeeper           gov.Keeper
	paramsKeeper        params.Keeper
}

func NewIrisApp(logger log.Logger, db dbm.DB, traceStore io.Writer, baseAppOptions ...func(*bam.BaseApp)) *IrisApp {
	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, cdc, logger, db, baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)

	// create your application object
	var app = &IrisApp{
		BaseApp:          bam.NewBaseApp(appName, cdc, logger, db),
		cdc:              cdc,
		keyMain:          sdk.NewKVStoreKey("main"),
		keyAccount:       sdk.NewKVStoreKey("acc"),
		keyStake:         sdk.NewKVStoreKey("stake"),
		keySlashing:      sdk.NewKVStoreKey("slashing"),
		keyGov:           sdk.NewKVStoreKey("gov"),
		keyFeeCollection: sdk.NewKVStoreKey("fee"),
		keyParams:        sdk.NewKVStoreKey("params"),
	}

	// define the accountMapper
	app.accountMapper = auth.NewAccountMapper(
		app.cdc,
		app.keyAccount,        // target store
		auth.ProtoBaseAccount, // prototype
	)

	// add handlers
	app.paramsKeeper = params.NewKeeper(app.cdc, app.keyParams)
	app.coinKeeper = bank.NewKeeper(app.accountMapper)
	app.stakeKeeper = stake.NewKeeper(app.cdc, app.keyStake, app.coinKeeper, app.RegisterCodespace(stake.DefaultCodespace))
	app.slashingKeeper = slashing.NewKeeper(app.cdc, app.keySlashing, app.stakeKeeper, app.RegisterCodespace(slashing.DefaultCodespace))
	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(app.cdc, app.keyFeeCollection, app.paramsKeeper.Getter())
	app.govKeeper = gov.NewKeeper(app.cdc, app.keyGov, app.paramsKeeper.Setter(), app.coinKeeper, app.stakeKeeper, app.RegisterCodespace(gov.DefaultCodespace))
	//app.govKeeper = gov.NewKeeper(app.cdc, app.keyGov, app.paramsKeeper.Setter(), app.coinKeeper, app.stakeKeeper, app.RegisterCodespace(gov.DefaultCodespace))

	// register message routes
	// need to update each module's msg type
	app.Router().
		AddRoute("bank", []*sdk.KVStoreKey{app.keyAccount}, bank.NewHandler(app.coinKeeper)).
		AddRoute("stake", []*sdk.KVStoreKey{app.keyStake, app.keyAccount}, stake.NewHandler(app.stakeKeeper)).
		AddRoute("slashing", []*sdk.KVStoreKey{app.keySlashing, app.keyStake}, slashing.NewHandler(app.slashingKeeper)).
		AddRoute("gov", []*sdk.KVStoreKey{app.keyGov, app.keyAccount, app.keyStake, app.keyParams}, gov.NewHandler(app.govKeeper))

	// initialize BaseApp
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper, app.feeCollectionKeeper))
	app.SetFeeRefundHandler(auth.NewFeeRefundHandler(app.accountMapper, app.feeCollectionKeeper))
	app.MountStoresIAVL(app.keyMain, app.keyAccount, app.keyStake, app.keySlashing, app.keyGov, app.keyFeeCollection, app.keyParams)
	var err error
	err = app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
}

// custom tx codec
func MakeCodec() *wire.Codec {
	var cdc = wire.NewCodec()
	bank.RegisterWire(cdc)
	stake.RegisterWire(cdc)
	slashing.RegisterWire(cdc)
	gov.RegisterWire(cdc)
	auth.RegisterWire(cdc)
	sdk.RegisterWire(cdc)
	wire.RegisterCrypto(cdc)
	return cdc
}

// application updates every end block
func (app *IrisApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	tags := slashing.BeginBlocker(ctx, req, app.slashingKeeper)

	return abci.ResponseBeginBlock{
		Tags: tags.ToKVPairs(),
	}
}

// application updates every end block
func (app *IrisApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	validatorUpdates := stake.EndBlocker(ctx, app.stakeKeeper)

	tags,_ := gov.EndBlocker(ctx, app.govKeeper)
	return abci.ResponseEndBlock{
		ValidatorUpdates: validatorUpdates,
		Tags:             tags,
	}
}

// custom logic for iris initialization
func (app *IrisApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	var genesisState GenesisState
	err := app.cdc.UnmarshalJSON(stateJSON, &genesisState)
	if err != nil {
		panic(err)
	}

	// load the accounts
	for _, gacc := range genesisState.Accounts {
		acc := gacc.ToAccount()
		acc.AccountNumber = app.accountMapper.GetNextAccountNumber(ctx)
		app.accountMapper.SetAccount(ctx, acc)
	}

	// load the initial stake information
	err = stake.InitGenesis(ctx, app.stakeKeeper, genesisState.StakeData)
	if err != nil {
		panic(err) // TODO https://github.com/cosmos/cosmos-sdk/issues/468
		// return sdk.ErrGenesisParse("").TraceCause(err, "")
	}

	gov.InitGenesis(ctx, app.govKeeper, gov.GenesisState{
		StartingProposalID: 1,
		DepositProcedure: gov.DepositProcedure{
			MinDeposit:       sdk.Coins{sdk.Coin{Denom: "iris", Amount: sdk.NewInt(int64(10)).Mul(gov.Pow10(18))}},
			MaxDepositPeriod: 1440,
		},
		VotingProcedure: gov.VotingProcedure{
			VotingPeriod: 5,
		},
		TallyingProcedure: gov.TallyingProcedure{
			Threshold:         sdk.NewRat(1, 2),
			Veto:              sdk.NewRat(1, 3),
			GovernancePenalty: sdk.NewRat(1, 100),
		},
	})

	feeTokenGensisConfig := auth.GenesisState{
		FeeTokenNative:    "iris",
		GasPriceThreshold: 20000000000, // 20(glue), 20*10^9, 1 glue = 10^9 lue/gas, 1 iris = 10^18 lue
	}

	auth.InitGenesis(ctx, app.paramsKeeper.Setter(), feeTokenGensisConfig)


	return abci.ResponseInitChain{}
}

// export the state of iris for a genesis file
func (app *IrisApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{})

	// iterate to get the accounts
	accounts := []GenesisAccount{}
	appendAccount := func(acc auth.Account) (stop bool) {
		account := NewGenesisAccountI(acc)
		accounts = append(accounts, account)
		return false
	}
	app.accountMapper.IterateAccounts(ctx, appendAccount)

	genState := GenesisState{
		Accounts:  accounts,
		StakeData: stake.WriteGenesis(ctx, app.stakeKeeper),
	}
	appState, err = wire.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}
	validators = stake.WriteValidators(ctx, app.stakeKeeper)
	return appState, validators, nil
}

func (app *IrisApp) replay() int64 {
	ctx := server.NewDefaultContext()
	ctx.Config.RootDir = viper.GetString(tmcli.HomeFlag)
	dbContext := node.DBContext{"state", ctx.Config}
	dbType := dbm.DBBackendType(dbContext.Config.DBBackend)
	stateDB := dbm.NewDB(dbContext.ID, dbType, dbContext.Config.DBDir())

	preState := sm.LoadPreState(stateDB)
	if preState.LastBlockHeight == 0 {
		panic(errors.New("can't replay the last block, last block height is 0"))
	}

	sm.SaveState(stateDB, preState)
	stateDB.Close()

	return preState.LastBlockHeight
}
