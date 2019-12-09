package clitest

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/irisnet/irishub/app"
	iconfig "github.com/irisnet/irishub/config"
	htlcmodule "github.com/irisnet/irishub/modules/htlc"
)

func init() {
	// set Bech32 config
	config := sdk.GetConfig()
	irisConfig := iconfig.GetConfig()
	config.SetBech32PrefixForAccount(irisConfig.GetBech32AccountAddrPrefix(), irisConfig.GetBech32AccountPubPrefix())
	config.SetBech32PrefixForValidator(irisConfig.GetBech32ValidatorAddrPrefix(), irisConfig.GetBech32ValidatorPubPrefix())
	config.SetBech32PrefixForConsensusNode(irisConfig.GetBech32ConsensusAddrPrefix(), irisConfig.GetBech32ConsensusPubPrefix())
	config.Seal()
}

func TestIrisCLIHTLC(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)

	flags := fmt.Sprintf("--home=%s --node=%v --chain-id=%v --output=json", f.IriscliHome, f.RPCAddr, f.ChainID)

	proc := f.GDStart()
	defer proc.Stop(false)

	tests.WaitForTMStart(f.Port)
	tests.WaitForNextNBlocksTM(2, f.Port)

	fooAddr := f.KeyAddress("foo")
	barAddr := f.KeyAddress("bar")

	fooAcc := f.QueryAccount(fooAddr, flags)
	fooCoin := fooAcc.Coins.AmountOf(sdk.DefaultBondDenom).String()
	require.Equal(t, "50000000", fooCoin)

	// testdata
	receiverOnOtherChain := "0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826"
	hashLock := "e8d4133e1a82c74e2746e78c19385706ea7958a0ca441a08dacfa10c48ce2561"
	secretHex := "5f5f5f6162636465666768696a6b6c6d6e6f707172737475767778797a5f5f5f"
	amount := "1000" + sdk.DefaultBondDenom
	timeLock := uint64(50)
	timestamp := uint64(1580000000)
	initSecret := []byte(nil)
	stateOpen := "open"
	stateCompleted := "completed"
	stateExpired := "expired"
	stateRefunded := "refunded"

	// create an htlc
	spStr := fmt.Sprintf("%s tx htlc create %v", f.IriscliBinary, flags)
	spStr += fmt.Sprintf(" --from=%s", "foo")
	spStr += fmt.Sprintf(" --to=%s", barAddr)
	spStr += fmt.Sprintf(" --receiver-on-other-chain=%s", receiverOnOtherChain)
	spStr += fmt.Sprintf(" --secret=%s", secretHex)
	spStr += fmt.Sprintf(" --amount=%s", amount)
	spStr += fmt.Sprintf(" --time-lock=%d", timeLock)
	spStr += fmt.Sprintf(" --timestamp=%d", timestamp)
	spStr += fmt.Sprintf(" --fees=%s", "3000"+sdk.DefaultBondDenom)
	spStr += " -y"

	require.True(t, executeWrite(t, spStr, client.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, f.Port)

	htlc := executeGetHTLC(t, fmt.Sprintf("%s query htlc htlc %s --output=json %v", f.IriscliBinary, strings.ToLower(strings.TrimSpace(hashLock)), flags))

	require.Equal(t, fooAddr, htlc.Sender)
	require.Equal(t, barAddr, htlc.To)
	require.Equal(t, receiverOnOtherChain, htlc.ReceiverOnOtherChain)
	require.Equal(t, amount, htlc.Amount.String())
	require.Equal(t, initSecret, htlc.Secret)
	require.Equal(t, timestamp, htlc.Timestamp)
	require.Equal(t, stateOpen, htlc.State.String())

	htlcAddr := supply.NewModuleAddress("htlc")
	htlcAcc := f.QueryAccount(htlcAddr, flags)
	htlcCoin := htlcAcc.Coins.AmountOf(sdk.DefaultBondDenom).String()
	require.Equal(t, "1000", htlcCoin)

	// claim an htlc
	spStr = fmt.Sprintf("%s tx htlc claim %v", f.IriscliBinary, flags)
	spStr += fmt.Sprintf(" --from=%s", "foo")
	spStr += fmt.Sprintf(" --hash-lock=%s", hashLock)
	spStr += fmt.Sprintf(" --secret=%s", secretHex)
	spStr += fmt.Sprintf(" --fees=%s", "3000"+sdk.DefaultBondDenom)
	spStr += " -y"

	require.True(t, executeWrite(t, spStr, client.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, f.Port)

	htlc = executeGetHTLC(t, fmt.Sprintf("%s query htlc htlc %s --output=json %v", f.IriscliBinary, strings.ToLower(strings.TrimSpace(hashLock)), flags))
	require.Equal(t, stateCompleted, htlc.State.String())

	htlcAcc = f.QueryAccount(htlcAddr, flags)
	require.Equal(t, "0", htlcAcc.GetCoins().AmountOf(sdk.DefaultBondDenom).String())

	barAcc := f.QueryAccount(barAddr, flags)
	barCoin := barAcc.Coins.AmountOf(sdk.DefaultBondDenom).String()
	require.Equal(t, "1000", barCoin)

	// testdata
	hashLock = "f054e34abd9ccc3cab12a5b797b8e9c053507f279e7e53fb3f9f44d178c94b20"
	timestamp = uint64(0)
	timeLock = uint64(50)

	// create an htlc
	spStr = fmt.Sprintf("%s tx htlc create %v", f.IriscliBinary, flags)
	spStr += fmt.Sprintf(" --from=%s", "foo")
	spStr += fmt.Sprintf(" --to=%s", barAddr)
	spStr += fmt.Sprintf(" --receiver-on-other-chain=%s", receiverOnOtherChain)
	spStr += fmt.Sprintf(" --secret=%s", secretHex)
	spStr += fmt.Sprintf(" --amount=%s", amount)
	spStr += fmt.Sprintf(" --time-lock=%d", timeLock)
	spStr += fmt.Sprintf(" --timestamp=%d", timestamp)
	spStr += fmt.Sprintf(" --fees=%s", "3000"+sdk.DefaultBondDenom)
	spStr += " -y"

	require.True(t, executeWrite(t, spStr, client.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, f.Port)

	htlc = executeGetHTLC(t, fmt.Sprintf("%s query htlc htlc %s --output=json %v", f.IriscliBinary, strings.ToLower(strings.TrimSpace(hashLock)), flags))
	require.Equal(t, fooAddr, htlc.Sender)
	require.Equal(t, barAddr, htlc.To)
	require.Equal(t, receiverOnOtherChain, htlc.ReceiverOnOtherChain)
	require.Equal(t, amount, htlc.Amount.String())
	require.Equal(t, initSecret, htlc.Secret)
	require.Equal(t, timestamp, htlc.Timestamp)
	require.Equal(t, stateOpen, htlc.State.String())

	htlcAcc = f.QueryAccount(htlcAddr, flags)
	htlcCoin = htlcAcc.Coins.AmountOf(sdk.DefaultBondDenom).String()
	require.Equal(t, "1000", htlcCoin)

	// refund an htlc and expect failure
	spStr = fmt.Sprintf("%s tx htlc refund %v", f.IriscliBinary, flags)
	spStr += fmt.Sprintf(" --from=%s", "foo")
	spStr += fmt.Sprintf(" --hash-lock=%s", hashLock)
	spStr += fmt.Sprintf(" --fees=%s", "3000"+sdk.DefaultBondDenom)
	spStr += " -y"

	require.True(t, executeWrite(t, spStr, client.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, f.Port)

	htlc = executeGetHTLC(t, fmt.Sprintf("%s query htlc htlc %s --output=json %v", f.IriscliBinary, strings.ToLower(strings.TrimSpace(hashLock)), flags))
	require.Equal(t, stateOpen, htlc.State.String())

	// refund an htlc and expect success
	tests.WaitForNextNBlocksTM(int64(timeLock), f.Port)

	htlc = executeGetHTLC(t, fmt.Sprintf("%s query htlc htlc %s --output=json %v", f.IriscliBinary, strings.ToLower(strings.TrimSpace(hashLock)), flags))
	require.Equal(t, stateExpired, htlc.State.String())

	require.True(t, executeWrite(t, spStr, client.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, f.Port)

	htlc = executeGetHTLC(t, fmt.Sprintf("%s query htlc htlc %s --output=json %v", f.IriscliBinary, strings.ToLower(strings.TrimSpace(hashLock)), flags))
	require.Equal(t, stateRefunded, htlc.State.String())

	htlcAcc = f.QueryAccount(htlcAddr, flags)
	require.Equal(t, "0", htlcAcc.GetCoins().AmountOf(sdk.DefaultBondDenom).String())
}

func executeGetHTLC(t *testing.T, cmdStr string) htlcmodule.HTLC {
	cdc := app.MakeCodec()
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var htlc htlcmodule.HTLC
	err := cdc.UnmarshalJSON([]byte(out), &htlc)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return htlc
}
