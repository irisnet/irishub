package types

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

// HTLC represents a HTLC
type HTLC struct {
	Sender               sdk.AccAddress `json:"sender"`                  // the initiator address
	Receiver             sdk.AccAddress `json:"receiver"`                // the recipient address
	ReceiverOnOtherChain []byte         `json:"receiver_in_other_chain"` // the recipient address on other chain
	OutAmount            sdk.Coin       `json:"out_amount"`              // the amount to be transferred
	InAmount             uint64         `json:"in_amount"`               // expected amount to be received from another HTLC
	Secret               []byte         `json:"secret"`                  // the random secret
	Timestamp            uint64         `json:"timestamp"`               // the time used to generate the hash lock together with secret
	ExpireHeight         uint64         `json:"expire_height"`           // the block height by which the HTLC expires
	State                uint8          `json:"state"`                   // the state of the HTLC(0:open,1:completed,2:expired)
}

// NewHTLC constructs a HTLC
func NewHTLC(sender sdk.AccAddress, receiver sdk.AccAddress, receiverOnOtherChain []byte, outAmount sdk.Coin, inAmount uint64, secret []byte, timestamp uint64, expireHeight uint64, state uint8) HTLC {
	return HTLC{
		Sender:               sender,
		Receiver:             receiver,
		ReceiverOnOtherChain: receiverOnOtherChain,
		OutAmount:            outAmount,
		InAmount:             inAmount,
		Secret:               secret,
		Timestamp:            timestamp,
		ExpireHeight:         expireHeight,
		State:                state,
	}
}

// GetSecretHashLock calculates the secret hash lock
func (h HTLC) GetSecretHashLock() []byte {
	return sdk.SHA256(append(h.Secret, sdk.Uint64ToBigEndian(h.Timestamp)...))
}

// String implements fmt.Stringer
func (h HTLC) String() string {
	return fmt.Sprintf(`HTLC:
	Sender:               %s
	Receiver:             %s
	ReceiverOnOtherChain: %v
	OutAmount:            %s
	InAmount:             %d
	Secret:               %s
	Timestamp:            %d
	ExpireHeight:         %d
	State:                %d`,
		h.Sender, h.Receiver, h.ReceiverOnOtherChain, h.OutAmount.String(), h.InAmount, hex.EncodeToString(h.Secret), h.Timestamp, h.ExpireHeight, h.State)
}
