package clitest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/cosmos/cosmos-sdk/tests"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/irisnet/irishub/app"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"encoding/hex"
	"net/http"
	"time"
)

func init() {
	irisHome, iriscliHome = getTestingHomeDirs()
}

func TestIrismon(t *testing.T) {
	tests.ExecuteT(t, fmt.Sprintf("iris --home=%s unsafe_reset_all", irisHome), "")
	executeWrite(t, fmt.Sprintf("iriscli keys delete --home=%s foo", iriscliHome), app.DefaultKeyPass)
	executeWrite(t, fmt.Sprintf("iriscli keys delete --home=%s bar", iriscliHome), app.DefaultKeyPass)
	chainID, _ := executeInit(t, fmt.Sprintf("iris init -o --name=foo --home=%s --home-client=%s", irisHome, iriscliHome))
	executeWrite(t, fmt.Sprintf("iriscli keys add --home=%s bar", iriscliHome), app.DefaultKeyPass)

	err := modifyGenesisFile(irisHome)
	require.NoError(t, err)

	// get a free port, also setup some common flags
	servAddr, port, err := server.FreeTCPAddr()
	require.NoError(t, err)
	flags := fmt.Sprintf("--home=%s --node=%v --chain-id=%v", iriscliHome, servAddr, chainID)

	// start iris server
	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iris start --home=%s --rpc.laddr=%v", irisHome, servAddr))

	defer proc.Stop(false)
	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(2, port)

	accountAddress, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show foo --output=json --home=%s", iriscliHome))
	validator := executeGetValidator(t, fmt.Sprintf("iriscli stake validator %s --output=json %v", accountAddress, flags))
	pk, err := sdk.GetValPubKeyBech32(validator.PubKey)
	pkHex := hex.EncodeToString(pk.Bytes())

	// get a free port
	_, port1, err := server.FreeTCPAddr()
	require.NoError(t, err)

	// start irismon
	proc1 := tests.GoExecuteTWithStdout(t, fmt.Sprintf("irismon --address=%s --chain-id=%s --account-address=%s --port=%s", pkHex, chainID, accountAddress.String(), port1))

	// wait 20s for irismon start
	time.Sleep(time.Second * 20)

	// irismon test
	resp, err := http.Get(fmt.Sprintf("http://0.0.0.0:%s", port1))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	defer proc1.Stop(false)
}
