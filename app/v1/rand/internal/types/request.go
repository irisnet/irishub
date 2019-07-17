package types

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

// Request represents a request for random number
type Request struct {
	Height   int64          `json:"height"`   // the height of the block in which the request tx is included
	Consumer sdk.AccAddress `json:"consumer"` // the address of the request account
}

// NewRequest constructs a request
func NewRequest(height int64, consumer sdk.AccAddress) Request {
	return Request{
		Height:   height,
		Consumer: consumer,
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
  Consumer:          %s`,
		r.Height, r.Consumer.String())
}

// GenerateRequestID generates a request id
func GenerateRequestID(request Request) string {
	reqID := make([]byte, 0)

	reqID = append(reqID, sdk.Uint64ToBigEndian(uint64(request.Height))...)
	reqID = append(reqID, []byte(request.Consumer)...)

	return hex.EncodeToString(sdk.SHA256(reqID))
}
