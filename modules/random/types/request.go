package types

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

// NewRequest constructs a request
func NewRequest(
	height int64,
	consumer sdk.AccAddress,
	txHash tmbytes.HexBytes,
	oracle bool,
	serviceFeeCap sdk.Coins,
	serviceContextID tmbytes.HexBytes,
) Request {
	return Request{
		Height:           height,
		Consumer:         consumer,
		TxHash:           txHash,
		Oracle:           oracle,
		ServiceFeeCap:    serviceFeeCap,
		ServiceContextID: serviceContextID,
	}
}

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
