package gov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/modules/record"
)

type RecordOutput struct {
	Filename     string         `json:"Filename"`     //  Filename of the File
	Filepath     string         `json:"Filepath"`     //  full path of the File
	Description  string         `json:"Description"`  //  Description of the File
	SubmitTime   int64          `json:"SubmitTime"`   //  File  submit unix timestamp
	OwnerAddress sdk.AccAddress `json:"OwnerAddress"` //  Address of the owner
	DataHash     string         `json:"DataHash"`     // ipfs hash of file
	DataSize     int64          `json:"DataSize"`     // File Size in bytes
	//PinedNode    string        `json:"PinedNode"` //pined node of ipfs
}

func ConvertRecordToRecordOutput(cliCtx context.CLIContext, r record.MsgSubmitFile) (RecordOutput, error) {

	// TODO : Currently we only copy values from record msg, we can call related methods later
	recordOutput := RecordOutput{
		Filename:     r.Filename,
		Filepath:     r.Filepath,
		Description:  r.Description,
		SubmitTime:   r.SubmitTime,
		OwnerAddress: r.OwnerAddress,
		DataHash:     r.DataHash,
		DataSize:     r.DataSize,
	}

	return recordOutput, nil
}
