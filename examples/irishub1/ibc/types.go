package ibc

import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/codec"
)

var (
	msgCdc *codec.Codec
)

func init() {
	msgCdc = codec.New()
}

type IBCSetMsg struct {
	Addr		sdk.AccAddress
}

var _ sdk.Msg = (*IBCSetMsg)(nil)

func NewIBCSetMsg(addr sdk.AccAddress) IBCSetMsg {
	return  IBCSetMsg{
		Addr:addr,
	}
}

func (msg  IBCSetMsg) Type() string {
	return "ibc-set"
}

func (msg  IBCSetMsg) Route() string {
	return "ibc-1"
}

func (msg  IBCSetMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg  IBCSetMsg) ValidateBasic() sdk.Error {
	return nil
}

func (msg  IBCSetMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Addr}
}

type IBCGetMsg struct {
	Addr		sdk.AccAddress
}

var _ sdk.Msg = (*IBCGetMsg)(nil)

func NewIBCGetMsg(addr sdk.AccAddress) IBCGetMsg {
	return  IBCGetMsg{
		Addr:addr,
	}
}

func (msg  IBCGetMsg) Type() string {
	return "ibc-get"
}

func (msg  IBCGetMsg) Route() string {
	return "ibc-1"
}

func (msg  IBCGetMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg  IBCGetMsg) ValidateBasic() sdk.Error {
	return nil
}

func (msg  IBCGetMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Addr}
}
