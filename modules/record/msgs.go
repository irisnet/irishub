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
	Description  string         // File description
	SubmitTime   int64          // File upload timestamp
	OwnerAddress sdk.AccAddress // Owner of file
	DataHash     string         // IPFS hash of file
	DataSize     int64          // File Size in bytes
	RecordId     string         // Record index ID
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
		Description:  description,
		SubmitTime:   submitTime,
		OwnerAddress: ownerAddress,
		DataHash:     dataHash,
		DataSize:     dataSize,
		RecordId:     string(KeyRecord(ownerAddress, dataHash)),
	}
}

// Implements Msg.
func (msg MsgSubmitFile) Type() string { return MsgType }

// Implements Msg.
func (msg MsgSubmitFile) ValidateBasic() sdk.Error {

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
	return fmt.Sprintf("MsgSubmitFile{%s, %d, %d}",
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
