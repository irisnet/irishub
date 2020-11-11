package types

// NewRandom constructs a Random
func NewRandom(requestTxHash string, height int64, value string) Random {
	return Random{
		RequestTxHash: requestTxHash,
		Height:        height,
		Value:         value,
	}
}
