package types

import (
	"fmt"
	"testing"

	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

var (
	senderAddr, _        = sdk.AccAddressFromBech32("faa128nh833v43sggcj65nk7khjka9dwngpl6j29hj")
	receiverAddr, _      = sdk.AccAddressFromBech32("faa1mrehjkgeg75nz2gk7lr7dnxvvtg4497jxss8hq")
	receiverOnOtherChain = []byte("receiverOnOtherChain")
	amount               = sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(10))
	secret               = []byte("___abcdefghijklmnopqrstuvwxyz___")
	timestamp            = uint64(1580000000)
	hashLock             = sdk.SHA256(append(secret, sdk.Uint64ToBigEndian(timestamp)...))
	timeLock             = uint64(50)
)

func TestNewMsgCreateHTLC(t *testing.T) {
	msg := NewMsgCreateHTLC(senderAddr, receiverAddr, receiverOnOtherChain, amount, hashLock, timestamp, timeLock)

	require.Equal(t, senderAddr, msg.Sender)
	require.Equal(t, receiverAddr, msg.Receiver)
	require.Equal(t, receiverOnOtherChain, msg.ReceiverOnOtherChain)
	require.Equal(t, amount, msg.Amount)
	require.Equal(t, hashLock, msg.HashLock)
	require.Equal(t, timestamp, msg.Timestamp)
	require.Equal(t, timeLock, msg.TimeLock)
}

func TestMsgCreateHTLCRoute(t *testing.T) {
	// build a MsgCreateHTLC
	msg := NewMsgCreateHTLC(senderAddr, receiverAddr, receiverOnOtherChain, amount, hashLock, timestamp, timeLock)
	require.Equal(t, RouterKey, msg.Route())
}

func TestMsgCreateHTLCValidation(t *testing.T) {
	emptyAddr := sdk.AccAddress{}
	errReceiverOnOtherChain := make([]byte, 33)
	errAmount := sdk.Coin{}
	errHashLock1 := []byte("xx")
	errHashLock2 := []byte("00")
	errTimeLock1 := uint64(49)
	errTimeLock2 := uint64(25481)

	testData := []struct {
		expectPass           bool
		sender               sdk.AccAddress
		receiver             sdk.AccAddress
		receiverOnOtherChain []byte
		amount               sdk.Coin
		hashLock             []byte
		timestamp            uint64
		timeLock             uint64
	}{
		// correct
		{true, senderAddr, receiverAddr, receiverOnOtherChain, amount, hashLock, timestamp, timeLock},
		// len(msg.Sender) == 0
		{false, emptyAddr, receiverAddr, receiverOnOtherChain, amount, hashLock, timestamp, timeLock},
		// len(msg.Receiver) == 0
		{false, senderAddr, emptyAddr, receiverOnOtherChain, amount, hashLock, timestamp, timeLock},
		// len(msg.ReceiverOnOtherChain) > MaxLengthForAddressOnOtherChain
		{false, senderAddr, receiverAddr, errReceiverOnOtherChain, amount, hashLock, timestamp, timeLock},
		// !msg.OutAmount.IsPositive()
		{false, senderAddr, receiverAddr, receiverOnOtherChain, errAmount, hashLock, timestamp, timeLock},
		// ValidateSecretHashLock(msg.SecretHashLock)
		{false, senderAddr, receiverAddr, receiverOnOtherChain, amount, errHashLock1, timestamp, timeLock},
		{false, senderAddr, receiverAddr, receiverOnOtherChain, amount, errHashLock2, timestamp, timeLock},
		// msg.TimeLock < MinTimeLock || msg.TimeLock > MaxTimeLock
		{false, senderAddr, receiverAddr, receiverOnOtherChain, amount, hashLock, timestamp, errTimeLock1},
		{false, senderAddr, receiverAddr, receiverOnOtherChain, amount, hashLock, timestamp, errTimeLock2},
	}

	for i, td := range testData {
		msg := NewMsgCreateHTLC(td.sender, td.receiver, td.receiverOnOtherChain, td.amount, td.hashLock, td.timestamp, td.timeLock)
		err := msg.ValidateBasic()
		if td.expectPass {
			require.Nil(t, err, "%d: %+v", i, err)
		} else {
			require.NotNil(t, err, "%d: %+v", i, err)
		}
	}
}

func TestMsgCreateHTLCGetSignBytes(t *testing.T) {
	msg := NewMsgCreateHTLC(senderAddr, receiverAddr, receiverOnOtherChain, amount, hashLock, timestamp, timeLock)
	res := msg.GetSignBytes()
	expected := `{"type":"irishub/htlc/MsgCreateHTLC","value":{"amount":{"amount":"10","denom":"iris-atto"},"hash_lock":"6NQTPhqCx04nRueMGThXBup5WKDKRBoI2s+hDEjOJWE=","receiver":"faa1mrehjkgeg75nz2gk7lr7dnxvvtg4497jxss8hq","receiver_on_other_chain":"cmVjZWl2ZXJPbk90aGVyQ2hhaW4=","sender":"faa128nh833v43sggcj65nk7khjka9dwngpl6j29hj","time_lock":"50","timestamp":"1580000000"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgCreateHTLCGetSigners(t *testing.T) {
	msg := NewMsgCreateHTLC(senderAddr, receiverAddr, receiverOnOtherChain, amount, hashLock, timestamp, timeLock)
	res := msg.GetSigners()
	expected := "[51E773C62CAC6084625AA4EDEB5E56E95AE9A03F]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

func TestNewMsgClaimHTLC(t *testing.T) {
	msg := NewMsgClaimHTLC(senderAddr, hashLock, secret)
	require.Equal(t, senderAddr, msg.Sender)
	require.Equal(t, secret, msg.Secret)
	require.Equal(t, hashLock, msg.HashLock)
}

func TestMsgClaimHTLCRoute(t *testing.T) {
	msg := NewMsgClaimHTLC(senderAddr, hashLock, secret)
	require.Equal(t, RouterKey, msg.Route())
}

func TestMsgClaimHTLCValidation(t *testing.T) {
	emptyAddr := sdk.AccAddress{}
	errSecret1 := []byte("xx")
	errSecret2 := []byte("00")
	errHashLock1 := []byte("xx")
	errHashLock2 := []byte("00")

	testData := []struct {
		expectPass bool
		sender     sdk.AccAddress
		secret     []byte
		hashLock   []byte
	}{
		// correct
		{true, senderAddr, secret, hashLock},
		// len(msg.Sender) == 0
		{false, emptyAddr, secret, hashLock},
		// ValidateSecret(msg.Secret)
		{false, senderAddr, errSecret1, hashLock},
		{false, senderAddr, errSecret2, hashLock},
		// ValidateSecretHashLock(msg.SecretHashLock)
		{false, senderAddr, secret, errHashLock1},
		{false, senderAddr, secret, errHashLock2},
	}

	for i, td := range testData {
		msg := NewMsgClaimHTLC(td.sender, td.hashLock, td.secret)
		err := msg.ValidateBasic()

		if td.expectPass {
			require.Nil(t, msg.ValidateBasic(), "%d: %+v", i, err)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "%d", i)
		}
	}
}

func TestMsgClaimHTLCGetSignBytes(t *testing.T) {
	msg := NewMsgClaimHTLC(senderAddr, hashLock, secret)
	res := msg.GetSignBytes()
	expected := `{"type":"irishub/htlc/MsgClaimHTLC","value":{"hash_lock":"6NQTPhqCx04nRueMGThXBup5WKDKRBoI2s+hDEjOJWE=","secret":"X19fYWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXpfX18=","sender":"faa128nh833v43sggcj65nk7khjka9dwngpl6j29hj"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgClaimHTLCGetSigners(t *testing.T) {
	msg := NewMsgClaimHTLC(senderAddr, hashLock, secret)
	res := msg.GetSigners()
	expected := "[51E773C62CAC6084625AA4EDEB5E56E95AE9A03F]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

func TestNewMsgRefundHTLC(t *testing.T) {
	msg := NewMsgRefundHTLC(senderAddr, hashLock)
	require.Equal(t, senderAddr, msg.Sender)
	require.Equal(t, hashLock, msg.HashLock)
}

func TestMsgRefundHTLCRoute(t *testing.T) {
	msg := NewMsgRefundHTLC(senderAddr, hashLock)
	require.Equal(t, RouterKey, msg.Route())
}

func TestMsgRefundHTLCValidation(t *testing.T) {
	emptyAddr := sdk.AccAddress{}
	errHashLock1 := []byte("xx")
	errHashLock2 := []byte("00")

	testData := []struct {
		expectPass bool
		sender     sdk.AccAddress
		hashLock   []byte
	}{
		// correct
		{true, senderAddr, hashLock},
		// len(msg.Sender) == 0
		{false, emptyAddr, hashLock},
		// ValidateSecretHashLock(msg.SecretHashLock)
		{false, senderAddr, errHashLock1},
		{false, senderAddr, errHashLock2},
	}

	for i, td := range testData {
		msg := NewMsgRefundHTLC(td.sender, td.hashLock)
		err := msg.ValidateBasic()

		if td.expectPass {
			require.Nil(t, msg.ValidateBasic(), "%d: %+v", i, err)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "%d", i)
		}
	}
}

func TestMsgRefundHTLCGetSignBytes(t *testing.T) {
	msg := NewMsgRefundHTLC(senderAddr, hashLock)
	res := msg.GetSignBytes()
	expected := `{"type":"irishub/htlc/MsgRefundHTLC","value":{"hash_lock":"6NQTPhqCx04nRueMGThXBup5WKDKRBoI2s+hDEjOJWE=","sender":"faa128nh833v43sggcj65nk7khjka9dwngpl6j29hj"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgRefundHTLCGetSigners(t *testing.T) {
	msg := NewMsgRefundHTLC(senderAddr, hashLock)
	res := msg.GetSigners()
	expected := "[51E773C62CAC6084625AA4EDEB5E56E95AE9A03F]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}
