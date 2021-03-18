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
	emptyAddr            = ""
	sender               = sdk.AccAddress(tmhash.SumTruncated([]byte("sender")))
	senderStr            = sender.String()
	recipient            = sdk.AccAddress(tmhash.SumTruncated([]byte("recipient")))
	recipientStr         = recipient.String()
	receiverOnOtherChain = "receiverOnOtherChain"
	senderOnOtherChain   = "senderOnOtherChain"
	amount               = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10)))
	secret               = tmbytes.HexBytes(tmhash.Sum([]byte("secret")))
	secretStr            = secret.String()
	timestamp            = uint64(1580000000)
	hashLock             = tmbytes.HexBytes(tmhash.Sum(append(secret, sdk.Uint64ToBigEndian(timestamp)...)))
	hashLockStr          = hashLock.String()
	id                   = tmbytes.HexBytes(tmhash.Sum(append(append(append(hashLock, sender...), recipient...), []byte(amount.String())...)))
	idStr                = id.String()
	timeLock             = uint64(50)
	transfer             = true
	notTransfer          = false
)

// TestNewMsgCreateHTLC tests constructor for MsgCreateHTLC
func TestNewMsgCreateHTLC(t *testing.T) {
	msg := NewMsgCreateHTLC(senderStr, recipientStr, receiverOnOtherChain, senderOnOtherChain, amount, hashLockStr, timestamp, timeLock, notTransfer)

	require.Equal(t, senderStr, msg.Sender)
	require.Equal(t, recipientStr, msg.To)
	require.Equal(t, receiverOnOtherChain, msg.ReceiverOnOtherChain)
	require.Equal(t, amount, msg.Amount)
	require.Equal(t, hashLockStr, msg.HashLock)
	require.Equal(t, timestamp, msg.Timestamp)
	require.Equal(t, timeLock, msg.TimeLock)
	require.Equal(t, notTransfer, msg.Transfer)
}

// TestMsgCreateHTLCRoute tests Route for MsgCreateHTLC
func TestMsgCreateHTLCRoute(t *testing.T) {
	msg := NewMsgCreateHTLC(senderStr, recipientStr, receiverOnOtherChain, senderOnOtherChain, amount, hashLockStr, timestamp, timeLock, notTransfer)
	require.Equal(t, "htlc", msg.Route())
}

// TestMsgCreateHTLCType tests Type for MsgCreateHTLC
func TestMsgCreateHTLCType(t *testing.T) {
	msg := NewMsgCreateHTLC(senderStr, recipientStr, receiverOnOtherChain, senderOnOtherChain, amount, hashLockStr, timestamp, timeLock, notTransfer)
	require.Equal(t, "create_htlc", msg.Type())
}

// TestMsgCreateHTLCValidation tests ValidateBasic for MsgCreateHTLC
func TestMsgCreateHTLCValidation(t *testing.T) {
	invalidReceiverOnOtherChain := strings.Repeat("r", 129)
	invalidSenderOnOtherChain := strings.Repeat("r", 129)
	invalidAmount := sdk.Coins{}
	invalidHashLock := "0x"
	invalidSmallTimeLock := uint64(49)
	invalidLargeTimeLock := uint64(25481)

	testMsgs := []MsgCreateHTLC{
		NewMsgCreateHTLC(senderStr, recipientStr, receiverOnOtherChain, senderOnOtherChain, amount, hashLockStr, timestamp, timeLock, notTransfer),             // valid htlc msg
		NewMsgCreateHTLC(senderStr, recipientStr, receiverOnOtherChain, senderOnOtherChain, amount, hashLockStr, timestamp, timeLock, transfer),                // valid htlt msg
		NewMsgCreateHTLC(emptyAddr, recipientStr, receiverOnOtherChain, senderOnOtherChain, amount, hashLockStr, timestamp, timeLock, notTransfer),             // missing sender
		NewMsgCreateHTLC(senderStr, emptyAddr, receiverOnOtherChain, senderOnOtherChain, amount, hashLockStr, timestamp, timeLock, notTransfer),                // missing recipient
		NewMsgCreateHTLC(senderStr, recipientStr, invalidReceiverOnOtherChain, senderOnOtherChain, amount, hashLockStr, timestamp, timeLock, notTransfer),      // too long receiver on other chain
		NewMsgCreateHTLC(senderStr, recipientStr, receiverOnOtherChain, invalidSenderOnOtherChain, amount, hashLockStr, timestamp, timeLock, notTransfer),      // too long sender on other chain
		NewMsgCreateHTLC(senderStr, recipientStr, receiverOnOtherChain, senderOnOtherChain, invalidAmount, hashLockStr, timestamp, timeLock, notTransfer),      // invalid amount
		NewMsgCreateHTLC(senderStr, recipientStr, receiverOnOtherChain, senderOnOtherChain, amount, invalidHashLock, timestamp, timeLock, notTransfer),         // invalid hash lock
		NewMsgCreateHTLC(senderStr, recipientStr, receiverOnOtherChain, senderOnOtherChain, amount, hashLockStr, timestamp, invalidSmallTimeLock, notTransfer), // too small time lock
		NewMsgCreateHTLC(senderStr, recipientStr, receiverOnOtherChain, senderOnOtherChain, amount, hashLockStr, timestamp, invalidLargeTimeLock, notTransfer), // too large time lock
	}

	testCases := []struct {
		msg     MsgCreateHTLC
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, "valid htlc"},
		{testMsgs[1], true, "valid htlt"},
		{testMsgs[2], false, "missing sender"},
		{testMsgs[3], false, "missing recipient"},
		{testMsgs[4], false, "too long receiver on other chain"},
		{testMsgs[5], false, "too long sender on other chain"},
		{testMsgs[6], false, "invalid amount"},
		{testMsgs[7], false, "invalid hash lock"},
		{testMsgs[8], false, "too small time lock"},
		{testMsgs[9], false, "too large time lock"},
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
	msg := NewMsgCreateHTLC(senderStr, recipientStr, receiverOnOtherChain, senderOnOtherChain, amount, hashLockStr, timestamp, timeLock, notTransfer)
	res := msg.GetSignBytes()

	expected := `{"type":"irismod/htlc/MsgCreateHTLC","value":{"amount":[{"amount":"10","denom":"stake"}],"hash_lock":"6F4ECE9B22CFC1CF39C9C73DD2D35867A8EC97C48A9C2F664FE5287865A18C2E","receiver_on_other_chain":"receiverOnOtherChain","sender":"cosmos1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgmr4lac","sender_on_other_chain":"senderOnOtherChain","time_lock":"50","timestamp":"1580000000","to":"cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgCreateHTLCGetSigners tests GetSigners for MsgCreateHTLC
func TestMsgCreateHTLCGetSigners(t *testing.T) {
	msg := NewMsgCreateHTLC(senderStr, recipientStr, receiverOnOtherChain, senderOnOtherChain, amount, hashLockStr, timestamp, timeLock, notTransfer)
	res := msg.GetSigners()

	expected := "[0A367B92CF0B037DFD89960EE832D56F7FC15168]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// TestNewMsgClaimHTLC tests constructor for MsgClaimHTLC
func TestNewMsgClaimHTLC(t *testing.T) {
	msg := NewMsgClaimHTLC(senderStr, idStr, secretStr)

	require.Equal(t, senderStr, msg.Sender)
	require.Equal(t, secretStr, msg.Secret)
	require.Equal(t, idStr, msg.Id)
}

// TestMsgClaimHTLCRoute tests Route for MsgClaimHTLC
func TestMsgClaimHTLCRoute(t *testing.T) {
	msg := NewMsgClaimHTLC(senderStr, idStr, secretStr)
	require.Equal(t, "htlc", msg.Route())
}

// TestMsgClaimHTLCType tests Type for MsgClaimHTLC
func TestMsgClaimHTLCType(t *testing.T) {
	msg := NewMsgClaimHTLC(senderStr, idStr, secret.String())
	require.Equal(t, "claim_htlc", msg.Type())
}

// TestMsgClaimHTLCValidation tests ValidateBasic for MsgClaimHTLC
func TestMsgClaimHTLCValidation(t *testing.T) {
	invalidID := "0x"
	invalidSecret := "0x"

	testMsgs := []MsgClaimHTLC{
		NewMsgClaimHTLC(senderStr, idStr, secret.String()),     // valid msg
		NewMsgClaimHTLC(emptyAddr, idStr, secret.String()),     // missing sender
		NewMsgClaimHTLC(senderStr, invalidID, secret.String()), // invalid id
		NewMsgClaimHTLC(senderStr, idStr, invalidSecret),       // invalid secret
	}

	testCases := []struct {
		msg     MsgClaimHTLC
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, "valid msg"},
		{testMsgs[1], false, "missing sender"},
		{testMsgs[2], false, "invalid id"},
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
	msg := NewMsgClaimHTLC(senderStr, idStr, secretStr)
	res := msg.GetSignBytes()

	expected := `{"type":"irismod/htlc/MsgClaimHTLC","value":{"id":"B94EFE2C859EDADE7F3F6CAF5D7A1CE388D65B9E63CB6CE0B824117F117695A7","secret":"2BB80D537B1DA3E38BD30361AA855686BDE0EACD7162FEF6A25FE97BF527A25B","sender":"cosmos1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgmr4lac"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgClaimHTLCGetSigners tests GetSigners for MsgClaimHTLC
func TestMsgClaimHTLCGetSigners(t *testing.T) {
	msg := NewMsgClaimHTLC(senderStr, idStr, secretStr)
	res := msg.GetSigners()

	expected := "[0A367B92CF0B037DFD89960EE832D56F7FC15168]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}
