package cli

import (
	"fmt"
	"testing"

	"regexp"
	"strings"

	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/tests"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

func TestIrisCLIService(t *testing.T) {
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

	// define service
	serviceName := "testService"
	serviceSchemas := `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`

	svcDefOutput, _ := tests.ExecuteT(t, fmt.Sprintf("iriscli service definition %s %v", serviceName, flags), "")
	require.Equal(t, "", svcDefOutput)

	sdStr := fmt.Sprintf("iriscli service define %v", flags)
	sdStr += fmt.Sprintf(" --from=%s", "foo")
	sdStr += fmt.Sprintf(" --name=%s", serviceName)
	sdStr += fmt.Sprintf(" --description=%s", "test")
	sdStr += fmt.Sprintf(" --tags=%s", "tag1,tag2")
	sdStr += fmt.Sprintf(" --author-description=%s", "foo")
	sdStr += fmt.Sprintf(" --schemas=%s", serviceSchemas)
	sdStr += fmt.Sprintf(" --fee=%s", "0.4iris")

	executeWrite(t, sdStr, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	fooAmt := getAmountFromCoinStr(fooCoin)

	if !(fooAmt > 49 && fooAmt < 50) {
		t.Error("Test Failed: (49, 50) expected, received: {}", fooAmt)
	}

	svcDef := executeGetServiceDefinition(t, fmt.Sprintf("iriscli service definition %s %v", serviceName, flags))
	require.Equal(t, serviceName, svcDef.Name)
	require.Equal(t, serviceSchemas, svcDef.Schemas)

	// bind service
	sdStr = fmt.Sprintf("iriscli service bind %v", flags)
	sdStr += fmt.Sprintf(" --service-name=%s", serviceName)
	sdStr += fmt.Sprintf(" --def-chain-id=%s", chainID)
	sdStr += fmt.Sprintf(" --bind-type=%s", "Local")
	sdStr += fmt.Sprintf(" --deposit=%s", "10iris")
	sdStr += fmt.Sprintf(" --prices=%s", "1iris")
	sdStr += fmt.Sprintf(" --avg-rsp-time=%d", 10000)
	sdStr += fmt.Sprintf(" --usable-time=%d", 10000)
	sdStr += fmt.Sprintf(" --fee=%s", "0.4iris")

	sdStrFoo := sdStr + fmt.Sprintf(" --from=%s", "foo")
	sdStrBar := sdStr + fmt.Sprintf(" --from=%s", "bar")

	executeWrite(t, sdStrFoo, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	fooAmt = getAmountFromCoinStr(fooCoin)

	if !(fooAmt > 39 && fooAmt < 40) {
		t.Error("Test Failed: (39, 40) expected, received: {}", fooAmt)
	}

	executeWrite(t, fmt.Sprintf("iriscli bank send --to=%s --from=%s --amount=20iris --fee=0.3iris %v", barAddr.String(), "foo", flags), sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	executeWrite(t, sdStrBar, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	barAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin := convertToIrisBaseAccount(t, barAcc)
	barAmt := getAmountFromCoinStr(barCoin)

	if !(barAmt > 9 && barAmt < 10) {
		t.Error("Test Failed: (9, 10) expected, received: {}", barAmt)
	}

	svcBinding := executeGetServiceBinding(t, fmt.Sprintf("iriscli service binding --service-name=%s --def-chain-id=%s --bind-chain-id=%s --provider=%s %v", serviceName, chainID, chainID, fooAddr.String(), flags))
	require.NotNil(t, svcBinding)

	svcBindings := executeGetServiceBindings(t, fmt.Sprintf("iriscli service bindings --service-name=%s --def-chain-id=%s %v", serviceName, chainID, flags))
	require.Equal(t, 2, len(svcBindings))

	// update binding
	ubStr := fmt.Sprintf("iriscli service update-binding %v", flags)
	ubStr += fmt.Sprintf(" --service-name=%s", serviceName)
	ubStr += fmt.Sprintf(" --def-chain-id=%s", chainID)
	ubStr += fmt.Sprintf(" --bind-type=%s", "Global")
	ubStr += fmt.Sprintf(" --deposit=%s", "1iris")
	ubStr += fmt.Sprintf(" --prices=%s", "0.1iris")
	ubStr += fmt.Sprintf(" --avg-rsp-time=%d", 99)
	ubStr += fmt.Sprintf(" --usable-time=%d", 99)
	ubStr += fmt.Sprintf(" --fee=%s", "0.4iris")
	ubStr += fmt.Sprintf(" --from=%s", "bar")

	executeWrite(t, ubStr, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	barAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin = convertToIrisBaseAccount(t, barAcc)
	barAmt = getAmountFromCoinStr(barCoin)

	if !(barAmt > 8 && barAmt < 9) {
		t.Error("Test Failed: (8, 9) expected, received: {}", barAmt)
	}

	svcBindings = executeGetServiceBindings(t, fmt.Sprintf("iriscli service bindings --service-name=%s --def-chain-id=%s %v", serviceName, chainID, flags))

	var totalDeposit sdk.Coins
	for _, bind := range svcBindings {
		totalDeposit = totalDeposit.Add(bind.Deposit)
	}
	require.Equal(t, "21000000000000000000iris-atto", totalDeposit.String())

	// disable binding
	executeWrite(t, fmt.Sprintf("iriscli service disable --def-chain-id=%s --service-name=%s --from=%s --fee=0.3iris %v", chainID, serviceName, "bar", flags), sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(12, port)

	// refund deposit
	executeWrite(t, fmt.Sprintf("iriscli service refund-deposit --service-name=%s --def-chain-id=%s --from=%s --fee=0.3iris %v", serviceName, chainID, "bar", flags), sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	barAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin = convertToIrisBaseAccount(t, barAcc)
	barAmt = getAmountFromCoinStr(barCoin)

	if !(barAmt > 18 && barAmt < 20) {
		t.Error("Test Failed: (18, 20) expected, received: {}", barAmt)
	}

	// call service
	caStr := fmt.Sprintf("iriscli service call %v", flags)
	caStr += fmt.Sprintf(" --def-chain-id=%s", chainID)
	caStr += fmt.Sprintf(" --service-name=%s", serviceName)
	caStr += fmt.Sprintf(" --bind-chain-id=%s", chainID)
	caStr += fmt.Sprintf(" --method-id=%d", 1)
	caStr += fmt.Sprintf(" --provider=%s", fooAddr.String())
	caStr += fmt.Sprintf(" --request-data=%s", "1234")
	caStr += fmt.Sprintf(" --service-fee=%s", "2iris")
	caStr += fmt.Sprintf(" --fee=%s", "0.4iris")
	caStr += fmt.Sprintf(" --from=%s", "bar")
	caStr += " --commit"

	_, outString, _ := executeWriteRetStdStreams(t, caStr, sdk.DefaultKeyPass)

	var digitsRegexp = regexp.MustCompile(`\"key\": \"request-id\",\n       \"value\": \".*\"`)
	requestTag := string(digitsRegexp.Find([]byte(outString)))
	requestId := strings.TrimSpace(strings.Split(requestTag, ":")[2])
	requestId = requestId[1 : len(requestId)-1]

	tests.WaitForNextNBlocksTM(2, port)

	svcRequests := executeGetServiceRequests(t, fmt.Sprintf("iriscli service requests --def-chain-id=%s --service-name=%s --bind-chain-id=%s --provider=%s %v", chainID, serviceName, chainID, fooAddr.String(), flags))
	require.Equal(t, 1, len(svcRequests))

	barAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin = convertToIrisBaseAccount(t, barAcc)
	barAmt = getAmountFromCoinStr(barCoin)

	if !(barAmt > 17 && barAmt < 19) {
		t.Error("Test Failed: (17, 19) expected, received: {}", barAmt)
	}

	// respond service
	reStr := fmt.Sprintf("iriscli service respond %v", flags)
	reStr += fmt.Sprintf(" --request-chain-id=%s", chainID)
	reStr += fmt.Sprintf(" --request-id=%s", requestId)
	reStr += fmt.Sprintf(" --response-data=%s", "1234")
	reStr += fmt.Sprintf(" --fee=%s", "0.4iris")
	reStr += fmt.Sprintf(" --from=%s", "foo")

	executeWrite(t, reStr, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(7, port)

	// fees
	fooFees := executeGetServiceFees(t, fmt.Sprintf("iriscli service fees %s %v", fooAddr.String(), flags))
	barFees := executeGetServiceFees(t, fmt.Sprintf("iriscli service fees %s %v", barAddr.String(), flags))

	require.Equal(t, "990000000000000000iris-atto", fooFees.IncomingFee.String())
	require.Nil(t, fooFees.ReturnedFee)
	require.Nil(t, barFees.ReturnedFee)
	require.Nil(t, barFees.IncomingFee)

	executeWrite(t, caStr, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(7, port)

	fooFees = executeGetServiceFees(t, fmt.Sprintf("iriscli service fees %s %v", fooAddr.String(), flags))
	barFees = executeGetServiceFees(t, fmt.Sprintf("iriscli service fees %s %v", barAddr.String(), flags))

	require.Equal(t, "990000000000000000iris-atto", fooFees.IncomingFee.String())
	require.Nil(t, fooFees.ReturnedFee)
	require.Equal(t, "1000000000000000000iris-atto", barFees.ReturnedFee.String())
	require.Nil(t, barFees.IncomingFee)

	svcBinding = executeGetServiceBinding(t, fmt.Sprintf("iriscli service binding --service-name=%s --def-chain-id=%s --bind-chain-id=%s --provider=%s %v", serviceName, chainID, chainID, fooAddr.String(), flags))
	require.NotNil(t, svcBinding)
	require.Equal(t, "10000000000000000000iris-atto", svcBinding.Deposit.String())
	require.Equal(t, true, svcBinding.Available)

	// refund fees
	executeWrite(t, fmt.Sprintf("iriscli service refund-fees %v --fee=%s --from=%s", flags, "0.4iris", "bar"), sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	barAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin = convertToIrisBaseAccount(t, barAcc)
	barAmt = getAmountFromCoinStr(barCoin)

	if !(barAmt > 17 && barAmt < 19) {
		t.Error("Test Failed: (17, 19) expected, received: {}", barAmt)
	}

	// withdraw fees
	executeWrite(t, fmt.Sprintf("iriscli service withdraw-fees %v --fee=%s --from=%s", flags, "0.4iris", "foo"), sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	fooAmt = getAmountFromCoinStr(fooCoin)

	if !(fooAmt > 19 && fooAmt < 21) {
		t.Error("Test Failed: (19, 21) expected, received: {}", fooAmt)
	}

	// withdraw tax
	barAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))

	wtStr := fmt.Sprintf("iriscli service withdraw-tax %v", flags)
	wtStr += fmt.Sprintf(" --withdraw-amount=%s", "0.001iris")
	wtStr += fmt.Sprintf(" --dest-address=%s", barAcc.Address)
	wtStr += fmt.Sprintf(" --fee=%s", "0.4iris")
	wtStr += fmt.Sprintf(" --from=%s", "foo")

	executeWrite(t, wtStr, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	barAcc1 := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))

	cliCtx := context.NewCLIContext()
	oldAmount, _ := cliCtx.ParseCoin(barAcc.Coins[0].String())
	newAmount, _ := cliCtx.ParseCoin(barAcc1.Coins[0].String())

	tax, _ := sdk.NewIntFromString("1000000000000000")
	require.Equal(t, oldAmount.Amount.Add(tax).String(), newAmount.Amount.String())
}
