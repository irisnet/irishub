package cli

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/irisnet/irishub/app/v3/service"
	"github.com/irisnet/irishub/tests"
	sdk "github.com/irisnet/irishub/types"
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

	// testing variables
	serviceName := "test"
	serviceDesc := "test"
	serviceTags := []string{"tag1", "tag2"}
	authorDesc := "author"
	serviceSchemas := `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`
	deposit := "10iris"
	priceAmt := 1 // 1iris
	price := fmt.Sprintf("%diris", priceAmt)
	pricing := fmt.Sprintf(`{"price":[{"denom":"iris-atto","amount":"%s"}]}`, sdk.NewIntWithDecimal(int64(priceAmt), 18).String())
	addedDeposit := "1iris"
	serviceFeeCap := "10iris"
	input := `{"pair":"iris-usdt"}`
	timeout := int64(5)
	repeatedFreq := uint64(20)
	repeatedTotal := int64(10)
	output := `{"last":100}`

	// define service
	svcDefOutput, _ := tests.ExecuteT(t, fmt.Sprintf("iriscli service definition %s %v", serviceName, flags), "")
	require.Equal(t, "", svcDefOutput)

	sdStr := fmt.Sprintf("iriscli service define %v", flags)
	sdStr += fmt.Sprintf(" --from=%s", "foo")
	sdStr += fmt.Sprintf(" --name=%s", serviceName)
	sdStr += fmt.Sprintf(" --description=%s", serviceDesc)
	sdStr += fmt.Sprintf(" --tags=%s", serviceTags)
	sdStr += fmt.Sprintf(" --author-description=%s", authorDesc)
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
	sbStr := fmt.Sprintf("iriscli service bind %v", flags)
	sbStr += fmt.Sprintf(" --service-name=%s", serviceName)
	sbStr += fmt.Sprintf(" --deposit=%s", deposit)
	sbStr += fmt.Sprintf(" --pricing=%s", pricing)
	sbStr += fmt.Sprintf(" --fee=%s", "0.4iris")

	sbStrFoo := sbStr + fmt.Sprintf(" --from=%s", "foo")
	sbStrBar := sbStr + fmt.Sprintf(" --from=%s", "bar")

	executeWrite(t, sbStrFoo, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	fooAmt = getAmountFromCoinStr(fooCoin)

	if !(fooAmt > 39 && fooAmt < 40) {
		t.Error("Test Failed: (39, 40) expected, received: {}", fooAmt)
	}

	executeWrite(t, fmt.Sprintf("iriscli bank send --to=%s --from=%s --amount=20iris --fee=0.3iris %v", barAddr.String(), "foo", flags), sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	executeWrite(t, sbStrBar, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	barAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin := convertToIrisBaseAccount(t, barAcc)
	barAmt := getAmountFromCoinStr(barCoin)

	if !(barAmt > 9 && barAmt < 10) {
		t.Error("Test Failed: (9, 10) expected, received: {}", barAmt)
	}

	svcBinding := executeGetServiceBinding(t, fmt.Sprintf("iriscli service binding %s %s %v", serviceName, fooAddr.String(), flags))
	require.Equal(t, serviceName, svcBinding.ServiceName)
	require.Equal(t, fooAddr, svcBinding.Provider)
	require.Equal(t, deposit, svcBinding.Deposit.MainUnitString())
	require.Equal(t, pricing, svcBinding.Pricing)
	require.True(t, svcBinding.Available)

	svcBindings := executeGetServiceBindings(t, fmt.Sprintf("iriscli service bindings %s %v", serviceName, flags))
	require.Equal(t, 2, len(svcBindings))

	// update binding
	ubStr := fmt.Sprintf("iriscli service update-binding %s %v", serviceName, flags)
	ubStr += fmt.Sprintf(" --deposit=%s", addedDeposit)
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

	svcBindings = executeGetServiceBindings(t, fmt.Sprintf("iriscli service bindings %s %v", serviceName, flags))

	var totalDeposit sdk.Coins
	for _, binding := range svcBindings {
		totalDeposit = totalDeposit.Add(binding.Deposit)
	}
	require.Equal(t, "21iris", totalDeposit.MainUnitString())

	// set withdrawal address
	swStr := fmt.Sprintf("iriscli service set-withdraw-addr %s %v", fooAddr.String(), flags)
	swStr += fmt.Sprintf(" --fee=%s", "0.4iris")
	swStr += fmt.Sprintf(" --from=%s", "bar")

	executeWrite(t, swStr, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	withdrawAddr := executeGetWithdrawAddress(t, fmt.Sprintf("iriscli service withdraw-addr %s %v", barAddr.String(), flags))
	require.Equal(t, fooAddr, withdrawAddr)

	// disable binding
	executeWrite(t, fmt.Sprintf("iriscli service disable %s --from=%s --fee=0.3iris %v", serviceName, "bar", flags), sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	svcBinding = executeGetServiceBinding(t, fmt.Sprintf("iriscli service binding %s %s %v", serviceName, barAddr.String(), flags))
	require.False(t, svcBinding.Available)
	require.False(t, svcBinding.DisabledTime.IsZero())

	// refund deposit
	tests.WaitForNextNBlocksTM(10, port)

	executeWrite(t, fmt.Sprintf("iriscli service refund-deposit %s --from=%s --fee=0.3iris %v", serviceName, "bar", flags), sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	svcBinding = executeGetServiceBinding(t, fmt.Sprintf("iriscli service binding %s %s %v", serviceName, barAddr.String(), flags))
	require.Equal(t, sdk.Coins(nil), svcBinding.Deposit)

	barAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin = convertToIrisBaseAccount(t, barAcc)
	barAmt = getAmountFromCoinStr(barCoin)

	if !(barAmt > 19 && barAmt < 20) {
		t.Error("Test Failed: (19, 20) expected, received: {}", barAmt)
	}

	// enable binding
	executeWrite(t, fmt.Sprintf("iriscli service enable %s --from=%s --fee=0.3iris --deposit=%s %v", serviceName, "bar", deposit, flags), sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	svcBinding = executeGetServiceBinding(t, fmt.Sprintf("iriscli service binding %s %s %v", serviceName, barAddr.String(), flags))
	require.True(t, svcBinding.Available)
	require.True(t, svcBinding.DisabledTime.IsZero())
	require.Equal(t, deposit, svcBinding.Deposit.MainUnitString())

	// call service
	csStr := fmt.Sprintf("iriscli service call %v", flags)
	csStr += fmt.Sprintf(" --service-name=%s", serviceName)
	csStr += fmt.Sprintf(" --providers=%s,%s", fooAddr.String(), barAddr.String())
	csStr += fmt.Sprintf(" --service-fee-cap=%s", serviceFeeCap)
	csStr += fmt.Sprintf(" --data=%s", input)
	csStr += fmt.Sprintf(" --timeout=%d", timeout)
	csStr += fmt.Sprintf(" --repeated=%v", true)
	csStr += fmt.Sprintf(" --frequency=%d", repeatedFreq)
	csStr += fmt.Sprintf(" --total=%d", repeatedTotal)
	csStr += fmt.Sprintf(" --fee=%s", "0.4iris")
	csStr += fmt.Sprintf(" --from=%s", "bar")
	csStr += " --commit"

	_, outString, _ := executeWriteRetStdStreams(t, csStr, sdk.DefaultKeyPass)

	var regExp = regexp.MustCompile(`\"key\": \"request-context-id\",\n.*?\"value\": \"(.*)\"`)
	requestContextID := string(regExp.FindSubmatch([]byte(outString))[1])

	tests.WaitForNextNBlocksTM(1, port)

	// query request context
	requestContext := executeGetRequestContext(t, fmt.Sprintf("iriscli service request-context %s %v", requestContextID, flags))

	require.Equal(t, serviceName, requestContext.ServiceName)
	require.Equal(t, []sdk.AccAddress{fooAddr, barAddr}, requestContext.Providers)
	require.Equal(t, barAddr, requestContext.Consumer)
	require.Equal(t, input, requestContext.Input)
	require.Equal(t, serviceFeeCap, requestContext.ServiceFeeCap.MainUnitString())
	require.Equal(t, timeout, requestContext.Timeout)
	require.Equal(t, false, requestContext.SuperMode)
	require.Equal(t, true, requestContext.Repeated)
	require.Equal(t, repeatedFreq, requestContext.RepeatedFrequency)
	require.Equal(t, repeatedTotal, requestContext.RepeatedTotal)
	require.Equal(t, uint64(1), requestContext.BatchCounter)
	require.Equal(t, uint16(2), requestContext.BatchRequestCount)
	require.Equal(t, uint16(0), requestContext.BatchResponseCount)
	require.Equal(t, service.BATCHRUNNING, requestContext.BatchState)
	require.Equal(t, service.RUNNING, requestContext.State)
	require.Equal(t, uint16(0), requestContext.ResponseThreshold)
	require.Equal(t, "", requestContext.ModuleName)

	// query requests by binding (foo)
	fooRequests := executeGetServiceRequests(t, fmt.Sprintf("iriscli service requests %s %s %v", serviceName, fooAddr.String(), flags))
	require.Equal(t, 1, len(fooRequests))

	require.Equal(t, serviceName, fooRequests[0].ServiceName)
	require.Equal(t, fooAddr, fooRequests[0].Provider)
	require.Equal(t, barAddr, fooRequests[0].Consumer)
	require.Equal(t, input, fooRequests[0].Input)
	require.Equal(t, price, fooRequests[0].ServiceFee.MainUnitString())
	require.Equal(t, false, fooRequests[0].SuperMode)
	require.Equal(t, requestContextID, hex.EncodeToString(fooRequests[0].RequestContextID))
	require.Equal(t, uint64(1), fooRequests[0].RequestContextBatchCounter)

	// query requests by binding (bar)
	barRequests := executeGetServiceRequests(t, fmt.Sprintf("iriscli service requests %s %s %v", serviceName, barAddr.String(), flags))
	require.Equal(t, 1, len(barRequests))

	require.Equal(t, serviceName, barRequests[0].ServiceName)
	require.Equal(t, barAddr, barRequests[0].Provider)
	require.Equal(t, barAddr, barRequests[0].Consumer)
	require.Equal(t, input, barRequests[0].Input)
	require.Equal(t, price, barRequests[0].ServiceFee.MainUnitString())
	require.Equal(t, false, barRequests[0].SuperMode)
	require.Equal(t, requestContextID, hex.EncodeToString(barRequests[0].RequestContextID))
	require.Equal(t, uint64(1), barRequests[0].RequestContextBatchCounter)

	// query requests by request context
	requests := executeGetServiceRequests(t, fmt.Sprintf("iriscli service requests %s %d %v", requestContextID, 1, flags))
	require.Equal(t, 2, len(requests))

	require.Equal(t, fooRequests[0], requests[0])
	require.Equal(t, barRequests[0], requests[1])

	// respond service (foo)

	fooRequestID := requestContextID + hex.EncodeToString(sdk.Uint64ToBigEndian(1)) + "0000"

	rsStr := fmt.Sprintf("iriscli service respond %v", flags)
	rsStr += fmt.Sprintf(" --request-id=%s", fooRequestID)
	rsStr += fmt.Sprintf(" --data=%s", output)
	rsStr += fmt.Sprintf(" --fee=%s", "0.4iris")
	rsStr += fmt.Sprintf(" --from=%s", "foo")

	executeWrite(t, rsStr, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(1, port)

	// query response (foo)
	fooResponse := executeGetServiceResponse(t, fmt.Sprintf("iriscli service response %s %v", fooRequestID, flags))

	require.Equal(t, fooAddr, fooResponse.Provider)
	require.Equal(t, barAddr, fooResponse.Consumer)
	require.Equal(t, output, fooResponse.Output)
	require.Equal(t, "", fooResponse.Error)
	require.Equal(t, requestContextID, hex.EncodeToString(fooResponse.RequestContextID))
	require.Equal(t, uint64(1), fooResponse.RequestContextBatchCounter)

	// query request context
	requestContext = executeGetRequestContext(t, fmt.Sprintf("iriscli service request-context %s %v", requestContextID, flags))

	require.Equal(t, uint64(1), requestContext.BatchCounter)
	require.Equal(t, uint16(2), requestContext.BatchRequestCount)
	require.Equal(t, uint16(1), requestContext.BatchResponseCount)
	require.Equal(t, service.BATCHRUNNING, requestContext.BatchState)
	require.Equal(t, service.RUNNING, requestContext.State)

	// respond service (bar)

	barRequestID := requestContextID + hex.EncodeToString(sdk.Uint64ToBigEndian(1)) + "0001"

	rsStr = fmt.Sprintf("iriscli service respond %v", flags)
	rsStr += fmt.Sprintf(" --request-id=%s", barRequestID)
	rsStr += fmt.Sprintf(" --data=%s", output)
	rsStr += fmt.Sprintf(" --fee=%s", "0.4iris")
	rsStr += fmt.Sprintf(" --from=%s", "bar")

	executeWrite(t, rsStr, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(1, port)

	// query response (bar)
	barResponse := executeGetServiceResponse(t, fmt.Sprintf("iriscli service response %s %v", barRequestID, flags))

	require.Equal(t, barAddr, barResponse.Provider)
	require.Equal(t, barAddr, barResponse.Consumer)
	require.Equal(t, output, barResponse.Output)
	require.Equal(t, "", barResponse.Error)
	require.Equal(t, requestContextID, hex.EncodeToString(barResponse.RequestContextID))
	require.Equal(t, uint64(1), barResponse.RequestContextBatchCounter)

	// query request context
	requestContext = executeGetRequestContext(t, fmt.Sprintf("iriscli service request-context %s %v", requestContextID, flags))

	require.Equal(t, uint64(1), requestContext.BatchCounter)
	require.Equal(t, uint16(2), requestContext.BatchRequestCount)
	require.Equal(t, uint16(2), requestContext.BatchResponseCount)
	require.Equal(t, service.BATCHCOMPLETED, requestContext.BatchState)
	require.Equal(t, service.RUNNING, requestContext.State)

	// query responses by request context
	responses := executeGetServiceResponses(t, fmt.Sprintf("iriscli service responses %s %d %v", requestContextID, 1, flags))
	require.Equal(t, 2, len(responses))

	require.Equal(t, fooResponse, responses[0])
	require.Equal(t, barResponse, responses[1])

	// pause the request context
	prcStr := fmt.Sprintf("iriscli service pause %s %v", requestContextID, flags)
	prcStr += fmt.Sprintf(" --fee=%s", "0.4iris")
	prcStr += fmt.Sprintf(" --from=%s", "bar")

	executeWrite(t, prcStr, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(1, port)

	// query request context
	requestContext = executeGetRequestContext(t, fmt.Sprintf("iriscli service request-context %s %v", requestContextID, flags))

	require.Equal(t, service.PAUSED, requestContext.State)

	// start the request context
	srcStr := fmt.Sprintf("iriscli service start %s %v", requestContextID, flags)
	srcStr += fmt.Sprintf(" --fee=%s", "0.4iris")
	srcStr += fmt.Sprintf(" --from=%s", "bar")

	executeWrite(t, srcStr, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(1, port)

	// query request context
	requestContext = executeGetRequestContext(t, fmt.Sprintf("iriscli service request-context %s %v", requestContextID, flags))

	require.Equal(t, service.RUNNING, requestContext.State)

	// kill the request context
	krcStr := fmt.Sprintf("iriscli service kill %s %v", requestContextID, flags)
	krcStr += fmt.Sprintf(" --fee=%s", "0.4iris")
	krcStr += fmt.Sprintf(" --from=%s", "bar")

	executeWrite(t, krcStr, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(1, port)

	// query request context
	requestContext = executeGetRequestContext(t, fmt.Sprintf("iriscli service request-context %s %v", requestContextID, flags))

	require.Equal(t, service.COMPLETED, requestContext.State)

	// query fees(bar)
	fees := executeGetServiceFees(t, fmt.Sprintf("iriscli service fees %s %v", barAddr.String(), flags))

	earnedFeesAmt := float64(priceAmt) * (1 - 0.01)
	require.Equal(t, fmt.Sprintf("%giris", earnedFeesAmt), fees.Coins.MainUnitString())

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	oldFooAmt := getAmountFromCoinStr(fooCoin)

	// withdraw fees(bar's withdrawal address is foo)
	executeWrite(t, fmt.Sprintf("iriscli service withdraw-fees %v --fee=0.4iris --from=%s", flags, "bar"), sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(1, port)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	newFooAmt := getAmountFromCoinStr(fooCoin)

	require.Equal(t, oldFooAmt+earnedFeesAmt, newFooAmt)

	// withdraw tax
	barAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin = convertToIrisBaseAccount(t, barAcc)
	oldBarAmt := getAmountFromCoinStr(barCoin)

	taxAmt := float64(priceAmt) * 0.01

	wtStr := fmt.Sprintf("iriscli service withdraw-tax %s %s %v", barAddr.String(), fmt.Sprintf("%giris", taxAmt), flags)
	wtStr += fmt.Sprintf(" --fee=%s", "0.4iris")
	wtStr += fmt.Sprintf(" --from=%s", "foo")

	executeWrite(t, wtStr, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	newBarAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin = convertToIrisBaseAccount(t, newBarAcc)
	newBarAmt := getAmountFromCoinStr(barCoin)

	require.Equal(t, oldBarAmt+taxAmt, newBarAmt)
}
