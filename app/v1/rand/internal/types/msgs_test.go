package types

import (
	"fmt"
	"testing"

	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

var (
	emptyAddr     sdk.AccAddress
	testAddr      = sdk.AccAddress([]byte("testAddr"))
	blockInterval = uint64(10)
)

func TestNewMsgRequestRand(t *testing.T) {
	msg := NewMsgRequestRand(testAddr, blockInterval)

	require.Equal(t, testAddr, msg.Consumer)
	require.Equal(t, blockInterval, msg.BlockInterval)
}

func TestMsgRequestRandRoute(t *testing.T) {
	// build a MsgRequestRand
	msg := NewMsgRequestRand(testAddr, blockInterval)

	require.Equal(t, "rand", msg.Route())
}

func TestMsgRequestRandValidation(t *testing.T) {
	testData := []struct {
		name          string
		consumer      sdk.AccAddress
		blockInterval uint64
		expectPass    bool
	}{
		{"empty consumer", emptyAddr, blockInterval, false},
		{"valid consumer", testAddr, blockInterval, true},
		{"invalid block interval", testAddr, 0, false},
	}

	for _, td := range testData {
		msg := NewMsgRequestRand(td.consumer, td.blockInterval)
		if td.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", td.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", td.name)
		}
	}
}

func TestMsgRequestRandGetSignBytes(t *testing.T) {
	var msg = NewMsgRequestRand(testAddr, blockInterval)
	res := msg.GetSignBytes()

	expected := "{\"type\":\"irishub/rand/MsgRequestRand\",\"value\":{\"block-interval\":\"10\",\"consumer\":\"faa1w3jhxazpv3j8yxhn3j0\"}}"
	require.Equal(t, expected, string(res))
}

func TestMsgRequestRandGetSigners(t *testing.T) {
	var msg = NewMsgRequestRand(testAddr, blockInterval)
	res := msg.GetSigners()

	expected := "[7465737441646472]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}
