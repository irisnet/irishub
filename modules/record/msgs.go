package record

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// name to idetify transaction types
const MsgType = "record"

//-----------------------------------------------------------
// MsgSubmitFile
type MsgSubmitFile struct {
	Filename    string         //  Filename of the File
	Filepath    string         //  full path of the File
	Description string         //  Description of the File
	FileType    string         //  Type of file
	Proposer    sdk.AccAddress //  Address of the proposer
	Amount      sdk.Coins      //  File upload paid by sender. Must be non negative.
	DataSize    int64          // File Size in bytes
}

//msg := record.NewMsgSubmitFile(filename, filepath, description, fileType, fromAddr, amount)
func NewMsgSubmitFile(filename string, filepath string, description string, fileType string, proposer sdk.AccAddress, amount sdk.Coins, dataSize int64) MsgSubmitFile {
	return MsgSubmitFile{
		Filename:    filename,
		Filepath:    filepath,
		Description: description,
		FileType:    fileType,
		Proposer:    proposer,
		Amount:      amount,
		DataSize:    dataSize,
	}
}

// Implements Msg.
func (msg MsgSubmitFile) Type() string { return MsgType }

// Implements Msg.
func (msg MsgSubmitFile) ValidateBasic() sdk.Error {
	if len(msg.Filename) == 0 {
		return ErrInvalidFilename(DefaultCodespace, msg.Filename) // TODO: Proper Error
	}
	if len(msg.Description) == 0 {
		return ErrInvalidDescription(DefaultCodespace, msg.Description) // TODO: Proper Error
	}

	if len(msg.Proposer) == 0 {
		return sdk.ErrInvalidAddress(msg.Proposer.String())
	}

	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}
	if !msg.Amount.IsNotNegative() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}

	return nil
}

func (msg MsgSubmitFile) String() string {
	return fmt.Sprintf("MsgSubmitFile{%s, %s, %s, %v}", msg.Filename, msg.Description, msg.FileType, msg.Amount)
}

// Implements Msg.
func (msg MsgSubmitFile) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg MsgSubmitFile) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgSubmitFile) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proposer}
}
