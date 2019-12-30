package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	AddGuardian `json:"add_guardian" yaml:"add_guardian"` // the guardian will be added
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
func (msg MsgAddProfiler) ValidateBasic() error {
	return msg.AddGuardian.ValidateBasic()
}

// GetSigners implements Msg.
func (msg MsgAddProfiler) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.AddedBy}
}

//______________________________________________________________________
// MsgDeleteProfiler - struct for delete a profiler
type MsgDeleteProfiler struct {
	DeleteGuardian `json:"delete_guardian" yaml:"delete_guardian"` // the guardian will be deleted
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
func (msg MsgDeleteProfiler) ValidateBasic() error {
	return msg.DeleteGuardian.ValidateBasic()
}

// GetSigners implements Msg.
func (msg MsgDeleteProfiler) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.DeletedBy}
}

//______________________________________________________________________
// MsgAddTrustee - struct for add a trustee
type MsgAddTrustee struct {
	AddGuardian `json:"add_guardian" yaml:"add_guardian"` // the guardian will be added
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
func (msg MsgAddTrustee) ValidateBasic() error {
	return msg.AddGuardian.ValidateBasic()
}

// GetSigners implements Msg.
func (msg MsgAddTrustee) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.AddedBy}
}

//______________________________________________________________________
// MsgDeleteTrustee - struct for delete a trustee
type MsgDeleteTrustee struct {
	DeleteGuardian `json:"delete_guardian" yaml:"delete_guardian"` // the guardian will be deleted
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
func (msg MsgDeleteTrustee) ValidateBasic() error {
	return msg.DeleteGuardian.ValidateBasic()
}

// GetSigners implements Msg.
func (msg MsgDeleteTrustee) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.DeletedBy}
}

//______________________________________________________________________

// AddGuardian defines the properties of add guardian message
type AddGuardian struct {
	Description string         `json:"description" yaml:"description"` // description of the guardian
	Address     sdk.AccAddress `json:"address" yaml:"address"`         // address added
	AddedBy     sdk.AccAddress `json:"added_by" yaml:"added_by"`       // address that initiated the tx
}

// DeleteGuardian defines the properties of delete guardian message
type DeleteGuardian struct {
	Address   sdk.AccAddress `json:"address" yaml:"address"`       // address deleted
	DeletedBy sdk.AccAddress `json:"deleted_by" yaml:"deleted_by"` // address that initiated the tx
}

// ValidateBasic validate the AddGuardian
func (g AddGuardian) ValidateBasic() error {
	if len(g.Description) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "description missing")
	}
	if len(g.Address) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "added address missing")
	}
	if len(g.AddedBy) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "operator address missing")
	}
	if err := g.EnsureLength(); err != nil {
		return err
	}
	return nil
}

// ValidateBasic validate the DeleteGuardian
func (g DeleteGuardian) ValidateBasic() error {
	if len(g.Address) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "deleted address missing")
	}
	if len(g.DeletedBy) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "operator address missing")
	}
	return nil
}

// EnsureLength validate the length of AddGuardian
func (g AddGuardian) EnsureLength() error {
	if len(g.Description) > 70 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid website length; got: %d, max: %d", len(g.Description), 70)
	}
	return nil
}
