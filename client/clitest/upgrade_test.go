package clitest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/tests"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/irisnet/irishub/app"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func init() {
	irisHome, iriscliHome = getTestingHomeDirs()
}

func TestIrisCLISoftwareUpgrade(t *testing.T) {
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
	barAddr, barPubKey := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show bar --output=json --home=%s", iriscliHome))
	barCeshPubKey := sdk.MustBech32ifyValPub(barPubKey)

	executeWrite(t, fmt.Sprintf("iriscli bank send %v --amount=10iris --to=%s --from=foo --gas=10000 --fee=0.04iris", flags, barAddr), app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	barAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin := convertToIrisBaseAccount(t, barAcc)
	require.Equal(t, "10iris", barCoin)

	fooAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin := convertToIrisBaseAccount(t, fooAcc)
	num := getAmuntFromCoinStr(t, fooCoin)

	if !(num > 999999999989 && num < 999999999990) {
		t.Error("Test Failed: (999999999989, 999999999990) expected, recieved: {}", num)
	}

	// create validator
	cvStr := fmt.Sprintf("iriscli stake create-validator %v", flags)
	cvStr += fmt.Sprintf(" --from=%s", "bar")
	cvStr += fmt.Sprintf(" --pubkey=%s", barCeshPubKey)
	cvStr += fmt.Sprintf(" --amount=%v", "2iris")
	cvStr += fmt.Sprintf(" --moniker=%v", "bar-vally")
	cvStr += fmt.Sprintf(" --fee=%s", "0.004iris")

	executeWrite(t, cvStr, app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	barAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin = convertToIrisBaseAccount(t, barAcc)
	num = getAmuntFromCoinStr(t, barCoin)

	if !(num > 7 && num < 8) {
		t.Error("Test Failed: (7, 8) expected, recieved: {}", num)
	}

	validator := executeGetValidator(t, fmt.Sprintf("iriscli stake validator %s --output=json %v", barAddr, flags))
	require.Equal(t, validator.Owner, barAddr)
	require.Equal(t, "2.0000000000", validator.Tokens)

	// unbond a single share
	unbondStr := fmt.Sprintf("iriscli stake unbond begin %v", flags)
	unbondStr += fmt.Sprintf(" --from=%s", "bar")
	unbondStr += fmt.Sprintf(" --address-validator=%s", barAddr)
	unbondStr += fmt.Sprintf(" --shares-amount=%v", "1")
	unbondStr += fmt.Sprintf(" --fee=%s", "0.004iris")

	success := executeWrite(t, unbondStr, app.DefaultKeyPass)
	require.True(t, success)
	tests.WaitForNextNBlocksTM(2, port)

	validator = executeGetValidator(t, fmt.Sprintf("iriscli stake validator %s --output=json %v", barAddr, flags))
	require.Equal(t, "1.0000000000", validator.Tokens)
}
