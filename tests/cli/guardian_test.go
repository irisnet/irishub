package cli

import (
	"fmt"
	"testing"

	"github.com/irisnet/irishub/tests"
	"github.com/stretchr/testify/require"
	"github.com/irisnet/irishub/app/v0"
)

func TestIrisCLIAddProfiler(t *testing.T) {
	chainID, servAddr, port := initializeFixtures(t)
	flags := fmt.Sprintf("--home=%s --node=%v --chain-id=%v", iriscliHome, servAddr, chainID)
	// start iris server
	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iris start --home=%s --rpc.laddr=%v", irisHome, servAddr))
	defer proc.Stop(false)
	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(2, port)
	fooAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show foo --output=json --home=%s", iriscliHome))
	barAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show bar --output=json --home=%s", iriscliHome))
	profilers := executeGetProfilers(t, fmt.Sprintf("iriscli guardian profilers %v", flags))
	require.Equal(t, 1, len(profilers))
	require.Equal(t, fooAddr, profilers[0].Addr)
	scStr := fmt.Sprintf("iriscli guardian add-profiler %v", flags)
	scStr += fmt.Sprintf(" --profiler-address=%s", barAddr)
	scStr += fmt.Sprintf(" --profiler-name=%s", "bar")
	scStr += fmt.Sprintf(" --fee=%s", "0.004iris")
	scStr += fmt.Sprintf(" --from=%s", "foo")
	executeWrite(t, scStr, v0.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)
	profilers = executeGetProfilers(t, fmt.Sprintf("iriscli guardian profilers %v", flags))
	for _, profiler := range profilers {
		if profiler.Name != "genesis" {
			require.Equal(t, barAddr, profiler.Addr)
			require.Equal(t, fooAddr, profiler.AddedAddr)
			require.Equal(t, "bar", profiler.Name)
		}
	}

	trustees := executeGetTrustees(t, fmt.Sprintf("iriscli guardian trustees %v", flags))
	for _, trustee := range trustees {
		require.Equal(t, fooAddr, trustee.Addr)
	}
}
