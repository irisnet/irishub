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

// String implements fmt.Stringer
func (r Rand) String() string {
	return fmt.Sprintf(`Rand:
  RequestTxHash:     %s
  Height:            %d
  Value:             %s`,
		hex.EncodeToString(r.RequestTxHash), r.Height, r.Value.Rat.FloatString(RandPrec))
}

// ReadableRand represents a shadow Rand intended for readable output
type ReadableRand struct {
	RequestTxHash string `json:"request_tx_hash"`
	Height        int64  `json:"height"`
	Value         string `json:"value"`
}

// MarshalJSON marshals rand to readable JSON
func (r Rand) MarshalJSON() ([]byte, error) {
	readableRand := ReadableRand{
		RequestTxHash: hex.EncodeToString(r.RequestTxHash),
		Height:        r.Height,
		Value:         r.Value.Rat.FloatString(RandPrec),
	}

	return msgCdc.MarshalJSON(readableRand)
}

// UnmarshalJSON unmarshals data to Rand
func (r *Rand) UnmarshalJSON(data []byte) error {
	var readableRand ReadableRand

	err := msgCdc.UnmarshalJSON(data, &readableRand)
	if err != nil {
		return err
	}

	txHash, err := hex.DecodeString(readableRand.RequestTxHash)
	if err != nil {
		return err
	}

	value, _ := sdk.NewRatFromDecimal(readableRand.Value, RandPrec)
	if err != nil {
		return err
	}

	rawRand := Rand{
		RequestTxHash: txHash,
		Height:        readableRand.Height,
		Value:         value,
	}

	*r = rawRand
	return nil
}
