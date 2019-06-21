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
	//TODO
	require.Equal(t, "", gatewayQuery)

	// define constant gateway fields
	moniker := "testgw"
	identity := "test-gateway-identity"
	details := "test-gateway"
	website := "https://www.test-gateway.io"

	// create a gateway
	spStr := fmt.Sprintf("iriscli asset create-gateway %v", flags)
	spStr += fmt.Sprintf(" --from=%s", "foo")
	spStr += fmt.Sprintf(" --moniker=%s", moniker)
	spStr += fmt.Sprintf(" --identity=%s", identity)
	spStr += fmt.Sprintf(" --details=%s", details)
	spStr += fmt.Sprintf(" --website=%s", website)
	spStr += fmt.Sprintf(" --create-fee=%s", "30iris")
	spStr += fmt.Sprintf(" --fee=%s", "0.4iris")

	require.True(t, executeWrite(t, spStr, sdk.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, port)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	num := getAmountFromCoinStr(fooCoin)

	// TODO: balance - create-fee
	if !(num > 41 && num < 45) {
		t.Error("Test Failed: (41, 45) expected, recieved:", num)
	}

	gateway := executeGetGateway(t, fmt.Sprintf("iriscli asset query-gateway --moniker=testgw --output=json %v", flags))
	require.Equal(t, moniker, gateway.Moniker)
	require.Equal(t, identity, gateway.Identity)
	require.Equal(t, details, gateway.Details)
	require.Equal(t, website, gateway.Website)

	gateways := executeGetGateways(t, fmt.Sprintf("iriscli asset query-gateways --owner=%s %v", fooAddr.String(), flags))
	require.Equal(t, 1, len(gateways))
}

func TestIrisCLIToken(t *testing.T) {
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

	gatewayToken, _ := tests.ExecuteT(t, fmt.Sprintf("iriscli asset query-token kitty %v", flags), "")
	//TODO
	require.Equal(t, "", gatewayToken)

	// define constant gateway fields
	family := "fungible"
	source := "native"
	symbol := "kitty"
	name := "kitty"
	initialSupply := 1000000

	// create a gateway
	spStr := fmt.Sprintf("iriscli asset issue-token %v", flags)
	spStr += fmt.Sprintf(" --from=%s", "foo")
	spStr += fmt.Sprintf(" --name=%s", name)
	spStr += fmt.Sprintf(" --symbol=%s", symbol)
	spStr += fmt.Sprintf(" --source=%s", source)
	spStr += fmt.Sprintf(" --family=%s", family)
	spStr += fmt.Sprintf(" --initial-supply=%d", initialSupply)
	spStr += fmt.Sprintf(" --fee=%s", "40iris")

	require.True(t, executeWrite(t, spStr, sdk.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, port)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))

	// TODO: make sure there are new tokens in fooAcc

	token := executeGetToken(t, fmt.Sprintf("iriscli asset query-token kitty --output=json %v", flags))
	require.Equal(t, name, token.Name)
	require.Equal(t, symbol, token.Symbol)
	require.Equal(t, family, token.Family.String())
	require.Equal(t, source, token.Source.String())
	//TODO:edit_token cli_test
}
