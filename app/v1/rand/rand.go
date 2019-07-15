package rand

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

// Rand represents a random number related to a height
type Rand struct {
	Height int64   `json:"height"` // the height of the block, of which the hash is used to generate the random number
	Value  sdk.Dec `json:"value"`  // the random number
}

// NewRand constructs a Rand
func NewRand(height int64, value sdk.Dec) Rand {
	return Rand{
		Height: height,
		Value:  value,
	}
}

// Validate checks if a rand is valid. Indended to validate random numbers exported to genesis
func (r Rand) Validate() sdk.Error {
	return nil
}

// String implements fmt.Stringer
func (r Rand) String() string {
	return fmt.Sprintf(`Rand:
  Height:            %s
  Value:             %s`,
		r.Height, r.Value.String())
}
