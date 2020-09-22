package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// QueryRecordParams defines QueryRecord params
type QueryRecordParams struct {
	RecordID []byte `json:"record_id"`
}

type RecordOutput struct {
	TxHash   string         `json:"tx_hash" yaml:"tx_hash"`
	Contents []Content      `json:"contents" yaml:"contents"`
	Creator  sdk.AccAddress `json:"creator" yaml:"creator"`
}
