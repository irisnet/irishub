package upgrade

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgSwitch struct {
	Title		string
	ProposalID	int64
	Voter		sdk.AccAddress
}

func NewMsgSwitch( title string, proposalID int64,voter sdk.AccAddress) MsgSwitch {
	return MsgSwitch{
		Title:title,
		ProposalID: proposalID,
		Voter:      voter,
	}
}

func (msg MsgSwitch) Type() string {
	return "record"
}

func (msg MsgSwitch) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg MsgSwitch) ValidateBasic() sdk.Error {
	if len(msg.Title) == 0 {
		return ErrInvalidTitle(DefaultCodespace, msg.Title) // TODO: Proper Error
	}
	if len(msg.Description) == 0 {
		return ErrInvalidDescription(DefaultCodespace, msg.Description) // TODO: Proper Error
	}
	if !validProposalType(msg.ProposalType) {
		return ErrInvalidProposalType(DefaultCodespace, msg.ProposalType)
	}
	if len(msg.Proposer) == 0 {
		return sdk.ErrInvalidAddress(msg.Proposer.String())
	}
	if !msg.InitialDeposit.IsValid() {
		return sdk.ErrInvalidCoins(msg.InitialDeposit.String())
	}
	if !msg.InitialDeposit.IsNotNegative() {
		return sdk.ErrInvalidCoins(msg.InitialDeposit.String())
	}
	return nil
	return nil
}

func (msg MsgSwitch) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Voter}
}