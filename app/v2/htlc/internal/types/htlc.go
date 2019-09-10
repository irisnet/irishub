package types

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/irisnet/irishub/types"
)

// the state of the HTLC
const (
	StateOpen      = uint8(0) // can claim
	StateCompleted = uint8(1) // claimed
	StateExpired   = uint8(2) // Expired
	StateRefunded  = uint8(3) // Refunded
)

// HTLC represents a HTLC
type HTLC struct {
	Sender               sdk.AccAddress `json:"sender"`                  // the initiator address
	Receiver             sdk.AccAddress `json:"receiver"`                // the recipient address
	ReceiverOnOtherChain []byte         `json:"receiver_on_other_chain"` // the recipient address on other chain
	OtherChainName       string         `json:"other_chain_name"`        // the name of the other chain
	OutAmount            sdk.Coin       `json:"out_amount"`              // the amount to be transferred
	InTokenName          string         `json:"in_token_name"`           // the name of the expected token to be received from other chain
	InAmount             uint64         `json:"in_amount"`               // expected amount to be received from other chain, which is of 8 decimals
	Secret               []byte         `json:"secret"`                  // the random secret which is 32 bytes length
	Timestamp            uint64         `json:"timestamp"`               // the time used to generate the hash lock together with secret if provided
	ExpireHeight         uint64         `json:"expire_height"`           // the block height by which the HTLC expires
	State                uint8          `json:"state"`                   // the state of the HTLC(0:open,1:completed,2:expired,3:refunded)
}

// NewHTLC constructs a HTLC
func NewHTLC(
	sender sdk.AccAddress,
	receiver sdk.AccAddress,
	receiverOnOtherChain []byte,
	outAmount sdk.Coin,
	inAmount uint64,
	secret []byte,
	timestamp uint64,
	expireHeight uint64,
	state uint8,
) HTLC {
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
	ReceiverOnOtherChain: %v
	OutAmount:            %s
	InAmount:             %d
	Secret:               %s
	Timestamp:            %d
	ExpireHeight:         %d
	State:                %d`,
		h.Sender,
		h.Receiver,
		h.ReceiverOnOtherChain,
		h.OutAmount.String(),
		h.InAmount,
		hex.EncodeToString(h.Secret),
		h.Timestamp,
		h.ExpireHeight,
		h.State,
	)
}
