package types

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

// String implements fmt.Stringer
func (r Request) String() string {
	return fmt.Sprintf(`Request:
  Height:            %d
  Consumer:          %s
  TxHash:            %s`,
		r.Height, r.Consumer.String(), hex.EncodeToString(r.TxHash))
}

// Requests is a set of requests
type Requests []Request

// String implements fmt.Stringer
func (rs Requests) String() string {
	if len(rs) == 0 {
		return "[]"
	}

	var str string
	for _, r := range rs {
		str += fmt.Sprintf("Request:\n  Height: %d, Consumer: %s, TxHash: %s", r.Height, r.Consumer.String(), hex.EncodeToString(r.TxHash))
	}

	return str
}

// GenerateRequestID generate a request id
func GenerateRequestID(r Request) []byte {
	reqID := make([]byte, 0)

	reqID = append(reqID, sdk.Uint64ToBigEndian(uint64(r.Height))...)
	reqID = append(reqID, []byte(r.Consumer)...)

	return SHA256(reqID)
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
