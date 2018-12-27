package cli

import (
	"fmt"
	"testing"

	"github.com/irisnet/irishub/tests"
	"github.com/stretchr/testify/require"
	"github.com/irisnet/irishub/app/v0"
	"github.com/irisnet/irishub/modules/guardian"
)

func TestIrisCLIAddProfiler(t *testing.T) {
	t.Parallel()
	chainID, servAddr, port, irisHome, iriscliHome, p2pAddr := initializeFixtures(t)

	flags := fmt.Sprintf("--home=%s --node=%v --chain-id=%v", iriscliHome, servAddr, chainID)

	// start iris server
	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iris start --home=%s --rpc.laddr=%v --p2p.laddr=%v", irisHome, servAddr, p2pAddr))
	defer proc.Stop(false)

	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(2, port)
	fooAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show foo --output=json --home=%s", iriscliHome))
	barAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show bar --output=json --home=%s", iriscliHome))

	// add profiler
	profilers := executeGetProfilers(t, fmt.Sprintf("iriscli guardian profilers %v", flags))
	require.Equal(t, 1, len(profilers))
	require.Equal(t, fooAddr, profilers[0].Address)

	paStr := fmt.Sprintf("iriscli guardian add-profiler %v", flags)
	paStr += fmt.Sprintf(" --address=%s", barAddr)
	paStr += fmt.Sprintf(" --description=%s", "bar")
	paStr += fmt.Sprintf(" --fee=%s", "0.004iris")
	paStr += fmt.Sprintf(" --from=%s", "foo")
	paStr += " --commit"
	require.True(t, executeWrite(t, paStr, v0.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, port)
	profilers = executeGetProfilers(t, fmt.Sprintf("iriscli guardian profilers %v", flags))
	require.Equal(t, 2, len(profilers))
	for _, profiler := range profilers {
		if profiler.AccountType != guardian.Genesis {
			require.Equal(t, barAddr, profiler.Address)
			require.Equal(t, fooAddr, profiler.AddedBy)
			require.Equal(t, "bar", profiler.Description)
			require.Equal(t, guardian.Ordinary, profiler.AccountType)
		}
	}

	// add trustee
	trustees := executeGetTrustees(t, fmt.Sprintf("iriscli guardian trustees %v", flags))
	require.Equal(t, 1, len(trustees))
	require.Equal(t, fooAddr, trustees[0].Address)

	taStr := fmt.Sprintf("iriscli guardian add-trustee %v", flags)
	taStr += fmt.Sprintf(" --description=%s", "bar")
	taStr += fmt.Sprintf(" --fee=%s", "0.004iris")
	taStr += fmt.Sprintf(" --address=%s", barAddr)
	taStr += fmt.Sprintf(" --from=%s", "foo")
	taStr += " --commit"
	require.True(t, executeWrite(t, taStr, v0.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, port)
	trustees = executeGetProfilers(t, fmt.Sprintf("iriscli guardian trustees %v", flags))
	require.Equal(t, 2, len(trustees))
	for _, trustee := range trustees {
		if trustee.AccountType != guardian.Genesis {
			require.Equal(t, barAddr, trustee.Address)
			require.Equal(t, fooAddr, trustee.AddedBy)
			require.Equal(t, "bar", trustee.Description)
			require.Equal(t, guardian.Ordinary, trustee.AccountType)
		}
	}

	// delete profiler
	pdStr := fmt.Sprintf("iriscli guardian delete-profiler %v", flags)
	pdStr += fmt.Sprintf(" --fee=%s", "0.004iris")
	pdStr += fmt.Sprintf(" --from=%s", "foo")
	pdStr += " --commit"

	pdbStr := pdStr + fmt.Sprintf(" --address=%s", barAddr)
	pdfStr := pdStr + fmt.Sprintf(" --address=%s", fooAddr)
	require.Equal(t, false, executeWrite(t, pdfStr, v0.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, port)
	require.True(t, executeWrite(t, pdbStr, v0.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, port)
	profilers = executeGetProfilers(t, fmt.Sprintf("iriscli guardian profilers %v", flags))
	require.Equal(t, 1, len(profilers))
	require.Equal(t, fooAddr, profilers[0].Address)

	// delete trustee
	tdStr := fmt.Sprintf("iriscli guardian delete-trustee %v", flags)
	tdStr += fmt.Sprintf(" --fee=%s", "0.004iris")
	tdStr += fmt.Sprintf(" --from=%s", "foo")
	tdStr += " --commit"

	tdbStr := tdStr + fmt.Sprintf(" --address=%s", barAddr)
	tdfStr := tdStr + fmt.Sprintf(" --address=%s", fooAddr)
	require.Equal(t, false, executeWrite(t, tdfStr, v0.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, port)
	require.True(t, executeWrite(t, tdbStr, v0.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, port)
	trustees = executeGetTrustees(t, fmt.Sprintf("iriscli guardian trustees %v", flags))
	require.Equal(t, 1, len(trustees))
	require.Equal(t, fooAddr, trustees[0].Address)
}
