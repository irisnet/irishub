package upgrade

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// name to idetify transaction types
const MsgRoute = "upgrade"

var _ sdk.Msg = MsgSwitch{}

type MsgSwitch struct {
	Title		string
	ProposalID	int64
	Voter		sdk.AccAddress
}

func NewMsgSwitch( title string, proposalID int64, voter sdk.AccAddress) MsgSwitch {
	return MsgSwitch{
		Title:title,
		ProposalID: proposalID,
		Voter:      voter,
	}
}

func (msg MsgSwitch) Route() string { return MsgRoute }
func (msg MsgSwitch) Type() string { return "submit_switch" }

func (msg MsgSwitch) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg MsgSwitch) ValidateBasic() sdk.Error {
	return nil
}

func (msg MsgSwitch) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Voter}
}
