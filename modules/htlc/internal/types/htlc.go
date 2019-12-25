package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HTLC represents an HTLC
type HTLC struct {
	Sender               sdk.AccAddress `json:"sender" yaml:"sender"`                                   // the initiator address
	To                   sdk.AccAddress `json:"to" yaml:"to"`                                           // the destination address
	ReceiverOnOtherChain string         `json:"receiver_on_other_chain" yaml:"receiver_on_other_chain"` // the claim receiving address on the other chain
	Amount               sdk.Coins      `json:"amount" yaml:"amount"`                                   // the amount to be transferred
	Secret               HTLCSecret     `json:"secret" yaml:"secret"`                                   // the random secret which is of 32 bytes
	Timestamp            uint64         `json:"timestamp" yaml:"timestamp"`                             // the timestamp, if provided, used to generate the hash lock together with secret
	ExpireHeight         uint64         `json:"expire_height" yaml:"expire_height"`                     // the block height by which the HTLC expires
	State                HTLCState      `json:"state" yaml:"state"`                                     // the state of the HTLC
}

// NewHTLC constructs an HTLC
func NewHTLC(
	sender sdk.AccAddress,
	to sdk.AccAddress,
	receiverOnOtherChain string,
	amount sdk.Coins,
	secret HTLCSecret,
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
			return SHA256(append(h.Secret, sdk.Uint64ToBigEndian(h.Timestamp)...))
		}
		return SHA256(h.Secret)
	}
	return nil
}

// Validate returns is the HTLC invalid
func (h HTLC) Validate(hashLock HTLCHashLock) error {
	if len(hashLock) != HashLockLength {
		return fmt.Errorf("the hash lock must be %d bytes long", HashLockLength)
	}
	if h.State != OPEN {
		return fmt.Errorf("htlc state must be OPEN")
	}
	if len(h.Sender) == 0 {
		return fmt.Errorf("the sender address must be specified")
	}
	if len(h.To) == 0 {
		return fmt.Errorf("the receiver address must be specified")
	}
	if len(h.ReceiverOnOtherChain) > MaxLengthForAddressOnOtherChain {
		return fmt.Errorf("the length of the receiver on other chain must be between [0,%d]", MaxLengthForAddressOnOtherChain)
	}
	if !h.Amount.IsValid() || !h.Amount.IsAllPositive() {
		return fmt.Errorf("invalid transferred amount: %s", h.Amount.String())
	}
	if len(h.Secret) != 0 {
		return fmt.Errorf("the secret length must be zero")
	}
	if h.ExpireHeight < 1 {
		return fmt.Errorf("expire height must be greater than 0")
	}
	return nil
}

// HTLCState represents the state of an HTLC
type HTLCState byte

// HTLCSecret represents the secret of an HTLC
type HTLCSecret []byte

// HTLCSecret represents the hash lock of an HTLC
type HTLCHashLock []byte

const (
	OPEN      HTLCState = 0x00 // claimable
	COMPLETED HTLCState = 0x01 // claimed
	EXPIRED   HTLCState = 0x02 // expired
)

var (
	// HTLCStateToStringMap is the map for HTLCState to state string
	HTLCStateToStringMap = map[HTLCState]string{
		OPEN:      "open",
		COMPLETED: "completed",
		EXPIRED:   "expired",
	}
	// StringToHTLCStateMap is the map for state string to HTLCState
	StringToHTLCStateMap = map[string]HTLCState{
		"open":      OPEN,
		"completed": COMPLETED,
		"expired":   EXPIRED,
	}
)

// HTLCStateFromString convert string to HTLCState
func HTLCStateFromString(str string) (HTLCState, error) {
	if state, ok := StringToHTLCStateMap[strings.ToLower(str)]; ok {
		return state, nil
	}
	return HTLCState(0xff), fmt.Errorf("'%s' is not a valid HTLC state", str)
}

// Format formats HTLCState
func (state HTLCState) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		_, _ = s.Write([]byte(fmt.Sprintf("%s", state.String())))
	default:
		_, _ = s.Write([]byte(fmt.Sprintf("%v", byte(state))))
	}
}

// String convert HTLCState to format string
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

// MarshalYAML implements the Marshaler of yaml
func (state HTLCState) MarshalYAML() (interface{}, error) {
	return state.String(), nil
}

// UnmarshalYAML implements the Unmarshaler of yaml
func (state *HTLCState) UnmarshalYAML(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	bz, err := HTLCStateFromString(s)
	if err != nil {
		return err
	}
	*state = bz
	return nil
}

// MarshalJSON marshals HTLCState to json
func (state HTLCState) MarshalJSON() ([]byte, error) {
	return json.Marshal(state.String())
}

// UnmarshalJSON unmarshals json to HTLCState
func (state *HTLCState) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return nil
	}
	bz, err := HTLCStateFromString(s)
	if err != nil {
		return err
	}
	*state = bz
	return nil
}

// String convert HTLCSecret to string
func (secret HTLCSecret) String() string {
	return hex.EncodeToString(secret)
}

// Marshal needed for protobuf compatibility
func (secret HTLCSecret) Marshal() ([]byte, error) {
	return secret, nil
}

// Unmarshal needed for protobuf compatibility
func (secret *HTLCSecret) Unmarshal(data []byte) error {
	*secret = data
	return nil
}

// MarshalYAML implements the Marshaler of yaml
func (secret HTLCSecret) MarshalYAML() (interface{}, error) {
	return secret.String(), nil
}

// UnmarshalYAML implements the Unmarshaler of yaml
func (secret *HTLCSecret) UnmarshalYAML(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return nil
	}
	bz, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	*secret = bz
	return nil
}

// MarshalJSON marshals HTLCSecret to json
func (secret HTLCSecret) MarshalJSON() ([]byte, error) {
	return json.Marshal(secret.String())
}

// UnmarshalJSON unmarshals json to HTLCSecret
func (secret *HTLCSecret) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	bz, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	*secret = bz
	return nil
}

// String convert HTLCHashLock to string
func (hashLock HTLCHashLock) String() string {
	return hex.EncodeToString(hashLock)
}

// Marshal needed for protobuf compatibility
func (hashLock HTLCHashLock) Marshal() ([]byte, error) {
	return hashLock, nil
}

// Unmarshal needed for protobuf compatibility
func (hashLock *HTLCHashLock) Unmarshal(data []byte) error {
	*hashLock = data
	return nil
}

// MarshalYAML implements the Marshaler of yaml
func (hashLock HTLCHashLock) MarshalYAML() (interface{}, error) {
	return hashLock.String(), nil
}

// UnmarshalYAML implements the Unmarshaler of yaml
func (hashLock *HTLCHashLock) UnmarshalYAML(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return nil
	}
	bz, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	*hashLock = bz
	return nil
}

// MarshalJSON marshals HTLCHashLock to json
func (hashLock HTLCHashLock) MarshalJSON() ([]byte, error) {
	return json.Marshal(hashLock.String())
}

// UnmarshalJSON unmarshals json to HTLCHashLock
func (hashLock *HTLCHashLock) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	bz, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	*hashLock = bz
	return nil
}
