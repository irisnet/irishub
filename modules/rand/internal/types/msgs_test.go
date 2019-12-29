package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/config"
)

var (
	emptyAddr     sdk.AccAddress
	testAddr      = sdk.AccAddress("testAddr")
	blockInterval = uint64(10)
)

func init() {
	sdk.GetConfig().SetBech32PrefixForAccount(config.GetConfig().GetBech32AccountAddrPrefix(), config.GetConfig().GetBech32AccountPubPrefix())
	sdk.GetConfig().SetBech32PrefixForValidator(config.GetConfig().GetBech32ValidatorAddrPrefix(), config.GetConfig().GetBech32ValidatorPubPrefix())
	sdk.GetConfig().SetBech32PrefixForConsensusNode(config.GetConfig().GetBech32ConsensusAddrPrefix(), config.GetConfig().GetBech32ConsensusPubPrefix())
}

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
		{"basic good", testAddr, blockInterval, true},
	}

	for _, td := range testData {
		msg := NewMsgRequestRand(td.consumer, td.blockInterval)
		if td.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", td.name)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", td.name)
		}
	}
}

func TestMsgRequestRandGetSignBytes(t *testing.T) {
	var msg = NewMsgRequestRand(testAddr, blockInterval)
	res := msg.GetSignBytes()

	expected := `{"type":"irishub/rand/MsgRequestRand","value":{"block_interval":"10","consumer":"faa1w3jhxazpv3j8yxhn3j0"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgRequestRandGetSigners(t *testing.T) {
	var msg = NewMsgRequestRand(testAddr, blockInterval)
	res := msg.GetSigners()

	expected := "[7465737441646472]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}
