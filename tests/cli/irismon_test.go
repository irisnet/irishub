package cli

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/irisnet/irishub/tests"
	"github.com/irisnet/irishub/server"
	sdk "github.com/irisnet/irishub/types"
	"encoding/hex"
	"net/http"
	"time"
)

func TestIrismon(t *testing.T) {
	t.SkipNow()

	chainID, servAddr, port := initializeFixtures(t)
	flags := fmt.Sprintf("--home=%s --node=%v --chain-id=%v", iriscliHome, servAddr, chainID)

	// start iris server
	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iris start --home=%s --rpc.laddr=%v", irisHome, servAddr))

	defer proc.Stop(false)
	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(2, port)

	accountAddress, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show foo --output=json --home=%s", iriscliHome))
	validator := executeGetValidator(t, fmt.Sprintf("iriscli stake validator %s --output=json %v", accountAddress, flags))
	pk, err := sdk.GetValPubKeyBech32(validator.ConsPubKey)
	pkHex := hex.EncodeToString(pk.Bytes())

	// get a free port
	_, port1, err := server.FreeTCPAddr()
	require.NoError(t, err)

	// start irismon
	proc1 := tests.GoExecuteTWithStdout(t, fmt.Sprintf("irismon --address=%s --chain-id=%s --account-address=%s --port=%s", pkHex, chainID, accountAddress.String(), port1))
	defer proc1.Stop(false)

	// wait 20s for irismon start
	time.Sleep(time.Second * 20)

	// irismon test
	resp, err := http.Get(fmt.Sprintf("http://localhost:%s", port1))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}
