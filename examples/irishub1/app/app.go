package app

import (
	"encoding/json"
	"io"
	"os"

	bam "github.com/irisnet/irishub/baseapp"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/stake"
	ibc1 "github.com/irisnet/irishub/examples/irishub1/ibc"
	"github.com/irisnet/irishub/modules/upgrade"

	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/viper"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/node"
	sm "github.com/tendermint/tendermint/state"
	bc "github.com/tendermint/tendermint/blockchain"
	"strings"
	"github.com/irisnet/irishub/modules/iparams"
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
	keyIBC           *sdk.KVStoreKey
	keyStake         *sdk.KVStoreKey
	keySlashing      *sdk.KVStoreKey
	keyGov           *sdk.KVStoreKey
	keyFeeCollection *sdk.KVStoreKey
	keyParams        *sdk.KVStoreKey
	keyUpgrade       *sdk.KVStoreKey

	// Manage getting and setting accounts
	accountMapper       auth.AccountMapper
	feeCollectionKeeper auth.FeeCollectionKeeper
	coinKeeper          bank.Keeper
	ibcMapper           ibc.Mapper
	ibc1Mapper          ibc1.Mapper
	stakeKeeper         stake.Keeper
	slashingKeeper      slashing.Keeper
	govKeeper           gov.Keeper
	paramsKeeper        params.Keeper
	upgradeKeeper       upgrade.Keeper
	iparamsKeeper       iparams.Keeper

	// fee manager
	feeManager  bam.FeeManager
}

func NewIrisApp(logger log.Logger, db dbm.DB, traceStore io.Writer, baseAppOptions ...func(*bam.BaseApp)) *IrisApp {
	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, cdc, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)

	// create your application object
	var app = &IrisApp{
		BaseApp:          bApp,
		cdc:              cdc,
		keyMain:          sdk.NewKVStoreKey("main"),
		keyAccount:       sdk.NewKVStoreKey("acc"),
		keyIBC:           sdk.NewKVStoreKey("ibc"),
		keyStake:         sdk.NewKVStoreKey("stake"),
		keySlashing:      sdk.NewKVStoreKey("slashing"),
		keyGov:           sdk.NewKVStoreKey("gov"),
		keyFeeCollection: sdk.NewKVStoreKey("fee"),
		keyParams:        sdk.NewKVStoreKey("params"),
		keyUpgrade:       sdk.NewKVStoreKey("upgrade"),
	}

	var lastHeight int64
	if viper.GetBool(FlagReplay) {
		lastHeight = app.replay()
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
	app.ibcMapper = ibc.NewMapper(app.cdc, app.keyIBC, app.RegisterCodespace(ibc.DefaultCodespace))
	app.ibc1Mapper = ibc1.NewMapper(app.cdc, app.keyIBC, app.RegisterCodespace(ibc1.DefaultCodespace))
	app.stakeKeeper = stake.NewKeeper(app.cdc, app.keyStake, app.coinKeeper, app.RegisterCodespace(stake.DefaultCodespace))
	app.slashingKeeper = slashing.NewKeeper(app.cdc, app.keySlashing, app.stakeKeeper, app.paramsKeeper.Getter(), app.RegisterCodespace(slashing.DefaultCodespace))
	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(app.cdc, app.keyFeeCollection)
	app.upgradeKeeper = upgrade.NewKeeper(app.cdc, app.keyUpgrade, app.stakeKeeper, app.iparamsKeeper.Setter())
	app.govKeeper = gov.NewKeeper(app.cdc, app.keyGov, app.paramsKeeper.Setter(), app.coinKeeper, app.stakeKeeper, app.RegisterCodespace(gov.DefaultCodespace))
	//app.govKeeper = gov.NewKeeper(app.cdc, app.keyGov, app.paramsKeeper.Setter(), app.coinKeeper, app.stakeKeeper, app.RegisterCodespace(gov.DefaultCodespace))

	// register message routes
	// need to update each module's msg type
	app.Router().
		AddRoute("bank", []*sdk.KVStoreKey{app.keyAccount}, bank.NewHandler(app.coinKeeper)).
		AddRoute("ibc", []*sdk.KVStoreKey{app.keyIBC, app.keyAccount}, ibc.NewHandler(app.ibcMapper, app.coinKeeper)).
		AddRoute("ibc-1", []*sdk.KVStoreKey{app.keyIBC, app.keyAccount}, ibc1.NewHandler(app.ibc1Mapper, app.coinKeeper)).
		AddRoute("stake", []*sdk.KVStoreKey{app.keyStake, app.keyAccount}, stake.NewHandler(app.stakeKeeper)).
		AddRoute("slashing", []*sdk.KVStoreKey{app.keySlashing, app.keyStake}, slashing.NewHandler(app.slashingKeeper)).
		AddRoute("gov", []*sdk.KVStoreKey{app.keyGov, app.keyAccount, app.keyStake, app.keyParams}, gov.NewHandler(app.govKeeper)).
		AddRoute("upgrade", []*sdk.KVStoreKey{app.keyUpgrade, app.keyStake}, upgrade.NewHandler(app.upgradeKeeper))

	app.feeManager = bam.NewFeeManager(app.iparamsKeeper.Getter())

	// initialize BaseApp
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper, app.feeCollectionKeeper))
	app.SetFeeRefundHandler(bam.NewFeeRefundHandler(app.accountMapper, app.feeCollectionKeeper, app.feeManager))
	app.SetFeePreprocessHandler(bam.NewFeePreprocessHandler(app.feeManager))
	app.MountStoresIAVL(app.keyMain, app.keyAccount, app.keyIBC, app.keyStake, app.keySlashing, app.keyGov, app.keyFeeCollection, app.keyParams, app.keyUpgrade)
	app.SetRunMsg(app.runMsgs)
	var err error
	if viper.GetBool(FlagReplay) {
		err = app.LoadVersion(lastHeight, app.keyMain)
	} else {
		err = app.LoadLatestVersion(app.keyMain)
	}
	if err != nil {
		cmn.Exit(err.Error())
	}

	upgrade.RegisterModuleList(app.Router())

	return app
}

// custom tx codec
func MakeCodec() *wire.Codec {
	var cdc = wire.NewCodec()
	ibc.RegisterWire(cdc)
	ibc1.RegisterWire(cdc)

	bank.RegisterWire(cdc)
	stake.RegisterWire(cdc)
	slashing.RegisterWire(cdc)
	gov.RegisterWire(cdc)
	auth.RegisterWire(cdc)
	upgrade.RegisterWire(cdc)
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

	tags := gov.EndBlocker(ctx, app.govKeeper)
	tags.AppendTags(upgrade.EndBlocker(ctx, app.upgradeKeeper))

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
	validators, err := stake.InitGenesis(ctx, app.stakeKeeper, genesisState.StakeData)
	if err != nil {
		panic(err)
	}

	minDeposit,err := IrisCt.ConvertToMinCoin(fmt.Sprintf("%d%s",10,denom))
	if err != nil {
		panic(err)
	}

	gov.InitGenesis(ctx, app.govKeeper, gov.GenesisState{
		StartingProposalID: 1,
		DepositProcedure: gov.DepositProcedure{
			MinDeposit:       sdk.Coins{minDeposit},
			MaxDepositPeriod: 1440,
		},
		VotingProcedure: gov.VotingProcedure{
			VotingPeriod: 30,
		},
		TallyingProcedure: gov.TallyingProcedure{
			Threshold:         sdk.NewRat(1, 2),
			Veto:              sdk.NewRat(1, 3),
			GovernancePenalty: sdk.NewRat(1, 100),
		},
	})

	feeTokenGensisConfig := bam.FeeGenesisStateConfig{
		FeeTokenNative:    IrisCt.MinUnit.Denom,
		GasPriceThreshold: 20000000000, // 20(glue), 20*10^9, 1 glue = 10^9 lue/gas, 1 iris = 10^18 lue
	}

	bam.InitGenesis(ctx, app.iparamsKeeper.Setter(), feeTokenGensisConfig)

	// load the address to pubkey map
	slashing.InitGenesis(ctx, app.slashingKeeper, genesisState.StakeData)

	upgrade.InitGenesis(ctx, app.upgradeKeeper, app.Router())

	return abci.ResponseInitChain{
		Validators: validators,
	}
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

// Iterates through msgs and executes them
func (app *IrisApp) runMsgs(ctx sdk.Context, msgs []sdk.Msg) (result sdk.Result) {
	// accumulate results
	logs := make([]string, 0, len(msgs))
	var data []byte   // NOTE: we just append them all (?!)
	var tags sdk.Tags // also just append them all
	var code sdk.ABCICodeType
	for msgIdx, msg := range msgs {
		// Match route.
		msgType, err := app.upgradeKeeper.GetMsgTypeInCurrentVersion(ctx, msg)
		fmt.Println("============ runMsgs() ===========  " + msgType)
		if err != nil {
			return err.Result()
		}

		handler := app.Router().Route(msgType)
		if handler == nil {
			return sdk.ErrUnknownRequest("Unrecognized Msg type: " + msgType).Result()
		}
		fmt.Println(msg)
		msgResult := handler(ctx, msg)

		// NOTE: GasWanted is determined by ante handler and
		// GasUsed by the GasMeter

		// Append Data and Tags
		data = append(data, msgResult.Data...)
		tags = append(tags, msgResult.Tags...)

		// Stop execution and return on first failed message.
		if !msgResult.IsOK() {
			logs = append(logs, fmt.Sprintf("Msg %d failed: %s", msgIdx, msgResult.Log))
			code = msgResult.Code
			break
		}

		// Construct usable logs in multi-message transactions.
		logs = append(logs, fmt.Sprintf("Msg %d: %s", msgIdx, msgResult.Log))
	}

	// Set the final gas values.
	result = sdk.Result{
		Code:    code,
		Data:    data,
		Log:     strings.Join(logs, "\n"),
		GasUsed: ctx.GasMeter().GasConsumed(),
		// TODO: FeeAmount/FeeDenom
		Tags: tags,
	}

	return result
}

func (app *IrisApp) replay() int64 {
	ctx := server.NewDefaultContext()
	ctx.Config.RootDir = viper.GetString(tmcli.HomeFlag)
	dbContext := node.DBContext{"state", ctx.Config}
	dbType := dbm.DBBackendType(dbContext.Config.DBBackend)
	stateDB := dbm.NewDB(dbContext.ID, dbType, dbContext.Config.DBDir())

	blockDBContext := node.DBContext{"blockstore", ctx.Config}
	blockStoreDB := dbm.NewDB(blockDBContext.ID, dbType, dbContext.Config.DBDir())
	blockStore := bc.NewBlockStore(blockStoreDB)

	defer func() {
		stateDB.Close()
		blockStoreDB.Close()
	} ()

	curState := sm.LoadState(stateDB)
	preState := sm.LoadPreState(stateDB)
	if curState.LastBlockHeight == preState.LastBlockHeight {
		panic(errors.New("there is no block now, can't replay"))
	}
	var loadHeight int64
	if blockStore.Height() == curState.LastBlockHeight {
		sm.SaveState(stateDB, preState)
		loadHeight = preState.LastBlockHeight
	} else if blockStore.Height() == curState.LastBlockHeight+1 {
		loadHeight = curState.LastBlockHeight
	} else {
		panic(errors.New("tendermint block store height should be at most one ahead of the its state height"))
	}

	return loadHeight
}
