package record

import (
	"encoding/binary"
	"testing"
	"time"

	"github.com/irisnet/irishub/simulation/mock"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// test ValidateBasic for MsgSubmitRecord
func TestMsgSubmitRecord(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(1, sdk.Coins{})

	normalDescription := "record description"
	emptyDescription := ""
	overFlowedDescription := createOverflowedData(normalDescription, UploadLimitOfDescription)

	normalRecordData := "record data"
	emptyRecordData := ""
	overflowedRecordData := createOverflowedData(normalRecordData, UploadLimitOfOnchain)

	tests := []struct {
		submitTime   int64
		ownerAddress sdk.AccAddress
		recordID     string
		description  string
		dataHash     string
		dataSize     int64
		data         string
		expectPass   bool
	}{
		// -------------------Data Field-------------------------
		{
			time.Now().Unix(),
			addrs[0],
			string(KeyRecord(getDataHash(normalRecordData))),
			normalDescription,
			getDataHash(normalRecordData),
			int64(binary.Size([]byte(normalRecordData))),
			normalRecordData,
			true,
		},
		{
			time.Now().Unix(),
			addrs[0],
			string(KeyRecord(getDataHash(emptyRecordData))),
			normalDescription,
			getDataHash(emptyRecordData),
			int64(binary.Size([]byte(emptyRecordData))),
			emptyRecordData,
			false,
		},
		{
			time.Now().Unix(),
			addrs[0],
			string(KeyRecord(getDataHash(overflowedRecordData))),
			normalDescription,
			getDataHash(overflowedRecordData),
			int64(binary.Size([]byte(overflowedRecordData))),
			overflowedRecordData,
			false,
		},
		// -------------------DataHash Field-------------------------
		{
			time.Now().Unix(),
			addrs[0],
			string(KeyRecord(getDataHash(normalRecordData))),
			normalDescription,
			"",
			int64(binary.Size([]byte(normalRecordData))),
			normalRecordData,
			false,
		},
		// -------------------OwnerAddress Field-------------------------
		{
			time.Now().Unix(),
			sdk.AccAddress{},
			string(KeyRecord(getDataHash(normalRecordData))),
			normalDescription,
			getDataHash(normalRecordData),
			int64(binary.Size([]byte(normalRecordData))),
			normalRecordData,
			false,
		},
		// -------------------Description Field--------------------------
		{
			time.Now().Unix(),
			addrs[0],
			string(KeyRecord(getDataHash(normalRecordData))),
			emptyDescription,
			getDataHash(normalRecordData),
			int64(binary.Size([]byte(normalRecordData))),
			normalRecordData,
			false,
		},
		{
			time.Now().Unix(),
			addrs[0],
			string(KeyRecord(getDataHash(normalRecordData))),
			overFlowedDescription,
			getDataHash(normalRecordData),
			int64(binary.Size([]byte(normalRecordData))),
			normalRecordData,
			false,
		},
	}

	for i, tc := range tests {
		msg := NewMsgSubmitRecord(tc.description, tc.submitTime, tc.ownerAddress, tc.dataHash, tc.dataSize, tc.data)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}
