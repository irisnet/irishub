package types

import (
	"encoding/hex"
	"fmt"

	"github.com/irisnet/irishub/app/v2/htlc"
	sdk "github.com/irisnet/irishub/types"
)

// HTLC represents an HTLC
type OutputHTLC struct {
	Sender               sdk.AccAddress `json:"sender"`                  // the initiator address
	To                   sdk.AccAddress `json:"to"`                      // the destination address
	ReceiverOnOtherChain string         `json:"receiver_on_other_chain"` // the claim receiving address on the other chain
	Amount               sdk.Coins      `json:"amount"`                  // the amount to be transferred
	Secret               string         `json:"secret"`                  // the random secret which is of 32 bytes
	Timestamp            uint64         `json:"timestamp"`               // the timestamp, if provided, used to generate the hash lock together with secret
	ExpireHeight         uint64         `json:"expire_height"`           // the block height by which the HTLC expires
	State                htlc.HTLCState `json:"state"`                   // the state of the HTLC
}

func NewOutputHTLC(h htlc.HTLC) OutputHTLC {
	return OutputHTLC{
		h.Sender,
		h.To,
		h.ReceiverOnOtherChain,
		h.Amount,
		hex.EncodeToString(h.Secret),
		h.Timestamp,
		h.ExpireHeight,
		h.State,
	}
}

// String implements fmt.Stringer
func (h OutputHTLC) String() string {
	return fmt.Sprintf(`HTLC:
	Sender:               %s
	Receiver:             %s
	ReceiverOnOtherChain: %s
	Amount:               %s
	Secret:               %s
	Timestamp:            %d
	ExpireHeight:         %d
	State:                %s`,
		h.Sender.String(),
		h.To.String(),
		h.ReceiverOnOtherChain,
		h.Amount.String(),
		h.Secret,
		h.Timestamp,
		h.ExpireHeight,
		h.State,
	)
}

// HumanString implements human
func (h OutputHTLC) HumanString(converter sdk.CoinsConverter) string {
	return fmt.Sprintf(`HTLC:
	Sender:               %s
	To:                   %s
	ReceiverOnOtherChain: %s
	Amount:               %s
	Secret:               %s
	Timestamp:            %d
	ExpireHeight:         %d
	State:                %s`,
		h.Sender,
		h.To,
		h.ReceiverOnOtherChain,
		converter.ToMainUnit(h.Amount),
		h.Secret,
		h.Timestamp,
		h.ExpireHeight,
		h.State,
	)
}
