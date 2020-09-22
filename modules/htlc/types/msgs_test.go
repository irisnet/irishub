package types

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/tmhash"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	sender               = sdk.AccAddress([]byte("sender"))
	recipient            = sdk.AccAddress([]byte("recipient"))
	receiverOnOtherChain = "receiverOnOtherChain"
	amount               = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10)))
	secret               = tmbytes.HexBytes(tmhash.Sum([]byte("secret")))
	timestamp            = uint64(1580000000)
	hashLock             = tmbytes.HexBytes(tmhash.Sum(append(secret, sdk.Uint64ToBigEndian(timestamp)...)))
	timeLock             = uint64(50)
)

// TestNewMsgCreateHTLC tests constructor for MsgCreateHTLC
func TestNewMsgCreateHTLC(t *testing.T) {
	msg := NewMsgCreateHTLC(sender, recipient, receiverOnOtherChain, amount, hashLock, timestamp, timeLock)

	require.Equal(t, sender, msg.Sender)
	require.Equal(t, recipient, msg.To)
	require.Equal(t, receiverOnOtherChain, msg.ReceiverOnOtherChain)
	require.Equal(t, amount, msg.Amount)
	require.Equal(t, hashLock, msg.HashLock)
	require.Equal(t, timestamp, msg.Timestamp)
	require.Equal(t, timeLock, msg.TimeLock)
}

// TestMsgCreateHTLCRoute tests Route for MsgCreateHTLC
func TestMsgCreateHTLCRoute(t *testing.T) {
	msg := NewMsgCreateHTLC(sender, recipient, receiverOnOtherChain, amount, hashLock, timestamp, timeLock)
	require.Equal(t, "htlc", msg.Route())
}

// TestMsgCreateHTLCType tests Type for MsgCreateHTLC
func TestMsgCreateHTLCType(t *testing.T) {
	msg := NewMsgCreateHTLC(sender, recipient, receiverOnOtherChain, amount, hashLock, timestamp, timeLock)
	require.Equal(t, "create_htlc", msg.Type())
}

// TestMsgCreateHTLCValidation tests ValidateBasic for MsgCreateHTLC
func TestMsgCreateHTLCValidation(t *testing.T) {
	emptyAddr := sdk.AccAddress{}

	invalidReceiverOnOtherChain := strings.Repeat("r", 129)
	invalidAmount := sdk.Coins{}
	invalidHashLock := []byte("0x")
	invalidSmallTimeLock := uint64(49)
	invalidLargeTimeLock := uint64(25481)

	testMsgs := []MsgCreateHTLC{
		NewMsgCreateHTLC(sender, recipient, receiverOnOtherChain, amount, hashLock, timestamp, timeLock),             // valid msg
		NewMsgCreateHTLC(emptyAddr, recipient, receiverOnOtherChain, amount, hashLock, timestamp, timeLock),          // missing sender
		NewMsgCreateHTLC(sender, emptyAddr, receiverOnOtherChain, amount, hashLock, timestamp, timeLock),             // missing recipient
		NewMsgCreateHTLC(sender, recipient, invalidReceiverOnOtherChain, amount, hashLock, timestamp, timeLock),      // too long receiver on other chain
		NewMsgCreateHTLC(sender, recipient, receiverOnOtherChain, invalidAmount, hashLock, timestamp, timeLock),      // invalid amount
		NewMsgCreateHTLC(sender, recipient, receiverOnOtherChain, amount, invalidHashLock, timestamp, timeLock),      // invalid hash lock
		NewMsgCreateHTLC(sender, recipient, receiverOnOtherChain, amount, hashLock, timestamp, invalidSmallTimeLock), // too small time lock
		NewMsgCreateHTLC(sender, recipient, receiverOnOtherChain, amount, hashLock, timestamp, invalidLargeTimeLock), // too large time lock
	}

	testCases := []struct {
		msg     MsgCreateHTLC
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], false, "missing sender"},
		{testMsgs[2], false, "missing recipient"},
		{testMsgs[3], false, "too long receiver on other chain"},
		{testMsgs[4], false, "invalid amount"},
		{testMsgs[5], false, "invalid hash lock"},
		{testMsgs[6], false, "too small time lock"},
		{testMsgs[7], false, "too large time lock"},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Msg %d failed: %v", i, err)
		} else {
			require.Error(t, err, "Invalid Msg %d passed: %s", i, tc.errMsg)
		}
	}
}

// TestMsgCreateHTLCGetSignBytes tests GetSignBytes for MsgCreateHTLC
func TestMsgCreateHTLCGetSignBytes(t *testing.T) {
	msg := NewMsgCreateHTLC(sender, recipient, receiverOnOtherChain, amount, hashLock, timestamp, timeLock)
	res := msg.GetSignBytes()

	expected := `{"type":"irismod/htlc/MsgCreateHTLC","value":{"amount":[{"amount":"10","denom":"stake"}],"hash_lock":"6F4ECE9B22CFC1CF39C9C73DD2D35867A8EC97C48A9C2F664FE5287865A18C2E","receiver_on_other_chain":"receiverOnOtherChain","sender":"cosmos1wdjkuer9wgh76ts6","time_lock":"50","timestamp":"1580000000","to":"cosmos1wfjkx6tsd9jkuaqhtdv59"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgCreateHTLCGetSigners tests GetSigners for MsgCreateHTLC
func TestMsgCreateHTLCGetSigners(t *testing.T) {
	msg := NewMsgCreateHTLC(sender, recipient, receiverOnOtherChain, amount, hashLock, timestamp, timeLock)
	res := msg.GetSigners()

	expected := "[73656E646572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// TestNewMsgClaimHTLC tests constructor for MsgClaimHTLC
func TestNewMsgClaimHTLC(t *testing.T) {
	msg := NewMsgClaimHTLC(sender, hashLock, secret)

	require.Equal(t, sender, msg.Sender)
	require.Equal(t, secret, msg.Secret)
	require.Equal(t, hashLock, msg.HashLock)
}

// TestMsgClaimHTLCRoute tests Route for MsgClaimHTLC
func TestMsgClaimHTLCRoute(t *testing.T) {
	msg := NewMsgClaimHTLC(sender, hashLock, secret)
	require.Equal(t, "htlc", msg.Route())
}

// TestMsgClaimHTLCType tests Type for MsgClaimHTLC
func TestMsgClaimHTLCType(t *testing.T) {
	msg := NewMsgClaimHTLC(sender, hashLock, secret)
	require.Equal(t, "claim_htlc", msg.Type())
}

// TestMsgClaimHTLCValidation tests ValidateBasic for MsgClaimHTLC
func TestMsgClaimHTLCValidation(t *testing.T) {
	emptyAddr := sdk.AccAddress{}

	invalidHashLock := []byte("0x")
	invalidSecret := []byte("0x")

	testMsgs := []MsgClaimHTLC{
		NewMsgClaimHTLC(sender, hashLock, secret),        // valid msg
		NewMsgClaimHTLC(emptyAddr, hashLock, secret),     // missing sender
		NewMsgClaimHTLC(sender, invalidHashLock, secret), // invalid hash lock
		NewMsgClaimHTLC(sender, hashLock, invalidSecret), // invalid secret
	}

	testCases := []struct {
		msg     MsgClaimHTLC
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], false, "missing sender"},
		{testMsgs[2], false, "invalid hash lock"},
		{testMsgs[3], false, "invalid secret"},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Msg %d failed: %v", i, err)
		} else {
			require.Error(t, err, "Invalid Msg %d passed: %s", i, tc.errMsg)
		}
	}
}

// TestMsgClaimHTLCGetSignBytes tests GetSignBytes for MsgClaimHTLC
func TestMsgClaimHTLCGetSignBytes(t *testing.T) {
	msg := NewMsgClaimHTLC(sender, hashLock, secret)
	res := msg.GetSignBytes()

	expected := `{"type":"irismod/htlc/MsgClaimHTLC","value":{"hash_lock":"6F4ECE9B22CFC1CF39C9C73DD2D35867A8EC97C48A9C2F664FE5287865A18C2E","secret":"2BB80D537B1DA3E38BD30361AA855686BDE0EACD7162FEF6A25FE97BF527A25B","sender":"cosmos1wdjkuer9wgh76ts6"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgClaimHTLCGetSigners tests GetSigners for MsgClaimHTLC
func TestMsgClaimHTLCGetSigners(t *testing.T) {
	msg := NewMsgClaimHTLC(sender, hashLock, secret)
	res := msg.GetSigners()

	expected := "[73656E646572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// TestNewMsgRefundHTLC tests constructor for MsgRefundHTLC
func TestNewMsgRefundHTLC(t *testing.T) {
	msg := NewMsgRefundHTLC(sender, hashLock)

	require.Equal(t, sender, msg.Sender)
	require.Equal(t, hashLock, msg.HashLock)
}

// TestMsgRefundHTLCRoute tests Route for MsgRefundHTLC
func TestMsgRefundHTLCRoute(t *testing.T) {
	msg := NewMsgRefundHTLC(sender, hashLock)
	require.Equal(t, "htlc", msg.Route())
}

// TestMsgRefundHTLCType tests Type for MsgRefundHTLC
func TestMsgRefundHTLCType(t *testing.T) {
	msg := NewMsgRefundHTLC(sender, hashLock)
	require.Equal(t, "refund_htlc", msg.Type())
}

// TestMsgRefundHTLCValidation tests ValidateBasic for MsgRefundHTLC
func TestMsgRefundHTLCValidation(t *testing.T) {
	emptyAddr := sdk.AccAddress{}

	invalidHashLock := []byte("0x")

	testMsgs := []MsgRefundHTLC{
		NewMsgRefundHTLC(sender, hashLock),        // valid msg
		NewMsgRefundHTLC(emptyAddr, hashLock),     // missing sender
		NewMsgRefundHTLC(sender, invalidHashLock), // invalid hash lock
	}

	testCases := []struct {
		msg     MsgRefundHTLC
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], false, "missing sender"},
		{testMsgs[2], false, "invalid hash lock"},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Msg %d failed: %v", i, err)
		} else {
			require.Error(t, err, "Invalid Msg %d passed: %s", i, tc.errMsg)
		}
	}
}

// TestMsgRefundHTLCGetSignBytes tests GetSignBytes for MsgRefundHTLC
func TestMsgRefundHTLCGetSignBytes(t *testing.T) {
	msg := NewMsgRefundHTLC(sender, hashLock)
	res := msg.GetSignBytes()

	expected := `{"type":"irismod/htlc/MsgRefundHTLC","value":{"hash_lock":"6F4ECE9B22CFC1CF39C9C73DD2D35867A8EC97C48A9C2F664FE5287865A18C2E","sender":"cosmos1wdjkuer9wgh76ts6"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgRefundHTLCGetSigners tests GetSigners for MsgRefundHTLC
func TestMsgRefundHTLCGetSigners(t *testing.T) {
	msg := NewMsgRefundHTLC(sender, hashLock)
	res := msg.GetSigners()

	expected := "[73656E646572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}
