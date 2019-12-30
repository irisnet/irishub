package types

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Request represents a request for a random number
type Request struct {
	Height   int64          `json:"height" yaml:"height"`     // the height of the block in which the request tx is included
	Consumer sdk.AccAddress `json:"consumer" yaml:"consumer"` // the request address
	TxHash   []byte         `json:"tx_hash" yaml:"tx_hash"`   // the request tx hash
}

// NewRequest constructs a request
func NewRequest(height int64, consumer sdk.AccAddress, txHash []byte) Request {
	return Request{
		Height:   height,
		Consumer: consumer,
		TxHash:   txHash,
	}
}

// Requests is a set of requests
type Requests []Request

// GenerateRequestID generate a request id
func GenerateRequestID(r Request) []byte {
	reqID := make([]byte, 0)

	reqID = append(reqID, sdk.Uint64ToBigEndian(uint64(r.Height))...)
	reqID = append(reqID, []byte(r.Consumer)...)

	return SHA256(reqID)
}

// CheckReqID checks if the given request id is valid
func CheckReqID(reqID string) error {
	if len(reqID) != 64 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("invalid request id: %s", reqID))
	}

	if _, err := hex.DecodeString(reqID); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("invalid request id: %s", reqID))
	}

	return nil
}
