package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/config"
)

// nolint: deadcode unused
var (
	sender, _    = sdk.AccAddressFromHex(crypto.AddressHash([]byte("sender")).String())
	recipient, _ = sdk.AccAddressFromHex(crypto.AddressHash([]byte("recipient")).String())

	amt = sdk.NewInt(100)

	denomETH = "eth"
	denomBTC = "btc"
	unidenom = FormatUniABSPrefix + denomBTC

	input             = sdk.NewCoin(denomETH, sdk.NewInt(1000))
	output            = sdk.NewCoin(denomBTC, sdk.NewInt(500))
	withdrawLiquidity = sdk.NewCoin(unidenom, sdk.NewInt(500))
	deadline          = int64(1580000000)

	emptyAddr sdk.AccAddress
	emptyTime int64
)

func init() {
	sdk.GetConfig().SetBech32PrefixForAccount(config.GetConfig().GetBech32AccountAddrPrefix(), config.GetConfig().GetBech32AccountPubPrefix())
	sdk.GetConfig().SetBech32PrefixForValidator(config.GetConfig().GetBech32ValidatorAddrPrefix(), config.GetConfig().GetBech32ValidatorPubPrefix())
	sdk.GetConfig().SetBech32PrefixForConsensusNode(config.GetConfig().GetBech32ConsensusAddrPrefix(), config.GetConfig().GetBech32ConsensusPubPrefix())
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
	require.Equal(t, TypeMsgSwapOrder, msg.Type())
}

func TestMsgSwapOrderGetSignBytes(t *testing.T) {
	msg := NewMsgSwapOrder(
		Input{Address: sender, Coin: input},
		Output{Address: sender, Coin: output},
		deadline,
		true,
	)
	res := msg.GetSignBytes()
	expected := `{"type":"irishub/coinswap/MsgSwapOrder","value":{"deadline":"1580000000","input":{"address":"faa1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgkwnkl5","coin":{"amount":"1000","denom":"eth"}},"is_buy_order":true,"output":{"address":"faa1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgkwnkl5","coin":{"amount":"500","denom":"btc"}}}}`
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
	expected := "[0A367B92CF0B037DFD89960EE832D56F7FC15168]"
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
				Input{Address: sender, Coin: sdk.NewCoin(denomETH, sdk.ZeroInt())},
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
				Output{Address: recipient, Coin: sdk.NewCoin(denomBTC, sdk.ZeroInt())},
				deadline,
				true,
			),
		}, {
			"swap and coin denomination are equal", false, NewMsgSwapOrder(Input{
				Address: sender, Coin: input},
				Output{Address: recipient, Coin: sdk.NewCoin(denomETH, amt)},
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
				require.NoError(t, err, tc.name)
			} else {
				require.Error(t, err, tc.name)
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
	require.Equal(t, TypeMsgAddLiquidity, msg.Type())
}

func TestMsgAddLiquidityGetSignBytes(t *testing.T) {
	msg := NewMsgAddLiquidity(input, amt, sdk.OneInt(), deadline, sender)
	res := msg.GetSignBytes()
	expected := `{"type":"irishub/coinswap/MsgAddLiquidity","value":{"deadline":"1580000000","exact_standard_amt":"100","max_token":{"amount":"1000","denom":"eth"},"min_liquidity":"1","sender":"faa1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgkwnkl5"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgAddLiquidityGetSigners(t *testing.T) {
	msg := NewMsgAddLiquidity(input, amt, sdk.OneInt(), deadline, sender)
	res := msg.GetSigners()
	expected := "[0A367B92CF0B037DFD89960EE832D56F7FC15168]"
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
		{"zero deposit coin", NewMsgAddLiquidity(sdk.NewCoin(denomBTC, sdk.ZeroInt()), amt, sdk.OneInt(), deadline, sender), false},
		{"invalid withdraw amount", NewMsgAddLiquidity(input, sdk.ZeroInt(), sdk.OneInt(), deadline, sender), false},
		{"deadline not initialized", NewMsgAddLiquidity(input, amt, sdk.OneInt(), emptyTime, sender), false},
		{"empty sender", NewMsgAddLiquidity(input, amt, sdk.OneInt(), deadline, emptyAddr), false},
		{"valid MsgAddLiquidity", NewMsgAddLiquidity(input, amt, sdk.OneInt(), deadline, sender), true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectPass {
				require.NoError(t, err, tc.name)
			} else {
				require.Error(t, err, tc.name)
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
	require.Equal(t, TypeMsgRemoveLiquidity, msg.Type())
}

func TestMsgRemoveLiquidityGetSignBytes(t *testing.T) {
	msg := NewMsgRemoveLiquidity(amt, withdrawLiquidity, sdk.OneInt(), deadline, sender)
	res := msg.GetSignBytes()
	expected := `{"type":"irishub/coinswap/MsgRemoveLiquidity","value":{"deadline":"1580000000","min_standard_amt":"1","min_token":"100","sender":"faa1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgkwnkl5","withdraw_liquidity":{"amount":"500","denom":"uni:btc"}}}`
	require.Equal(t, expected, string(res))
}

func TestMsgRemoveLiquidityGetSigners(t *testing.T) {
	msg := NewMsgRemoveLiquidity(amt, withdrawLiquidity, sdk.OneInt(), deadline, sender)
	res := msg.GetSigners()
	expected := "[0A367B92CF0B037DFD89960EE832D56F7FC15168]"
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
				require.NoError(t, err, tc.name)
			} else {
				require.Error(t, err, tc.name)
			}
		})
	}

}
