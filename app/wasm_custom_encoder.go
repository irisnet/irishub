package app

import (
	"encoding/base64"
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
	Router string `json:"router"`
	Data   string `json:"data"`
}

// Encode implement wasm module CustomEncoder
func (encoder WasmCustomEncoder) Encode(sender sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error) {
	var msg MsgWasmCustom
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "data type not match []MsgWasmCustom, json.Unmarshal failed")
	}

	sysMsg, err := encoder.registry.Resolve(msg.Router)
	if err != nil {
		return nil, sdkerrors.Wrapf(wasm.ErrInvalidMsg, err.Error())
	}

	msgByte, err := base64.StdEncoding.DecodeString(msg.Data)
	if err != nil {
		return nil, sdkerrors.Wrapf(wasm.ErrInvalidMsg, err.Error())
	}

	if err := json.Unmarshal(msgByte, &sysMsg); err != nil {
		return nil, sdkerrors.Wrapf(wasm.ErrInvalidMsg, "data not match %v, json.Unmarshal failed", sysMsg)
	}

	m, ok := sysMsg.(sdk.Msg)
	if !ok {
		return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "not implement sdk.Msg")
	}
	return []sdk.Msg{m}, nil
}

// NewMessageEncoders overide the wasm module CustomEncoder
func NewMessageEncoders(registry types.InterfaceRegistry) *wasm.MessageEncoders {
	return &wasm.MessageEncoders{
		Custom: WasmCustomEncoder{registry: registry}.Encode,
	}
}
