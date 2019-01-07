package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v0"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/irisnet/irishub/app/v1"
)

const (
	appName          = "IrisApp"
	FlagReplay       = "replay"
	FlagReplayHeight = "replay_height"
	//Keep snapshot every at syncable height
	DefaultSyncableHeight = 10000
)

// default home directories for expected binaries
var (
	DefaultLCDHome  = os.ExpandEnv("$HOME/.irislcd")
	DefaultCLIHome  = os.ExpandEnv("$HOME/.iriscli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.iris")
)

// Extended ABCI application
type IrisApp struct {
	*BaseApp
}

func NewIrisApp(logger log.Logger, db dbm.DB, traceStore io.Writer, baseAppOptions ...func(*BaseApp)) *IrisApp {
	cdc := MakeCodec()

	bApp := NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)

	// create your application object
	var app = &IrisApp{
		BaseApp: bApp,
	}
	engine := protocol.NewProtocolEngine(cdc)
	app.SetProtocolEngine(&engine)
	app.MountStoresIAVL(engine.GetKVStoreKeys())
	app.MountStoresTransient(engine.GetTransientStoreKeys())

	var err error
	if viper.GetInt64(FlagReplayHeight) > 0 {
		replayHeight := viper.GetInt64(FlagReplayHeight)
		loadHeight := app.replayToHeight(replayHeight, app.Logger)
		app.Logger.Info(fmt.Sprintf("Load store at %d, start to replay to %d", loadHeight, replayHeight))
		err = app.LoadVersion(loadHeight, protocol.KeyMain, true)
	} else {
		err = app.LoadLatestVersion(engine.GetKeyMain())
	}
	if err != nil {
		cmn.Exit(err.Error())
	}

	protocol0 := v0.NewProtocolVersion0(cdc)
	engine.Add(protocol0)
	protocol1 := v1.NewProtocolVersion1(cdc)
	engine.Add(protocol1)

	engine.LoadCurrentProtocol(app.GetKVStore(protocol.KeyProtocol))
	bApp.txDecoder = auth.DefaultTxDecoder(engine.GetCurrentProtocol().GetCodec())

	return app
}

// custom tx codec
func MakeCodec() *codec.Codec {
	return v1.MakeCodec()
}

// export the state of iris for a genesis file
func (app *IrisApp) ExportAppStateAndValidators(forZeroHeight bool) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
	return app.Engine.GetCurrentProtocol().ExportAppStateAndValidators(ctx, forZeroHeight)
}

func (app *IrisApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, protocol.KeyMain, false)
}

func (app *IrisApp) replayToHeight(replayHeight int64, logger log.Logger) int64 {
	loadHeight := int64(0)
	logger.Info("Please make sure the replay height is less than block height")
	if replayHeight >= DefaultSyncableHeight {
		loadHeight = replayHeight - replayHeight%DefaultSyncableHeight
	} else {
		// version 1 will always be kept
		loadHeight = 1
	}
	logger.Info("This replay operation will change the application store, please spare your node home directory first")
	logger.Info("Confirm that:(y/n)")
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		cmn.Exit(err.Error())
	}
	confirm := strings.TrimSpace(input)
	if confirm != "y" && confirm != "yes" {
		cmn.Exit("Abort replay operation")
	}
	return loadHeight
}
