package guardian

import (
	sdk "github.com/irisnet/irishub/types"
)

const MsgType = "guardian"

//______________________________________________________________________
// MsgAddProfiler - struct for add a profiler
type MsgAddProfiler struct {
	AddGuardian
}

func NewMsgAddProfiler(description string, address, addedAddress sdk.AccAddress) MsgAddProfiler {
	return MsgAddProfiler{
		AddGuardian: AddGuardian{
			Description:  description,
			Address:      address,
			AddedAddress: addedAddress,
		},
	}
}
func (msg MsgAddProfiler) Route() string { return MsgType }
func (msg MsgAddProfiler) Type() string  { return "guardian add-profiler" }
func (msg MsgAddProfiler) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgAddProfiler) ValidateBasic() sdk.Error {
	return msg.AddGuardian.ValidateBasic()
}

func (msg MsgAddProfiler) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.AddedAddress}
}

//______________________________________________________________________
// MsgDeleteProfiler - struct for delete a profiler
type MsgDeleteProfiler struct {
	DeleteGuardian
}

func NewMsgDeleteProfiler(address, deletedAddress sdk.AccAddress) MsgDeleteProfiler {
	return MsgDeleteProfiler{
		DeleteGuardian: DeleteGuardian{
			Address:        address,
			DeletedAddress: deletedAddress,
		},
	}
}
func (msg MsgDeleteProfiler) Route() string { return MsgType }
func (msg MsgDeleteProfiler) Type() string  { return "guardian delete-profiler" }
func (msg MsgDeleteProfiler) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgDeleteProfiler) ValidateBasic() sdk.Error {
	return msg.DeleteGuardian.ValidateBasic()
}

func (msg MsgDeleteProfiler) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.DeletedAddress}
}

//______________________________________________________________________
// MsgAddTrustee - struct for add a trustee
type MsgAddTrustee struct {
	AddGuardian
}

func NewMsgAddTrustee(description string, accountType AccountType, address, addedAddress sdk.AccAddress) MsgAddTrustee {
	return MsgAddTrustee{
		AddGuardian: AddGuardian{
			Description:  description,
			Address:      address,
			AddedAddress: addedAddress,
		},
	}
}
func (msg MsgAddTrustee) Route() string { return MsgType }
func (msg MsgAddTrustee) Type() string  { return "guardian add-trustee" }
func (msg MsgAddTrustee) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgAddTrustee) ValidateBasic() sdk.Error {
	return msg.AddGuardian.ValidateBasic()
}

func (msg MsgAddTrustee) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.AddedAddress}
}

//______________________________________________________________________
// MsgDeleteTrustee - struct for delete a trustee
type MsgDeleteTrustee struct {
	DeleteGuardian
}

func NewMsgDeleteTrustee(address, deletedAddress sdk.AccAddress) MsgDeleteTrustee {
	return MsgDeleteTrustee{
		DeleteGuardian: DeleteGuardian{
			Address:        address,
			DeletedAddress: deletedAddress,
		},
	}
}
func (msg MsgDeleteTrustee) Route() string { return MsgType }
func (msg MsgDeleteTrustee) Type() string  { return "guardian delete-trustee" }
func (msg MsgDeleteTrustee) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgDeleteTrustee) ValidateBasic() sdk.Error {
	return msg.DeleteGuardian.ValidateBasic()
}

func (msg MsgDeleteTrustee) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.DeletedAddress}
}

//______________________________________________________________________

type AddGuardian struct {
	Description  string         `json:"description"`
	Address      sdk.AccAddress `json:"address"`
	AddedAddress sdk.AccAddress `json:"added_address"`
}

type DeleteGuardian struct {
	Address        sdk.AccAddress `json:"address"`
	DeletedAddress sdk.AccAddress `json:"deleted_address"`
}

func (g AddGuardian) ValidateBasic() sdk.Error {
	if len(g.Description) == 0 {
		return ErrInvalidDescription(DefaultCodespace)
	}
	if len(g.Address) == 0 {
		return sdk.ErrInvalidAddress(g.AddedAddress.String())
	}
	if len(g.AddedAddress) == 0 {
		return sdk.ErrInvalidAddress(g.AddedAddress.String())
	}
	return nil
}

func (g DeleteGuardian) ValidateBasic() sdk.Error {
	if len(g.Address) == 0 {
		return sdk.ErrInvalidAddress(g.Address.String())
	}
	if len(g.DeletedAddress) == 0 {
		return sdk.ErrInvalidAddress(g.DeletedAddress.String())
	}
	return nil
}
