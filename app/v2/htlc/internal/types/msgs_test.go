package types

import (
	"encoding/hex"
	"fmt"
	"testing"

	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

var (
	emptyAddr            sdk.AccAddress
	senderAddr           = sdk.AccAddress([]byte("sender"))
	receiverAddr         = sdk.AccAddress([]byte("receiver"))
	receiverOnOtherChain = []byte("receiverOnOtherChain")
	outAmount            = sdk.NewCoin("iris", sdk.NewInt(10))
	inAmount             = uint64(100)
	secret               = []byte("secret")
	timestamp            = uint64(1580000000)
	secretHashLock       = hex.EncodeToString(sdk.SHA256(append(secret, sdk.Uint64ToBigEndian(timestamp)...)))
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
	// TODO
}

func TestMsgCreateHTLCGetSignBytes(t *testing.T) {
	msg := NewMsgCreateHTLC(senderAddr, receiverAddr, receiverOnOtherChain, outAmount, inAmount, secretHashLock, timestamp, timeLock)
	res := msg.GetSignBytes()

	expected := `{"type":"irishub/htlc/MsgCreateHTLC","value":{"in_amount":"100","out_amount":{"amount":"10","denom":"iris"},"receiver":"faa1wfjkxetfwejhytm7qdx","receiver_on_other_chain":"cmVjZWl2ZXJPbk90aGVyQ2hhaW4=","secret_hash_lock":"20af01f5feea9aaf81d3633fdd6edcc32263a3cf908e4fada5a64d305d37d59f","sender":"faa1wdjkuer9wg8a5rfv","time_lock":"50","timestamp":"1580000000"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgCreateHTLCGetSigners(t *testing.T) {
	msg := NewMsgCreateHTLC(senderAddr, receiverAddr, receiverOnOtherChain, outAmount, inAmount, secretHashLock, timestamp, timeLock)
	res := msg.GetSigners()

	expected := "[73656E646572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}
