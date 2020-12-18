package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

var _ sdk.Msg = &MsgSend{}

// NewMsgSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgSend(in []banktypes.Input, out []banktypes.Output) MsgSend {
	return MsgSend{Inputs: in, Outputs: out}
}

// Implements Msg.
// nolint
func (msg MsgSend) Route() string { return ModuleName }
func (msg MsgSend) Type() string  { return "send" }

// Implements Msg.
func (msg MsgSend) ValidateBasic() error {
	// this just makes sure all the inputs and outputs are properly formatted,
	// not that they actually have the money inside
	if len(msg.Inputs) == 0 {
		return banktypes.ErrNoInputs
	}

	if len(msg.Outputs) == 0 {
		return banktypes.ErrNoOutputs
	}

	return banktypes.ValidateInputsOutputs(msg.Inputs, msg.Outputs)
}

// Implements Msg.
func (msg MsgSend) GetSignBytes() []byte {
	var inputs, outputs []json.RawMessage
	for _, input := range msg.Inputs {
		inputs = append(inputs, ModuleCdc.LegacyAmino.MustMarshalJSON(input))
	}
	for _, output := range msg.Outputs {
		outputs = append(outputs, ModuleCdc.LegacyAmino.MustMarshalJSON(output))
	}
	b, err := ModuleCdc.LegacyAmino.MarshalJSON(struct {
		Inputs  []json.RawMessage `json:"inputs"`
		Outputs []json.RawMessage `json:"outputs"`
	}{
		Inputs:  inputs,
		Outputs: outputs,
	})
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgSend) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.Inputs))
	for i, in := range msg.Inputs {
		addr, _ := sdk.AccAddressFromBech32(in.Address)
		addrs[i] = addr
	}
	return addrs
}
