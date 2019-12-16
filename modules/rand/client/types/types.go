package types

import "fmt"

// ReadableRand represents a shadow Rand intended for readable output
type ReadableRand struct {
	RequestTxHash string `json:"request_tx_hash" yaml:"request_tx_hash"`
	Height        int64  `json:"height" yaml:"height"`
	Value         string `json:"value" yaml:"value"`
}

// String implements fmt.Stringer
func (rr ReadableRand) String() string {
	return fmt.Sprintf(`Rand:
  RequestTxHash:     %s
  Height:            %d
  Value:             %s`,
		rr.RequestTxHash, rr.Height, rr.Value)
}
