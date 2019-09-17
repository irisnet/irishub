package cli

import (
	"fmt"
	"github.com/irisnet/irishub/tests"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIrisCLIBankSend(t *testing.T) {
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

	executeWrite(t, fmt.Sprintf("iriscli bank send %v --amount=10iris --to=%s --from=foo --gas=10000 --fee=0.3iris", flags, barAddr), sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	barAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin := convertToIrisBaseAccount(t, barAcc)
	require.Equal(t, "10iris", barCoin)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	num := getAmountFromCoinStr(fooCoin)

	if !(num > 39 && num < 40) {
		t.Error("Test Failed: (39, 40) expected, received: {}", num)
	}

	// test autosequencing
	executeWrite(t, fmt.Sprintf("iriscli bank send %v --amount=10iris --to=%s --from=foo --gas=10000 --fee=0.3iris", flags, barAddr), sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	barAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin = convertToIrisBaseAccount(t, barAcc)
	require.Equal(t, "20iris", barCoin)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	num = getAmountFromCoinStr(fooCoin)

	if !(num > 29 && num < 30) {
		t.Error("Test Failed: (29, 30) expected, received: {}", num)
	}

	// test memo
	executeWrite(t, fmt.Sprintf("iriscli bank send %v --amount=10iris --to=%s --from=foo --memo 'testmemo' --gas=10000 --fee=0.3iris", flags, barAddr), sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	barAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin = convertToIrisBaseAccount(t, barAcc)
	require.Equal(t, "30iris", barCoin)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	num = getAmountFromCoinStr(fooCoin)

	if !(num > 19 && num < 20) {
		t.Error("Test Failed: (69, 70) expected, received: {}", num)
	}
}

func TestIrisCLIBankTokenStatsById(t *testing.T) {
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

	executeWrite(t, fmt.Sprintf("iriscli asset issue-token %v --family=fungible --source=native  --symbol=kitty --name=eeee --initial-supply=1000 --decimal=0 --from=foo --fee=0.6iris", flags), sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	executeWrite(t, fmt.Sprintf("iriscli bank burn --from=foo --amount=10kitty %v --fee=0.6iris", flags), sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	tokenStats := executeGetTokenStatsForAsset(t, fmt.Sprintf("iriscli bank token-stats %v kitty", flags))
	require.Nil(t, tokenStats.LooseTokens)
	require.Nil(t, tokenStats.BondedTokens)
	require.Contains(t, tokenStats.TotalSupply, sdk.NewCoin("kitty-min", sdk.NewIntWithDecimal(990, 0)))
	require.Contains(t, tokenStats.BurnedTokens, sdk.NewCoin("kitty-min", sdk.NewIntWithDecimal(10, 0)))
}
