package types

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

// Request represents a request for a random number
type Request struct {
	Height   int64          `json:"height"`   // the height of the block in which the request tx is included
	Consumer sdk.AccAddress `json:"consumer"` // the request address
	TxHash   []byte         `json:"txhash"`   // the request tx hash
}

// NewRequest constructs a request
func NewRequest(height int64, consumer sdk.AccAddress, txHash []byte) Request {
	return Request{
		Height:   height,
		Consumer: consumer,
		TxHash:   txHash,
	}
}

// Validate checks if a request is valid. Indended to validate requests exported to genesis
func (r Request) Validate() sdk.Error {
	return nil
}

// String implements fmt.Stringer
func (r Request) String() string {
	return fmt.Sprintf(`Request:
  Height:            %d
  Consumer:          %s
  TxHash:            %s`,
		r.Height, r.Consumer.String(), string(r.TxHash))
}

// GenerateRequestID generate a request id
func GenerateRequestID(r Request) string {
	id := make([]byte, 0)

	id = append(id, sdk.Uint64ToBigEndian(uint64(r.Height))...)
	id = append(id, []byte(r.Consumer)...)

	return hex.EncodeToString(sdk.SHA256(id))
}

// CheckReqID checks if the given request id is valid
func CheckReqID(reqID string) sdk.Error {
	if len(reqID) != 64 {
		return ErrInvalidReqID(DefaultCodespace, fmt.Sprintf("invalid request id: %s", reqID))
	}

	if _, err := hex.DecodeString(reqID); err != nil {
		return ErrInvalidReqID(DefaultCodespace, fmt.Sprintf("invalid request id: %s", reqID))
	}

	return nil
}
