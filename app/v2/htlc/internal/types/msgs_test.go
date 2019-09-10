package types

import (
	"encoding/hex"
	"fmt"
	"testing"

	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

var (
	senderAddr, _        = sdk.AccAddressFromBech32("faa128nh833v43sggcj65nk7khjka9dwngpl6j29hj")
	receiverAddr, _      = sdk.AccAddressFromBech32("faa1mrehjkgeg75nz2gk7lr7dnxvvtg4497jxss8hq")
	receiverOnOtherChain = []byte("receiverOnOtherChain")
	outAmount            = sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(10))
	inAmount             = uint64(100)
	secret               = "5f5f5f6162636465666768696a6b6c6d6e6f707172737475767778797a5f5f5f"
	secretStr            = "___abcdefghijklmnopqrstuvwxyz___"
	timestamp            = uint64(1580000000)
	secretHashLock       = hex.EncodeToString(sdk.SHA256(append([]byte(secretStr), sdk.Uint64ToBigEndian(timestamp)...)))
	timeLock             = uint64(50)
)

func TestNewMsgCreateHTLC(t *testing.T) {
	msg := NewMsgCreateHTLC(senderAddr, receiverAddr, receiverOnOtherChain, outAmount, inAmount, secretHashLock, timestamp, timeLock)

	require.Equal(t, senderAddr, msg.Sender)
	require.Equal(t, receiverAddr, msg.Receiver)
	require.Equal(t, receiverOnOtherChain, msg.ReceiverOnOtherChain)
	require.Equal(t, outAmount, msg.OutAmount)
	require.Equal(t, inAmount, msg.InAmount)
	require.Equal(t, secretHashLock, msg.SecretHashLock)
	require.Equal(t, timestamp, msg.Timestamp)
	require.Equal(t, timeLock, msg.TimeLock)
}

func TestMsgCreateHTLCRoute(t *testing.T) {
	// build a MsgCreateHTLC
	msg := NewMsgCreateHTLC(senderAddr, receiverAddr, receiverOnOtherChain, outAmount, inAmount, secretHashLock, timestamp, timeLock)
	require.Equal(t, "htlc", msg.Route())
}

func TestMsgCreateHTLCValidation(t *testing.T) {
	emptyAddr := sdk.AccAddress{}
	errReceiverOnOtherChain := make([]byte, 33)
	errOutAmount := sdk.Coin{}
	errSecretHashLock1 := "xx"
	errSecretHashLock2 := "00"
	errTimeLock1 := uint64(49)
	errTimeLock2 := uint64(25481)

	testData := []struct {
		expectPass           bool
		sender               sdk.AccAddress
		receiver             sdk.AccAddress
		receiverOnOtherChain []byte
		outAmount            sdk.Coin
		inAmount             uint64
		secretHashLock       string
		timestamp            uint64
		timeLock             uint64
	}{
		// correct
		{true, senderAddr, receiverAddr, receiverOnOtherChain, outAmount, inAmount, secretHashLock, timestamp, timeLock},
		// len(msg.Sender) == 0
		{false, emptyAddr, receiverAddr, receiverOnOtherChain, outAmount, inAmount, secretHashLock, timestamp, timeLock},
		// len(msg.Receiver) == 0
		{false, senderAddr, emptyAddr, receiverOnOtherChain, outAmount, inAmount, secretHashLock, timestamp, timeLock},
		// len(msg.ReceiverOnOtherChain) > MaxLengthForAddressOnOtherChain
		{false, senderAddr, receiverAddr, errReceiverOnOtherChain, outAmount, inAmount, secretHashLock, timestamp, timeLock},
		// !msg.OutAmount.IsPositive()
		{false, senderAddr, receiverAddr, receiverOnOtherChain, errOutAmount, inAmount, secretHashLock, timestamp, timeLock},
		// ValidateSecretHashLock(msg.SecretHashLock)
		{false, senderAddr, receiverAddr, receiverOnOtherChain, outAmount, inAmount, errSecretHashLock1, timestamp, timeLock},
		{false, senderAddr, receiverAddr, receiverOnOtherChain, outAmount, inAmount, errSecretHashLock2, timestamp, timeLock},
		// msg.TimeLock < MinTimeLock || msg.TimeLock > MaxTimeLock
		{false, senderAddr, receiverAddr, receiverOnOtherChain, outAmount, inAmount, secretHashLock, timestamp, errTimeLock1},
		{false, senderAddr, receiverAddr, receiverOnOtherChain, outAmount, inAmount, secretHashLock, timestamp, errTimeLock2},
	}

	for i, td := range testData {
		msg := NewMsgCreateHTLC(td.sender, td.receiver, td.receiverOnOtherChain, td.outAmount, td.inAmount, td.secretHashLock, td.timestamp, td.timeLock)
		err := msg.ValidateBasic()
		if td.expectPass {
			require.Nil(t, err, "%d: %+v", i, err)
		} else {
			require.NotNil(t, err, "%d: %+v", i, err)
		}
	}
}

func TestMsgCreateHTLCGetSignBytes(t *testing.T) {
	msg := NewMsgCreateHTLC(senderAddr, receiverAddr, receiverOnOtherChain, outAmount, inAmount, secretHashLock, timestamp, timeLock)
	res := msg.GetSignBytes()
	expected := `{"type":"irishub/htlc/MsgCreateHTLC","value":{"in_amount":"100","out_amount":{"amount":"10","denom":"iris-atto"},"receiver":"faa1mrehjkgeg75nz2gk7lr7dnxvvtg4497jxss8hq","receiver_on_other_chain":"cmVjZWl2ZXJPbk90aGVyQ2hhaW4=","secret_hash_lock":"e8d4133e1a82c74e2746e78c19385706ea7958a0ca441a08dacfa10c48ce2561","sender":"faa128nh833v43sggcj65nk7khjka9dwngpl6j29hj","time_lock":"50","timestamp":"1580000000"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgCreateHTLCGetSigners(t *testing.T) {
	msg := NewMsgCreateHTLC(senderAddr, receiverAddr, receiverOnOtherChain, outAmount, inAmount, secretHashLock, timestamp, timeLock)
	res := msg.GetSigners()
	expected := "[51E773C62CAC6084625AA4EDEB5E56E95AE9A03F]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

func TestNewMsgClaimHTLC(t *testing.T) {
	msg := NewMsgClaimHTLC(senderAddr, secret, secretHashLock)
	require.Equal(t, senderAddr, msg.Sender)
	require.Equal(t, secret, msg.Secret)
	require.Equal(t, secretHashLock, msg.SecretHashLock)
}

func TestMsgClaimHTLCRoute(t *testing.T) {
	msg := NewMsgClaimHTLC(senderAddr, secret, secretHashLock)
	require.Equal(t, "htlc", msg.Route())
}

func TestMsgClaimHTLCValidation(t *testing.T) {
	emptyAddr := sdk.AccAddress{}
	errSecret1 := "xx"
	errSecret2 := "00"
	errSecretHashLock1 := "xx"
	errSecretHashLock2 := "00"

	testData := []struct {
		expectPass     bool
		sender         sdk.AccAddress
		secret         string
		secretHashLock string
	}{
		// correct
		{true, senderAddr, secret, secretHashLock},
		// len(msg.Sender) == 0
		{false, emptyAddr, secret, secretHashLock},
		// ValidateSecret(msg.Secret)
		{false, senderAddr, errSecret1, secretHashLock},
		{false, senderAddr, errSecret2, secretHashLock},
		// ValidateSecretHashLock(msg.SecretHashLock)
		{false, senderAddr, secret, errSecretHashLock1},
		{false, senderAddr, secret, errSecretHashLock2},
	}

	for i, td := range testData {
		msg := NewMsgClaimHTLC(td.sender, td.secret, td.secretHashLock)
		err := msg.ValidateBasic()

		if td.expectPass {
			require.Nil(t, msg.ValidateBasic(), "%d: %+v", i, err)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "%d", i)
		}
	}
}

func TestMsgClaimHTLCGetSignBytes(t *testing.T) {
	msg := NewMsgClaimHTLC(senderAddr, secret, secretHashLock)
	res := msg.GetSignBytes()
	expected := `{"secret":"5f5f5f6162636465666768696a6b6c6d6e6f707172737475767778797a5f5f5f","secret_hash_lock":"e8d4133e1a82c74e2746e78c19385706ea7958a0ca441a08dacfa10c48ce2561","sender":"faa128nh833v43sggcj65nk7khjka9dwngpl6j29hj"}`
	require.Equal(t, expected, string(res))
}

func TestMsgClaimHTLCGetSigners(t *testing.T) {
	msg := NewMsgClaimHTLC(senderAddr, secret, secretHashLock)
	res := msg.GetSigners()
	expected := "[51E773C62CAC6084625AA4EDEB5E56E95AE9A03F]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

func TestNewMsgRefundHTLC(t *testing.T) {
	msg := NewMsgRefundHTLC(senderAddr, secretHashLock)
	require.Equal(t, senderAddr, msg.Sender)
	require.Equal(t, secretHashLock, msg.SecretHashLock)
}

func TestMsgRefundHTLCRoute(t *testing.T) {
	msg := NewMsgRefundHTLC(senderAddr, secretHashLock)
	require.Equal(t, "htlc", msg.Route())
}

func TestMsgRefundHTLCValidation(t *testing.T) {
	emptyAddr := sdk.AccAddress{}
	errSecretHashLock1 := "xx"
	errSecretHashLock2 := "00"

	testData := []struct {
		expectPass     bool
		sender         sdk.AccAddress
		secretHashLock string
	}{
		// correct
		{true, senderAddr, secretHashLock},
		// len(msg.Sender) == 0
		{false, emptyAddr, secretHashLock},
		// ValidateSecretHashLock(msg.SecretHashLock)
		{false, senderAddr, errSecretHashLock1},
		{false, senderAddr, errSecretHashLock2},
	}

	for i, td := range testData {
		msg := NewMsgRefundHTLC(td.sender, td.secretHashLock)
		err := msg.ValidateBasic()

		if td.expectPass {
			require.Nil(t, msg.ValidateBasic(), "%d: %+v", i, err)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "%d", i)
		}
	}
}

func TestMsgRefundHTLCGetSignBytes(t *testing.T) {
	msg := NewMsgRefundHTLC(senderAddr, secretHashLock)
	res := msg.GetSignBytes()
	expected := `{"secret_hash_lock":"e8d4133e1a82c74e2746e78c19385706ea7958a0ca441a08dacfa10c48ce2561","sender":"faa128nh833v43sggcj65nk7khjka9dwngpl6j29hj"}`
	require.Equal(t, expected, string(res))
}

func TestMsgRefundHTLCGetSigners(t *testing.T) {
	msg := NewMsgRefundHTLC(senderAddr, secretHashLock)
	res := msg.GetSigners()
	expected := "[51E773C62CAC6084625AA4EDEB5E56E95AE9A03F]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}
