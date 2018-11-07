package clitest

import (
	"testing"
	"github.com/cosmos/cosmos-sdk/tests"
	"fmt"
	"github.com/irisnet/irishub/app"
	"io/ioutil"
	"os"
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

	serviceName := "testService"

	serviceQuery, _ := tests.ExecuteT(t, fmt.Sprintf("iriscli iservice definition --name=%s --def-chain-id=%s %v", serviceName, chainID, flags), "")
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
	sdStr += fmt.Sprintf(" --name=%s", serviceName)
	sdStr += fmt.Sprintf(" --service-description=%s", "test")
	sdStr += fmt.Sprintf(" --tags=%s", "tag1 tag2")
	sdStr += fmt.Sprintf(" --author-description=%s", "foo")
	sdStr += fmt.Sprintf(" --broadcast=%s", "Broadcast")
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

	serviceDef := executeGetServiceDefinition(t, fmt.Sprintf("iriscli iservice definition --name=%s --def-chain-id=%s %v", serviceName, chainID, flags))
	require.Equal(t, serviceName, serviceDef.Name)

	// method test
	require.Equal(t, "SayHello", serviceDef.Methods[0].Name)
	require.Equal(t, "sayHello", serviceDef.Methods[0].Description)
	require.Equal(t, "NoCached", serviceDef.Methods[0].OutputCached.String())
	require.Equal(t, "NoPrivacy", serviceDef.Methods[0].OutputPrivacy.String())
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
