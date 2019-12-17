package types

// ReadableRand represents a shadow Rand intended for readable output
type ReadableRand struct {
	RequestTxHash string `json:"request_tx_hash" yaml:"request_tx_hash"`
	Height        int64  `json:"height" yaml:"height"`
	Value         string `json:"value" yaml:"value"`
}
