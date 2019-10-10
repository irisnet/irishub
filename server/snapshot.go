package server

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/tendermint/tendermint/consensus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	bc "github.com/tendermint/tendermint/blockchain"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

const flagTmpDir = "tmp-dir"
const pathSeparator = string(os.PathSeparator)

// SnapshotCmd delete historical block data and index data
func SnapshotCmd(ctx *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "snapshot current block information and drop other block",
		RunE: func(cmd *cobra.Command, args []string) error {
			home := viper.GetString("home")
			emptyState, err := isEmptyState(home)
			if err != nil || emptyState {
				fmt.Println("WARNING: State is not initialized.")
				return nil
			}
			srcDir := filepath.Join(home, "data")

			targetDir := viper.GetString(flagTmpDir)
			if len(targetDir) == 0 {
				targetDir = filepath.Join(home, "data.bak")
			}

			if err = dumpData(srcDir, targetDir); err != nil {
				os.RemoveAll(targetDir)
				fmt.Println(fmt.Sprintf("FAILED: %s", err.Error()))
				return err
			}
			fmt.Println(fmt.Sprintf("snapshot file is stored in [%s]", targetDir))
			return nil
		},
	}
	cmd.Flags().String(flagTmpDir, "", "snapshot file storage directory")
	return cmd
}

func loadDb(name, path string) *dbm.GoLevelDB {
	db, err := dbm.NewGoLevelDB(name, path)
	if err != nil {
		panic("load db failed")
	}
	return db
}

func dumpData(home, targetDir string) error {
	//save last block and flush disk
	lastHeight := snapshotBlock(home, targetDir)
	if err := snapshotCsWAL(home, targetDir, lastHeight); err != nil {
		return err
	}

	//copy application
	appDir := filepath.Join(home, "application.db")
	appTargetDir := filepath.Join(targetDir, "application.db")
	if err := copyDir(appDir, appTargetDir); err != nil {
		return err
	}

	//copy state
	stateDir := filepath.Join(home, "state.db")
	targetStateDir := filepath.Join(targetDir, "state.db")
	if err := copyDir(stateDir, targetStateDir); err != nil {
		return err
	}

	//copy evidence.db
	evidenceDir := filepath.Join(home, "evidence.db")
	evidenceTargetDir := filepath.Join(targetDir, "evidence.db")
	return copyDir(evidenceDir, evidenceTargetDir)
}

func snapshotBlock(home, targetDir string) int64 {
	originDb := loadDb("blockstore", home)
	defer originDb.Close()

	originStore := bc.NewBlockStore(originDb)
	height := originStore.Height()

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
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("cannot replay height %d. WAL does not contain #ENDHEIGHT for %d", height, height-1)
	}
	defer gr.Close() // nolint: errcheck

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
		targetWAL.Write(msg.Msg)
	}
	targetWAL.WriteSync(consensus.EndHeightMessage{height})
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
					fmt.Println(err)
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
