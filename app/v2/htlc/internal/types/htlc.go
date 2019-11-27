package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	sdk "github.com/irisnet/irishub/types"
)

// HTLC represents an HTLC
type HTLC struct {
	Sender               sdk.AccAddress `json:"sender"`                  // the initiator address
	To                   sdk.AccAddress `json:"to"`                      // the destination address
	ReceiverOnOtherChain string         `json:"receiver_on_other_chain"` // the claim receiving address on the other chain
	Amount               sdk.Coins      `json:"amount"`                  // the amount to be transferred
	Secret               []byte         `json:"secret"`                  // the random secret which is of 32 bytes
	Timestamp            uint64         `json:"timestamp"`               // the timestamp, if provided, used to generate the hash lock together with secret
	ExpireHeight         uint64         `json:"expire_height"`           // the block height by which the HTLC expires
	State                HTLCState      `json:"state"`                   // the state of the HTLC
}

// NewHTLC constructs an HTLC
func NewHTLC(
	sender sdk.AccAddress,
	to sdk.AccAddress,
	receiverOnOtherChain string,
	amount sdk.Coins,
	secret []byte,
	timestamp uint64,
	expireHeight uint64,
	state HTLCState,
) HTLC {
	return HTLC{
		Sender:               sender,
		To:                   to,
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
	if h.State == COMPLETED {
		if h.Timestamp > 0 {
			return sdk.SHA256(append(h.Secret, sdk.Uint64ToBigEndian(h.Timestamp)...))
		}

		return sdk.SHA256(h.Secret)
	}

	return nil
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
	State:                %s`,
		h.Sender.String(),
		h.To.String(),
		h.ReceiverOnOtherChain,
		h.Amount.String(),
		hex.EncodeToString(h.Secret),
		h.Timestamp,
		h.ExpireHeight,
		h.State,
	)
}

// HumanString implements human
func (h HTLC) HumanString(converter sdk.CoinsConverter) string {
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
		hex.EncodeToString(h.Secret),
		h.Timestamp,
		h.ExpireHeight,
		h.State,
	)
}

// HTLCState represents the state of an HTLC
type HTLCState byte

const (
	OPEN      HTLCState = 0x00 // claimable
	COMPLETED HTLCState = 0x01 // claimed
	EXPIRED   HTLCState = 0x02 // expired
	REFUNDED  HTLCState = 0x03 // refunded
)

var (
	HTLCStateToStringMap = map[HTLCState]string{
		OPEN:      "open",
		COMPLETED: "completed",
		EXPIRED:   "expired",
		REFUNDED:  "refunded",
	}
	StringToHTLCStateMap = map[string]HTLCState{
		"open":      OPEN,
		"completed": COMPLETED,
		"expired":   EXPIRED,
		"refunded":  REFUNDED,
	}
)

func HTLCStateFromString(str string) (HTLCState, error) {
	if state, ok := StringToHTLCStateMap[strings.ToLower(str)]; ok {
		return state, nil
	}
	return HTLCState(0xff), fmt.Errorf("'%s' is not a valid HTLC state", str)
}

func (state HTLCState) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(fmt.Sprintf("%s", state.String())))
	default:
		s.Write([]byte(fmt.Sprintf("%v", byte(state))))
	}
}

func (state HTLCState) String() string {
	return HTLCStateToStringMap[state]
}

// Marshal needed for protobuf compatibility
func (state HTLCState) Marshal() ([]byte, error) {
	return []byte{byte(state)}, nil
}

// Unmarshal needed for protobuf compatibility
func (state *HTLCState) Unmarshal(data []byte) error {
	*state = HTLCState(data[0])
	return nil
}

// Marshals to JSON using string
func (state HTLCState) MarshalJSON() ([]byte, error) {
	return json.Marshal(state.String())
}

// Unmarshals from JSON
func (state *HTLCState) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil
	}

	bz, err := HTLCStateFromString(s)
	if err != nil {
		return err
	}
	*state = bz
	return nil
}
