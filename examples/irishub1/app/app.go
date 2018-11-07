package app

import (
	"encoding/json"
	"io"
	"os"

	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/stake"
	bam "github.com/irisnet/irishub/baseapp"
	ibc1 "github.com/irisnet/irishub/examples/irishub1/ibc"
	"github.com/irisnet/irishub/modules/gov"
	"github.com/irisnet/irishub/modules/gov/params"
	"github.com/irisnet/irishub/iparam"
	"github.com/irisnet/irishub/modules/record"
	"github.com/irisnet/irishub/modules/upgrade"
	"github.com/irisnet/irishub/modules/upgrade/params"
	"github.com/irisnet/irishub/modules/iservice"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	bc "github.com/tendermint/tendermint/blockchain"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	sm "github.com/tendermint/tendermint/state"
	tmtypes "github.com/tendermint/tendermint/types"
	"strings"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"sort"
	"github.com/irisnet/irishub/modules/iservice/params"
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
	cdc *codec.Codec

	// keys to access the substores
	keyMain          *sdk.KVStoreKey
	keyAccount       *sdk.KVStoreKey
	keyIBC           *sdk.KVStoreKey
	keyStake         *sdk.KVStoreKey
	tkeyStake        *sdk.TransientStoreKey
	keySlashing      *sdk.KVStoreKey
	keyMint          *sdk.KVStoreKey
	keyDistr         *sdk.KVStoreKey
	tkeyDistr        *sdk.TransientStoreKey
	keyGov           *sdk.KVStoreKey
	keyFeeCollection *sdk.KVStoreKey
	keyParams        *sdk.KVStoreKey
	tkeyParams       *sdk.TransientStoreKey
	keyUpgrade       *sdk.KVStoreKey
	keyIservice      *sdk.KVStoreKey
	keyRecord        *sdk.KVStoreKey

	// Manage getting and setting accounts
	accountMapper       auth.AccountKeeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	bankKeeper          bank.Keeper
	ibcMapper           ibc.Mapper
	ibc1Mapper          ibc1.Mapper
	stakeKeeper         stake.Keeper
	slashingKeeper      slashing.Keeper
	mintKeeper          mint.Keeper
	distrKeeper         distr.Keeper
	govKeeper           gov.Keeper
	paramsKeeper        params.Keeper
	upgradeKeeper       upgrade.Keeper
	iserviceKeeper      iservice.Keeper
	recordKeeper        record.Keeper

	// fee manager
	feeManager bam.FeeManager
}

func NewIrisApp(logger log.Logger, db dbm.DB, traceStore io.Writer, baseAppOptions ...func(*bam.BaseApp)) *IrisApp {
	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)

	// create your application object
	var app = &IrisApp{
		BaseApp:          bApp,
		cdc:              cdc,
		keyMain:          sdk.NewKVStoreKey("main"),
		keyAccount:       sdk.NewKVStoreKey("acc"),
		keyIBC:           sdk.NewKVStoreKey("ibc"),
		keyStake:         sdk.NewKVStoreKey("stake"),
		tkeyStake:        sdk.NewTransientStoreKey("transient_stake"),
		keyMint:          sdk.NewKVStoreKey("mint"),
		keyDistr:         sdk.NewKVStoreKey("distr"),
		tkeyDistr:        sdk.NewTransientStoreKey("transient_distr"),
		keySlashing:      sdk.NewKVStoreKey("slashing"),
		keyGov:           sdk.NewKVStoreKey("gov"),
		keyRecord:        sdk.NewKVStoreKey("record"),
		keyFeeCollection: sdk.NewKVStoreKey("fee"),
		keyParams:        sdk.NewKVStoreKey("params"),
		tkeyParams:       sdk.NewTransientStoreKey("transient_params"),
		keyUpgrade:       sdk.NewKVStoreKey("upgrade"),
		keyIservice:      sdk.NewKVStoreKey("iservice"),
	}

	var lastHeight int64
	if viper.GetBool(FlagReplay) {
		lastHeight = app.replay()
	}

	// define the AccountKeeper
	app.accountMapper = auth.NewAccountKeeper(
		app.cdc,
		app.keyAccount,        // target store
		auth.ProtoBaseAccount, // prototype
	)

	// add handlers
	app.bankKeeper = bank.NewBaseKeeper(app.accountMapper)
	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(
		app.cdc,
		app.keyFeeCollection,
	)
	app.paramsKeeper = params.NewKeeper(
		app.cdc,
		app.keyParams, app.tkeyParams,
	)
	app.ibcMapper = ibc.NewMapper(
		app.cdc,
		app.keyIBC, app.RegisterCodespace(ibc.DefaultCodespace),
	)
	app.ibc1Mapper = ibc1.NewMapper(
		app.cdc,
		app.keyIBC, app.RegisterCodespace(ibc1.DefaultCodespace),
	)
	app.stakeKeeper = stake.NewKeeper(
		app.cdc,
		app.keyStake, app.tkeyStake,
		app.bankKeeper, app.paramsKeeper.Subspace(stake.DefaultParamspace),
		app.RegisterCodespace(stake.DefaultCodespace),
	)
	app.mintKeeper = mint.NewKeeper(app.cdc, app.keyMint,
		app.paramsKeeper.Subspace(mint.DefaultParamspace),
		app.stakeKeeper, app.feeCollectionKeeper,
	)
	app.distrKeeper = distr.NewKeeper(
		app.cdc,
		app.keyDistr,
		app.paramsKeeper.Subspace(distr.DefaultParamspace),
		app.bankKeeper, app.stakeKeeper, app.feeCollectionKeeper,
		app.RegisterCodespace(stake.DefaultCodespace),
	)
	app.slashingKeeper = slashing.NewKeeper(
		app.cdc,
		app.keySlashing,
		app.stakeKeeper, app.paramsKeeper.Subspace(slashing.DefaultParamspace),
		app.RegisterCodespace(slashing.DefaultCodespace),
	)
	app.upgradeKeeper = upgrade.NewKeeper(
		app.cdc,
		app.keyUpgrade, app.stakeKeeper,
	)

	app.govKeeper = gov.NewKeeper(
		app.cdc,
		app.keyGov,
		app.bankKeeper, app.stakeKeeper,
		app.RegisterCodespace(gov.DefaultCodespace),
	)

	app.recordKeeper = record.NewKeeper(
		app.cdc,
		app.keyRecord,
		app.RegisterCodespace(record.DefaultCodespace),
	)
	app.iserviceKeeper = iservice.NewKeeper(
		app.cdc,
		app.keyIservice,
		app.bankKeeper,
		app.RegisterCodespace(iservice.DefaultCodespace),
	)

	// register the staking hooks
	app.stakeKeeper = app.stakeKeeper.WithHooks(
		NewHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()))

	// register message routes
	// need to update each module's msg type
	app.Router().
		AddRoute("bank", []*sdk.KVStoreKey{app.keyAccount}, bank.NewHandler(app.bankKeeper)).
		AddRoute("ibc", []*sdk.KVStoreKey{app.keyIBC, app.keyAccount}, ibc.NewHandler(app.ibcMapper, app.bankKeeper)).
		AddRoute("ibc-1", []*sdk.KVStoreKey{app.keyIBC, app.keyAccount}, ibc1.NewHandler(app.ibc1Mapper, app.bankKeeper)).
		AddRoute("stake", []*sdk.KVStoreKey{app.keyStake, app.keyAccount}, stake.NewHandler(app.stakeKeeper)).
		AddRoute("slashing", []*sdk.KVStoreKey{app.keySlashing, app.keyStake}, slashing.NewHandler(app.slashingKeeper)).
		AddRoute("distr",  []*sdk.KVStoreKey{app.keyDistr}, distr.NewHandler(app.distrKeeper)).
		AddRoute("gov", []*sdk.KVStoreKey{app.keyGov, app.keyAccount, app.keyStake, app.keyParams}, gov.NewHandler(app.govKeeper)).
		AddRoute("upgrade", []*sdk.KVStoreKey{app.keyUpgrade, app.keyStake}, upgrade.NewHandler(app.upgradeKeeper)).
		AddRoute("record", []*sdk.KVStoreKey{app.keyRecord}, record.NewHandler(app.recordKeeper)).
		AddRoute("iservice", []*sdk.KVStoreKey{app.keyIservice}, iservice.NewHandler(app.iserviceKeeper))

	app.QueryRouter().
		AddRoute("gov", gov.NewQuerier(app.govKeeper)).
		AddRoute("stake", stake.NewQuerier(app.stakeKeeper, app.cdc))


	app.feeManager = bam.NewFeeManager(app.paramsKeeper.Subspace("Fee"))

	// initialize BaseApp
	app.MountStoresIAVL(app.keyMain, app.keyAccount, app.keyIBC, app.keyStake, app.keySlashing, app.keyGov, app.keyMint, app.keyDistr,
		app.keyFeeCollection, app.keyParams, app.keyUpgrade, app.keyRecord, app.keyIservice)
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper, app.feeCollectionKeeper))
	app.MountStoresTransient(app.tkeyParams, app.tkeyStake, app.tkeyDistr)
	app.SetFeeRefundHandler(bam.NewFeeRefundHandler(app.accountMapper, app.feeCollectionKeeper, app.feeManager))
	app.SetFeePreprocessHandler(bam.NewFeePreprocessHandler(app.feeManager))
	app.SetEndBlocker(app.EndBlocker)
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
	app.upgradeKeeper.RefreshVersionList(app.GetKVStore(app.keyUpgrade))

	iparam.SetParamReadWriter(app.paramsKeeper.Subspace(iparam.SignalParamspace).WithTypeTable(
		params.NewTypeTable(
			upgradeparams.CurrentUpgradeProposalIdParameter.GetStoreKey(), int64((0)),
			upgradeparams.ProposalAcceptHeightParameter.GetStoreKey(), int64(0),
			upgradeparams.SwitchPeriodParameter.GetStoreKey(), int64(0),
		)),
		&upgradeparams.CurrentUpgradeProposalIdParameter,
		&upgradeparams.ProposalAcceptHeightParameter,
		&upgradeparams.SwitchPeriodParameter)

	iparam.SetParamReadWriter(app.paramsKeeper.Subspace(iparam.GovParamspace).WithTypeTable(
		params.NewTypeTable(
			govparams.DepositProcedureParameter.GetStoreKey(), govparams.DepositProcedure{},
			govparams.VotingProcedureParameter.GetStoreKey(), govparams.VotingProcedure{},
			govparams.TallyingProcedureParameter.GetStoreKey(), govparams.TallyingProcedure{},
		)),
		&govparams.DepositProcedureParameter,
		&govparams.VotingProcedureParameter,
		&govparams.TallyingProcedureParameter)

	iparam.RegisterGovParamMapping(
		&govparams.DepositProcedureParameter,
		&govparams.VotingProcedureParameter,
		&govparams.TallyingProcedureParameter)

	iparam.SetParamReadWriter(app.paramsKeeper.Subspace(iparam.ServiceParamspace).WithTypeTable(
		params.NewTypeTable(
			iserviceparams.MaxRequestTimeoutParameter.GetStoreKey(), int64(0),
			iserviceparams.MinProviderDepositParameter.GetStoreKey(), sdk.Coins{},
		)),
		&iserviceparams.MaxRequestTimeoutParameter,
		&iserviceparams.MinProviderDepositParameter)

	return app
}

// custom tx codec
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	ibc.RegisterCodec(cdc)
	ibc1.RegisterCodec(cdc)

	bank.RegisterCodec(cdc)
	stake.RegisterCodec(cdc)
	distr.RegisterCodec(cdc)
	slashing.RegisterCodec(cdc)
	gov.RegisterCodec(cdc)
	record.RegisterCodec(cdc)
	upgrade.RegisterCodec(cdc)
	iservice.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

// application updates every end block
func (app *IrisApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	tags := slashing.BeginBlocker(ctx, req, app.slashingKeeper)

	// distribute rewards from previous block
	distr.BeginBlocker(ctx, req, app.distrKeeper)

	// mint new tokens for this new block
	mint.BeginBlocker(ctx, app.mintKeeper)

	return abci.ResponseBeginBlock{
		Tags: tags.ToKVPairs(),
	}
}

// application updates every end block
func (app *IrisApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	tags := gov.EndBlocker(ctx, app.govKeeper)
	validatorUpdates := stake.EndBlocker(ctx, app.stakeKeeper)
	tags.AppendTags(upgrade.EndBlocker(ctx, app.upgradeKeeper))
	// Add these new validators to the addr -> pubkey map.
	app.slashingKeeper.AddValidators(ctx, validatorUpdates)
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
	gov.InitGenesis(ctx, app.govKeeper, genesisState.GovData)

	feeTokenGensisConfig := bam.FeeGenesisStateConfig{
		FeeTokenNative:    IrisCt.MinUnit.Denom,
		GasPriceThreshold: 20000000000, // 20(glue), 20*10^9, 1 glue = 10^9 lue/gas, 1 iris = 10^18 lue
	}

	bam.InitGenesis(ctx, app.feeManager, feeTokenGensisConfig)

	// load the address to pubkey map
	slashing.InitGenesis(ctx, app.slashingKeeper, genesisState.SlashingData, genesisState.StakeData)
	mint.InitGenesis(ctx, app.mintKeeper, genesisState.MintData)
	distr.InitGenesis(ctx, app.distrKeeper, genesisState.DistrData)
	err = IrisValidateGenesisState(genesisState)
	if err != nil {
		panic(err) // TODO find a way to do this w/o panics
	}

	if len(genesisState.GenTxs) > 0 {
		for _, genTx := range genesisState.GenTxs {
			var tx auth.StdTx
			err = app.cdc.UnmarshalJSON(genTx, &tx)
			if err != nil {
				panic(err)
			}
			bz := app.cdc.MustMarshalBinary(tx)
			res := app.BaseApp.DeliverTx(bz)
			if !res.IsOK() {
				panic(res.Log)
			}
		}

		validators = app.stakeKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	}
	app.slashingKeeper.AddValidators(ctx, validators)

	// sanity check
	if len(req.Validators) > 0 {
		if len(req.Validators) != len(validators) {
			panic(fmt.Errorf("len(RequestInitChain.Validators) != len(validators) (%d != %d) ", len(req.Validators), len(validators)))
		}
		sort.Sort(abci.ValidatorUpdates(req.Validators))
		sort.Sort(abci.ValidatorUpdates(validators))
		for i, val := range validators {
			if !val.Equal(req.Validators[i]) {
				panic(fmt.Errorf("validators[%d] != req.Validators[%d] ", i, i))
			}
		}
	}

	upgrade.InitGenesis(ctx, app.upgradeKeeper, app.Router(), genesisState.UpgradeData)
	iservice.InitGenesis(ctx, genesisState.IserviceData)

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
	genState := NewGenesisState(
		accounts,
		stake.WriteGenesis(ctx, app.stakeKeeper),
		mint.WriteGenesis(ctx, app.mintKeeper),
		distr.WriteGenesis(ctx, app.distrKeeper),
		gov.WriteGenesis(ctx, app.govKeeper),
		upgrade.WriteGenesis(ctx, app.upgradeKeeper),
		slashing.GenesisState{}, // TODO create write methods
	)
	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}
	validators = stake.WriteValidators(ctx, app.stakeKeeper)
	return appState, validators, nil
}

// Iterates through msgs and executes them
func (app *IrisApp) runMsgs(ctx sdk.Context, msgs []sdk.Msg, mode bam.RunTxMode) (result sdk.Result) {
	// accumulate results
	logs := make([]string, 0, len(msgs))
	var data []byte   // NOTE: we just append them all (?!)
	var tags sdk.Tags // also just append them all
	var code sdk.ABCICodeType
	for msgIdx, msg := range msgs {
		// Match route.
		var msgType string
		var err sdk.Error
		if ctx.BlockHeight() != 0 {
			msgType, err = app.upgradeKeeper.GetMsgTypeInCurrentVersion(ctx, msg)

			if err != nil {
				return err.Result()
			}

		} else {
			msgType = msg.Route()
		}

		handler := app.Router().Route(msgType)
		if handler == nil {
			return sdk.ErrUnknownRequest("Unrecognized Msg type: " + msgType).Result()
		}

		var msgResult sdk.Result
		if mode != bam.RunTxModeCheck {
			msgResult = handler(ctx, msg)
		}

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
	}()

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

//______________________________________________________________________________________________

// Combined Staking Hooks
type Hooks struct {
	dh distr.Hooks
	sh slashing.Hooks
}

func NewHooks(dh distr.Hooks, sh slashing.Hooks) Hooks {
	return Hooks{dh, sh}
}

var _ sdk.StakingHooks = Hooks{}

// nolint
func (h Hooks) OnValidatorCreated(ctx sdk.Context, addr sdk.ValAddress) {
	h.dh.OnValidatorCreated(ctx, addr)
}
func (h Hooks) OnValidatorModified(ctx sdk.Context, addr sdk.ValAddress) {
	h.dh.OnValidatorModified(ctx, addr)
}
func (h Hooks) OnValidatorRemoved(ctx sdk.Context, addr sdk.ValAddress) {
	h.dh.OnValidatorRemoved(ctx, addr)
}
func (h Hooks) OnValidatorBonded(ctx sdk.Context, addr sdk.ConsAddress, operator sdk.ValAddress) {
	h.dh.OnValidatorBonded(ctx, addr, operator)
	h.sh.OnValidatorBonded(ctx, addr, operator)
}
func (h Hooks) OnValidatorPowerDidChange(ctx sdk.Context, addr sdk.ConsAddress, operator sdk.ValAddress) {
	h.dh.OnValidatorPowerDidChange(ctx, addr, operator)
	h.sh.OnValidatorPowerDidChange(ctx, addr, operator)
}
func (h Hooks) OnValidatorBeginUnbonding(ctx sdk.Context, addr sdk.ConsAddress, operator sdk.ValAddress) {
	h.dh.OnValidatorBeginUnbonding(ctx, addr, operator)
	h.sh.OnValidatorBeginUnbonding(ctx, addr, operator)
}
func (h Hooks) OnDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.dh.OnDelegationCreated(ctx, delAddr, valAddr)
}
func (h Hooks) OnDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.dh.OnDelegationSharesModified(ctx, delAddr, valAddr)
}
func (h Hooks) OnDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.dh.OnDelegationRemoved(ctx, delAddr, valAddr)
}
