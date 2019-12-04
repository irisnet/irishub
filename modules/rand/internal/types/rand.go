package types

import (
	"encoding/hex"
	"fmt"
	"math/big"
)

// Rand represents a random number with related data
type Rand struct {
	RequestTxHash []byte   `json:"request_tx_hash"` // the original request tx hash
	Height        int64    `json:"height"`          // the height of the block where the random number is generated
	Value         *big.Rat `json:"value"`           // the actual random number
}

// NewRand constructs a Rand
func NewRand(requestTxHash []byte, height int64, value *big.Rat) Rand {
	return Rand{
		RequestTxHash: requestTxHash,
		Height:        height,
		Value:         value,
	}
}

// String implements fmt.Stringer
func (r Rand) String() string {
	return fmt.Sprintf(`Rand:
  RequestTxHash:     %s
  Height:            %d
  Value:             %s`,
		hex.EncodeToString(r.RequestTxHash), r.Height, r.Value.FloatString(RandPrec))
}
