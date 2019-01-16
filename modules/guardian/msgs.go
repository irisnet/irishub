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

func NewMsgAddProfiler(description string, address, addedBy sdk.AccAddress) MsgAddProfiler {
	return MsgAddProfiler{
		AddGuardian: AddGuardian{
			Description: description,
			Address:     address,
			AddedBy:     addedBy,
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
	return []sdk.AccAddress{msg.AddedBy}
}

//______________________________________________________________________
// MsgDeleteProfiler - struct for delete a profiler
type MsgDeleteProfiler struct {
	DeleteGuardian
}

func NewMsgDeleteProfiler(address, deletedBy sdk.AccAddress) MsgDeleteProfiler {
	return MsgDeleteProfiler{
		DeleteGuardian: DeleteGuardian{
			Address:   address,
			DeletedBy: deletedBy,
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
	return []sdk.AccAddress{msg.DeletedBy}
}

//______________________________________________________________________
// MsgAddTrustee - struct for add a trustee
type MsgAddTrustee struct {
	AddGuardian
}

func NewMsgAddTrustee(description string, address, addedAddress sdk.AccAddress) MsgAddTrustee {
	return MsgAddTrustee{
		AddGuardian: AddGuardian{
			Description: description,
			Address:     address,
			AddedBy:     addedAddress,
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
	return []sdk.AccAddress{msg.AddedBy}
}

//______________________________________________________________________
// MsgDeleteTrustee - struct for delete a trustee
type MsgDeleteTrustee struct {
	DeleteGuardian
}

func NewMsgDeleteTrustee(address, deletedBy sdk.AccAddress) MsgDeleteTrustee {
	return MsgDeleteTrustee{
		DeleteGuardian: DeleteGuardian{
			Address:   address,
			DeletedBy: deletedBy,
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
	return []sdk.AccAddress{msg.DeletedBy}
}

//______________________________________________________________________

type AddGuardian struct {
	Description string         `json:"description"`
	Address     sdk.AccAddress `json:"address"`  // address added
	AddedBy     sdk.AccAddress `json:"added_by"` // address that initiated the tx
}

type DeleteGuardian struct {
	Address   sdk.AccAddress `json:"address"`    // address deleted
	DeletedBy sdk.AccAddress `json:"deleted_by"` // address that initiated the tx
}

func (g AddGuardian) ValidateBasic() sdk.Error {
	if len(g.Description) == 0 {
		return ErrInvalidDescription(DefaultCodespace)
	}
	if len(g.Address) == 0 {
		return sdk.ErrInvalidAddress(g.Address.String())
	}
	if len(g.AddedBy) == 0 {
		return sdk.ErrInvalidAddress(g.AddedBy.String())
	}
	return nil
}

func (g DeleteGuardian) ValidateBasic() sdk.Error {
	if len(g.Address) == 0 {
		return sdk.ErrInvalidAddress(g.Address.String())
	}
	if len(g.DeletedBy) == 0 {
		return sdk.ErrInvalidAddress(g.DeletedBy.String())
	}
	return nil
}
