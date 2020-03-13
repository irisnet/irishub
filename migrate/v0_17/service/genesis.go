package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	cmn "github.com/tendermint/tendermint/libs/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Params            Params                    `json:"params"`             // service params
	Definitions       []ServiceDefinition       `json:"definitions"`        // service definitions
	Bindings          []ServiceBinding          `json:"bindings"`           // service bindings
	WithdrawAddresses map[string]sdk.AccAddress `json:"withdraw_addresses"` // withdraw addresses
	RequestContexts   map[string]RequestContext `json:"request_contexts"`   // request contexts
}

type Params struct {
	MaxRequestTimeout    int64         `json:"max_request_timeout"`
	MinDepositMultiple   int64         `json:"min_deposit_multiple"`
	MinDeposit           sdk.Coins     `json:"min_deposit"`
	ServiceFeeTax        sdk.Dec       `json:"service_fee_tax"`
	SlashFraction        sdk.Dec       `json:"slash_fraction"`
	ComplaintRetrospect  time.Duration `json:"complaint_retrospect"`
	ArbitrationTimeLimit time.Duration `json:"arbitration_time_limit"`
	TxSizeLimit          uint64        `json:"tx_size_limit"`
}

type ServiceDefinition struct {
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	Tags              []string       `json:"tags"`
	Author            sdk.AccAddress `json:"author"`
	AuthorDescription string         `json:"author_description"`
	Schemas           string         `json:"schemas"`
}

type ServiceBinding struct {
	ServiceName  string         `json:"service_name"`
	Provider     sdk.AccAddress `json:"provider"`
	Deposit      sdk.Coins      `json:"deposit"`
	Pricing      string         `json:"pricing"`
	Available    bool           `json:"available"`
	DisabledTime time.Time      `json:"disabled_time"`
}

type RequestContext struct {
	ServiceName        string                   `json:"service_name"`
	Providers          []sdk.AccAddress         `json:"providers"`
	Consumer           sdk.AccAddress           `json:"consumer"`
	Input              string                   `json:"input"`
	ServiceFeeCap      sdk.Coins                `json:"service_fee_cap"`
	Timeout            int64                    `json:"timeout"`
	SuperMode          bool                     `json:"super_mode"`
	Repeated           bool                     `json:"repeated"`
	RepeatedFrequency  uint64                   `json:"repeated_frequency"`
	RepeatedTotal      int64                    `json:"repeated_total"`
	BatchCounter       uint64                   `json:"batch_counter"`
	BatchRequestCount  uint16                   `json:"batch_request_count"`
	BatchResponseCount uint16                   `json:"batch_response_count"`
	BatchState         RequestContextBatchState `json:"batch_state"`
	State              RequestContextState      `json:"state"`
	ResponseThreshold  uint16                   `json:"response_threshold"`
	ModuleName         string                   `json:"module_name"`
}

type CompactRequest struct {
	RequestContextID           cmn.HexBytes
	RequestContextBatchCounter uint64
	Provider                   sdk.AccAddress
	ServiceFee                 sdk.Coins
	RequestHeight              int64
}

type RequestContextBatchState byte

const (
	BATCHRUNNING   RequestContextBatchState = 0x00 // running
	BATCHCOMPLETED RequestContextBatchState = 0x01 // completed
)

var (
	RequestContextBatchStateToStringMap = map[RequestContextBatchState]string{
		BATCHRUNNING:   "running",
		BATCHCOMPLETED: "completed",
	}
	StringToRequestContextBatchStateMap = map[string]RequestContextBatchState{
		"running":   BATCHRUNNING,
		"completed": BATCHCOMPLETED,
	}
)

func RequestContextBatchStateFromString(str string) (RequestContextBatchState, error) {
	if state, ok := StringToRequestContextBatchStateMap[strings.ToLower(str)]; ok {
		return state, nil
	}
	return RequestContextBatchState(0xff), fmt.Errorf("'%s' is not a valid request context batch state", str)
}

func (state RequestContextBatchState) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(state.String()))
	default:
		s.Write([]byte(fmt.Sprintf("%v", byte(state))))
	}
}

func (state RequestContextBatchState) String() string {
	return RequestContextBatchStateToStringMap[state]
}

// Marshal needed for protobuf compatibility
func (state RequestContextBatchState) Marshal() ([]byte, error) {
	return []byte{byte(state)}, nil
}

// Unmarshal needed for protobuf compatibility
func (state *RequestContextBatchState) Unmarshal(data []byte) error {
	*state = RequestContextBatchState(data[0])
	return nil
}

// Marshals to JSON using string
func (state RequestContextBatchState) MarshalJSON() ([]byte, error) {
	return json.Marshal(state.String())
}

// Unmarshals from JSON
func (state *RequestContextBatchState) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return nil
	}

	bz, err := RequestContextBatchStateFromString(s)
	if err != nil {
		return err
	}

	*state = bz
	return nil
}

type RequestContextState byte

const (
	RUNNING   RequestContextState = 0x00 // running
	PAUSED    RequestContextState = 0x01 // paused
	COMPLETED RequestContextState = 0x02 // completed
)

var (
	RequestContextStateToStringMap = map[RequestContextState]string{
		RUNNING:   "running",
		PAUSED:    "paused",
		COMPLETED: "completed",
	}
	StringToRequestContextStateMap = map[string]RequestContextState{
		"running":   RUNNING,
		"paused":    PAUSED,
		"completed": COMPLETED,
	}
)

func RequestContextStateFromString(str string) (RequestContextState, error) {
	if state, ok := StringToRequestContextStateMap[strings.ToLower(str)]; ok {
		return state, nil
	}
	return RequestContextState(0xff), fmt.Errorf("'%s' is not a valid request context state", str)
}

func (state RequestContextState) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(state.String()))
	default:
		s.Write([]byte(fmt.Sprintf("%v", byte(state))))
	}
}

func (state RequestContextState) String() string {
	return RequestContextStateToStringMap[state]
}

// Marshal needed for protobuf compatibility
func (state RequestContextState) Marshal() ([]byte, error) {
	return []byte{byte(state)}, nil
}

// Unmarshal needed for protobuf compatibility
func (state *RequestContextState) Unmarshal(data []byte) error {
	*state = RequestContextState(data[0])
	return nil
}

// Marshals to JSON using string
func (state RequestContextState) MarshalJSON() ([]byte, error) {
	return json.Marshal(state.String())
}

// Unmarshals from JSON
func (state *RequestContextState) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return nil
	}

	bz, err := RequestContextStateFromString(s)
	if err != nil {
		return err
	}

	*state = bz
	return nil
}
