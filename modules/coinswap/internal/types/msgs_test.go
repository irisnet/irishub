package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/config"
)

// nolint: deadcode unused
var (
	sender    sdk.AccAddress
	recipient sdk.AccAddress

	amt = sdk.NewInt(100)

	denom0   = "atom"
	denom1   = "btc"
	unidenom = FormatUniABSPrefix + "btc"

	input             = sdk.NewCoin(denom0, sdk.NewInt(1000))
	output            = sdk.NewCoin(denom1, sdk.NewInt(500))
	withdrawLiquidity = sdk.NewCoin(unidenom, sdk.NewInt(500))
	deadline          = int64(1580000000)

	emptyAddr sdk.AccAddress
	emptyTime int64
)

func init() {
	sdk.GetConfig().SetBech32PrefixForAccount(config.GetConfig().GetBech32AccountAddrPrefix(), config.GetConfig().GetBech32AccountPubPrefix())
	sdk.GetConfig().SetBech32PrefixForValidator(config.GetConfig().GetBech32ValidatorAddrPrefix(), config.GetConfig().GetBech32ValidatorPubPrefix())
	sdk.GetConfig().SetBech32PrefixForConsensusNode(config.GetConfig().GetBech32ConsensusAddrPrefix(), config.GetConfig().GetBech32ConsensusPubPrefix())

	sender, _ = sdk.AccAddressFromBech32("faa128nh833v43sggcj65nk7khjka9dwngpl6j29hj")
	recipient, _ = sdk.AccAddressFromBech32("faa1mrehjkgeg75nz2gk7lr7dnxvvtg4497jxss8hq")
}

// ----------------------------------------------
// test MsgSwapOrder
// ----------------------------------------------

func TestNewMsgSwapOrder(t *testing.T) {
	inputTmp := Input{Address: sender, Coin: input}
	outputTmp := Output{Address: sender, Coin: output}
	msg := NewMsgSwapOrder(
		inputTmp,
		outputTmp,
		deadline,
		true,
	)
	require.Equal(t, inputTmp, msg.Input)
	require.Equal(t, outputTmp, msg.Output)
	require.Equal(t, deadline, msg.Deadline)
	require.Equal(t, true, msg.IsBuyOrder)
}

func TestMsgSwapOrderRoute(t *testing.T) {
	msg := NewMsgSwapOrder(
		Input{Address: sender, Coin: input},
		Output{Address: sender, Coin: output},
		deadline,
		true,
	)
	require.Equal(t, RouterKey, msg.Route())
}

func TestMsgSwapOrderType(t *testing.T) {
	msg := NewMsgSwapOrder(
		Input{Address: sender, Coin: input},
		Output{Address: sender, Coin: output},
		deadline,
		true,
	)
	require.Equal(t, MsgTypeSwapOrder, msg.Type())
}

func TestMsgSwapOrderGetSignBytes(t *testing.T) {
	msg := NewMsgSwapOrder(
		Input{Address: sender, Coin: input},
		Output{Address: sender, Coin: output},
		deadline,
		true,
	)
	res := msg.GetSignBytes()
	expected := `{"type":"irishub/coinswap/MsgSwapOrder","value":{"deadline":"1580000000","input":{"address":"faa128nh833v43sggcj65nk7khjka9dwngpl6j29hj","coin":{"amount":"1000","denom":"atom"}},"is_buy_order":true,"output":{"address":"faa128nh833v43sggcj65nk7khjka9dwngpl6j29hj","coin":{"amount":"500","denom":"btc"}}}}`
	require.Equal(t, expected, string(res))
}

func TestMsgSwapOrderGetSigners(t *testing.T) {
	msg := NewMsgSwapOrder(
		Input{Address: sender, Coin: input},
		Output{Address: sender, Coin: output},
		deadline,
		true,
	)
	res := msg.GetSigners()
	expected := "[51E773C62CAC6084625AA4EDEB5E56E95AE9A03F]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// test ValidateBasic for MsgSwapOrder
func TestMsgSwapOrderValidation(t *testing.T) {
	tests := []struct {
		name       string
		expectPass bool
		msg        MsgSwapOrder
	}{
		{
			"no input coin",
			false,
			NewMsgSwapOrder(
				Input{Address: sender},
				Output{recipient, output},
				deadline,
				true,
			),
		}, {
			"zero input coin",
			false,
			NewMsgSwapOrder(
				Input{Address: sender, Coin: sdk.NewCoin(denom0, sdk.ZeroInt())},
				Output{Address: recipient, Coin: output},
				deadline,
				true,
			),
		}, {
			"no output coin",
			false,
			NewMsgSwapOrder(
				Input{Address: sender, Coin: input},
				Output{Address: recipient, Coin: sdk.Coin{}},
				deadline,
				false,
			),
		}, {
			"zero output coin",
			false,
			NewMsgSwapOrder(
				Input{Address: sender, Coin: input},
				Output{Address: recipient, Coin: sdk.NewCoin(denom1, sdk.ZeroInt())},
				deadline,
				true,
			),
		}, {
			"swap and coin denomination are equal", false, NewMsgSwapOrder(Input{
				Address: sender, Coin: input},
				Output{Address: recipient, Coin: sdk.NewCoin(denom0, amt)},
				deadline,
				true,
			),
		}, {
			"deadline not initialized",
			false,
			NewMsgSwapOrder(
				Input{Address: sender, Coin: input},
				Output{Address: recipient, Coin: output},
				emptyTime,
				true,
			),
		}, {
			"no sender",
			false,
			NewMsgSwapOrder(
				Input{Address: emptyAddr, Coin: input},
				Output{Address: recipient, Coin: output},
				deadline,
				true,
			),
		}, {
			"no recipient",
			true,
			NewMsgSwapOrder(
				Input{Address: sender, Coin: input},
				Output{Address: emptyAddr, Coin: output},
				deadline,
				true,
			),
		}, {
			"valid MsgSwapOrder",
			true,
			NewMsgSwapOrder(
				Input{Address: sender, Coin: input},
				Output{Address: recipient, Coin: output},
				deadline,
				true,
			),
		}, {
			"sender and recipient are same",
			true,
			NewMsgSwapOrder(
				Input{Address: sender, Coin: input},
				Output{Address: sender, Coin: output},
				deadline,
				true,
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectPass {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}

// ----------------------------------------------
// test MsgAddLiquidity
// ----------------------------------------------

func TestNewMsgAddLiquidity(t *testing.T) {
	msg := NewMsgAddLiquidity(input, amt, sdk.OneInt(), deadline, sender)
	require.Equal(t, input, msg.MaxToken)
	require.Equal(t, amt, msg.ExactStandardAmt)
	require.Equal(t, sdk.OneInt(), msg.MinLiquidity)
	require.Equal(t, deadline, msg.Deadline)
	require.Equal(t, sender, msg.Sender)
}

func TestMsgAddLiquidityRoute(t *testing.T) {
	msg := NewMsgAddLiquidity(input, amt, sdk.OneInt(), deadline, sender)
	require.Equal(t, RouterKey, msg.Route())
}

func TestMsgAddLiquidityType(t *testing.T) {
	msg := NewMsgAddLiquidity(input, amt, sdk.OneInt(), deadline, sender)
	require.Equal(t, MsgTypeAddLiquidity, msg.Type())
}

func TestMsgAddLiquidityGetSignBytes(t *testing.T) {
	msg := NewMsgAddLiquidity(input, amt, sdk.OneInt(), deadline, sender)
	res := msg.GetSignBytes()
	expected := `{"type":"irishub/coinswap/MsgAddLiquidity","value":{"deadline":"1580000000","exact_standard_amt":"100","max_token":{"amount":"1000","denom":"atom"},"min_liquidity":"1","sender":"faa128nh833v43sggcj65nk7khjka9dwngpl6j29hj"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgAddLiquidityGetSigners(t *testing.T) {
	msg := NewMsgAddLiquidity(input, amt, sdk.OneInt(), deadline, sender)
	res := msg.GetSigners()
	expected := "[51E773C62CAC6084625AA4EDEB5E56E95AE9A03F]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// test ValidateBasic for MsgAddLiquidity
func TestMsgAddLiquidityValidation(t *testing.T) {
	tests := []struct {
		name       string
		msg        MsgAddLiquidity
		expectPass bool
	}{
		{"no deposit coin", NewMsgAddLiquidity(sdk.Coin{}, amt, sdk.OneInt(), deadline, sender), false},
		{"zero deposit coin", NewMsgAddLiquidity(sdk.NewCoin(denom1, sdk.ZeroInt()), amt, sdk.OneInt(), deadline, sender), false},
		{"invalid withdraw amount", NewMsgAddLiquidity(input, sdk.ZeroInt(), sdk.OneInt(), deadline, sender), false},
		{"deadline not initialized", NewMsgAddLiquidity(input, amt, sdk.OneInt(), emptyTime, sender), false},
		{"empty sender", NewMsgAddLiquidity(input, amt, sdk.OneInt(), deadline, emptyAddr), false},
		{"valid MsgAddLiquidity", NewMsgAddLiquidity(input, amt, sdk.OneInt(), deadline, sender), true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectPass {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}

// ----------------------------------------------
// test MsgRemoveLiquidity
// ----------------------------------------------

func TestNewMsgRemoveLiquidity(t *testing.T) {
	msg := NewMsgRemoveLiquidity(amt, withdrawLiquidity, sdk.OneInt(), deadline, sender)
	require.Equal(t, amt, msg.MinToken)
	require.Equal(t, withdrawLiquidity, msg.WithdrawLiquidity)
	require.Equal(t, sdk.OneInt(), msg.MinStandardAmt)
	require.Equal(t, deadline, msg.Deadline)
	require.Equal(t, sender, msg.Sender)
}

func TestMsgRemoveLiquidityRoute(t *testing.T) {
	msg := NewMsgRemoveLiquidity(amt, withdrawLiquidity, sdk.OneInt(), deadline, sender)
	require.Equal(t, RouterKey, msg.Route())
}

func TestMsgRemoveLiquidityType(t *testing.T) {
	msg := NewMsgRemoveLiquidity(amt, withdrawLiquidity, sdk.OneInt(), deadline, sender)
	require.Equal(t, MsgTypeRemoveLiquidity, msg.Type())
}

func TestMsgRemoveLiquidityGetSignBytes(t *testing.T) {
	msg := NewMsgRemoveLiquidity(amt, withdrawLiquidity, sdk.OneInt(), deadline, sender)
	res := msg.GetSignBytes()
	expected := `{"type":"irishub/coinswap/MsgRemoveLiquidity","value":{"deadline":"1580000000","min_standard_amt":"1","min_token":"100","sender":"faa128nh833v43sggcj65nk7khjka9dwngpl6j29hj","withdraw_liquidity":{"amount":"500","denom":"uni:btc"}}}`
	require.Equal(t, expected, string(res))
}

func TestMsgRemoveLiquidityGetSigners(t *testing.T) {
	msg := NewMsgRemoveLiquidity(amt, withdrawLiquidity, sdk.OneInt(), deadline, sender)
	res := msg.GetSigners()
	expected := "[51E773C62CAC6084625AA4EDEB5E56E95AE9A03F]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// test ValidateBasic for MsgRemoveLiquidity
func TestMsgRemoveLiquidityValidation(t *testing.T) {
	tests := []struct {
		name       string
		msg        MsgRemoveLiquidity
		expectPass bool
	}{
		{"no withdraw coin", NewMsgRemoveLiquidity(amt, sdk.Coin{}, sdk.OneInt(), deadline, sender), false},
		{"zero withdraw coin", NewMsgRemoveLiquidity(amt, sdk.NewCoin(unidenom, sdk.ZeroInt()), sdk.OneInt(), deadline, sender), false},
		{"invalid minimum token amount", NewMsgRemoveLiquidity(sdk.NewInt(-100), withdrawLiquidity, sdk.OneInt(), deadline, sender), false},
		{"invalid minimum standard token amount", NewMsgRemoveLiquidity(amt, withdrawLiquidity, sdk.NewInt(-100), deadline, sender), false},
		{"deadline not initialized", NewMsgRemoveLiquidity(amt, withdrawLiquidity, sdk.OneInt(), emptyTime, sender), false},
		{"empty sender", NewMsgRemoveLiquidity(amt, withdrawLiquidity, sdk.OneInt(), deadline, emptyAddr), false},
		{"valid MsgRemoveLiquidity", NewMsgRemoveLiquidity(amt, withdrawLiquidity, sdk.OneInt(), deadline, sender), true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectPass {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
		})
	}

}
