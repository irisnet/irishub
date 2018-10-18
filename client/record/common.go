package record

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/modules/record"
)

const (
	UploadLimitOfOnchain = 1024        //upload limit on chain in bytes(1K currently)
	UploadLimitOfIpfs    = 1024 * 1024 //upload limit on chain in bytes(1M currently)
	IpfsHashLength       = 46
)

type RecordOutput struct {
	SubmitTime   string         `json:"submit_time"` // File upload timestamp
	OwnerAddress sdk.AccAddress `json:"owner_addr"`  // Owner of file
	RecordID     string         `json:"record_id"`   // Record index ID
	Description  string         `json:"description"` // Data/file description
	DataHash     string         `json:"data_hash"`   // Data/file hash
	DataSize     int64          `json:"data_size"`   // Data/file Size in bytes
	Data         string         `json:"data"`        // Onchain data
}

func ConvertRecordToRecordOutput(cliCtx context.CLIContext, r record.MsgSubmitFile) (RecordOutput, error) {

	utcTime := time.Unix(r.SubmitTime, 0).Format("2006-01-02 15:04:05")

	recordOutput := RecordOutput{
		SubmitTime:   utcTime,
		OwnerAddress: r.OwnerAddress,
		RecordID:     r.RecordID,
		Description:  r.Description,
		DataHash:     r.DataHash,
		DataSize:     r.DataSize,
		Data:         r.Data,
	}

	return recordOutput, nil
}
