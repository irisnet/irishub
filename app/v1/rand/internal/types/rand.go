package types

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

// Rand represents a random number with related data
type Rand struct {
	Request `json:"request"` // the original request
	Height  int64            `json:"height"` // the height of the block used to generate the random number
	Value   sdk.Dec          `json:"value"`  // the actual random number
}

// NewRand constructs a Rand
func NewRand(request Request, height int64, value sdk.Dec) Rand {
	return Rand{
		Request: request,
		Height:  height,
		Value:   value,
	}
}

// Validate checks if a rand is valid. Indended to validate random numbers exported to genesis
func (r Rand) Validate() sdk.Error {
	return nil
}

// String implements fmt.Stringer
func (r Rand) String() string {
	return fmt.Sprintf(`Rand:
  Request:           %s
  Height:            %d,
  Value:             %s`,
		r.Request.String(), r.Height, r.Value.String())
}
