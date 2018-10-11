package clitest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/tests"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/irisnet/irishub/app"
)

func init() {
	irisHome, iriscliHome = getTestingHomeDirs()
}

func TestIrisCLIBankSend(t *testing.T) {
	tests.ExecuteT(t, fmt.Sprintf("iris --home=%s unsafe_reset_all", irisHome), "")
	executeWrite(t, fmt.Sprintf("iriscli keys delete --home=%s foo", iriscliHome), app.DefaultKeyPass)
	executeWrite(t, fmt.Sprintf("iriscli keys delete --home=%s bar", iriscliHome), app.DefaultKeyPass)
	chainID := executeInit(t, fmt.Sprintf("iris init -o --name=foo --home=%s --home-client=%s", irisHome, iriscliHome))
	executeWrite(t, fmt.Sprintf("iriscli keys add --home=%s bar", iriscliHome), app.DefaultKeyPass)

	err := modifyGenesisFile(t, irisHome)
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

	fooAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show foo --output=json --home=%s", iriscliHome))
	barAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show bar --output=json --home=%s", iriscliHome))

	fooAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin := convertToIrisBaseAccount(t, fooAcc)
	require.Equal(t, "1000000000000iris", fooCoin)

	executeWrite(t, fmt.Sprintf("iriscli bank send %v --amount=10iris --to=%s --from=foo --gas=10000 --fee=0.04iris", flags, barAddr), app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	barAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin := convertToIrisBaseAccount(t, barAcc)
	require.Equal(t, "10iris", barCoin)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	num := getAmuntFromCoinStr(t, fooCoin)

	if !(num > 999999999989 && num < 999999999990) {
		t.Error("Test Failed: (999999999989, 999999999990) expected, recieved: {}", num)
	}

	// test autosequencing
	executeWrite(t, fmt.Sprintf("iriscli bank send %v --amount=10iris --to=%s --from=foo --gas=10000 --fee=0.04iris", flags, barAddr), app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	barAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin = convertToIrisBaseAccount(t, barAcc)
	require.Equal(t, "20iris", barCoin)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	num = getAmuntFromCoinStr(t, fooCoin)

	if !(num > 999999999979 && num < 999999999980) {
		t.Error("Test Failed: (999999999979, 999999999980) expected, recieved: {}", num)
	}

	// test memo
	executeWrite(t, fmt.Sprintf("iriscli bank send %v --amount=10iris --to=%s --from=foo --memo 'testmemo' --gas=10000 --fee=0.04iris", flags, barAddr), app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	barAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin = convertToIrisBaseAccount(t, barAcc)
	require.Equal(t, "30iris", barCoin)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	num = getAmuntFromCoinStr(t, fooCoin)

	if !(num > 999999999969 && num < 999999999970) {
		t.Error("Test Failed: (999999999969, 999999999970) expected, recieved: {}", num)
	}
}
