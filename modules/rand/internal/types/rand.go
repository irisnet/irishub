package types

// Rand represents a random number with related data
type Rand struct {
	RequestTxHash []byte `json:"request_tx_hash" yaml:"request_tx_hash"` // the original request tx hash
	Height        int64  `json:"height" yaml:"height"`                   // the height of the block where the random number is generated
	Value         string `json:"value" yaml:"value"`                     // the actual random number
}

// NewRand constructs a Rand
func NewRand(requestTxHash []byte, height int64, value string) Rand {
	return Rand{
		RequestTxHash: requestTxHash,
		Height:        height,
		Value:         value,
	}
}

// ReadableRand for client use
type ReadableRand struct {
	RequestTxHash string `json:"request_tx_hash" yaml:"request_tx_hash"` // the original request tx hash
	Height        int64  `json:"height" yaml:"height"`                   // the height of the block where the random number is generated
	Value         string `json:"value" yaml:"value"`                     // the actual random number
}
