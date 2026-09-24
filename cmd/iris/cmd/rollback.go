package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"cosmossdk.io/store/rootmulti"
	storetypes "cosmossdk.io/store/types"
	cmtDB "github.com/cometbft/cometbft-db"
	cmtcfg "github.com/cometbft/cometbft/config"
	cmtstate "github.com/cometbft/cometbft/state"
	cmtstore "github.com/cometbft/cometbft/store"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"

	"github.com/irisnet/irishub/v4/app"
)

// Keep the SDK's single-step rollback unchanged unless --height is supplied.
func addRollbackHeightFlag(cmd *cobra.Command) {
	legacyRun := cmd.RunE
	cmd.Flags().Int64("height", 0, "target committed height (requires --hard; all later blocks are removed)")
	cmd.Long += `
With --hard --height H, roll both stores back to committed height H and remove
every later local block. H must be positive, retained by the application and
block stores, and no greater than the current consensus state height. The node
must be stopped. Back up the entire node data directory first.
The consensus WAL is moved into a backup directory to avoid replaying old
consensus messages. Validator signing state is NOT reset. Coordinate signer
recovery and peer isolation separately before restarting a rolled-back network.
An interrupted operation must be rerun with the same height before restarting.
Snapshots and transaction indexes are not rolled back; disable stale snapshot
serving and rebuild or archive auxiliary indexes before serving the new history.
The command refuses to remove committed or pending local blocks containing
evidence because the separate evidence database is not reverted. Resolve
evidence state separately.
`
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if !cmd.Flags().Changed("height") {
			return legacyRun(cmd, args)
		}
		if len(args) != 0 {
			return fmt.Errorf("rollback takes no positional arguments")
		}
		hard, err := cmd.Flags().GetBool("hard")
		if err != nil {
			return err
		}
		height, err := cmd.Flags().GetInt64("height")
		if err != nil {
			return err
		}
		if !hard || height <= 0 {
			return fmt.Errorf("--height requires --hard and a positive target height")
		}
		ctx := server.GetServerContextFromCmd(cmd)
		bs, ss, err := openRollbackStores(ctx.Config)
		if err != nil {
			return err
		}
		defer bs.Close()
		defer ss.Close()
		db, err := dbm.NewDB("application", server.GetAppDBBackend(ctx.Viper), filepath.Join(ctx.Config.RootDir, "data"))
		if err != nil {
			return err
		}
		// Do not load latest: after an interrupted multistore rollback individual
		// trees may already be at H while the root metadata still points above H.
		irisApp := app.NewIrisApp(ctx.Logger, db, nil, false, ctx.Viper)
		defer irisApp.Close()
		hash, err := rollbackStoresToHeight(bs, ss, irisApp.CommitMultiStore(), db, height, func() error {
			return archiveRollbackWAL(ctx.Config.Consensus.WalFile(), height, cmd.OutOrStdout())
		})
		if err != nil {
			return fmt.Errorf("rollback to %d failed; keep the node stopped and retry the same target: %w", height, err)
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "Rolled back application, consensus state and block store to height %d and hash %X\n", height, hash)
		return err
	}
}

func openRollbackStores(cfg *cmtcfg.Config) (*cmtstore.BlockStore, cmtstate.Store, error) {
	// Do not silently create empty databases when --home is wrong.
	for _, name := range []string{"blockstore.db", "state.db"} {
		if _, err := os.Stat(filepath.Join(cfg.DBDir(), name)); err != nil {
			return nil, nil, err
		}
	}
	blockDB, err := cmtDB.NewDB("blockstore", cmtDB.BackendType(cfg.DBBackend), cfg.DBDir())
	if err != nil {
		return nil, nil, err
	}
	stateDB, err := cmtDB.NewDB("state", cmtDB.BackendType(cfg.DBBackend), cfg.DBDir())
	if err != nil {
		_ = blockDB.Close()
		return nil, nil, err
	}
	return cmtstore.NewBlockStore(blockDB), cmtstate.NewStore(stateDB, cmtstate.StoreOptions{
		DiscardABCIResponses: cfg.Storage.DiscardABCIResponses,
	}), nil
}

func rollbackStoresToHeight(bs *cmtstore.BlockStore, ss cmtstate.Store, ms storetypes.CommitMultiStore, appDB dbm.DB, target int64, beforeWrite func() error) ([]byte, error) {
	state, err := ss.Load()
	if err != nil {
		return nil, err
	}
	if state.IsEmpty() || target <= 0 || target < state.InitialHeight || target > state.LastBlockHeight {
		return nil, fmt.Errorf("invalid target %d for consensus height %d (initial height %d)", target, state.LastBlockHeight, state.InitialHeight)
	}
	if bs.Height() != state.LastBlockHeight && bs.Height() != state.LastBlockHeight+1 {
		return nil, fmt.Errorf("block height %d must equal consensus height %d or its pending successor", bs.Height(), state.LastBlockHeight)
	}
	appVersion := rootmulti.GetLatestVersion(appDB)
	if target < bs.Base() || target > appVersion {
		return nil, fmt.Errorf("target %d unavailable: block base %d, application latest version %d", target, bs.Base(), appVersion)
	}
	// Preflight every input the native one-step rollback will read, before
	// deleting any block. A pruned intermediate validator/parameter record must
	// not leave a half-completed rollback merely because it was discovered late.
	for h := state.LastBlockHeight; h >= target; h-- {
		if bs.LoadBlockMeta(h) == nil {
			return nil, fmt.Errorf("block metadata at height %d is unavailable", h)
		}
		if h > target {
			if _, err := ss.LoadValidators(h - 1); err != nil {
				return nil, fmt.Errorf("validators at height %d: %w", h-1, err)
			}
			if _, err := ss.LoadConsensusParams(h); err != nil {
				return nil, fmt.Errorf("consensus parameters at height %d: %w", h, err)
			}
		}
	}
	if err := ensureNoDiscardedEvidence(bs, target+1, bs.Height()); err != nil {
		return nil, err
	}
	expectedHash := state.AppHash
	if target < state.LastBlockHeight {
		expectedHash = bs.LoadBlockMeta(target + 1).Header.AppHash
	}
	// Loading all mounted trees proves H is actually retained; root commit
	// metadata alone can outlive pruned tree versions. This does not delete any
	// versions. It also permits recovery from a partially completed app rollback.
	if err := ms.LoadVersion(target); err != nil {
		return nil, fmt.Errorf("application version %d unavailable: %w", target, err)
	}
	if err := verifyLoadedAppVersion(ms, target, expectedHash); err != nil {
		return nil, fmt.Errorf("application version %d does not match consensus: %w", target, err)
	}
	if beforeWrite != nil {
		if err := beforeWrite(); err != nil {
			return nil, err
		}
	}
	// CometBFT saves the previous state before deleting its block. If interrupted
	// between those writes, its native pending-block branch safely finishes that
	// deletion on retry. Application rollback comes last so it cannot get ahead
	// of block deletion. If interrupted, LoadVersion(H) above can still reopen H.
	for bs.Height() > target {
		if _, _, err := cmtstate.Rollback(bs, ss, true); err != nil {
			return nil, err
		}
	}
	if err := ms.RollbackToVersion(target); err != nil {
		return nil, fmt.Errorf("application rollback: %w", err)
	}
	finalState, err := ss.Load()
	if err != nil {
		return nil, err
	}
	if finalState.LastBlockHeight != target || bs.Height() != target || rootmulti.GetLatestVersion(appDB) != target ||
		!bytes.Equal(finalState.AppHash, ms.LastCommitID().Hash) {
		return nil, fmt.Errorf("post-rollback stores are inconsistent at target %d", target)
	}
	return finalState.AppHash, nil
}

func ensureNoDiscardedEvidence(bs *cmtstore.BlockStore, first, last int64) error {
	for height := first; height <= last; height++ {
		block := bs.LoadBlock(height)
		if block == nil {
			return fmt.Errorf("block at height %d is unavailable for evidence preflight", height)
		}
		if count := len(block.Evidence.Evidence); count != 0 {
			return fmt.Errorf("cannot roll back through block %d containing %d evidence item(s): evidence database state requires separate recovery", height, count)
		}
	}
	return nil
}

func verifyLoadedAppVersion(ms storetypes.CommitMultiStore, target int64, expectedHash []byte) error {
	infoReader, ok := ms.(interface {
		GetCommitInfo(int64) (*storetypes.CommitInfo, error)
	})
	if !ok {
		return fmt.Errorf("multistore does not expose commit info")
	}
	keyLister, ok := ms.(interface {
		StoreKeysByName() map[string]storetypes.StoreKey
	})
	if !ok {
		return fmt.Errorf("multistore does not expose mounted store keys")
	}
	historical, err := infoReader.GetCommitInfo(target)
	if err != nil {
		return fmt.Errorf("load historical commit info: %w", err)
	}
	if historical.Version != target || !bytes.Equal(historical.CommitID().Hash, expectedHash) {
		return fmt.Errorf("historical commit info does not match consensus")
	}

	historicalStores := make(map[string]storetypes.CommitID, len(historical.StoreInfos))
	for _, storeInfo := range historical.StoreInfos {
		if _, exists := historicalStores[storeInfo.Name]; exists {
			return fmt.Errorf("duplicate store %q in historical commit info", storeInfo.Name)
		}
		historicalStores[storeInfo.Name] = storeInfo.CommitId
	}
	actual := storetypes.CommitInfo{Version: target}
	for name, key := range keyLister.StoreKeysByName() {
		storeType := ms.GetStore(key).GetStoreType()
		if storeType == storetypes.StoreTypeTransient || storeType == storetypes.StoreTypeMemory {
			continue
		}
		store := ms.GetCommitKVStore(key)
		expected, exists := historicalStores[name]
		if !exists {
			return fmt.Errorf("mounted persistent store %q is missing from historical commit info", name)
		}
		commitID := store.LastCommitID()
		if commitID.Version != target || expected.Version != target || !bytes.Equal(commitID.Hash, expected.Hash) {
			return fmt.Errorf("loaded store %q does not match historical version %d", name, target)
		}
		actual.StoreInfos = append(actual.StoreInfos, storetypes.StoreInfo{Name: name, CommitId: commitID})
		delete(historicalStores, name)
	}
	if len(historicalStores) != 0 {
		return fmt.Errorf("historical commit info contains stores not mounted by this binary: %v", mapKeys(historicalStores))
	}
	if !bytes.Equal(actual.CommitID().Hash, expectedHash) {
		return fmt.Errorf("loaded stores do not reproduce the historical application hash")
	}
	return nil
}

func mapKeys[V any](values map[string]V) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// Comet's WAL group consists of its head and numeric-suffix rotated files.
// Archive only these exact regular files, retaining unrelated files and signer
// state. On failure the node stays stopped; rerunning archives remaining files.
func archiveRollbackWAL(walPath string, target int64, out io.Writer) error {
	dir, base := filepath.Dir(walPath), filepath.Base(walPath)
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	var paths []string
	for _, entry := range entries {
		name := entry.Name()
		if name != base {
			suffix, ok := strings.CutPrefix(name, base+".")
			if !ok || len(suffix) < 3 || strings.Trim(suffix, "0123456789") != "" {
				continue
			}
		}
		path := filepath.Join(dir, name)
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("WAL artifact is not a regular file: %s", path)
		}
		paths = append(paths, path)
	}
	if len(paths) == 0 {
		return nil
	}
	backup, err := os.MkdirTemp(dir, fmt.Sprintf("rollback-wal-%d-", target))
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "Archiving consensus WAL to %s\n", backup); err != nil {
		return err
	}
	for _, path := range paths {
		if err := os.Rename(path, filepath.Join(backup, filepath.Base(path))); err != nil {
			return fmt.Errorf("archive WAL %s (earlier files retained in %s): %w", path, backup, err)
		}
	}
	// Persist the renames before changing either database. A crash must not
	// resurrect an old WAL head after the stores have been rolled back.
	for _, path := range []string{backup, dir} {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		err = file.Sync()
		closeErr := file.Close()
		if err != nil {
			return err
		}
		if closeErr != nil {
			return closeErr
		}
	}
	return nil
}
