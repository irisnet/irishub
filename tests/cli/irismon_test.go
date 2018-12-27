package cli

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/irisnet/irishub/server"
	"github.com/irisnet/irishub/tests"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

func TestIrismon(t *testing.T) {
	t.Parallel()
	chainID, servAddr, port, irisHome, iriscliHome, p2pAddr := initializeFixtures(t)
	// start iris server
	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iris start --home=%s --rpc.laddr=%v --p2p.laddr=%v", irisHome, servAddr, p2pAddr))

	defer proc.Stop(false)
	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(2, port)

	accAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show foo --output=json --home=%s", iriscliHome))
	valAddr := hex.EncodeToString(sdk.ValAddress(accAddr))

	// get a free port
	_, port1, err := server.FreeTCPAddr()
	require.NoError(t, err)

	// start irismon
	println(fmt.Sprintf("iristool monitor --validator-address=%s --chain-id=%s --account-address=%s --port=%s --node=http://localhost:%s", valAddr, chainID, accAddr.String(), port1, port))
	proc1 := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iristool monitor --validator-address=%s --chain-id=%s --account-address=%s --port=%s --node=http://localhost:%s", valAddr, chainID, accAddr.String(), port1, port))
	defer proc1.Stop(false)

	// wait 20s for irismon start
	time.Sleep(time.Second * 20)

	// irismon test
	resp, err := http.Get(fmt.Sprintf("http://localhost:%s", port1))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}
