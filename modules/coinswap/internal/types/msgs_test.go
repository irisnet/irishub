package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// nolint: deadcode unused
var (
	amt = sdk.NewInt(100)

	senderPk    = ed25519.GenPrivKey().PubKey()
	recipientPk = ed25519.GenPrivKey().PubKey()
	sender      = sdk.AccAddress(senderPk.Address())
	recipient   = sdk.AccAddress(recipientPk.Address())

	denom0   = "atom-min"
	denom1   = "btc-min"
	unidenom = FormatUniABSPrefix + "btc-min"

	input             sdk.Coin
	output            sdk.Coin
	withdrawLiquidity sdk.Coin
	deadline          = time.Now().Unix()

	emptyAddr sdk.AccAddress
	emptyTime int64
)

func init() {
	_ = sdk.RegisterDenom(denom0, sdk.NewDecWithPrec(1, sdk.Precision))
	_ = sdk.RegisterDenom(denom1, sdk.NewDecWithPrec(1, sdk.Precision))
	_ = sdk.RegisterDenom(unidenom, sdk.NewDecWithPrec(1, sdk.Precision))

	input = sdk.NewCoin(denom0, sdk.NewInt(1000))
	output = sdk.NewCoin(denom1, sdk.NewInt(500))
	withdrawLiquidity = sdk.NewCoin(unidenom, sdk.NewInt(500))
}

// test ValidateBasic for MsgSwapOrder
func TestMsgSwapOrder(t *testing.T) {
	tests := []struct {
		name       string
		msg        MsgSwapOrder
		expectPass bool
	}{
		{"no input coin", NewMsgSwapOrder(Input{Address: sender}, Output{recipient, output}, deadline, true), false},
		{"zero input coin", NewMsgSwapOrder(Input{
			Address: sender,
			Coin:    sdk.NewCoin(denom0, sdk.ZeroInt()),
		}, Output{
			Address: recipient,
			Coin:    output,
		}, deadline, true), false},
		{"no output coin", NewMsgSwapOrder(Input{
			Address: sender,
			Coin:    input,
		}, Output{
			Address: recipient,
			Coin:    sdk.Coin{},
		}, deadline, false), false},
		{"zero output coin", NewMsgSwapOrder(Input{
			Address: sender,
			Coin:    input,
		}, Output{
			Address: recipient,
			Coin:    sdk.NewCoin(denom1, sdk.ZeroInt()),
		}, deadline, true), false},
		{"swap and coin denomination are equal", NewMsgSwapOrder(Input{
			Address: sender,
			Coin:    input,
		}, Output{
			Address: recipient,
			Coin:    sdk.NewCoin(denom0, amt),
		}, deadline, true), false},
		{"deadline not initialized", NewMsgSwapOrder(Input{
			Address: sender,
			Coin:    input,
		}, Output{
			Address: recipient,
			Coin:    output,
		}, emptyTime, true), false},
		{"no sender", NewMsgSwapOrder(Input{
			Address: emptyAddr,
			Coin:    input,
		}, Output{
			Address: recipient,
			Coin:    output,
		}, deadline, true), false},
		{"no recipient", NewMsgSwapOrder(Input{
			Address: sender,
			Coin:    input,
		}, Output{
			Address: emptyAddr,
			Coin:    output,
		}, deadline, true), true},
		{"valid MsgSwapOrder", NewMsgSwapOrder(Input{
			Address: sender,
			Coin:    input,
		}, Output{
			Address: recipient,
			Coin:    output,
		}, deadline, true), true},
		{"sender and recipient are same", NewMsgSwapOrder(Input{
			Address: sender,
			Coin:    input,
		}, Output{
			Address: sender,
			Coin:    output,
		}, deadline, true), true},
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

// test ValidateBasic for MsgAddLiquidity
func TestMsgAddLiquidity(t *testing.T) {
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

// test ValidateBasic for MsgRemoveLiquidity
func TestMsgRemoveLiquidity(t *testing.T) {
	tests := []struct {
		name       string
		msg        MsgRemoveLiquidity
		expectPass bool
	}{
		{"no withdraw coin", NewMsgRemoveLiquidity(amt, sdk.Coin{}, sdk.OneInt(), deadline, sender), false},
		{"zero withdraw coin", NewMsgRemoveLiquidity(amt, sdk.NewCoin(unidenom, sdk.ZeroInt()), sdk.OneInt(), deadline, sender), false},
		{"invalid minimum token amount", NewMsgRemoveLiquidity(sdk.NewInt(-100), withdrawLiquidity, sdk.OneInt(), deadline, sender), false},
		{"invalid minimum iris amount", NewMsgRemoveLiquidity(amt, withdrawLiquidity, sdk.NewInt(-100), deadline, sender), false},
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
