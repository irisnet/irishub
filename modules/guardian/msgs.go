package guardian

import (
	sdk "github.com/irisnet/irishub/types"
	"regexp"
)

const MsgType = "guardian"

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
func (msg MsgAddProfiler) Type() string  { return "profiling add-profiler" }
func (msg MsgAddProfiler) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgAddProfiler) ValidateBasic() sdk.Error {
	if !validName(msg.Name) {
		return ErrInvalidProfilerName(DefaultCodespace, msg.Name)
	}
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

//______________________________________________________________________

func validName(name string) bool {
	if len(name) == 0 || len(name) > 128 {
		return false
	}

	// Must contain alphanumeric characters, _ and - only
	reg := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	return !reg.Match([]byte(name))
}
