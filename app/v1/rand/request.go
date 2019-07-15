package rand

import (
	"crypto/sha256"
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

// generateRequestID generates a request id
func generateRequestID(request Request) string {
	reqIDBytes := make([]byte, 0)

	reqIDBytes = append(reqIDBytes, sdk.Uint64ToBigEndian(uint64(request.Height))...)
	reqIDBytes = append(reqIDBytes, []byte(request.Consumer)...)

	hash := sha256.Sum256(reqIDBytes)
	return hex.EncodeToString(hash[:])
}
