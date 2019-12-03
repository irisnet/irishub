package server

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/irisnet/irishub/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	bc "github.com/tendermint/tendermint/blockchain"
	"github.com/tendermint/tendermint/consensus"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	tmsm "github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

const flagTmpDir = "tmp-dir"
const pathSeparator = string(os.PathSeparator)

// SnapshotCmd delete historical block data and index data
func SnapshotCmd(ctx *Context, cdc *codec.Codec, appReset AppReset) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "snapshot the latest information and drop the others",
		RunE: func(cmd *cobra.Command, args []string) error {
			home := viper.GetString(tmcli.HomeFlag)
			emptyState, err := isEmptyState(home)
			if err != nil || emptyState {
				ctx.Logger.Error("State is not initialized")
				return nil
			}

			targetDir := viper.GetString(flagTmpDir)
			if len(targetDir) == 0 {
				targetDir = filepath.Join(home, "data.bak")
			}
			dataDir := filepath.Join(home, "data")

			if err = snapshot(ctx, cdc, dataDir, targetDir, appReset); err != nil {
				_ = os.RemoveAll(targetDir)
				ctx.Logger.Error("snapshot file is created failed")
				return err
			}
			ctx.Logger.Info("snapshot file is created successful", "location", targetDir)
			return nil
		},
	}
	cmd.Flags().String(flagTmpDir, "", "Snapshot file storage directory")
	// core flags for the ABCI application
	cmd.Flags().String(flagTraceStore, "", "Enable KVStore tracing to an output file")
	cmd.Flags().String(flagPruning, "syncable", "Pruning strategy: syncable, nothing, everything")
	cmd.Flags().Bool(flagCheckInvariant, false, "Enable invariant check on mainnet, ignore this flag on testnet")
	return cmd
}

func loadDb(name, path string) *dbm.GoLevelDB {
	db, err := dbm.NewGoLevelDB(name, path)
	if err != nil {
		panic("load db failed")
	}
	return db
}

func snapshot(ctx *Context, cdc *codec.Codec, dataDir, targetDir string, appReset AppReset) error {
	blockDB := loadDb("blockstore", dataDir)
	blockStore := bc.NewBlockStore(blockDB)

	stateDB := loadDb("state", dataDir)
	state := tmsm.LoadState(stateDB)

	defer func() {
		blockDB.Close()
		stateDB.Close()
	}()
	if blockStore.Height() != state.LastBlockHeight {
		if err := reset(ctx, appReset, state.LastBlockHeight); err != nil {
			return err
		}
	}

	//save local current block and flush disk
	snapshotBlock(blockStore, targetDir, state.LastBlockHeight)
	//save local current block height state
	snapshotState(cdc, stateDB, targetDir)
	//save local current block height consensus data
	_ = snapshotCsWAL(dataDir, targetDir, state.LastBlockHeight)

	//copy application
	appDir := filepath.Join(dataDir, "application.db")
	appTargetDir := filepath.Join(targetDir, "application.db")
	if err := copyDir(appDir, appTargetDir); err != nil {
		return err
	}

	//copy evidence.db
	evidenceDir := filepath.Join(dataDir, "evidence.db")
	evidenceTargetDir := filepath.Join(targetDir, "evidence.db")
	return copyDir(evidenceDir, evidenceTargetDir)
}

func snapshotState(cdc *codec.Codec, tmDB *dbm.GoLevelDB, targetDir string) {
	targetDb := loadDb("state", targetDir)
	defer targetDb.Close()

	state := tmsm.LoadState(tmDB)

	saveValidatorsInfo(cdc, tmDB, targetDb, state.LastBlockHeight)
	saveConsensusParamsInfo(cdc, tmDB, targetDb, state.LastBlockHeight)
	tmsm.SaveState(targetDb, state)
}

func snapshotBlock(originStore *bc.BlockStore, targetDir string, height int64) int64 {
	targetDb := loadDb("blockstore", targetDir)
	defer targetDb.Close()

	bsj := bc.BlockStoreStateJSON{Height: height - 1}
	bsj.Save(targetDb)
	targetStore := bc.NewBlockStore(targetDb)

	block := originStore.LoadBlock(height)
	seenCommit := originStore.LoadSeenCommit(height)
	partSet := block.MakePartSet(types.BlockPartSizeBytes)
	targetStore.SaveBlock(block, partSet, seenCommit)
	return height
}

func snapshotCsWAL(home, targetDir string, height int64) error {
	walTargetDir := filepath.Join(targetDir, "cs.wal", "wal")
	targetWAL, err := consensus.NewWAL(walTargetDir)

	walSourceDir := filepath.Join(home, "cs.wal", "wal")
	sourceWAL, err := consensus.NewWAL(walSourceDir)
	if err != nil {
		return errors.New("failed to open WAL for consensus state")
	}

	gr, found, err := sourceWAL.SearchForEndHeight(height, &consensus.WALSearchOptions{IgnoreDataCorruptionErrors: true})

	if err != nil || !found {
		return fmt.Errorf("cannot replay height %d. WAL does not contain #ENDHEIGHT for %d", height, height-1)
	}

	defer func() {
		if err = gr.Close(); err != nil {
			return
		}
	}()

	var msg *consensus.TimedWALMessage
	dec := consensus.NewWALDecoder(gr)
	for {
		msg, err = dec.Decode()
		if err == io.EOF {
			break
		} else if consensus.IsDataCorruptionError(err) {
			return fmt.Errorf("data has been corrupted in last height %d of consensus WAL", height)
		} else if err != nil {
			return err
		}
		if err := targetWAL.Write(msg.Msg); err != nil {
			return err
		}
	}
	_ = targetWAL.WriteSync(consensus.EndHeightMessage{Height: height})
	return nil
}

func copyDir(srcPath string, destPath string) error {
	if _, err := os.Stat(srcPath); err != nil {
		return err
	}

	return filepath.Walk(srcPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		if f.IsDir() {
			return nil
		}
		path = strings.Replace(path, fmt.Sprintf("\\%s", pathSeparator), pathSeparator, -1)
		destNewPath := strings.Replace(path, srcPath, destPath, -1)
		_, err = copyFile(path, destNewPath)
		return err
	})
}

func copyFile(src, dest string) (w int64, err error) {
	srcFile, err := os.Open(src)
	defer srcFile.Close()
	if err != nil {
		return
	}

	destSplitPathDirs := strings.Split(dest, pathSeparator)

	destSplitPath := ""
	for index, dir := range destSplitPathDirs {
		if index < len(destSplitPathDirs)-1 {
			destSplitPath = destSplitPath + dir + pathSeparator
			if b, _ := pathExists(destSplitPath); b == false {
				err := os.Mkdir(destSplitPath, os.ModePerm)
				if err != nil {
					return 0, err
				}
			}
		}
	}
	dstFile, err := os.Create(dest)
	defer dstFile.Close()
	if err != nil {
		return
	}
	return io.Copy(dstFile, srcFile)
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func loadValidatorsInfo(cdc *codec.Codec, db dbm.DB, height int64) *tmsm.ValidatorsInfo {
	buf := db.Get(calcValidatorsKey(height))
	if len(buf) == 0 {
		return nil
	}

	v := new(tmsm.ValidatorsInfo)
	err := cdc.UnmarshalBinaryBare(buf, v)
	if err != nil {
		return v
	}
	return v
}

func saveValidatorsInfo(cdc *codec.Codec, originDb, targetDb dbm.DB, height int64) {
	valInfo := loadValidatorsInfo(cdc, originDb, height)
	if valInfo.LastHeightChanged > height {
		panic("LastHeightChanged cannot be greater than ValidatorsInfo height")
	}
	if valInfo.ValidatorSet == nil {
		valInfo = loadValidatorsInfo(cdc, originDb, valInfo.LastHeightChanged)
	}
	targetDb.Set(calcValidatorsKey(valInfo.LastHeightChanged), valInfo.Bytes())
}

func loadConsensusParamsInfo(cdc *codec.Codec, db dbm.DB, height int64) *tmsm.ConsensusParamsInfo {
	buf := db.Get(calcConsensusParamsKey(height))
	if len(buf) == 0 {
		return nil
	}

	paramsInfo := new(tmsm.ConsensusParamsInfo)
	err := cdc.UnmarshalBinaryBare(buf, paramsInfo)
	if err != nil {
		return paramsInfo
	}
	return paramsInfo
}

func saveConsensusParamsInfo(cdc *codec.Codec, originDb, targetDb dbm.DB, height int64) {
	consensusParamsInfo := loadConsensusParamsInfo(cdc, originDb, height)
	if consensusParamsInfo.ConsensusParams.Equals(&types.ConsensusParams{}) {
		consensusParamsInfo = loadConsensusParamsInfo(cdc, originDb, consensusParamsInfo.LastHeightChanged)
	}
	paramsInfo := &tmsm.ConsensusParamsInfo{
		LastHeightChanged: consensusParamsInfo.LastHeightChanged,
	}
	targetDb.Set(calcConsensusParamsKey(consensusParamsInfo.LastHeightChanged), paramsInfo.Bytes())
}

func calcValidatorsKey(height int64) []byte {
	return []byte(fmt.Sprintf("validatorsKey:%v", height))
}

func calcConsensusParamsKey(height int64) []byte {
	return []byte(fmt.Sprintf("consensusParamsKey:%v", height))
}

func reset(ctx *Context, appReset AppReset, height int64) error {
	cfg := ctx.Config
	home := cfg.RootDir
	traceWriterFile := viper.GetString(flagTraceStore)

	db, _ := openDB(home)
	traceWriter, _ := openTraceWriter(traceWriterFile)

	err := appReset(ctx, ctx.Logger, db, traceWriter, height)
	if err != nil {
		return err
	}
	return nil
}
