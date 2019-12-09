package clitest

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/tests"
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/modules/guardian"
	"github.com/stretchr/testify/require"
)

func TestIrisCLIAddProfiler(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)

	fooAddr := f.KeyAddress(keyFoo)
	barAddr := f.KeyAddress(keyBar)

	success, _, stderr := f.TxAddProfiler(fooAddr.String(), barAddr.String(), "test")
	require.True(f.T, success)
	require.Empty(f.T, stderr)

	res := f.QueryProfiler()
	require.Len(f.T, res, 1)
	require.Equal(f.T, barAddr, res[0].Address)

	success, _, stderr = f.TxDeleteProfiler(fooAddr.String(), barAddr.String())
	require.True(f.T, success)
	require.Empty(f.T, stderr)

	res = f.QueryProfiler()
	require.Len(f.T, res, 0)

	// Cleanup testing directories
	f.Cleanup()
}

func TestIrisCLIAddTrustee(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)

	fooAddr := f.KeyAddress(keyFoo)
	barAddr := f.KeyAddress(keyBar)

	success, _, stderr := f.TxAddTrustee(fooAddr.String(), barAddr.String(), "test")
	require.True(f.T, success)
	require.Empty(f.T, stderr)

	res := f.QueryTrustee()
	require.Len(f.T, res, 1)
	require.Equal(f.T, barAddr, res[0].Address)

	success, _, stderr = f.TxDeleteTrustee(fooAddr.String(), barAddr.String())
	require.True(f.T, success)
	require.Empty(f.T, stderr)

	res = f.QueryTrustee()
	require.Len(f.T, res, 0)

	// Cleanup testing directories
	f.Cleanup()
}

//___________________________________________________________________________________
// iriscli tx guardian

// TxAddProfiler is iriscli tx guardian add-profiler
func (f *Fixtures) TxAddProfiler(from, address, description string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx guardian add-profiler %v --from=%s --address=%s --description=%s", f.IriscliBinary, f.Flags(), from, address, description)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxAddTrustee is iriscli tx guardian add-trustee
func (f *Fixtures) TxAddTrustee(from, address, description string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx guardian add-trustee %v --from=%s --address=%s --description=%s", f.IriscliBinary, f.Flags(), from, address, description)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxDeleteProfiler is iriscli tx guardian delete-profiler
func (f *Fixtures) TxDeleteProfiler(from, address string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx guardian delete-profiler %v --from=%s --address=%s", f.IriscliBinary, f.Flags(), from, address)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxDeleteTrustee is iriscli tx guardian delete-trustee
func (f *Fixtures) TxDeleteTrustee(from, address string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx guardian  delete-trustee %v --from=%s --address=%s", f.IriscliBinary, f.Flags(), from, address)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// QueryProfiler is iriscli query guardian profilers
func (f *Fixtures) QueryProfiler() (result []guardian.Guardian) {
	cmd := fmt.Sprintf("%s query guardian profilers --output=%s", f.IriscliBinary, "json")
	out, _ := tests.ExecuteT(f.T, cmd, "")
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &result)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return result
}

// QueryTrustee is iriscli query guardian profilers
func (f *Fixtures) QueryTrustee() (result []guardian.Guardian) {
	cmd := fmt.Sprintf("%s query guardian trustees --output=%s", f.IriscliBinary, "json")
	out, _ := tests.ExecuteT(f.T, cmd, "")
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &result)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return result
}
