package profiling

import (
	sdk "github.com/irisnet/irishub/types"
)

const MsgType = "profiling"

//______________________________________________________________________
// MsgAddProfiler - struct for add a profiler
type MsgAddProfiler struct {
	Profiler
}

func NewMsgAddProfiler(addr, addedAddr sdk.AccAddress, name string) MsgAddProfiler {
	return MsgAddProfiler{
		Profiler{
			Name:      name,
			Addr:      addr,
			AddedAddr: addedAddr,
		},
	}
}
func (msg MsgAddProfiler) Route() string { return MsgType }
func (msg MsgAddProfiler) Type() string  { return "service add-profiler" }
func (msg MsgAddProfiler) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}
func (msg MsgAddProfiler) ValidateBasic() sdk.Error {
	if len(msg.Addr) == 0 {
		return sdk.ErrInvalidAddress(msg.Addr.String())
	}
	if len(msg.AddedAddr) == 0 {
		return sdk.ErrInvalidAddress(msg.AddedAddr.String())
	}
	return nil
}
func (msg MsgAddProfiler) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.AddedAddr}
}
