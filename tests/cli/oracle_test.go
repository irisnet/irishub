package cli

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/irisnet/irishub/app/v3/service"
	"github.com/irisnet/irishub/tests"
	sdk "github.com/irisnet/irishub/types"
)

func TestIrisCLIOracle(t *testing.T) {
	t.Parallel()
	chainID, servAddr, port, irisHome, iriscliHome, p2pAddr := initializeFixtures(t)

	flags := fmt.Sprintf("--home=%s --node=%v --chain-id=%v --output=json", iriscliHome, servAddr, chainID)

	// start iris server
	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iris start --home=%s --rpc.laddr=%v --p2p.laddr=%v", irisHome, servAddr, p2pAddr))

	defer func() {
		_ = proc.Stop(false)
	}()
	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(2, port)

	fooAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show foo --output=json --home=%s", iriscliHome))
	_, _ = executeGetAddrPK(t, fmt.Sprintf("iriscli keys show bar --output=json --home=%s", iriscliHome))

	fooAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin := convertToIrisBaseAccount(t, fooAcc)
	require.Equal(t, "50iris", fooCoin)

	// testing variables
	feedName := "test-feed"
	description := "feed-usdt"
	aggregateFunc := "avg"
	valueJsonPath := "last"
	latestHistory := uint64(2)

	serviceName := "test"
	serviceDesc := "test"
	serviceTags := []string{"tag1", "tag2"}
	authorDesc := "author"
	serviceSchemas := `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`
	deposit := "10iris"
	priceAmt := 1 // 1iris
	pricing := fmt.Sprintf(`{"price":[{"denom":"iris-atto","amount":"%s"}]}`, sdk.NewIntWithDecimal(int64(priceAmt), 18).String())
	serviceFeeCap := "10iris"
	input := `{"pair":"iris-usdt"}`
	timeout := int64(5)
	repeatedFreq := uint64(10)
	repeatedTotal := int64(1)
	responseThreshold := uint16(1)
	result := `{"code":200,"message":""}`
	output := `{"last":100}`

	// define service
	svcDefOutput, _ := tests.ExecuteT(t, fmt.Sprintf("iriscli service definition %s %v", serviceName, flags), "")
	require.Equal(t, "", svcDefOutput)

	defineService(t, flags, serviceName, serviceDesc, serviceTags, authorDesc, serviceSchemas, port)

	// bind service
	_ = bindService(t, flags, serviceName, deposit, pricing, port)

	//create feed
	cfCmdStr := fmt.Sprintf("iriscli oracle create %v", flags)
	cfCmdStr += fmt.Sprintf(" --from=%s", "foo")
	cfCmdStr += fmt.Sprintf(" --fee=%s", "0.4iris")
	cfCmdStr += fmt.Sprintf(" --feed-name=%s", feedName)
	cfCmdStr += fmt.Sprintf(" --description=%s", description)
	cfCmdStr += fmt.Sprintf(" --latest-history=%d", latestHistory)
	cfCmdStr += fmt.Sprintf(" --service-name=%s", serviceName)
	cfCmdStr += fmt.Sprintf(" --input=%s", input)
	cfCmdStr += fmt.Sprintf(" --providers=%s", fooAddr.String())
	cfCmdStr += fmt.Sprintf(" --service-fee-cap=%s", serviceFeeCap)
	cfCmdStr += fmt.Sprintf(" --timeout=%d", timeout)
	cfCmdStr += fmt.Sprintf(" --frequency=%d", repeatedFreq)
	cfCmdStr += fmt.Sprintf(" --total=%d", repeatedTotal)
	cfCmdStr += fmt.Sprintf(" --threshold=%d", responseThreshold)
	cfCmdStr += fmt.Sprintf(" --aggregate-func=%s", aggregateFunc)
	cfCmdStr += fmt.Sprintf(" --value-json-path=%s", valueJsonPath)
	cfCmdStr += " --commit"

	executeWrite(t, cfCmdStr, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(1, port)

	//query feed
	qfCmdStr := fmt.Sprintf("iriscli oracle query-feed %s %v", feedName, flags)
	feed := executeGetFeed(t, qfCmdStr)
	require.Equal(t, feedName, feed.Feed.FeedName)
	require.Equal(t, description, feed.Feed.Description)
	require.Equal(t, aggregateFunc, feed.Feed.AggregateFunc)
	require.Equal(t, valueJsonPath, feed.Feed.ValueJsonPath)
	require.Equal(t, latestHistory, feed.Feed.LatestHistory)
	require.Equal(t, fooAddr, feed.Feed.Creator)
	require.Equal(t, serviceName, feed.ServiceName)
	require.EqualValues(t, []sdk.AccAddress{fooAddr}, feed.Providers)
	require.Equal(t, input, feed.Input)
	require.Equal(t, timeout, feed.Timeout)
	require.Equal(t, serviceFeeCap, feed.ServiceFeeCap.MainUnitString())
	require.Equal(t, repeatedFreq, feed.RepeatedFrequency)
	require.Equal(t, repeatedTotal, feed.RepeatedTotal)
	require.Equal(t, responseThreshold, feed.ResponseThreshold)
	require.Equal(t, service.PAUSED, feed.State)

	//edit feed
	description = "feed-eth"
	//TODO
	efCmdStr := fmt.Sprintf("iriscli oracle edit %s %v", feedName, flags)
	efCmdStr += fmt.Sprintf(" --from=%s", "foo")
	efCmdStr += fmt.Sprintf(" --fee=%s", "0.4iris")
	efCmdStr += fmt.Sprintf(" --description=%s", description)
	executeWrite(t, efCmdStr, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(1, port)
	//query feed
	feed = executeGetFeed(t, qfCmdStr)
	require.Equal(t, feedName, feed.Feed.FeedName)
	require.Equal(t, description, feed.Feed.Description)

	//start feed
	sfCmdStr := fmt.Sprintf("iriscli oracle start %s %v", feedName, flags)
	sfCmdStr += fmt.Sprintf(" --from=%s", "foo")
	sfCmdStr += fmt.Sprintf(" --fee=%s", "0.4iris")
	sfCmdStr += " --commit"
	executeWrite(t, sfCmdStr, sdk.DefaultKeyPass)

	for {
		// query requests by binding (foo)
		fooRequests := executeGetServiceRequests(t, fmt.Sprintf("iriscli service requests %s %s %v", serviceName, fooAddr.String(), flags))
		if len(fooRequests) == 1 {
			//respond service
			respondService(t, fooRequests[0].ID.String(), flags, result, output, port)
			goto verifyValue
		}
		tests.WaitForNextNBlocksTM(1, port)
	}

verifyValue:
	//query value
	qvCmdStr := fmt.Sprintf("iriscli oracle query-value %s %v", feedName, flags)
	values := executeGetFeedValue(t, qvCmdStr)
	require.Equal(t, 1, len(values))
	require.Equal(t, "100.00000000", values[0].Data)
}

func respondService(t *testing.T, requestID string, flags string, result string, output string, port string) (string, string) {
	rsStr := fmt.Sprintf("iriscli service respond %v", flags)
	rsStr += fmt.Sprintf(" --request-id=%s", requestID)
	rsStr += fmt.Sprintf(" --result=%s", result)
	rsStr += fmt.Sprintf(" --data=%s", output)
	rsStr += fmt.Sprintf(" --fee=%s", "0.4iris")
	rsStr += fmt.Sprintf(" --from=%s", "foo")

	executeWrite(t, rsStr, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(1, port)
	return requestID, rsStr
}

func bindService(t *testing.T, flags string, serviceName string, deposit string, pricing string, port string) string {
	sbStr := fmt.Sprintf("iriscli service bind %v", flags)
	sbStr += fmt.Sprintf(" --service-name=%s", serviceName)
	sbStr += fmt.Sprintf(" --deposit=%s", deposit)
	sbStr += fmt.Sprintf(" --pricing=%s", pricing)
	sbStr += fmt.Sprintf(" --fee=%s", "0.4iris")

	sbStrFoo := sbStr + fmt.Sprintf(" --from=%s", "foo")
	sbStrBar := sbStr + fmt.Sprintf(" --from=%s", "bar")

	executeWrite(t, sbStrFoo, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)
	return sbStrBar
}

func defineService(t *testing.T, flags string, serviceName string, serviceDesc string, serviceTags []string, authorDesc string, serviceSchemas string, port string) {
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
}
