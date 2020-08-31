package types

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var (
	emptyAddr     sdk.AccAddress
	testAddr      = sdk.AccAddress([]byte("testAddr"))
	blockInterval = uint64(10)
	serviceFeeCap = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000000000000000000)))
)

func TestNewMsgRequestRandom(t *testing.T) {
	msg := NewMsgRequestRandom(testAddr, blockInterval, false, serviceFeeCap)

	require.Equal(t, testAddr, msg.Consumer)
	require.Equal(t, blockInterval, msg.BlockInterval)
}

func TestMsgRequestRandomRoute(t *testing.T) {
	// build a MsgRequestRandom
	msg := NewMsgRequestRandom(testAddr, blockInterval, false, serviceFeeCap)

	require.Equal(t, "rand", msg.Route())
}

func TestMsgRequestRandomValidation(t *testing.T) {
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
		msg := NewMsgRequestRandom(td.consumer, td.blockInterval, td.oracle, td.serviceFeeCap)
		if td.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", td.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", td.name)
		}
	}
}

func TestMsgRequestRandomGetSignBytes(t *testing.T) {
	var msg = NewMsgRequestRandom(testAddr, blockInterval, true, serviceFeeCap)
	res := msg.GetSignBytes()

	expected := "{\"type\":\"irishub/rand/MsgRequestRandom\",\"value\":{\"block_interval\":\"10\",\"consumer\":\"cosmos1w3jhxazpv3j8y5jww2c\",\"oracle\":true,\"service_fee_cap\":[{\"amount\":\"1000000000000000000\",\"denom\":\"iris-atto\"}]}}"
	require.Equal(t, expected, string(res))
}

func TestMsgRequestRandomGetSigners(t *testing.T) {
	var msg = NewMsgRequestRandom(testAddr, blockInterval, false, serviceFeeCap)
	res := msg.GetSigners()

	expected := "[7465737441646472]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}
