package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/store/iavl"
	"cosmossdk.io/store/metrics"
	"cosmossdk.io/store/rootmulti"
	storetypes "cosmossdk.io/store/types"
	cmtDB "github.com/cometbft/cometbft-db"
	cmtcfg "github.com/cometbft/cometbft/config"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmtstate "github.com/cometbft/cometbft/state"
	cmtstore "github.com/cometbft/cometbft/store"
	cmttypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

type rollbackFixture struct {
	blocks *cmtstore.BlockStore
	states cmtstate.Store
	app    *rootmulti.Store
	db     dbm.DB
	keys   []*storetypes.KVStoreKey
	saved  []cmtstate.State
}

func newRollbackFixture(t *testing.T) *rollbackFixture {
	t.Helper()
	return newRollbackFixtureWithDBs(t, cmtDB.NewMemDB(), cmtDB.NewMemDB(), dbm.NewMemDB())
}

func newRollbackFixtureWithDBs(t *testing.T, blockDB, stateDB cmtDB.DB, appDB dbm.DB) *rollbackFixture {
	t.Helper()
	f := &rollbackFixture{
		blocks: cmtstore.NewBlockStore(blockDB),
		states: cmtstate.NewStore(stateDB, cmtstate.StoreOptions{}),
		db:     appDB,
		keys:   []*storetypes.KVStoreKey{storetypes.NewKVStoreKey("first"), storetypes.NewKVStoreKey("second")},
	}
	f.app = f.newAppStore()
	require.NoError(t, f.app.LoadLatestVersion())
	vals, _ := cmttypes.RandValidatorSet(3, 10)
	state, err := cmtstate.MakeGenesisState(&cmttypes.GenesisDoc{
		GenesisTime: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		ChainID:     "rollback-test", InitialHeight: 1,
		ConsensusParams: cmttypes.DefaultConsensusParams(),
		Validators: []cmttypes.GenesisValidator{
			{PubKey: vals.Validators[0].PubKey, Power: 10},
			{PubKey: vals.Validators[1].PubKey, Power: 10},
			{PubKey: vals.Validators[2].PubKey, Power: 10},
		},
	})
	require.NoError(t, err)
	require.NoError(t, f.states.Save(state))
	f.saved = append(f.saved, state)
	var previousBlockID cmttypes.BlockID
	for h := int64(1); h <= 6; h++ {
		lastCommit := testCommit(h-1, previousBlockID)
		block, err := state.MakeBlock(h, nil, lastCommit, nil, state.Validators.GetProposer().Address)
		require.NoError(t, err)
		parts, err := block.MakePartSet(cmttypes.BlockPartSizeBytes)
		require.NoError(t, err)
		blockID := cmttypes.BlockID{Hash: block.Hash(), PartSetHeader: parts.Header()}
		f.blocks.SaveBlock(block, parts, testCommit(h, blockID))
		previousBlockID = blockID
		for _, key := range f.keys {
			f.app.GetKVStore(key).Set([]byte("height"), []byte(fmt.Sprint(h)))
		}
		id := f.app.Commit()
		state.LastBlockHeight = h
		state.LastBlockID = blockID
		state.LastBlockTime = block.Time
		state.LastValidators = state.Validators.Copy()
		state.Validators = state.NextValidators.Copy()
		state.NextValidators = state.NextValidators.Copy()
		if h == 3 {
			changed := state.NextValidators.Validators[0].Copy()
			changed.VotingPower = 20
			require.NoError(t, state.NextValidators.UpdateWithChangeSet([]*cmttypes.Validator{changed}))
			state.LastHeightValidatorsChanged = h + 2
		}
		state.NextValidators.IncrementProposerPriority(1)
		if h == 4 {
			state.ConsensusParams.Block.MaxBytes = 1000000
			state.LastHeightConsensusParamsChanged = h + 1
		}
		state.AppHash = id.Hash
		require.NoError(t, f.states.Save(state))
		f.saved = append(f.saved, state.Copy())
	}
	return f
}

func testCommit(height int64, blockID cmttypes.BlockID) *cmttypes.Commit {
	commit := &cmttypes.Commit{Height: height, BlockID: blockID}
	if height > 0 {
		commit.Signatures = []cmttypes.CommitSig{cmttypes.NewCommitSigAbsent()}
	}
	return commit
}

func TestRollbackDiskStoresReopen(t *testing.T) {
	cfg := cmtcfg.DefaultConfig().SetRoot(t.TempDir())
	blockDB, err := cmtDB.NewDB("blockstore", cmtDB.GoLevelDBBackend, cfg.DBDir())
	require.NoError(t, err)
	stateDB, err := cmtDB.NewDB("state", cmtDB.GoLevelDBBackend, cfg.DBDir())
	require.NoError(t, err)
	appDB, err := dbm.NewDB("application", dbm.GoLevelDBBackend, cfg.DBDir())
	require.NoError(t, err)
	f := newRollbackFixtureWithDBs(t, blockDB, stateDB, appDB)
	_, err = rollbackStoresToHeight(f.blocks, f.states, f.app, f.db, 2, nil)
	require.NoError(t, err)
	require.NoError(t, f.blocks.Close())
	require.NoError(t, f.states.Close())
	require.NoError(t, f.db.Close())
	// Reopen actual on-disk databases, with no in-memory store or tree state.
	f.blocks, f.states, err = openRollbackStores(cfg)
	require.NoError(t, err)
	t.Cleanup(func() { _ = f.blocks.Close(); _ = f.states.Close(); _ = f.db.Close() })
	f.db, err = dbm.NewDB("application", dbm.GoLevelDBBackend, cfg.DBDir())
	require.NoError(t, err)
	f.app = f.newAppStore()
	require.NoError(t, f.app.LoadLatestVersion())
	f.assertHeight(t, 2)
	for h := int64(3); h <= 6; h++ {
		// A block must also be absent from the hash lookup after reopening.
		require.Nil(t, f.blocks.LoadBlockByHash(f.saved[h].LastBlockID.Hash))
	}
}

func (f *rollbackFixture) newAppStore() *rootmulti.Store {
	return f.newAppStoreWithKeys(f.keys)
}

func (f *rollbackFixture) newAppStoreWithKeys(keys []*storetypes.KVStoreKey) *rootmulti.Store {
	ms := rootmulti.NewStore(f.db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	for _, key := range keys {
		ms.MountStoreWithDB(key, storetypes.StoreTypeIAVL, nil)
	}
	return ms
}

func (f *rollbackFixture) assertHeight(t *testing.T, target int64) {
	t.Helper()
	s, err := f.states.Load()
	require.NoError(t, err)
	require.Equal(t, target, s.LastBlockHeight)
	require.Equal(t, target, f.blocks.Height())
	require.Equal(t, target, f.app.LatestVersion())
	require.Equal(t, f.saved[target].AppHash, s.AppHash)
	require.Equal(t, s.AppHash, f.app.LastCommitID().Hash)
	require.Equal(t, f.saved[target].ConsensusParams, s.ConsensusParams)
	require.Equal(t, f.saved[target].LastValidators, s.LastValidators)
	require.Equal(t, f.saved[target].Validators, s.Validators)
	require.Equal(t, f.saved[target].NextValidators, s.NextValidators)
	for h := target + 1; h <= 7; h++ {
		require.Nil(t, f.blocks.LoadBlock(h))
		require.Nil(t, f.blocks.LoadBlockMeta(h))
	}
	for _, key := range f.keys {
		require.Equal(t, []byte(fmt.Sprint(target)), f.app.GetKVStore(key).Get([]byte("height")))
	}
}

func TestRollbackToHeight(t *testing.T) {
	for _, pending := range []bool{false, true} {
		t.Run(fmt.Sprintf("pending=%t", pending), func(t *testing.T) {
			f := newRollbackFixture(t)
			if pending {
				s := f.saved[6]
				block, err := s.MakeBlock(7, nil, testCommit(6, s.LastBlockID), nil, s.Validators.GetProposer().Address)
				require.NoError(t, err)
				parts, err := block.MakePartSet(cmttypes.BlockPartSizeBytes)
				require.NoError(t, err)
				blockID := cmttypes.BlockID{Hash: block.Hash(), PartSetHeader: parts.Header()}
				f.blocks.SaveBlock(block, parts, testCommit(7, blockID))
			}
			hash, err := rollbackStoresToHeight(f.blocks, f.states, f.app, f.db, 2, nil)
			require.NoError(t, err)
			require.Equal(t, f.saved[2].AppHash, hash)
			f.assertHeight(t, 2)
			// A completed operation can be retried without stepping back again.
			_, err = rollbackStoresToHeight(f.blocks, f.states, f.app, f.db, 2, nil)
			require.NoError(t, err)
			f.assertHeight(t, 2)
			// Reopen the app exactly as startup does and commit a different H+1.
			f.app = f.newAppStore()
			require.NoError(t, f.app.LoadLatestVersion())
			f.app.GetKVStore(f.keys[0]).Set([]byte("height"), []byte("new-branch"))
			id := f.app.Commit()
			require.EqualValues(t, 3, id.Version)
			require.NotEqual(t, f.saved[3].AppHash, id.Hash)
		})
	}
}

func TestRollbackPendingBlockOnly(t *testing.T) {
	f := newRollbackFixture(t)
	s := f.saved[6]
	block, err := s.MakeBlock(7, nil, testCommit(6, s.LastBlockID), nil, s.Validators.GetProposer().Address)
	require.NoError(t, err)
	parts, err := block.MakePartSet(cmttypes.BlockPartSizeBytes)
	require.NoError(t, err)
	blockID := cmttypes.BlockID{Hash: block.Hash(), PartSetHeader: parts.Header()}
	f.blocks.SaveBlock(block, parts, testCommit(7, blockID))
	_, err = rollbackStoresToHeight(f.blocks, f.states, f.app, f.db, 6, nil)
	require.NoError(t, err)
	f.assertHeight(t, 6)
}

func TestRollbackPreflightDoesNotDeleteBlocks(t *testing.T) {
	for _, scenario := range []string{"zero", "negative", "future", "pruned-app", "wrong-hash", "missing-validators", "archive-failure"} {
		t.Run(scenario, func(t *testing.T) {
			f := newRollbackFixture(t)
			target := int64(2)
			var before func() error
			switch scenario {
			case "zero":
				target = 0
			case "negative":
				target = -1
			case "future":
				target = 7
			case "pruned-app":
				// No mounted tree retains H even though root commit info does.
				for _, key := range f.keys {
					require.NoError(t, f.app.GetCommitKVStore(key).(*iavl.Store).DeleteVersionsTo(3))
				}
			case "wrong-hash":
				// At the tip, expected hash comes from Comet state.
				target = 6
				state, err := f.states.Load()
				require.NoError(t, err)
				state.AppHash = []byte("invalid")
				require.NoError(t, f.states.Save(state))
			case "missing-validators":
				f.states = missingRollbackValidators{Store: f.states}
			case "archive-failure":
				before = func() error { return errors.New("archive failed") }
			}
			_, err := rollbackStoresToHeight(f.blocks, f.states, f.app, f.db, target, before)
			require.Error(t, err)
			require.EqualValues(t, 6, f.blocks.Height())
			state, err := f.states.Load()
			require.NoError(t, err)
			require.EqualValues(t, 6, state.LastBlockHeight)
			require.EqualValues(t, 6, rootmulti.GetLatestVersion(f.db))
		})
	}
}

func TestRollbackRejectsHistoricalStoreSetMismatchBeforeMutation(t *testing.T) {
	f := newRollbackFixture(t)
	// The target history contains both persistent stores, but the current binary
	// mounts only the first one. rootmulti.LoadVersion accepts this extra historic
	// StoreInfo, while RollbackToVersion would rewrite metadata without it.
	currentApp := f.newAppStoreWithKeys(f.keys[:1])
	_, err := rollbackStoresToHeight(f.blocks, f.states, currentApp, f.db, 2, nil)
	require.ErrorContains(t, err, "stores not mounted")
	require.EqualValues(t, 6, f.blocks.Height())
	state, err := f.states.Load()
	require.NoError(t, err)
	require.EqualValues(t, 6, state.LastBlockHeight)
	require.EqualValues(t, 6, rootmulti.GetLatestVersion(f.db))
	info, err := currentApp.GetCommitInfo(2)
	require.NoError(t, err)
	require.Len(t, info.StoreInfos, 2)
}

func TestRollbackRejectsDiscardedEvidenceBeforeMutation(t *testing.T) {
	f := newRollbackFixture(t)
	state := f.saved[6]
	vals := state.Validators
	voteA := &cmttypes.Vote{
		Type:             cmtproto.PrecommitType,
		Height:           1,
		BlockID:          cmttypes.BlockID{Hash: bytes.Repeat([]byte{1}, 32), PartSetHeader: cmttypes.PartSetHeader{Total: 1, Hash: bytes.Repeat([]byte{1}, 32)}},
		Timestamp:        time.Date(2026, 1, 1, 0, 0, 1, 0, time.UTC),
		ValidatorAddress: vals.Validators[0].Address,
		ValidatorIndex:   0,
		Signature:        []byte{1},
	}
	voteB := &cmttypes.Vote{
		Type:             cmtproto.PrecommitType,
		Height:           1,
		BlockID:          cmttypes.BlockID{Hash: bytes.Repeat([]byte{2}, 32), PartSetHeader: cmttypes.PartSetHeader{Total: 1, Hash: bytes.Repeat([]byte{2}, 32)}},
		Timestamp:        voteA.Timestamp,
		ValidatorAddress: vals.Validators[0].Address,
		ValidatorIndex:   0,
		Signature:        []byte{2},
	}
	evidence, err := cmttypes.NewDuplicateVoteEvidence(voteA, voteB, voteA.Timestamp, vals)
	require.NoError(t, err)
	block, err := state.MakeBlock(7, nil, testCommit(6, state.LastBlockID), []cmttypes.Evidence{evidence}, vals.GetProposer().Address)
	require.NoError(t, err)
	parts, err := block.MakePartSet(cmttypes.BlockPartSizeBytes)
	require.NoError(t, err)
	blockID := cmttypes.BlockID{Hash: block.Hash(), PartSetHeader: parts.Header()}
	f.blocks.SaveBlock(block, parts, testCommit(7, blockID))

	_, err = rollbackStoresToHeight(f.blocks, f.states, f.app, f.db, 2, nil)
	require.ErrorContains(t, err, "containing 1 evidence item")
	require.EqualValues(t, 7, f.blocks.Height())
	stateAfter, err := f.states.Load()
	require.NoError(t, err)
	require.EqualValues(t, 6, stateAfter.LastBlockHeight)
	require.EqualValues(t, 6, rootmulti.GetLatestVersion(f.db))
}

type missingRollbackValidators struct{ cmtstate.Store }

func (s missingRollbackValidators) LoadValidators(h int64) (*cmttypes.ValidatorSet, error) {
	if h == 2 {
		return nil, errors.New("pruned validators")
	}
	return s.Store.LoadValidators(h)
}

type interruptedRollbackState struct {
	cmtstate.Store
	interrupted bool
}

func (s *interruptedRollbackState) Save(state cmtstate.State) error {
	if err := s.Store.Save(state); err != nil {
		return err
	}
	if !s.interrupted {
		s.interrupted = true
		return errors.New("interrupted after state save, before block deletion")
	}
	return nil
}

func TestRollbackResumeAfterConsensusSave(t *testing.T) {
	f := newRollbackFixture(t)
	ss := &interruptedRollbackState{Store: f.states}
	_, err := rollbackStoresToHeight(f.blocks, ss, f.app, f.db, 2, nil)
	require.Error(t, err)
	state, err := ss.Load()
	require.NoError(t, err)
	require.EqualValues(t, 5, state.LastBlockHeight)
	require.EqualValues(t, 6, f.blocks.Height())
	_, err = rollbackStoresToHeight(f.blocks, ss, f.app, f.db, 2, nil)
	require.NoError(t, err)
	f.assertHeight(t, 2)
}

type interruptedRollbackApp struct{ *rootmulti.Store }

func (s interruptedRollbackApp) RollbackToVersion(target int64) error {
	return errors.New("interrupted before application rollback")
}

func TestRollbackResumePartialApplication(t *testing.T) {
	f := newRollbackFixture(t)
	_, err := rollbackStoresToHeight(f.blocks, f.states, interruptedRollbackApp{f.app}, f.db, 2, nil)
	require.Error(t, err)
	require.EqualValues(t, 2, f.blocks.Height())
	// Simulate a crash after the first IAVL tree has discarded future versions,
	// before remaining trees and root metadata have been updated.
	require.NoError(t, f.app.GetCommitKVStore(f.keys[0]).(*iavl.Store).LoadVersionForOverwriting(2))
	require.EqualValues(t, 6, rootmulti.GetLatestVersion(f.db))
	f.app = f.newAppStore()
	_, err = rollbackStoresToHeight(f.blocks, f.states, f.app, f.db, 2, nil)
	require.NoError(t, err)
	f.assertHeight(t, 2)
}

func TestRollbackHeightFlag(t *testing.T) {
	for _, args := range [][]string{nil, {"--hard"}, {"--height", "2"}, {"--hard", "--height", "0"}, {"--hard", "--height", "-1"}, {"--hard", "--height", "2", "extra"}} {
		t.Run(fmt.Sprint(args), func(t *testing.T) {
			legacy := false
			cmd := &cobra.Command{Use: "rollback", RunE: func(*cobra.Command, []string) error { legacy = true; return nil }}
			cmd.Flags().Bool("hard", false, "")
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)
			addRollbackHeightFlag(cmd)
			cmd.SetArgs(args)
			err := cmd.Execute()
			if len(args) <= 1 {
				require.NoError(t, err)
				require.True(t, legacy)
			} else {
				require.Error(t, err)
				require.False(t, legacy)
			}
		})
	}
}

func TestArchiveRollbackWAL(t *testing.T) {
	root := t.TempDir()
	cfg := cmtcfg.DefaultConfig().SetRoot(root)
	wal := cfg.Consensus.WalFile()
	require.NoError(t, os.MkdirAll(filepath.Dir(wal), 0700))
	for _, path := range []string{wal, wal + ".000", wal + ".001", wal + ".CORRUPTED", filepath.Join(root, "data", "priv_validator_state.json")} {
		require.NoError(t, os.WriteFile(path, []byte(path), 0600))
	}
	var out bytes.Buffer
	require.NoError(t, archiveRollbackWAL(wal, 2, &out))
	backups, err := filepath.Glob(filepath.Join(filepath.Dir(wal), "rollback-wal-2-*"))
	require.NoError(t, err)
	require.Len(t, backups, 1)
	for _, path := range []string{wal, wal + ".000", wal + ".001"} {
		_, err := os.Stat(path)
		require.True(t, os.IsNotExist(err))
		data, err := os.ReadFile(filepath.Join(backups[0], filepath.Base(path)))
		require.NoError(t, err)
		require.Equal(t, []byte(path), data)
	}
	_, err = os.Stat(wal + ".CORRUPTED")
	require.NoError(t, err)
	_, err = os.Stat(filepath.Join(root, "data", "priv_validator_state.json"))
	require.NoError(t, err)
	require.NoError(t, archiveRollbackWAL(wal, 2, io.Discard))
	// Never move a symlink or a directory masquerading as a WAL group member.
	require.NoError(t, os.Symlink(wal+".CORRUPTED", wal))
	require.Error(t, archiveRollbackWAL(wal, 2, io.Discard))
}
