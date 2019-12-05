package clitest

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/config"
	htlcmodule "github.com/irisnet/irishub/modules/htlc"
)

// TODO: fix
func TestIrisCLIHTLC(t *testing.T) {
	t.Parallel()
	fixtures := InitFixtures(t)

	flags := fmt.Sprintf("--home=%s --node=%v --chain-id=%v --output=json", fixtures.IriscliHome, fixtures.RPCAddr, fixtures.ChainID)

	// start iris server
	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iris start --home=%s --rpc.laddr=%v --p2p.laddr=%v", fixtures.IrisdHome, fixtures.RPCAddr, fixtures.P2PAddr))

	defer proc.Stop(false)
	tests.WaitForTMStart(fixtures.Port)
	tests.WaitForNextNBlocksTM(2, fixtures.Port)

	fooAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show foo --output=json --home=%s", fixtures.IriscliHome))
	barAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show bar --output=json --home=%s", fixtures.IriscliHome))

	fooAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin := fooAcc.Coins.AmountOf(config.Iris).String()
	require.Equal(t, "50iris", fooCoin)

	// testdata
	receiverOnOtherChain := "0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826"
	hashLock := "e8d4133e1a82c74e2746e78c19385706ea7958a0ca441a08dacfa10c48ce2561"
	secretHex := "5f5f5f6162636465666768696a6b6c6d6e6f707172737475767778797a5f5f5f"
	amount := "10000000000000000000iris"
	amountIris := "10000000000000000000iris"
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

	require.True(t, executeWrite(t, spStr, config.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, fixtures.Port)

	htlc := executeGetHTLC(t, fmt.Sprintf("iriscli htlc query-htlc %s --output=json %v", strings.ToLower(strings.TrimSpace(hashLock)), flags))
	require.Equal(t, fooAddr, htlc.Sender)
	require.Equal(t, barAddr, htlc.To)
	require.Equal(t, receiverOnOtherChain, htlc.ReceiverOnOtherChain)
	require.Equal(t, amount, htlc.Amount.String())
	require.Equal(t, initSecret, htlc.Secret)
	require.Equal(t, timestamp, htlc.Timestamp)
	require.Equal(t, stateOpen, htlc.State.String())

	htlcAddr := supply.NewModuleAddress("htlc")
	htlcAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", htlcAddr, flags))
	htlcCoin := htlcAcc.Coins.AmountOf(config.Iris).String()
	require.Equal(t, "10iris", htlcCoin)

	// claim an htlc
	spStr = fmt.Sprintf("iriscli htlc claim %v", flags)
	spStr += fmt.Sprintf(" --from=%s", "foo")
	spStr += fmt.Sprintf(" --hash-lock=%s", hashLock)
	spStr += fmt.Sprintf(" --secret=%s", secretHex)
	spStr += fmt.Sprintf(" --fee=%s", "0.4iris")

	require.True(t, executeWrite(t, spStr, config.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, fixtures.Port)

	htlc = executeGetHTLC(t, fmt.Sprintf("iriscli htlc query-htlc %s --output=json %v", strings.ToLower(strings.TrimSpace(hashLock)), flags))
	require.Equal(t, stateCompleted, htlc.State.String())

	htlcAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", htlcAddr, flags))
	require.Equal(t, "0", htlcAcc.GetCoins().AmountOf(config.Iris).String())

	barAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin := barAcc.Coins.AmountOf(config.Iris).String()
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

	require.True(t, executeWrite(t, spStr, config.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, fixtures.Port)

	htlc = executeGetHTLC(t, fmt.Sprintf("iriscli htlc query-htlc %s --output=json %v", strings.ToLower(strings.TrimSpace(hashLock)), flags))
	require.Equal(t, fooAddr, htlc.Sender)
	require.Equal(t, barAddr, htlc.To)
	require.Equal(t, receiverOnOtherChain, htlc.ReceiverOnOtherChain)
	require.Equal(t, amount, htlc.Amount.String())
	require.Equal(t, initSecret, htlc.Secret)
	require.Equal(t, timestamp, htlc.Timestamp)
	require.Equal(t, stateOpen, htlc.State.String())

	htlcAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", htlcAddr, flags))
	htlcCoin = htlcAcc.Coins.AmountOf(config.Iris).String()
	require.Equal(t, "10iris", htlcCoin)

	// refund an htlc and expect failure
	spStr = fmt.Sprintf("iriscli htlc refund %v", flags)
	spStr += fmt.Sprintf(" --from=%s", "foo")
	spStr += fmt.Sprintf(" --hash-lock=%s", hashLock)
	spStr += fmt.Sprintf(" --fee=%s", "0.4iris")

	require.Zero(t, executeWrite(t, spStr, config.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, fixtures.Port)

	htlc = executeGetHTLC(t, fmt.Sprintf("iriscli htlc query-htlc %s --output=json %v", strings.ToLower(strings.TrimSpace(hashLock)), flags))
	require.Equal(t, stateOpen, htlc.State.String())

	// refund an htlc and expect success
	tests.WaitForNextNBlocksTM(int64(timeLock), fixtures.Port)

	htlc = executeGetHTLC(t, fmt.Sprintf("iriscli htlc query-htlc %s --output=json %v", strings.ToLower(strings.TrimSpace(hashLock)), flags))
	require.Equal(t, stateExpired, htlc.State.String())

	require.True(t, executeWrite(t, spStr, config.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, fixtures.Port)

	htlc = executeGetHTLC(t, fmt.Sprintf("iriscli htlc query-htlc %s --output=json %v", strings.ToLower(strings.TrimSpace(hashLock)), flags))
	require.Equal(t, stateRefunded, htlc.State.String())

	htlcAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", htlcAddr, flags))
	require.Equal(t, "0", htlcAcc.GetCoins().AmountOf(config.Iris).String())
}

func executeGetAddrPK(t *testing.T, cmdStr string) (sdk.AccAddress, crypto.PubKey) {
	cdc := app.MakeCodec()
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var ko keys.KeyOutput
	_ = cdc.UnmarshalJSON([]byte(out), &ko)

	pk, err := sdk.GetAccPubKeyBech32(ko.PubKey)
	require.NoError(t, err)

	accAddr, err := sdk.AccAddressFromBech32(ko.Address)
	require.NoError(t, err)

	return accAddr, pk
}

func executeGetAccount(t *testing.T, cmdStr string) (acc auth.BaseAccount) {
	cdc := app.MakeCodec()
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var initRes map[string]json.RawMessage
	err := json.Unmarshal([]byte(out), &initRes)
	require.NoError(t, err, "out %v, err %v", out, err)

	err = cdc.UnmarshalJSON([]byte(out), &acc)
	require.NoError(t, err, "acc %v, err %v", string(out), err)

	return acc
}

func executeGetHTLC(t *testing.T, cmdStr string) htlcmodule.HTLC {
	cdc := app.MakeCodec()
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var htlc htlcmodule.HTLC
	err := cdc.UnmarshalJSON([]byte(out), &htlc)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return htlc
}
