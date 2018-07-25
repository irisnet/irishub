package upgrade

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgSwitch struct {
	Title		string
	ProposalID	int
	Voter		sdk.AccAddress
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
	return nil
}

func (msg MsgSwitch) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Voter}
}
