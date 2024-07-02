package types

import (
	"fmt"
	"testing"

	"github.com/cometbft/cometbft/crypto/tmhash"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var (
	emptyAddr     = ""
	testAddr      = sdk.AccAddress(tmhash.SumTruncated([]byte("testAddr"))).String()
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

	require.Equal(t, "random", msg.Route())
}

func TestMsgRequestRandomValidation(t *testing.T) {
	testData := []struct {
		name          string
		consumer      string
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
	msg := NewMsgRequestRandom(testAddr, blockInterval, true, serviceFeeCap)
	res := msg.GetSignBytes()

	expected := fmt.Sprintf(
		"{\"type\":\"irismod/random/MsgRequestRandom\",\"value\":{\"block_interval\":\"10\",\"consumer\":\"cosmos133ee7f22kzn7khtdw8d72dgyre0txe5zll7d5w\",\"oracle\":true,\"service_fee_cap\":[{\"amount\":\"1000000000000000000\",\"denom\":\"%s\"}]}}",
		sdk.DefaultBondDenom,
	)
	require.Equal(t, expected, string(res))
}

func TestMsgRequestRandomGetSigners(t *testing.T) {
	msg := NewMsgRequestRandom(testAddr, blockInterval, false, serviceFeeCap)
	res := msg.GetSigners()

	expected := "[8C739F254AB0A7EB5D6D71DBE535041E5EB36682]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}
