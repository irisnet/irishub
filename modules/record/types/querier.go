package types

// QueryRecordParams defines QueryRecord params
type QueryRecordParams struct {
	RecordID []byte `json:"record_id"`
}

type RecordOutput struct {
	TxHash   string    `json:"tx_hash" yaml:"tx_hash"`
	Contents []Content `json:"contents" yaml:"contents"`
	Creator  string    `json:"creator" yaml:"creator"`
}
