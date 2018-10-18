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
	Filename     string         //  Filename of the File
	Filepath     string         //  full path of the File
	Description  string         //  Description of the File
	SubmitTime   int64          //  File  submit unix timestamp
	OwnerAddress sdk.AccAddress //  Address of the owner
	DataHash     string         // ipfs hash of file
	DataSize     int64          // File Size in bytes
	RecordId     string         // Record key of kvstore
	PinedNode    string         //pined node of ipfs
}

func NewMsgSubmitFile(filename string,
	filepath string,
	description string,
	submitTime int64,
	ownerAddress sdk.AccAddress,
	dataHash string,
	dataSize int64,
	pinedNode string) MsgSubmitFile {
	return MsgSubmitFile{
		Filename:     filename,
		Filepath:     filepath,
		Description:  description,
		SubmitTime:   submitTime,
		OwnerAddress: ownerAddress,
		DataHash:     dataHash,
		DataSize:     dataSize,
		RecordId:     string(KeyRecord(ownerAddress, dataHash)),
		PinedNode:    pinedNode,
	}
}

// Implements Msg.
func (msg MsgSubmitFile) Type() string { return MsgType }

// Implements Msg.
func (msg MsgSubmitFile) ValidateBasic() sdk.Error {
	if len(msg.Filename) == 0 {
		return ErrInvalidFilename(DefaultCodespace, msg.Filename)
	}

	if len(msg.Description) == 0 {
		return ErrInvalidDescription(DefaultCodespace, msg.Description)
	}

	if len(msg.DataHash) == 0 {
		return ErrFailUploadFile(DefaultCodespace, msg.DataHash)
	}

	if len(msg.OwnerAddress) == 0 {
		return sdk.ErrInvalidAddress(msg.OwnerAddress.String())
	}

	return nil
}

func (msg MsgSubmitFile) String() string {
	return fmt.Sprintf("MsgSubmitFile{%s, %s, %d, %d}",
		msg.Filename,
		msg.OwnerAddress,
		msg.DataSize,
		msg.SubmitTime,
	)
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
	return []sdk.AccAddress{msg.OwnerAddress}
}
