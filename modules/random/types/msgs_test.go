package types

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/irisnet/irishub/address"
)

var (
	emptyAddr     sdk.AccAddress
	testAddr      = sdk.AccAddress("testAddr")
	blockInterval = uint64(10)
)

func init() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(address.Bech32PrefixAccAddr, address.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(address.Bech32PrefixValAddr, address.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(address.Bech32PrefixConsAddr, address.Bech32PrefixConsPub)
	config.Seal()
}

func TestNewMsgRequestRandom(t *testing.T) {
	msg := NewMsgRequestRandom(testAddr, blockInterval)

	require.Equal(t, testAddr, msg.Consumer)
	require.Equal(t, blockInterval, msg.BlockInterval)
}

func TestMsgRequestRandomRoute(t *testing.T) {
	// build a MsgRequestRandom
	msg := NewMsgRequestRandom(testAddr, blockInterval)

	require.Equal(t, "rand", msg.Route())
}

func TestMsgRequestRandomValidation(t *testing.T) {
	testData := []struct {
		name          string
		consumer      sdk.AccAddress
		blockInterval uint64
		expectPass    bool
	}{
		{"empty consumer", emptyAddr, blockInterval, false},
		{"basic good", testAddr, blockInterval, true},
	}

	for _, td := range testData {
		msg := NewMsgRequestRandom(td.consumer, td.blockInterval)
		if td.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", td.name)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", td.name)
		}
	}
}

func TestMsgRequestRandomGetSignBytes(t *testing.T) {
	var msg = NewMsgRequestRandom(testAddr, blockInterval)
	res := msg.GetSignBytes()

	expected := `{"type":"irishub/rand/MsgRequestRandom","value":{"block_interval":"10","consumer":"faa1w3jhxazpv3j8yxhn3j0"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgRequestRandomGetSigners(t *testing.T) {
	var msg = NewMsgRequestRandom(testAddr, blockInterval)
	res := msg.GetSigners()

	expected := "[7465737441646472]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}
