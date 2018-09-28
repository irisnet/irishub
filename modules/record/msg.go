package record

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// name to idetify transaction types
const MsgType = "record"

// MsgRecord
type MsgRecord struct {
}

func NewMsgRecord(tx string, file string, owner sdk.AccAddress) MsgRecord {
	return MsgRecord{}
}

func (msg MsgRecord) Type() string { return MsgType }

func (msg MsgRecord) ValidateBasic() sdk.Error {
	// TO DO
	return nil
}

func (msg MsgRecord) String() string {
	return ""
}

func (msg MsgRecord) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg MsgRecord) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgRecord) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}
