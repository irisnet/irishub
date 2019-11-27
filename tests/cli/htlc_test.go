package cli

import (
	"fmt"
	"strings"
	"testing"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/tests"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

func TestIrisCLIHTLC(t *testing.T) {
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

	// testdata
	receiverOnOtherChain := "0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826"
	hashLock := "e8d4133e1a82c74e2746e78c19385706ea7958a0ca441a08dacfa10c48ce2561"
	secretHex := "5f5f5f6162636465666768696a6b6c6d6e6f707172737475767778797a5f5f5f"
	amount := "10000000000000000000iris-atto"
	amountIris := "10000000000000000000iris-atto"
	timeLock := uint64(50)
	timestamp := uint64(1580000000)
	initSecret := []byte(nil)
	stateOpen := "open"
	stateCompleted := "completed"
	stateExpired := "expired"
	stateRefunded := "refunded"

	// create an htlc
	spStr := fmt.Sprintf("iriscli htlc create %v", flags)
	spStr += fmt.Sprintf(" --from=%s", "foo")
	spStr += fmt.Sprintf(" --to=%s", barAddr)
	spStr += fmt.Sprintf(" --receiver-on-other-chain=%s", receiverOnOtherChain)
	spStr += fmt.Sprintf(" --secret=%s", secretHex)
	spStr += fmt.Sprintf(" --amount=%s", amountIris)
	spStr += fmt.Sprintf(" --time-lock=%d", timeLock)
	spStr += fmt.Sprintf(" --timestamp=%d", timestamp)
	spStr += fmt.Sprintf(" --fee=%s", "0.4iris")

	require.True(t, executeWrite(t, spStr, sdk.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, port)

	htlc := executeGetHtlc(t, fmt.Sprintf("iriscli htlc query-htlc %s --output=json %v", strings.ToLower(strings.TrimSpace(hashLock)), flags))
	require.Equal(t, fooAddr, htlc.Sender)
	require.Equal(t, barAddr, htlc.To)
	require.Equal(t, receiverOnOtherChain, htlc.ReceiverOnOtherChain)
	require.Equal(t, amount, htlc.Amount.String())
	require.Equal(t, initSecret, htlc.Secret)
	require.Equal(t, timestamp, htlc.Timestamp)
	require.Equal(t, stateOpen, htlc.State.String())

	htlcAddr := auth.HTLCLockedCoinsAccAddr
	htlcAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", htlcAddr, flags))
	htlcCoin := convertToIrisBaseAccount(t, htlcAcc)
	require.Equal(t, "10iris", htlcCoin)

	// claim an htlc
	spStr = fmt.Sprintf("iriscli htlc claim %v", flags)
	spStr += fmt.Sprintf(" --from=%s", "foo")
	spStr += fmt.Sprintf(" --hash-lock=%s", hashLock)
	spStr += fmt.Sprintf(" --secret=%s", secretHex)
	spStr += fmt.Sprintf(" --fee=%s", "0.4iris")

	require.True(t, executeWrite(t, spStr, sdk.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, port)

	htlc = executeGetHtlc(t, fmt.Sprintf("iriscli htlc query-htlc %s --output=json %v", strings.ToLower(strings.TrimSpace(hashLock)), flags))
	require.Equal(t, stateCompleted, htlc.State.String())

	htlcAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", htlcAddr, flags))
	require.Equal(t, "0", htlcAcc.GetCoins().AmountOf(sdk.IrisAtto).String())

	barAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin := convertToIrisBaseAccount(t, barAcc)
	require.Equal(t, "10iris", barCoin)

	// testdata
	hashLock = "f054e34abd9ccc3cab12a5b797b8e9c053507f279e7e53fb3f9f44d178c94b20"
	timestamp = uint64(0)
	timeLock = uint64(50)

	// create an htlc
	spStr = fmt.Sprintf("iriscli htlc create %v", flags)
	spStr += fmt.Sprintf(" --from=%s", "foo")
	spStr += fmt.Sprintf(" --to=%s", barAddr)
	spStr += fmt.Sprintf(" --receiver-on-other-chain=%s", receiverOnOtherChain)
	spStr += fmt.Sprintf(" --secret=%s", secretHex)
	spStr += fmt.Sprintf(" --amount=%s", amountIris)
	spStr += fmt.Sprintf(" --time-lock=%d", timeLock)
	spStr += fmt.Sprintf(" --timestamp=%d", timestamp)
	spStr += fmt.Sprintf(" --fee=%s", "0.4iris")

	require.True(t, executeWrite(t, spStr, sdk.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, port)

	htlc = executeGetHtlc(t, fmt.Sprintf("iriscli htlc query-htlc %s --output=json %v", strings.ToLower(strings.TrimSpace(hashLock)), flags))
	require.Equal(t, fooAddr, htlc.Sender)
	require.Equal(t, barAddr, htlc.To)
	require.Equal(t, receiverOnOtherChain, htlc.ReceiverOnOtherChain)
	require.Equal(t, amount, htlc.Amount.String())
	require.Equal(t, initSecret, htlc.Secret)
	require.Equal(t, timestamp, htlc.Timestamp)
	require.Equal(t, stateOpen, htlc.State.String())

	htlcAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", htlcAddr, flags))
	htlcCoin = convertToIrisBaseAccount(t, htlcAcc)
	require.Equal(t, "10iris", htlcCoin)

	// refund an htlc and expect failure
	spStr = fmt.Sprintf("iriscli htlc refund %v", flags)
	spStr += fmt.Sprintf(" --from=%s", "foo")
	spStr += fmt.Sprintf(" --hash-lock=%s", hashLock)
	spStr += fmt.Sprintf(" --fee=%s", "0.4iris")

	require.Zero(t, executeWrite(t, spStr, sdk.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, port)

	htlc = executeGetHtlc(t, fmt.Sprintf("iriscli htlc query-htlc %s --output=json %v", strings.ToLower(strings.TrimSpace(hashLock)), flags))
	require.Equal(t, stateOpen, htlc.State.String())

	// refund an htlc and expect success
	tests.WaitForNextNBlocksTM(int64(timeLock), port)

	htlc = executeGetHtlc(t, fmt.Sprintf("iriscli htlc query-htlc %s --output=json %v", strings.ToLower(strings.TrimSpace(hashLock)), flags))
	require.Equal(t, stateExpired, htlc.State.String())

	require.True(t, executeWrite(t, spStr, sdk.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, port)

	htlc = executeGetHtlc(t, fmt.Sprintf("iriscli htlc query-htlc %s --output=json %v", strings.ToLower(strings.TrimSpace(hashLock)), flags))
	require.Equal(t, stateRefunded, htlc.State.String())

	htlcAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", htlcAddr, flags))
	require.Equal(t, "0", htlcAcc.GetCoins().AmountOf(sdk.IrisAtto).String())
}
