package cli

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgRecord struct {
	OwnerAddress string
	SubmitTime   string
	DataHash     string
	DataSize     int
	PinedNode    string
}

var tmp_msg = MsgRecord{

	OwnerAddress: "this is owner address",
	SubmitTime:   "this is submit time",
	DataHash:     "this is data hash",
	DataSize:     1000,
	PinedNode:    "this is pinednode",
}

func (msg MsgRecord) Type() string { return "" }

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
	return []byte("")
}

func (msg MsgRecord) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}
