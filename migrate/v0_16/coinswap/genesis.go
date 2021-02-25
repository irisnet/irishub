package coinswap

import (
	"math/big"
)

type GenesisState struct {
	Params Params `json:"params"`
}

type Params struct {
	Fee Rat `json:"fee"`
}

type Rat struct {
	*big.Rat `json:"rat"`
}

// Requires a valid JSON string - strings quotes and calls UnmarshalText
func (r *Rat) UnmarshalAmino(text string) (err error) {
	tempRat := big.NewRat(0, 1)
	err = tempRat.UnmarshalText([]byte(text))
	if err != nil {
		return err
	}
	r.Rat = tempRat
	return nil
}
