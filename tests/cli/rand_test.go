package cli

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/irisnet/irishub/app/v3/rand"
	"github.com/irisnet/irishub/tests"
	sdk "github.com/irisnet/irishub/types"
)

func TestIrisCLIRand(t *testing.T) {
	t.Parallel()
	chainID, servAddr, port, irisHome, iriscliHome, p2pAddr := initializeFixtures(t)

	flags := fmt.Sprintf("--home=%s --node=%v --chain-id=%v --output=json", iriscliHome, servAddr, chainID)

	// start iris server
	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iris start --home=%s --rpc.laddr=%v --p2p.laddr=%v", irisHome, servAddr, p2pAddr))

	defer func() { _ = proc.Stop(false) }()
	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(2, port)

	fooAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show foo --output=json --home=%s", iriscliHome))
	barAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show bar --output=json --home=%s", iriscliHome))

	fooAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin := convertToIrisBaseAccount(t, fooAcc)
	require.Equal(t, "50iris", fooCoin)

	executeWrite(t, fmt.Sprintf("iriscli bank send --to=%s --from=%s --amount=20iris --fee=0.3iris %v", barAddr.String(), "foo", flags), sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	barAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin := convertToIrisBaseAccount(t, barAcc)
	require.Equal(t, "20iris", barCoin)

	// service data
	serviceName := "random"
	serviceDesc := "random"
	serviceTags := []string{"tag1", "tag2"}
	authorDesc := "author"
	serviceSchemas := `{"input":{"type":"object","properties":{}},"output":{"type":"object","properties":{"seed":{"description":"seed","type":"string","pattern":"^[0-9a-fA-F]{64}$"}}},"error":{"type":"string"}}`
	deposit := "10iris"
	priceAmt := 1 // 1iris
	pricing := fmt.Sprintf(`{"price":[{"denom":"iris-atto","amount":"%s"}]}`, sdk.NewIntWithDecimal(int64(priceAmt), 18).String())
	price := fmt.Sprintf("%diris", priceAmt)
	input := `{}`
	// addedDeposit := "1iris"
	// timeout := int64(5)
	result := `{"code":200,"message":""}`
	output := `{"seed":"3132333435363738393031323334353637383930313233343536373839303132"}`

	// random test data
	blockInterval := int64(5)
	oracle := "true"
	serviceFeeCap := "10iris"

	// define service (foo)
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

	svcDef := executeGetServiceDefinition(t, fmt.Sprintf("iriscli service definition %s %v", serviceName, flags))
	require.Equal(t, serviceName, svcDef.Name)
	require.Equal(t, serviceSchemas, svcDef.Schemas)

	// bind service (foo)
	sbStr := fmt.Sprintf("iriscli service bind %v", flags)
	sbStr += fmt.Sprintf(" --service-name=%s", serviceName)
	sbStr += fmt.Sprintf(" --deposit=%s", deposit)
	sbStr += fmt.Sprintf(" --pricing=%s", pricing)
	sbStr += fmt.Sprintf(" --fee=%s", "0.4iris")
	sbStr += sbStr + fmt.Sprintf(" --from=%s", "foo")

	executeWrite(t, sbStr, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	svcBinding := executeGetServiceBinding(t, fmt.Sprintf("iriscli service binding %s %s %v", serviceName, fooAddr.String(), flags))
	require.Equal(t, serviceName, svcBinding.ServiceName)
	require.Equal(t, fooAddr, svcBinding.Provider)
	require.Equal(t, deposit, svcBinding.Deposit.MainUnitString())
	require.Equal(t, pricing, svcBinding.Pricing)
	require.True(t, svcBinding.Available)

	svcBindings := executeGetServiceBindings(t, fmt.Sprintf("iriscli service bindings %s %v", serviceName, flags))
	require.Equal(t, 1, len(svcBindings))

	// request random (bar)
	rrStr := fmt.Sprintf("iriscli rand request-rand %v", flags)
	rrStr += fmt.Sprintf(" --from=%s", "bar")
	rrStr += fmt.Sprintf(" --fee=%s", "0.4iris")
	rrStr += fmt.Sprintf(" --block-interval=%d", blockInterval)
	rrStr += fmt.Sprintf(" --oracle=%s", oracle)
	rrStr += fmt.Sprintf(" --service-fee-cap=%s", serviceFeeCap)
	rrStr += " --commit"

	success, out, _ := executeWriteRetStdStreams(t, rrStr, sdk.DefaultKeyPass)
	require.True(t, success)

	var regExp = regexp.MustCompile(`\"key\": \"rand-height\",\n.*?\"value\": \"(.*)\"`)
	heightString := string(regExp.FindSubmatch([]byte(out))[1])
	height, err := strconv.ParseInt(heightString, 10, 64)
	require.NoError(t, err)

	tests.WaitForNextNBlocksTM(2, port)

	// query rand requests by height
	randRequests := executeGetRandRequests(t, fmt.Sprintf("iriscli rand query-queue --queue-height=%d %s %v", height, fooAddr.String(), flags))
	require.Equal(t, 1, len(randRequests))

	randReqID := hex.EncodeToString(rand.GenerateRequestID(randRequests[0]))
	tests.WaitForNextNBlocksTM(5, port)

	// query service requests by binding
	serviceRequests := executeGetServiceRequests(t, fmt.Sprintf("iriscli service requests %s %s %v", serviceName, fooAddr.String(), flags))
	require.Equal(t, 1, len(serviceRequests))
	require.Equal(t, serviceName, serviceRequests[0].ServiceName)
	require.Equal(t, fooAddr, serviceRequests[0].Provider)
	require.Equal(t, barAddr, serviceRequests[0].Consumer)
	require.Equal(t, input, serviceRequests[0].Input)
	require.Equal(t, price, serviceRequests[0].ServiceFee.MainUnitString())
	require.Equal(t, false, serviceRequests[0].SuperMode)
	require.Equal(t, uint64(1), serviceRequests[0].RequestContextBatchCounter)

	// respond service request (foo)
	reqID := hex.EncodeToString(serviceRequests[0].ID)

	srStr := fmt.Sprintf("iriscli service respond %v", flags)
	srStr += fmt.Sprintf(" --from=%s", "foo")
	srStr += fmt.Sprintf(" --fee=%s", "0.4iris")
	srStr += fmt.Sprintf(" --request-id=%s", reqID)
	srStr += fmt.Sprintf(" --result=%s", result)
	srStr += fmt.Sprintf(" --data=%s", output)

	success = executeWrite(t, srStr, sdk.DefaultKeyPass)
	require.True(t, success)

	tests.WaitForNextNBlocksTM(2, port)

	// query random by request id
	rand := executeGetRand(t, fmt.Sprintf("iriscli rand query-rand --request-id=%s %v", randReqID, flags))
	require.True(t, len(rand.RequestTxHash) > 0)
	require.True(t, rand.Height > 0)
	require.True(t, rand.Value.String() != "")
}
