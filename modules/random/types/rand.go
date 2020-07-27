package types

// NewRandom constructs a Random
func NewRandom(requestTxHash []byte, height int64, value string) Random {
	return Random{
		RequestTxHash: requestTxHash,
		Height:        height,
		Value:         value,
	}
}

// ReadableRandom for client use
type ReadableRandom struct {
	RequestTxHash string `json:"request_tx_hash" yaml:"request_tx_hash"` // the original request tx hash
	Height        int64  `json:"height" yaml:"height"`                   // the height of the block where the random number is generated
	Value         string `json:"value" yaml:"value"`                     // the actual random number
}
