package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/config"
)

var (
	senderAddr, _        = sdk.AccAddressFromHex(crypto.AddressHash([]byte("sender")).String())
	toAddr, _            = sdk.AccAddressFromHex(crypto.AddressHash([]byte("to")).String())
	receiverOnOtherChain = "receiverOnOtherChain"
	amount               = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10)))
	secret               = HTLCSecret("___abcdefghijklmnopqrstuvwxyz___")
	timestamp            = uint64(1580000000)
	hashLock             = HTLCHashLock(SHA256(append(secret, sdk.Uint64ToBigEndian(timestamp)...)))
	timeLock             = uint64(50)
)

func init() {
	sdk.GetConfig().SetBech32PrefixForAccount(config.GetConfig().GetBech32AccountAddrPrefix(), config.GetConfig().GetBech32AccountPubPrefix())
	sdk.GetConfig().SetBech32PrefixForValidator(config.GetConfig().GetBech32ValidatorAddrPrefix(), config.GetConfig().GetBech32ValidatorPubPrefix())
	sdk.GetConfig().SetBech32PrefixForConsensusNode(config.GetConfig().GetBech32ConsensusAddrPrefix(), config.GetConfig().GetBech32ConsensusPubPrefix())
}

// ----------------------------------------------
// test MsgCreateHTLC
// ----------------------------------------------

func TestNewMsgCreateHTLC(t *testing.T) {
	msg := NewMsgCreateHTLC(senderAddr, toAddr, receiverOnOtherChain, amount, hashLock, timestamp, timeLock)
	require.Equal(t, senderAddr, msg.Sender)
	require.Equal(t, toAddr, msg.To)
	require.Equal(t, receiverOnOtherChain, msg.ReceiverOnOtherChain)
	require.Equal(t, amount, msg.Amount)
	require.Equal(t, hashLock, msg.HashLock)
	require.Equal(t, timestamp, msg.Timestamp)
	require.Equal(t, timeLock, msg.TimeLock)
}

func TestMsgCreateHTLCRoute(t *testing.T) {
	// build a MsgCreateHTLC
	msg := NewMsgCreateHTLC(senderAddr, toAddr, receiverOnOtherChain, amount, hashLock, timestamp, timeLock)
	require.Equal(t, RouterKey, msg.Route())
}

func TestMsgCreateHTLCType(t *testing.T) {
	// build a MsgCreateHTLC
	msg := NewMsgCreateHTLC(senderAddr, toAddr, receiverOnOtherChain, amount, hashLock, timestamp, timeLock)
	require.Equal(t, TypeMsgCreateHTLC, msg.Type())
}

func TestMsgCreateHTLCGetSignBytes(t *testing.T) {
	msg := NewMsgCreateHTLC(senderAddr, toAddr, receiverOnOtherChain, amount, hashLock, timestamp, timeLock)
	res := msg.GetSignBytes()
	expected := `{"type":"irishub/htlc/MsgCreateHTLC","value":{"amount":[{"amount":"10","denom":"stake"}],"hash_lock":"e8d4133e1a82c74e2746e78c19385706ea7958a0ca441a08dacfa10c48ce2561","receiver_on_other_chain":"receiverOnOtherChain","sender":"faa1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgkwnkl5","time_lock":"50","timestamp":"1580000000","to":"faa1vcl2r0llu5pc70cv7enlznzz2lhl2tthnqea0f"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgCreateHTLCGetSigners(t *testing.T) {
	msg := NewMsgCreateHTLC(senderAddr, toAddr, receiverOnOtherChain, amount, hashLock, timestamp, timeLock)
	res := msg.GetSigners()
	expected := "[0A367B92CF0B037DFD89960EE832D56F7FC15168]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

func TestMsgCreateHTLCValidation(t *testing.T) {
	emptyAddr := sdk.AccAddress{}
	errReceiverOnOtherChain := string(make([]byte, 129))
	errAmount := sdk.Coins{}
	errHashLock1 := HTLCHashLock("xx")
	errHashLock2 := HTLCHashLock("00")
	errTimeLock1 := uint64(49)
	errTimeLock2 := uint64(25481)

	testData := []struct {
		expectPass           bool
		sender               sdk.AccAddress
		to                   sdk.AccAddress
		receiverOnOtherChain string
		amount               sdk.Coins
		hashLock             HTLCHashLock
		timestamp            uint64
		timeLock             uint64
	}{
		// correct
		{true, senderAddr, toAddr, receiverOnOtherChain, amount, hashLock, timestamp, timeLock},
		// len(msg.Sender) == 0
		{false, emptyAddr, toAddr, receiverOnOtherChain, amount, hashLock, timestamp, timeLock},
		// len(msg.Recipient) == 0
		{false, senderAddr, emptyAddr, receiverOnOtherChain, amount, hashLock, timestamp, timeLock},
		// len(msg.ToOnOtherChain) > MaxLengthForAddressOnOtherChain
		{false, senderAddr, toAddr, errReceiverOnOtherChain, amount, hashLock, timestamp, timeLock},
		// !msg.OutAmount.IsPositive()
		{false, senderAddr, toAddr, receiverOnOtherChain, errAmount, hashLock, timestamp, timeLock},
		// ValidateSecretHashLock(msg.SecretHashLock)
		{false, senderAddr, toAddr, receiverOnOtherChain, amount, errHashLock1, timestamp, timeLock},
		{false, senderAddr, toAddr, receiverOnOtherChain, amount, errHashLock2, timestamp, timeLock},
		// msg.TimeLock < MinTimeLock || msg.TimeLock > MaxTimeLock
		{false, senderAddr, toAddr, receiverOnOtherChain, amount, hashLock, timestamp, errTimeLock1},
		{false, senderAddr, toAddr, receiverOnOtherChain, amount, hashLock, timestamp, errTimeLock2},
	}

	for i, td := range testData {
		t.Run(string(i), func(t *testing.T) {
			msg := NewMsgCreateHTLC(td.sender, td.to, td.receiverOnOtherChain, td.amount, td.hashLock, td.timestamp, td.timeLock)
			err := msg.ValidateBasic()
			if td.expectPass {
				require.NoError(t, err, "%d: %+v", i, err)
			} else {
				require.Error(t, err, "%d: %+v", i, err)
			}
		})
	}
}

// ----------------------------------------------
// test MsgClaimHTLC
// ----------------------------------------------

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

func TestMsgClaimHTLCType(t *testing.T) {
	msg := NewMsgClaimHTLC(senderAddr, hashLock, secret)
	require.Equal(t, TypeMsgClaimHTLC, msg.Type())
}

func TestMsgClaimHTLCGetSignBytes(t *testing.T) {
	msg := NewMsgClaimHTLC(senderAddr, hashLock, secret)
	res := msg.GetSignBytes()
	expected := `{"type":"irishub/htlc/MsgClaimHTLC","value":{"hash_lock":"e8d4133e1a82c74e2746e78c19385706ea7958a0ca441a08dacfa10c48ce2561","secret":"5f5f5f6162636465666768696a6b6c6d6e6f707172737475767778797a5f5f5f","sender":"faa1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgkwnkl5"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgClaimHTLCGetSigners(t *testing.T) {
	msg := NewMsgClaimHTLC(senderAddr, hashLock, secret)
	res := msg.GetSigners()
	expected := "[0A367B92CF0B037DFD89960EE832D56F7FC15168]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

func TestMsgClaimHTLCValidation(t *testing.T) {
	emptyAddr := sdk.AccAddress{}
	errSecret1 := HTLCSecret("xx")
	errSecret2 := HTLCSecret("00")
	errHashLock1 := HTLCHashLock("xx")
	errHashLock2 := HTLCHashLock("00")

	testData := []struct {
		expectPass bool
		sender     sdk.AccAddress
		secret     HTLCSecret
		hashLock   HTLCHashLock
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
		t.Run(string(i), func(t *testing.T) {
			msg := NewMsgClaimHTLC(td.sender, td.hashLock, td.secret)
			err := msg.ValidateBasic()

			if td.expectPass {
				require.NoError(t, msg.ValidateBasic(), "%d: %+v", i, err)
			} else {
				require.Error(t, msg.ValidateBasic(), "%d", i)
			}
		})
	}
}

// ----------------------------------------------
// test MsgRefundHTLC
// ----------------------------------------------

func TestNewMsgRefundHTLC(t *testing.T) {
	msg := NewMsgRefundHTLC(senderAddr, hashLock)
	require.Equal(t, senderAddr, msg.Sender)
	require.Equal(t, hashLock, msg.HashLock)
}

func TestMsgRefundHTLCRoute(t *testing.T) {
	msg := NewMsgRefundHTLC(senderAddr, hashLock)
	require.Equal(t, RouterKey, msg.Route())
}

func TestMsgRefundHTLCType(t *testing.T) {
	msg := NewMsgRefundHTLC(senderAddr, hashLock)
	require.Equal(t, TypeMsgRefundHTLC, msg.Type())
}

func TestMsgRefundHTLCGetSignBytes(t *testing.T) {
	msg := NewMsgRefundHTLC(senderAddr, hashLock)
	res := msg.GetSignBytes()
	expected := `{"type":"irishub/htlc/MsgRefundHTLC","value":{"hash_lock":"e8d4133e1a82c74e2746e78c19385706ea7958a0ca441a08dacfa10c48ce2561","sender":"faa1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgkwnkl5"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgRefundHTLCGetSigners(t *testing.T) {
	msg := NewMsgRefundHTLC(senderAddr, hashLock)
	res := msg.GetSigners()
	expected := "[0A367B92CF0B037DFD89960EE832D56F7FC15168]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

func TestMsgRefundHTLCValidation(t *testing.T) {
	emptyAddr := sdk.AccAddress{}
	errHashLock1 := HTLCHashLock("xx")
	errHashLock2 := HTLCHashLock("00")

	testData := []struct {
		expectPass bool
		sender     sdk.AccAddress
		hashLock   HTLCHashLock
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
		t.Run(string(i), func(t *testing.T) {
			msg := NewMsgRefundHTLC(td.sender, td.hashLock)
			err := msg.ValidateBasic()
			if td.expectPass {
				require.NoError(t, msg.ValidateBasic(), "%d: %+v", i, err)
			} else {
				require.Error(t, msg.ValidateBasic(), "%d", i)
			}
		})
	}
}
