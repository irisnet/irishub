package types

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

// Rand represents a random number related to a height
type Rand struct {
	Height   int64          `json:"height"`   // the height of the block used to generate the random number
	Consumer sdk.AccAddress `json:"consumer"` // the address requesting the random number
	Value    sdk.Dec        `json:"value"`    // the random number
}

// NewRand constructs a Rand
func NewRand(height int64, consumer sdk.AccAddress, value sdk.Dec) Rand {
	return Rand{
		Height:   height,
		Consumer: consumer,
		Value:    value,
	}
}

// Validate checks if a rand is valid. Indended to validate random numbers exported to genesis
func (r Rand) Validate() sdk.Error {
	return nil
}

// String implements fmt.Stringer
func (r Rand) String() string {
	return fmt.Sprintf(`Rand:
  Height:            %d,
  Consumer:          %s,
  Value:             %s`,
		r.Height, r.Consumer.String(), r.Value.String())
}
