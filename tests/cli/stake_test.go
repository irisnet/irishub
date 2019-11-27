package cli

import (
	"fmt"
	"testing"

	"github.com/irisnet/irishub/tests"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

func TestIrisCLIStakeCreateValidator(t *testing.T) {
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

	irisHomeB, _ := getTestingHomeDirsB()
	executeInit(t, fmt.Sprintf("iris init -o --moniker=foo --home=%s", irisHomeB))
	barCeshPubKey := executeGetValidatorPK(t, fmt.Sprintf("iris tendermint show-validator --home=%s", irisHomeB))

	executeWrite(t, fmt.Sprintf("iriscli bank send %v --amount=10iris --to=%s --from=foo --gas=10000 --fee=0.3iris", flags, barAddr), sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	barAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin := convertToIrisBaseAccount(t, barAcc)
	require.Equal(t, "10iris", barCoin)

	fooAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin := convertToIrisBaseAccount(t, fooAcc)
	num := getAmountFromCoinStr(fooCoin)

	if !(num > 39 && num < 40) {
		t.Error("Test Failed: (39, 40) expected, received: {}", num)
	}

	// create validator
	cvStr := fmt.Sprintf("iriscli stake create-validator %v", flags)
	cvStr += fmt.Sprintf(" --from=%s", "bar")
	cvStr += fmt.Sprintf(" --pubkey=%s", barCeshPubKey)
	cvStr += fmt.Sprintf(" --amount=%v", "2iris")
	cvStr += fmt.Sprintf(" --moniker=%v", "bar-vally")
	cvStr += fmt.Sprintf(" --fee=%s", "0.4iris")
	cvStr += fmt.Sprintf(" --commission-rate=%s", "0.1")

	executeWrite(t, cvStr, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	barAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin = convertToIrisBaseAccount(t, barAcc)
	num = getAmountFromCoinStr(barCoin)

	if !(num > 7 && num < 8) {
		t.Error("Test Failed: (7, 8) expected, received: {}", num)
	}

	valAddr := sdk.ValAddress(barAddr).String()
	validator := executeGetValidator(t, fmt.Sprintf("iriscli stake validator %s --output=json %v", valAddr, flags))
	require.Equal(t, valAddr, validator.OperatorAddr.String())
	require.Equal(t, "2.0000000000000000000000000000", validator.Tokens)

	// unbond a single share
	unbondStr := fmt.Sprintf("iriscli stake unbond %v", flags)
	unbondStr += fmt.Sprintf(" --from=%s", "bar")
	unbondStr += fmt.Sprintf(" --address-validator=%s", valAddr)
	unbondStr += fmt.Sprintf(" --shares-amount=%v", "1")
	unbondStr += fmt.Sprintf(" --fee=%s", "0.4iris")
	unbondStr += fmt.Sprintf(" --gas=%s", "40000")

	success := executeWrite(t, unbondStr, sdk.DefaultKeyPass)
	require.True(t, success)
	tests.WaitForNextNBlocksTM(2, port)

	validator = executeGetValidator(t, fmt.Sprintf("iriscli stake validator %s --output=json %v", valAddr, flags))
	require.Equal(t, "1.0000000000000000000000000000", validator.Tokens)
}
