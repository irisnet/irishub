package clitest

import (
	"fmt"
	"os"
	"testing"
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/tests"
	"github.com/irisnet/irishub/app"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func init() {
	irisHome, iriscliHome = getTestingHomeDirs()
}

func TestIrisCLIIserviceDefine(t *testing.T) {
	chainID, servAddr, port := initializeFixtures(t)
	flags := fmt.Sprintf("--home=%s --node=%v --chain-id=%v", iriscliHome, servAddr, chainID)

	// start iris server
	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iris start --home=%s --rpc.laddr=%v", irisHome, servAddr))

	defer proc.Stop(false)
	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(2, port)

	fooAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show foo --output=json --home=%s", iriscliHome))
	barAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show bar --output=json --home=%s", iriscliHome))

	serviceName := "testService"

	serviceQuery, _ := tests.ExecuteT(t, fmt.Sprintf("iriscli iservice definition --service-name=%s --def-chain-id=%s %v", serviceName, chainID, flags), "")
	require.Equal(t, "", serviceQuery)

	fooAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin := convertToIrisBaseAccount(t, fooAcc)
	num := getAmountFromCoinStr(fooCoin)
	require.Equal(t, "50iris", fooCoin)

	// iservice define
	fileName := iriscliHome + string(os.PathSeparator) + "test.proto"
	defer tests.ExecuteT(t, fmt.Sprintf("rm -f %s", fileName), "")
	ioutil.WriteFile(fileName, []byte(idlContent), 0644)
	sdStr := fmt.Sprintf("iriscli iservice define %v", flags)
	sdStr += fmt.Sprintf(" --from=%s", "foo")
	sdStr += fmt.Sprintf(" --service-name=%s", serviceName)
	sdStr += fmt.Sprintf(" --service-description=%s", "test")
	sdStr += fmt.Sprintf(" --tags=%s", "tag1,tag2")
	sdStr += fmt.Sprintf(" --author-description=%s", "foo")
	sdStr += fmt.Sprintf(" --messaging=%s", "Multicast")
	sdStr += fmt.Sprintf(" --file=%s", fileName)
	sdStr += fmt.Sprintf(" --fee=%s", "0.004iris")

	executeWrite(t, sdStr, app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	num = getAmountFromCoinStr(fooCoin)

	if !(num > 49 && num < 50) {
		t.Error("Test Failed: (49, 50) expected, recieved: {}", num)
	}

	serviceDef := executeGetServiceDefinition(t, fmt.Sprintf("iriscli iservice definition --service-name=%s --def-chain-id=%s %v", serviceName, chainID, flags))
	require.Equal(t, serviceName, serviceDef.Name)

	// method test
	require.Equal(t, "SayHello", serviceDef.Methods[0].Name)
	require.Equal(t, "sayHello", serviceDef.Methods[0].Description)
	require.Equal(t, "NoCached", serviceDef.Methods[0].OutputCached.String())
	require.Equal(t, "NoPrivacy", serviceDef.Methods[0].OutputPrivacy.String())

	// binding test
	sdStr = fmt.Sprintf("iriscli iservice bind %v", flags)
	sdStr += fmt.Sprintf(" --service-name=%s", serviceName)
	sdStr += fmt.Sprintf(" --def-chain-id=%s", chainID)
	sdStr += fmt.Sprintf(" --bind-type=%s", "Local")
	sdStr += fmt.Sprintf(" --deposit=%s", "10iris")
	sdStr += fmt.Sprintf(" --prices=%s", "1iris")
	sdStr += fmt.Sprintf(" --avg-rsp-time=%d", 10000)
	sdStr += fmt.Sprintf(" --usable-time=%d", 10000)
	sdStr += fmt.Sprintf(" --expiration=%d", -1)
	sdStr += fmt.Sprintf(" --fee=%s", "0.004iris")

	sdStrFoo := sdStr + fmt.Sprintf(" --from=%s", "foo")
	sdStrBar := sdStr + fmt.Sprintf(" --from=%s", "bar")

	executeWrite(t, sdStrFoo, app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	num = getAmountFromCoinStr(fooCoin)

	if !(num > 39 && num < 40) {
		t.Error("Test Failed: (39, 40) expected, recieved: {}", num)
	}

	executeWrite(t, fmt.Sprintf("iriscli bank send --to=%s --from=%s --amount=20iris --fee=0.004iris %v", barAddr.String(), "foo", flags), app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)
	executeWrite(t, sdStrBar, app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)
	barAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin := convertToIrisBaseAccount(t, barAcc)
	barNum := getAmountFromCoinStr(barCoin)

	if !(barNum > 9 && barNum < 10) {
		t.Error("Test Failed: (9, 10) expected, recieved: {}", barNum)
	}

	serviceBinding := executeGetServiceBinding(t, fmt.Sprintf("iriscli iservice binding --service-name=%s --def-chain-id=%s --bind-chain-id=%s --provider=%s %v", serviceName, chainID, chainID, fooAddr.String(), flags))
	require.NotNil(t, serviceBinding)

	serviceBindings := executeGetServiceBindings(t, fmt.Sprintf("iriscli iservice bindings --service-name=%s --def-chain-id=%s %v", serviceName, chainID, flags))
	require.Equal(t, 2, len(serviceBindings))

	// binding update test
	sdStr = fmt.Sprintf("iriscli iservice update-binding %v", flags)
	sdStr += fmt.Sprintf(" --service-name=%s", serviceName)
	sdStr += fmt.Sprintf(" --def-chain-id=%s", chainID)
	sdStr += fmt.Sprintf(" --bind-type=%s", "Global")
	sdStr += fmt.Sprintf(" --deposit=%s", "1iris")
	sdStr += fmt.Sprintf(" --prices=%s", "5iris")
	sdStr += fmt.Sprintf(" --avg-rsp-time=%d", 99)
	sdStr += fmt.Sprintf(" --usable-time=%d", 99)
	sdStr += fmt.Sprintf(" --expiration=%d", 1)
	sdStr += fmt.Sprintf(" --fee=%s", "0.004iris")
	sdStr += fmt.Sprintf(" --from=%s", "bar")
	executeWrite(t, sdStr, app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)
	barAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin = convertToIrisBaseAccount(t, barAcc)
	barNum = getAmountFromCoinStr(barCoin)

	if !(barNum > 8 && barNum < 9) {
		t.Error("Test Failed: (8, 9) expected, recieved: {}", barNum)
	}
	serviceBindings = executeGetServiceBindings(t, fmt.Sprintf("iriscli iservice bindings --service-name=%s --def-chain-id=%s %v", serviceName, chainID, flags))
	var totalDeposit sdk.Coins
	for _, bind := range serviceBindings {
		totalDeposit = totalDeposit.Plus(bind.Deposit)
	}
	require.Equal(t, "21000000000000000000iris-atto", totalDeposit.String())

	// refund-deposit test
	tests.WaitForNextNBlocksTM(8, port)
	executeWrite(t, fmt.Sprintf("iriscli iservice refund-deposit --service-name=%s --def-chain-id=%s --from=%s --fee=0.004iris %v", serviceName, chainID, "bar", flags), app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)
	barAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin = convertToIrisBaseAccount(t, barAcc)
	barNum = getAmountFromCoinStr(barCoin)
	if !(barNum > 19 && barNum < 20) {
		t.Error("Test Failed: (19, 20) expected, recieved: {}", barNum)
	}
}

const idlContent = `
	syntax = "proto3";

	// The greeting service definition.
	service Greeter {
		//@Attribute description:sayHello
		//@Attribute output_privacy:NoPrivacy
		//@Attribute output_cached:NoCached
		rpc SayHello (HelloRequest) returns (HelloReply) {}
	}

	// The request message containing the user's name.
	message HelloRequest {
		string name = 1;
	}

	// The response message containing the greetings
	message HelloReply {
		string message = 1;
	}`
