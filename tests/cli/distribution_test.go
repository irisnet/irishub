package cli

import (
	"fmt"
	"testing"

	"github.com/irisnet/irishub/tests"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

func TestIrisCLIDistribution(t *testing.T) {
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
	num := getAmountFromCoinStr(fooCoin)

	valAddr := sdk.ValAddress(fooAddr).String()

	ddiList := executeGetDelegatorDistrInfo(t, fmt.Sprintf("iriscli distribution delegator-distr-info %s %s", fooAddr, flags))
	require.Equal(t, 1, len(ddiList))
	require.Equal(t, int64(0), ddiList[0].DelPoolWithdrawalHeight)
	require.Equal(t, fooAddr, ddiList[0].DelegatorAddr)
	require.Equal(t, valAddr, ddiList[0].ValOperatorAddr.String())

	ddi := executeGetDelegationDistrInfo(t, fmt.Sprintf("iriscli distribution delegation-distr-info --address-delegator=%s --address-validator=%s %s", fooAddr, valAddr, flags))
	require.Equal(t, int64(0), ddi.DelPoolWithdrawalHeight)
	require.Equal(t, fooAddr, ddi.DelegatorAddr)
	require.Equal(t, valAddr, ddi.ValOperatorAddr.String())

	vdi := executeGetValidatorDistrInfo(t, fmt.Sprintf("iriscli distribution validator-distr-info %s %s", valAddr, flags))
	require.Equal(t, valAddr, vdi.OperatorAddr.String())
	require.Equal(t, int64(0), vdi.FeePoolWithdrawalHeight)
	numDelPool := getAmountFromCoinStr(vdi.DelPool)
	numValCommission := getAmountFromCoinStr(vdi.ValCommission)
	require.True(t, numDelPool > numValCommission)

	executeWrite(t, fmt.Sprintf("iriscli distribution withdraw-rewards --from=foo --fee=0.4iris %s", flags), sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	numNew := getAmountFromCoinStr(fooCoin)
	require.True(t, numNew > num)
}

func TestIrisCLIWithdrawReward(t *testing.T) {
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

	valAddr := sdk.ValAddress(fooAddr).String()

	vdi := executeGetValidatorDistrInfo(t, fmt.Sprintf("iriscli distribution validator-distr-info %s %s", valAddr, flags))
	require.Equal(t, valAddr, vdi.OperatorAddr.String())
	require.Equal(t, int64(0), vdi.FeePoolWithdrawalHeight)
	numDelPool := getAmountFromCoinStr(vdi.DelPool)
	numValCommission := getAmountFromCoinStr(vdi.ValCommission)
	require.True(t, numDelPool > numValCommission)

	executeWrite(t, fmt.Sprintf("iriscli distribution withdraw-rewards --only-from-validator=%s --from=foo --fee=0.4iris %s", valAddr, flags), sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	num := getAmountFromCoinStr(fooCoin)
	require.True(t, num > 10)

	vdi = executeGetValidatorDistrInfo(t, fmt.Sprintf("iriscli distribution validator-distr-info %s %s", valAddr, flags))
	require.Equal(t, valAddr, vdi.OperatorAddr.String())
	numValCommission = getAmountFromCoinStr(vdi.ValCommission)
	require.True(t, numValCommission > 0)

	executeWrite(t, fmt.Sprintf("iriscli distribution withdraw-rewards --is-validator=true --from=foo --fee=0.4iris %s", flags), sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	numNew := getAmountFromCoinStr(fooCoin)

	if numNew <= num {
		t.Error("Test Failed: if --is-validator is true, more reward should be return")
	}
}
