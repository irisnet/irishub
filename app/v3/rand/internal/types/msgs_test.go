package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/irisnet/irishub/types"
)

var (
	emptyAddr     sdk.AccAddress
	testAddr      = sdk.AccAddress([]byte("testAddr"))
	blockInterval = uint64(10)
)

func TestNewMsgRequestRand(t *testing.T) {
	msg := NewMsgRequestRand(testAddr, blockInterval, false)

	require.Equal(t, testAddr, msg.Consumer)
	require.Equal(t, blockInterval, msg.BlockInterval)
}

func TestMsgRequestRandRoute(t *testing.T) {
	// build a MsgRequestRand
	msg := NewMsgRequestRand(testAddr, blockInterval, false)

	require.Equal(t, "rand", msg.Route())
}

func TestMsgRequestRandValidation(t *testing.T) {
	testData := []struct {
		name          string
		consumer      sdk.AccAddress
		blockInterval uint64
		oracle        bool
		expectPass    bool
	}{
		{"empty consumer", emptyAddr, blockInterval, false, false},
		{"basic good", testAddr, blockInterval, false, true},
	}

	for _, td := range testData {
		msg := NewMsgRequestRand(td.consumer, td.blockInterval, td.oracle)
		if td.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", td.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", td.name)
		}
	}
}

func TestMsgRequestRandGetSignBytes(t *testing.T) {
	var msg = NewMsgRequestRand(testAddr, blockInterval, false)
	res := msg.GetSignBytes()

	expected := "{\"type\":\"irishub/rand/MsgRequestRand\",\"value\":{\"block-interval\":\"10\",\"consumer\":\"faa1w3jhxazpv3j8yxhn3j0\"}}"
	require.Equal(t, expected, string(res))
}

func TestMsgRequestRandGetSigners(t *testing.T) {
	var msg = NewMsgRequestRand(testAddr, blockInterval, false)
	res := msg.GetSigners()

	expected := "[7465737441646472]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}
