package gov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/modules/record"
)

type RecordOutput struct {
	Description  string         `json:"Description"`  // File description
	SubmitTime   int64          `json:"SubmitTime"`   // File upload timestamp
	OwnerAddress sdk.AccAddress `json:"OwnerAddress"` // Owner of file
	DataHash     string         `json:"DataHash"`     // IPFS hash of file
	DataSize     int64          `json:"DataSize"`     // File Size in bytes
	RecordId     string         `json:"RecordId"`     // Record index ID
}

func ConvertRecordToRecordOutput(cliCtx context.CLIContext, r record.MsgSubmitFile) (RecordOutput, error) {
	recordOutput := RecordOutput{
		Description:  r.Description,
		SubmitTime:   r.SubmitTime,
		OwnerAddress: r.OwnerAddress,
		DataHash:     r.DataHash,
		DataSize:     r.DataSize,
		RecordId:     r.RecordId,
	}

	return recordOutput, nil
}
