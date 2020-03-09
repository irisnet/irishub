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
	serviceFeeCap = sdk.NewCoins(sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(1000000000000000000)))
)

func TestNewMsgRequestRand(t *testing.T) {
	msg := NewMsgRequestRand(testAddr, blockInterval, false, serviceFeeCap)

	require.Equal(t, testAddr, msg.Consumer)
	require.Equal(t, blockInterval, msg.BlockInterval)
}

func TestMsgRequestRandRoute(t *testing.T) {
	// build a MsgRequestRand
	msg := NewMsgRequestRand(testAddr, blockInterval, false, serviceFeeCap)

	require.Equal(t, "rand", msg.Route())
}

func TestMsgRequestRandValidation(t *testing.T) {
	testData := []struct {
		name          string
		consumer      sdk.AccAddress
		blockInterval uint64
		oracle        bool
		serviceFeeCap sdk.Coins
		expectPass    bool
	}{
		{"empty consumer", emptyAddr, blockInterval, false, serviceFeeCap, false},
		{"basic good", testAddr, blockInterval, false, serviceFeeCap, true},
	}

	for _, td := range testData {
		msg := NewMsgRequestRand(td.consumer, td.blockInterval, td.oracle, td.serviceFeeCap)
		if td.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", td.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", td.name)
		}
	}
}

func TestMsgRequestRandGetSignBytes(t *testing.T) {
	var msg = NewMsgRequestRand(testAddr, blockInterval, false, serviceFeeCap)
	res := msg.GetSignBytes()

	expected := "{\"type\":\"irishub/rand/MsgRequestRand\",\"value\":{\"block_interval\":\"10\",\"consumer\":\"faa1w3jhxazpv3j8yxhn3j0\",\"oracle\":false,\"service_fee_cap\":[{\"amount\":\"1000000000000000000\",\"denom\":\"iris-atto\"}]}}"
	require.Equal(t, expected, string(res))
}

func TestMsgRequestRandGetSigners(t *testing.T) {
	var msg = NewMsgRequestRand(testAddr, blockInterval, false, serviceFeeCap)
	res := msg.GetSigners()

	expected := "[7465737441646472]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}
