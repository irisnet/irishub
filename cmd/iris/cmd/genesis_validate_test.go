package cmd

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"testing"

	"cosmossdk.io/log"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client/flags"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/irisnet/irishub/v5/app"
)

func TestGenesisValidateLegacyLSM(t *testing.T) {
	application := app.NewIrisApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, simtestutil.AppOptionsMap{})
	state := application.DefaultGenesis()
	var staking map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(state[stakingtypes.ModuleName], &staking))
	// These disabled-feature locks are removed by the full app preflight.
	// The module-only fallback deliberately rejects them, so this checks that
	// the actual CLI command reaches application-wide validation.
	staking["tokenize_share_locks"] = json.RawMessage(`[{"address":"unused","status":"LOCKED"}]`)
	staking["total_liquid_staked_tokens"] = json.RawMessage(`"0"`)
	staking["tokenize_share_records"] = json.RawMessage(`[]`)
	state[stakingtypes.ModuleName], _ = json.Marshal(staking)
	data, err := json.Marshal(state)
	require.NoError(t, err)
	genesis := genutiltypes.NewAppGenesisWithVersion("lsm-test-1", data)
	require.NoError(t, genesis.ValidateAndComplete())
	path := filepath.Join(t.TempDir(), "genesis.json")
	require.NoError(t, genesis.SaveAs(path))
	command := NewRootCmd()
	// server/cmd.Execute installs this flag around NewRootCmd in the binary.
	command.PersistentFlags().String(flags.FlagHome, t.TempDir(), "node home")
	var output bytes.Buffer
	command.SetOut(&output)
	command.SetErr(&output)
	command.SetArgs([]string{"genesis", "validate-genesis", path, "--home", t.TempDir()})
	require.NoError(t, command.Execute())
	require.Contains(t, output.String(), "is a valid genesis file")
	// A record without matching bank principal must fail through the same CLI.
	staking["tokenize_share_records"] = json.RawMessage(`[{"id":"1","owner":"invalid","module_account":"tokenizeshare_1","validator":"invalid"}]`)
	state[stakingtypes.ModuleName], _ = json.Marshal(staking)
	genesis.AppState, _ = json.Marshal(state)
	require.NoError(t, genesis.SaveAs(path))
	command = NewRootCmd()
	command.PersistentFlags().String(flags.FlagHome, t.TempDir(), "node home")
	command.SetOut(&output)
	command.SetErr(&output)
	command.SetArgs([]string{"genesis", "validate-genesis", path, "--home", t.TempDir()})
	require.Error(t, command.Execute())
}
