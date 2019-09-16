package cli

import (
	"fmt"
	"strings"
	"testing"

	"github.com/irisnet/irishub/tests"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

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

	family := "fungible"
	source := "native"
	symbol := "btc"
	name := "Bitcoin"
	initialSupply := 2000000000
	decimal := 18
	canonicalSymbol := "Btc"
	minUnitAlias := "Satoshi"
	gateway := "ABC"

	// issue a token
	spStr := fmt.Sprintf("iriscli asset issue-token %v", flags)
	spStr += fmt.Sprintf(" --from=%s", "foo")
	spStr += fmt.Sprintf(" --family=%s", family)
	spStr += fmt.Sprintf(" --source=%s", source)
	spStr += fmt.Sprintf(" --symbol=%s", symbol)
	spStr += fmt.Sprintf(" --name=%s", name)
	spStr += fmt.Sprintf(" --initial-supply=%d", initialSupply)
	spStr += fmt.Sprintf(" --decimal=%d", decimal)
	spStr += fmt.Sprintf(" --canonical-symbol=%s", canonicalSymbol)
	spStr += fmt.Sprintf(" --min-unit-alias=%s", minUnitAlias)
	spStr += fmt.Sprintf(" --gateway=%s", gateway)
	spStr += fmt.Sprintf(" --fee=%s", "0.4iris")

	require.True(t, executeWrite(t, spStr, sdk.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, port)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	amt := getAmountFromCoinStr(fooCoin)

	// 30iris is used to issue tokens
	if !(amt > 19 && amt < 20) {
		t.Error("Test Failed: (19, 20) expected, received:", amt)
	}

	token := executeGetToken(t, fmt.Sprintf("iriscli asset query-token %s --output=json %v", strings.ToLower(strings.TrimSpace(symbol)), flags))
	require.Equal(t, strings.ToLower(strings.TrimSpace(family)), token.Family.String())
	require.Equal(t, strings.ToLower(strings.TrimSpace(source)), token.Source.String())
	require.Equal(t, strings.ToLower(strings.TrimSpace(symbol)), token.Symbol)
	require.Equal(t, strings.TrimSpace(name), token.Name)
	require.Equal(t, strings.ToLower(strings.TrimSpace(minUnitAlias)), token.MinUnitAlias)
	require.Equal(t, sdk.NewIntWithDecimal(int64(initialSupply), decimal), token.InitialSupply)
	require.Equal(t, uint8(decimal), token.Decimal)
	require.Equal(t, "", token.CanonicalSymbol) // ignored by native token
	require.Equal(t, "", token.Gateway)         // ignored by native token

}

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
	barAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show bar --output=json --home=%s", iriscliHome))

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
	cgStr := fmt.Sprintf("iriscli asset create-gateway %v", flags)
	cgStr += fmt.Sprintf(" --from=%s", "foo")
	cgStr += fmt.Sprintf(" --moniker=%s", moniker)
	cgStr += fmt.Sprintf(" --identity=%s", identity)
	cgStr += fmt.Sprintf(" --details=%s", details)
	cgStr += fmt.Sprintf(" --website=%s", website)
	cgStr += fmt.Sprintf(" --fee=%s", "0.4iris")

	require.True(t, executeWrite(t, cgStr, sdk.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, port)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	num := getAmountFromCoinStr(fooCoin)

	// TODO: balance - create-fee
	if !(num > 41 && num < 45) {
		t.Error("Test Failed: (41, 45) expected, received:", num)
	}

	gateway := executeGetGateway(t, fmt.Sprintf("iriscli asset query-gateway --moniker=testgw --output=json %v", flags))
	require.Equal(t, moniker, gateway.Moniker)
	require.Equal(t, identity, gateway.Identity)
	require.Equal(t, details, gateway.Details)
	require.Equal(t, website, gateway.Website)

	gateways := executeGetGateways(t, fmt.Sprintf("iriscli asset query-gateways --owner=%s %v", fooAddr.String(), flags))
	require.Equal(t, 1, len(gateways))

	// transfer the gateway owner
	tgStr := fmt.Sprintf("iriscli asset transfer-gateway-owner %v", flags)
	tgStr += fmt.Sprintf(" --from=%s", "foo")
	tgStr += fmt.Sprintf(" --moniker=%s", moniker)
	tgStr += fmt.Sprintf(" --to=%s", barAddr.String())
	tgStr += fmt.Sprintf(" --fee=%s", "0.4iris")

	require.True(t, executeWrite(t, tgStr, sdk.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, port)

	gateway = executeGetGateway(t, fmt.Sprintf("iriscli asset query-gateway --moniker=%s %v", moniker, flags))
	require.Equal(t, barAddr.String(), gateway.Owner.String())

	gateways = executeGetGateways(t, fmt.Sprintf("iriscli asset query-gateways --owner=%s %v", barAddr.String(), flags))
	require.Equal(t, 1, len(gateways))
	require.Equal(t, moniker, gateways[0].Moniker)
	require.Equal(t, identity, gateways[0].Identity)
	require.Equal(t, details, gateways[0].Details)
	require.Equal(t, website, gateways[0].Website)

	gateways = executeGetGateways(t, fmt.Sprintf("iriscli asset query-gateways --owner=%s %v", fooAddr.String(), flags))
	require.Equal(t, 0, len(gateways))
}

func TestIrisCLIEditToken(t *testing.T) {
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

	family := "fungible"
	source := "native"
	symbol := "AbcdefgH"
	name := "Bitcoin"
	initialSupply := 100000000
	decimal := 18
	canonicalSymbol := "Btc"
	minUnitAlias := "Satoshi"
	gateway := "ABC"

	// issue a token
	spStr := fmt.Sprintf("iriscli asset issue-token %v", flags)
	spStr += fmt.Sprintf(" --from=%s", "foo")
	spStr += fmt.Sprintf(" --family=%s", family)
	spStr += fmt.Sprintf(" --source=%s", source)
	spStr += fmt.Sprintf(" --symbol=%s", symbol)
	spStr += fmt.Sprintf(" --name=%s", name)
	spStr += fmt.Sprintf(" --initial-supply=%d", initialSupply)
	spStr += fmt.Sprintf(" --decimal=%d", decimal)
	spStr += fmt.Sprintf(" --canonical-symbol=%s", canonicalSymbol)
	spStr += fmt.Sprintf(" --min-unit-alias=%s", minUnitAlias)
	spStr += fmt.Sprintf(" --gateway=%s", gateway)
	spStr += fmt.Sprintf(" --fee=%s", "0.4iris")

	require.True(t, executeWrite(t, spStr, sdk.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, port)

	//
	editTokenStr := fmt.Sprintf("iriscli asset edit-token %s", strings.ToLower(strings.TrimSpace(symbol)))
	editTokenStr += fmt.Sprintf(" --from=%s", "foo")
	editTokenStr += fmt.Sprintf(" --name=%s", "BTC_Token")
	editTokenStr += fmt.Sprintf(" --canonical-symbol=%s", "BTC1")
	editTokenStr += fmt.Sprintf(" --max-supply=%d", 200000000)
	editTokenStr += fmt.Sprintf(" --mintable=%v", true)
	editTokenStr += fmt.Sprintf(" --fee=%s %v", "0.4iris", flags)
	require.True(t, executeWrite(t, editTokenStr, sdk.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, port)

	token := executeGetToken(t, fmt.Sprintf("iriscli asset query-token %s --output=json %v", strings.ToLower(strings.TrimSpace(symbol)), flags))

	require.Equal(t, "BTC_Token", token.Name)
	require.Equal(t, "", token.CanonicalSymbol)
	require.Equal(t, sdk.NewIntWithDecimal(int64(200000000), decimal), token.MaxSupply)
	require.Equal(t, true, token.Mintable)
}
