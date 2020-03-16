package clitest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/tests"

	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/modules/guardian"
)

func TestIrisCLIAddProfiler(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)

	fooAddr := f.KeyAddress(keyFoo)
	barAddr := f.KeyAddress(keyBar)

	// start iris server
	proc := f.GDStart()
	defer proc.Stop(false)

	description := "test"

	success, _, _ := f.TxAddProfiler(fooAddr.String(), barAddr.String(), description, "-y")
	require.True(f.T, success)

	tests.WaitForNextNBlocksTM(1, f.Port)
	// Ensure transaction tags can be queried
	searchResult := f.QueryTxs(1, 50, "message.action=add_profiler", fmt.Sprintf("message.sender=%s", fooAddr))
	require.Len(t, searchResult.Txs, 1)

	expGuardian := guardian.NewGuardian(description, guardian.Ordinary, barAddr, fooAddr)

	res := f.QueryProfilers()
	require.NotEmpty(f.T, res)
	require.Contains(f.T, res, expGuardian)

	success, _, _ = f.TxDeleteProfiler(fooAddr.String(), barAddr.String(), "-y")
	require.True(f.T, success)

	tests.WaitForNextNBlocksTM(1, f.Port)
	// Ensure transaction tags can be queried
	searchResult = f.QueryTxs(1, 50, "message.action=delete_profiler", fmt.Sprintf("message.sender=%s", fooAddr))
	require.Len(t, searchResult.Txs, 1)

	res = f.QueryProfilers()
	require.NotEmpty(f.T, res)
	require.NotContains(f.T, res, expGuardian)

	// Cleanup testing directories
	f.Cleanup()
}

func TestIrisCLIAddTrustee(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)

	// start iris server
	proc := f.GDStart()
	defer proc.Stop(false)

	fooAddr := f.KeyAddress(keyFoo)
	barAddr := f.KeyAddress(keyBar)

	description := "test"

	success, _, _ := f.TxAddTrustee(fooAddr.String(), barAddr.String(), description, "-y")
	require.True(f.T, success)

	expGuardian := guardian.NewGuardian(description, guardian.Ordinary, barAddr, fooAddr)

	tests.WaitForNextNBlocksTM(1, f.Port)
	// Ensure transaction tags can be queried
	searchResult := f.QueryTxs(1, 50, "message.action=add_trustee", fmt.Sprintf("message.sender=%s", fooAddr))
	require.Len(t, searchResult.Txs, 1)

	res := f.QueryTrustees()
	require.NotEmpty(f.T, res)
	require.Contains(f.T, res, expGuardian)

	success, _, _ = f.TxDeleteTrustee(fooAddr.String(), barAddr.String(), "-y")
	require.True(f.T, success)

	tests.WaitForNextNBlocksTM(1, f.Port)
	// Ensure transaction tags can be queried
	searchResult = f.QueryTxs(1, 50, "message.action=delete_trustee", fmt.Sprintf("message.sender=%s", fooAddr))
	require.Len(t, searchResult.Txs, 1)

	res = f.QueryTrustees()
	require.NotEmpty(f.T, res)
	require.NotContains(f.T, res, expGuardian)

	// Cleanup testing directories
	f.Cleanup()
}

//___________________________________________________________________________________
// iriscli tx guardian

// TxAddProfiler is iriscli tx guardian add-profiler
func (f *Fixtures) TxAddProfiler(from, address, description string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx guardian add-profiler %v --keyring-backend=test --from=%s --address=%s --description=%s", f.IriscliBinary, f.Flags(), from, address, description)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), keys.DefaultKeyPass)
}

// TxAddTrustee is iriscli tx guardian add-trustee
func (f *Fixtures) TxAddTrustee(from, address, description string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx guardian add-trustee %v --keyring-backend=test --from=%s --address=%s --description=%s", f.IriscliBinary, f.Flags(), from, address, description)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), keys.DefaultKeyPass)
}

// TxDeleteProfiler is iriscli tx guardian delete-profiler
func (f *Fixtures) TxDeleteProfiler(from, address string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx guardian delete-profiler %v --keyring-backend=test --from=%s --address=%s", f.IriscliBinary, f.Flags(), from, address)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), keys.DefaultKeyPass)
}

// TxDeleteTrustee is iriscli tx guardian delete-trustee
func (f *Fixtures) TxDeleteTrustee(from, address string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx guardian delete-trustee %v --keyring-backend=test --from=%s --address=%s", f.IriscliBinary, f.Flags(), from, address)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), keys.DefaultKeyPass)
}

// QueryProfiler is iriscli query guardian profilers
func (f *Fixtures) QueryProfilers() (result guardian.Profilers) {
	cmd := fmt.Sprintf("%s query guardian profilers --output=%s %v", f.IriscliBinary, "json", f.Flags())
	out, _ := tests.ExecuteT(f.T, cmd, "")
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &result)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return
}

// QueryTrustee is iriscli query guardian profilers
func (f *Fixtures) QueryTrustees() (result guardian.Trustees) {
	cmd := fmt.Sprintf("%s query guardian trustees --output=%s %v", f.IriscliBinary, "json", f.Flags())
	out, _ := tests.ExecuteT(f.T, cmd, "")
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &result)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return
}
