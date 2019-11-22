package app

import (
	"fmt"

	"github.com/irisnet/irishub/server"
	"github.com/spf13/viper"
	bc "github.com/tendermint/tendermint/blockchain"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	sm "github.com/tendermint/tendermint/state"
	dbm "github.com/tendermint/tm-db"
)

func Replay(logger log.Logger) int64 {
	ctx := server.NewDefaultContext()
	ctx.Config.RootDir = viper.GetString(tmcli.HomeFlag)
	dbContext := node.DBContext{"state", ctx.Config}
	dbType := dbm.DBBackendType(dbContext.Config.DBBackend)
	stateDB := dbm.NewDB(dbContext.ID, dbType, dbContext.Config.DBDir())
	defer stateDB.Close()

	blockDBContext := node.DBContext{"blockstore", ctx.Config}
	blockStoreDB := dbm.NewDB(blockDBContext.ID, dbType, dbContext.Config.DBDir())
	defer blockStoreDB.Close()
	blockStore := bc.NewBlockStore(blockStoreDB)

	curState := sm.LoadState(stateDB)
	preState := sm.LoadPreState(stateDB)
	if curState.LastBlockHeight == preState.LastBlockHeight && preState.LastBlockHeight == 0 {
		panic(fmt.Errorf("there is no block now, can't replay"))
	}
	var loadHeight int64
	if blockStore.Height() == curState.LastBlockHeight {
		logger.Info(fmt.Sprintf("Blockstore height equals to current state height %d", curState.LastBlockHeight))
		logger.Info("Just reset state DB to last height")
		sm.SaveState(stateDB, preState)
		loadHeight = preState.LastBlockHeight
	} else if blockStore.Height() == curState.LastBlockHeight+1 {
		logger.Info(fmt.Sprintf("Blockstore height %d, current state height %d", blockStore.Height(), curState.LastBlockHeight))
		logger.Info(fmt.Sprintf("Retreat block %d in block store and reset state DB to last height", blockStore.Height()))
		blockStore.RetreatLastBlock()
		sm.SaveState(stateDB, preState)
		loadHeight = preState.LastBlockHeight
	} else if blockStore.Height() == curState.LastBlockHeight+2 && curState.LastBlockHeight == preState.LastBlockHeight {
		logger.Info(fmt.Sprintf("Blockstore height %d, current state height %d, pre-state height %d", blockStore.Height(), curState.LastBlockHeight, preState.LastBlockHeight))
		logger.Info("State store has already been retreat to last height, only retreat block store to last height")
		blockStore.RetreatLastBlock()
		loadHeight = preState.LastBlockHeight
	} else {
		panic(fmt.Errorf("unexpected situation: tendermint block store height %d, tendermint state store height %d", blockStore.Height(), curState.LastBlockHeight))
	}

	return loadHeight
}
