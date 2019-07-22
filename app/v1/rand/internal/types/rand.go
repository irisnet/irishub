package types

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

// Rand represents a random number with related data
type Rand struct {
	RequestTxHash []byte  `json:"request_tx_hash"` // the original request tx hash
	Height        int64   `json:"height"`          // the height of the block used to generate the random number
	Value         sdk.Rat `json:"value"`           // the actual random number
}

// NewRand constructs a Rand
func NewRand(requestTxHash []byte, height int64, value sdk.Rat) Rand {
	return Rand{
		RequestTxHash: requestTxHash,
		Height:        height,
		Value:         value,
	}
}

// Validate checks if a rand is valid. Indended to validate random numbers exported to genesis
func (r Rand) Validate() sdk.Error {
	return nil
}

// String implements fmt.Stringer
func (r Rand) String() string {
	return fmt.Sprintf(`Rand:
  RequestTxHash:     %s
  Height:            %d,
  Value:             %s`,
		hex.EncodeToString(r.RequestTxHash), r.Height, r.Value.Rat.FloatString(RandPrec))
}
