package cli

import (
	"fmt"
	"testing"

	"github.com/irisnet/irishub/tests"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

func TestIrisCLIGateway(t *testing.T) {
	t.Parallel()
	chainID, servAddr, port, irisHome, iriscliHome, p2pAddr := initializeFixtures(t)

	flags := fmt.Sprintf("--home=%s --node=%v --chain-id=%v --output=json", iriscliHome, servAddr, chainID)

	// start iris server
	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iris start --home=%s --rpc.laddr=%v --p2p.laddr=%v", irisHome, servAddr, p2pAddr))

	defer proc.Stop(false)
	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(2, port)

	fooAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show foo --output=json --home=%s", iriscliHome))

	fooAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin := convertToIrisBaseAccount(t, fooAcc)
	require.Equal(t, "50iris", fooCoin)

	gatewayQuery, _ := tests.ExecuteT(t, fmt.Sprintf("iriscli asset query-gateway --moniker=uniquenm %v", flags), "")
	require.Equal(t, "null", gatewayQuery)

	// TODO: delete this comment after realizing gateway creation

	// define constant gateway fields
	moniker := "testgw"
	identity := "test gateway identity"
	details := "test gateway"
	website := "https://www.test-gateway.io"

	// create a gateway
	spStr := fmt.Sprintf("iriscli asset create-gateway %v", flags)
	spStr += fmt.Sprintf(" --from=%s", "foo")
	spStr += fmt.Sprintf(" --moniker=%s", moniker)
	spStr += fmt.Sprintf(" --identity=%s", identity)
	spStr += fmt.Sprintf(" --details=%s", details)
	spStr += fmt.Sprintf(" --website=%s", website)
	spStr += fmt.Sprintf(" --fee=%s", "0.4iris")

	executeWrite(t, spStr, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	num := getAmountFromCoinStr(fooCoin)

	if !(num > 44 && num < 45) {
		t.Error("Test Failed: (44, 45) expected, recieved:", num)
	}

	gateway := executeGetGateway(t, fmt.Sprintf("iriscli asset query-gateway --moniker=testgw --output=json %v", flags))
	require.Equal(t, moniker, gateway.Moniker)
	require.Equal(t, identity, gateway.Identity)
	require.Equal(t, details, gateway.Details)
	require.Equal(t, website, gateway.Website)

	gateways := executeGetGateways(t, fmt.Sprintf("iriscli asset query-gateways --owner=foo %v", flags))
	require.Equal(t, 1, len(gateways))
}
