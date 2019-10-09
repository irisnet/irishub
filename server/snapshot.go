package server

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	bc "github.com/tendermint/tendermint/blockchain"
	dbm "github.com/tendermint/tm-db"
)

const flagTmpDir = "tmp-dir"

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
	snapshotBlock(home, targetDir)

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
	partSet := block.MakePartSet(65536)
	targetStore.SaveBlock(block, partSet, seenCommit)
	return height
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
		path = strings.Replace(path, "\\", string(os.PathSeparator), -1)
		destNewPath := strings.Replace(path, srcPath, destPath, -1)
		_, err = copyFile(path, destNewPath)
		return err
	})
}

func copyFile(src, dest string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()
	destSplitPathDirs := strings.Split(dest, string(os.PathSeparator))

	destSplitPath := ""
	for index, dir := range destSplitPathDirs {
		if index < len(destSplitPathDirs)-1 {
			destSplitPath = destSplitPath + dir + string(os.PathSeparator)
			if b, _ := pathExists(destSplitPath); b == false {
				err := os.Mkdir(destSplitPath, os.ModePerm)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	dstFile, err := os.Create(dest)
	if err != nil {
		return
	}

	defer dstFile.Close()
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
