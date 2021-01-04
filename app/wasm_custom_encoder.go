package app

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/CosmWasm/wasmd/x/wasm"
)

// WasmCustomEncoder implement wasm module CustomEncoder
type WasmCustomEncoder struct {
	registry types.InterfaceRegistry
}

// MsgWasmCustom implement wasm module CustomEncoder
type MsgWasmCustom struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

// Encode implement wasm module CustomEncoder
func (encoder WasmCustomEncoder) Encode(sender sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error) {
	var msgs []MsgWasmCustom
	if err := json.Unmarshal(data, &msgs); err != nil {
		return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "data type not match []MsgWasmCustom, json.Unmarshal failed")
	}

	var moduleMsgs = make([]sdk.Msg, len(msgs))
	for i, msg := range msgs {
		sysMsg, err := encoder.registry.Resolve(msg.Type)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(msg.Value, &sysMsg); err != nil {
			return nil, sdkerrors.Wrapf(wasm.ErrInvalidMsg, "data not match %v, json.Unmarshal failed", sysMsg)
		}

		if m, ok := sysMsg.(sdk.Msg); ok {
			moduleMsgs[i] = m
		} else {
			return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "not implement sdk.Msg")
		}
	}
	return moduleMsgs, nil
}

// NewMessageEncoders overide the wasm module CustomEncoder
func NewMessageEncoders(registry types.InterfaceRegistry) *wasm.MessageEncoders {
	return &wasm.MessageEncoders{
		Custom: WasmCustomEncoder{registry: registry}.Encode,
	}
}
