package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgAddProfiler    = "add_profiler"    // type for MsgAddProfiler
	TypeMsgDeleteProfiler = "delete_profiler" // type for MsgDeleteProfiler
	TypeMsgAddTrustee     = "add_trustee"     // type for MsgAddTrustee
	TypeMsgDeleteTrustee  = "delete_trustee"  // type for MsgDeleteTrustee
)

var (
	_ sdk.Msg = MsgAddProfiler{}
	_ sdk.Msg = MsgAddTrustee{}
	_ sdk.Msg = MsgDeleteProfiler{}
	_ sdk.Msg = MsgDeleteTrustee{}
)

// MsgAddProfiler - struct for add a profiler
type MsgAddProfiler struct {
	AddGuardian
}

// NewMsgAddProfiler constructs a MsgAddProfiler
func NewMsgAddProfiler(description string, address, addedBy sdk.AccAddress) MsgAddProfiler {
	return MsgAddProfiler{
		AddGuardian: AddGuardian{
			Description: description,
			Address:     address,
			AddedBy:     addedBy,
		},
	}
}

// Route implements Msg.
func (msg MsgAddProfiler) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgAddProfiler) Type() string { return TypeMsgAddProfiler }

// GetSignBytes implements Msg.
func (msg MsgAddProfiler) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgAddProfiler) ValidateBasic() sdk.Error {
	return msg.AddGuardian.ValidateBasic()
}

// GetSigners implements Msg.
func (msg MsgAddProfiler) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.AddedBy}
}

//______________________________________________________________________
// MsgDeleteProfiler - struct for delete a profiler
type MsgDeleteProfiler struct {
	DeleteGuardian
}

// NewMsgDeleteProfiler constructs a MsgDeleteProfiler
func NewMsgDeleteProfiler(address, deletedBy sdk.AccAddress) MsgDeleteProfiler {
	return MsgDeleteProfiler{
		DeleteGuardian: DeleteGuardian{
			Address:   address,
			DeletedBy: deletedBy,
		},
	}
}

// Route implements Msg.
func (msg MsgDeleteProfiler) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgDeleteProfiler) Type() string { return TypeMsgDeleteProfiler }

// GetSignBytes implements Msg.
func (msg MsgDeleteProfiler) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// RoValidateBasicute implements Msg.
func (msg MsgDeleteProfiler) ValidateBasic() sdk.Error {
	return msg.DeleteGuardian.ValidateBasic()
}

// GetSigners implements Msg.
func (msg MsgDeleteProfiler) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.DeletedBy}
}

//______________________________________________________________________
// MsgAddTrustee - struct for add a trustee
type MsgAddTrustee struct {
	AddGuardian
}

// NewMsgAddTrustee constructs a MsgAddTrustee
func NewMsgAddTrustee(description string, address, addedAddress sdk.AccAddress) MsgAddTrustee {
	return MsgAddTrustee{
		AddGuardian: AddGuardian{
			Description: description,
			Address:     address,
			AddedBy:     addedAddress,
		},
	}
}

// Route implements Msg.
func (msg MsgAddTrustee) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgAddTrustee) Type() string { return TypeMsgAddTrustee }

// GetSignBytes implements Msg.
func (msg MsgAddTrustee) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgAddTrustee) ValidateBasic() sdk.Error {
	return msg.AddGuardian.ValidateBasic()
}

// GetSigners implements Msg.
func (msg MsgAddTrustee) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.AddedBy}
}

//______________________________________________________________________
// MsgDeleteTrustee - struct for delete a trustee
type MsgDeleteTrustee struct {
	DeleteGuardian
}

// NewMsgDeleteTrustee constructs a MsgDeleteTrustee
func NewMsgDeleteTrustee(address, deletedBy sdk.AccAddress) MsgDeleteTrustee {
	return MsgDeleteTrustee{
		DeleteGuardian: DeleteGuardian{
			Address:   address,
			DeletedBy: deletedBy,
		},
	}
}

// Route implements Msg.
func (msg MsgDeleteTrustee) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgDeleteTrustee) Type() string { return TypeMsgDeleteTrustee }

// GetSignBytes implements Msg.
func (msg MsgDeleteTrustee) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgDeleteTrustee) ValidateBasic() sdk.Error {
	return msg.DeleteGuardian.ValidateBasic()
}

// GetSigners implements Msg.
func (msg MsgDeleteTrustee) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.DeletedBy}
}

//______________________________________________________________________

// AddGuardian
type AddGuardian struct {
	Description string         `json:"description" yaml:"description"` //
	Address     sdk.AccAddress `json:"address" yaml:"address"`         // address added
	AddedBy     sdk.AccAddress `json:"added_by" yaml:"added_by"`       // address that initiated the tx
}

// DeleteGuardian
type DeleteGuardian struct {
	Address   sdk.AccAddress `json:"address" yaml:"address"`       // address deleted
	DeletedBy sdk.AccAddress `json:"deleted_by" yaml:"deleted_by"` // address that initiated the tx
}

// ValidateBasic
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
	if err := g.EnsureLength(); err != nil {
		return err
	}
	return nil
}

// ValidateBasic
func (g DeleteGuardian) ValidateBasic() sdk.Error {
	if len(g.Address) == 0 {
		return sdk.ErrInvalidAddress(g.Address.String())
	}
	if len(g.DeletedBy) == 0 {
		return sdk.ErrInvalidAddress(g.DeletedBy.String())
	}
	return nil
}

// EnsureLength
func (g AddGuardian) EnsureLength() sdk.Error {
	if len(g.Description) > 70 {
		return sdk.NewError(DefaultCodespace, CodeInvalidGuardian, "description", len(g.Description), 70)
	}
	return nil
}
