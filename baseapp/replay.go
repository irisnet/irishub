package baseapp

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/irisnet/irishub/server"
	"github.com/spf13/viper"
	bc "github.com/tendermint/tendermint/blockchain"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	sm "github.com/tendermint/tendermint/state"
)

const (
	DefaultSyncableHeight = 10000
)

func Replay(logger log.Logger) int64 {
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
		panic(fmt.Errorf("there is no block now, can't replay"))
	}
	var loadHeight int64
	if blockStore.Height() == curState.LastBlockHeight {
		logger.Info(fmt.Sprintf("blockstore height equals to current state height %d", curState.LastBlockHeight))
		logger.Info("Just reset state DB to last height")
		sm.SaveState(stateDB, preState)
		loadHeight = preState.LastBlockHeight
	} else if blockStore.Height() == curState.LastBlockHeight+1 {
		logger.Info(fmt.Sprintf("blockstore height %d, current state height %d", blockStore.Height(), curState.LastBlockHeight))
		logger.Info(fmt.Sprintf("Retreat block %d in block store and reset state DB to last height", blockStore.Height()))
		blockStore.RetreatLastBlock()
		sm.SaveState(stateDB, preState)
		loadHeight = preState.LastBlockHeight
	} else {
		panic(fmt.Errorf("tendermint block store height should be at most one ahead of the its state height"))
	}

	return loadHeight
}

func ReplayToHeight(replayHeight int64, logger log.Logger) int64 {
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
