package types

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

// the states of the HTLC
const (
	StateOpen      = uint8(0) // not claimed
	StateCompleted = uint8(1) // claimed
	StateExpired   = uint8(2) // Expired
	StateRefunded  = uint8(3) // Refunded
)

// HTLC represents a HTLC
type HTLC struct {
	Sender               sdk.AccAddress `json:"sender"`                  // the initiator address
	Receiver             sdk.AccAddress `json:"receiver"`                // the recipient address
	ReceiverOnOtherChain []byte         `json:"receiver_on_other_chain"` // the recipient address on other chain
	Amount               sdk.Coin       `json:"amount"`                  // the amount to be transferred
	Secret               []byte         `json:"secret"`                  // the random secret which is of 32 bytes
	Timestamp            uint64         `json:"timestamp"`               // the timestamp, if provided, used to generate the hash lock together with secret
	ExpireHeight         uint64         `json:"expire_height"`           // the block height by which the HTLC expires
	State                uint8          `json:"state"`                   // the state of the HTLC
}

// NewHTLC constructs a HTLC
func NewHTLC(
	sender sdk.AccAddress,
	receiver sdk.AccAddress,
	receiverOnOtherChain []byte,
	amount sdk.Coin,
	secret []byte,
	timestamp uint64,
	expireHeight uint64,
	state uint8,
) HTLC {
	return HTLC{
		Sender:               sender,
		Receiver:             receiver,
		ReceiverOnOtherChain: receiverOnOtherChain,
		Amount:               amount,
		Secret:               secret,
		Timestamp:            timestamp,
		ExpireHeight:         expireHeight,
		State:                state,
	}
}

// GetHashLock calculates the hash lock
func (h HTLC) GetHashLock() []byte {
	if h.Timestamp > 0 {
		return sdk.SHA256(append(h.Secret, sdk.Uint64ToBigEndian(h.Timestamp)...))
	}

	return sdk.SHA256(h.Secret)
}

// String implements fmt.Stringer
func (h HTLC) String() string {
	return fmt.Sprintf(`HTLC:
	Sender:               %s
	Receiver:             %s
	ReceiverOnOtherChain: %s
	Amount:               %s
	Secret:               %s
	Timestamp:            %d
	ExpireHeight:         %d
	State:                %d`,
		h.Sender,
		h.Receiver,
		hex.EncodeToString(h.ReceiverOnOtherChain),
		h.Amount.String(),
		hex.EncodeToString(h.Secret),
		h.Timestamp,
		h.ExpireHeight,
		h.State,
	)
}
