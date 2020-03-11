package types

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	cmn "github.com/tendermint/tendermint/libs/common"

	sdk "github.com/irisnet/irishub/types"
)

// RequestContext defines a context which holds request-related data
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

// NewRequestContext creates a new RequestContext instance
func NewRequestContext(
	serviceName string,
	providers []sdk.AccAddress,
	consumer sdk.AccAddress,
	input string,
	serviceFeeCap sdk.Coins,
	timeout int64,
	superMode bool,
	repeated bool,
	repeatedFrequency uint64,
	repeatedTotal int64,
	batchCounter uint64,
	batchRequestCount,
	batchResponseCount uint16,
	batchState RequestContextBatchState,
	state RequestContextState,
	responseThreshold uint16,
	moduleName string,
) RequestContext {
	return RequestContext{
		ServiceName:        serviceName,
		Providers:          providers,
		Consumer:           consumer,
		Input:              input,
		ServiceFeeCap:      serviceFeeCap,
		Timeout:            timeout,
		SuperMode:          superMode,
		Repeated:           repeated,
		RepeatedFrequency:  repeatedFrequency,
		RepeatedTotal:      repeatedTotal,
		BatchCounter:       batchCounter,
		BatchRequestCount:  batchRequestCount,
		BatchResponseCount: batchResponseCount,
		BatchState:         batchState,
		State:              state,
		ResponseThreshold:  responseThreshold,
		ModuleName:         moduleName,
	}
}

// String implements Stringer
func (rc RequestContext) String() string {
	providers := ""

	for _, p := range rc.Providers {
		providers += p.String() + ","
	}

	if len(providers) > 0 {
		providers = providers[0 : len(providers)-1]
	}

	return fmt.Sprintf(`RequestContext:
	ServiceName:             %s
	Providers:               %s
	Consumer:                %s
	Input:                   %s
	ServiceFeeCap:           %s
	Timeout:                 %d 
	SuperMode:               %v
	Repeated:                %v
	RepeatedFrequency:       %d
	RepeatedTotal:           %d
	BatchCounter:            %d
	BatchRequestCount:       %d
	BatchResponseCount:      %d
	BatchState:              %s
	State:                   %s
	ResponseThreshold:       %d
	ModuleName:              %s`,
		rc.ServiceName, providers, rc.Consumer, rc.Input, rc.ServiceFeeCap.String(),
		rc.Timeout, rc.SuperMode, rc.Repeated, rc.RepeatedFrequency, rc.RepeatedTotal,
		rc.BatchCounter, rc.BatchRequestCount, rc.BatchResponseCount, rc.BatchState, rc.State,
		rc.ResponseThreshold, rc.ModuleName,
	)
}

// HumanString implements human stringer
func (rc RequestContext) HumanString(converter sdk.CoinsConverter) string {
	providers := ""

	for _, p := range rc.Providers {
		providers += p.String() + ","
	}

	if len(providers) > 0 {
		providers = providers[0 : len(providers)-1]
	}

	return fmt.Sprintf(`RequestContext:
	ServiceName:             %s
	Providers:               %s
	Consumer:                %s
	Input:                   %s
	ServiceFeeCap:           %s
	Timeout:                 %d 
	SuperMode:               %v
	Repeated:                %v
	RepeatedFrequency:       %d
	RepeatedTotal:           %d
	BatchCounter:            %d
	BatchRequestCount:       %d
	BatchResponseCount:      %d
	BatchState:              %s
	State:                   %s
	ResponseThreshold:       %d
	ModuleName:              %s`,
		rc.ServiceName, providers, rc.Consumer, rc.Input, converter.ToMainUnit(rc.ServiceFeeCap),
		rc.Timeout, rc.SuperMode, rc.Repeated, rc.RepeatedFrequency, rc.RepeatedTotal,
		rc.BatchCounter, rc.BatchRequestCount, rc.BatchResponseCount, rc.BatchState, rc.State,
		rc.ResponseThreshold, rc.ModuleName,
	)
}

// CompactRequest defines a compact request with a request context ID
type CompactRequest struct {
	RequestContextID           cmn.HexBytes
	RequestContextBatchCounter uint64
	Provider                   sdk.AccAddress
	ServiceFee                 sdk.Coins
	RequestHeight              int64
}

// NewCompactRequest creates a new CompactRequest instance
func NewCompactRequest(
	requestContextID cmn.HexBytes,
	batchCounter uint64,
	provider sdk.AccAddress,
	serviceFee sdk.Coins,
	requestHeight int64,
) CompactRequest {
	return CompactRequest{
		RequestContextID:           requestContextID,
		RequestContextBatchCounter: batchCounter,
		Provider:                   provider,
		ServiceFee:                 serviceFee,
		RequestHeight:              requestHeight,
	}
}

// Request defines a request which contains the detailed request data
type Request struct {
	ServiceName                string         `json:"service_name"`
	Provider                   sdk.AccAddress `json:"provider"`
	Consumer                   sdk.AccAddress `json:"consumer"`
	Input                      string         `json:"input"`
	ServiceFee                 sdk.Coins      `json:"service_fee"`
	SuperMode                  bool           `json:"super_mode"`
	RequestHeight              int64          `json:"request_height"`
	ExpirationHeight           int64          `json:"expiration_height"`
	RequestContextID           cmn.HexBytes   `json:"request_context_id"`
	RequestContextBatchCounter uint64         `json:"request_context_batch_counter"`
}

// NewRequest creates a new Request instance
func NewRequest(
	serviceName string,
	provider,
	consumer sdk.AccAddress,
	input string,
	serviceFee sdk.Coins,
	superMode bool,
	requestHeight int64,
	expirationHeight int64,
	requestContextID cmn.HexBytes,
	batchCounter uint64,
) Request {
	return Request{
		ServiceName:                serviceName,
		Provider:                   provider,
		Consumer:                   consumer,
		Input:                      input,
		ServiceFee:                 serviceFee,
		SuperMode:                  superMode,
		RequestHeight:              requestHeight,
		ExpirationHeight:           expirationHeight,
		RequestContextID:           requestContextID,
		RequestContextBatchCounter: batchCounter,
	}
}

// String implements Stringer
func (r Request) String() string {
	return fmt.Sprintf(`Request:
	ServiceName:             %s
	Provider:                %s
	Consumer:                %s
	Input:                   %s
	ServiceFee:              %s
	SuperMode:               %v 
	RequestHeight:           %d
	ExpirationHeight:        %d
	RequestContextID:        %s
	BatchCounter:            %d`,
		r.ServiceName, r.Provider, r.Consumer, r.Input, r.ServiceFee.String(),
		r.SuperMode, r.RequestHeight, r.ExpirationHeight,
		r.RequestContextID.String(), r.RequestContextBatchCounter,
	)
}

// HumanString implements human stringer
func (r Request) HumanString(converter sdk.CoinsConverter) string {
	return fmt.Sprintf(`Request:
	ServiceName:             %s
	Provider:                %s
	Consumer:                %s
	Input:                   %s
	ServiceFee:              %s
	SuperMode:               %v 
	RequestHeight:           %d
	ExpirationHeight:        %d
	RequestContextID:        %s
	BatchCounter:            %d`,
		r.ServiceName, r.Provider, r.Consumer, r.Input, converter.ToMainUnit(r.ServiceFee),
		r.SuperMode, r.RequestHeight, r.ExpirationHeight,
		r.RequestContextID.String(), r.RequestContextBatchCounter,
	)
}

// Requests represents a set of requests
type Requests []Request

// String implements Stringer
func (rs Requests) String() string {
	if len(rs) == 0 {
		return "[]"
	}

	var str string
	for _, r := range rs {
		str += r.String() + "\n"
	}

	return str
}

// HumanString implements human stringer
func (rs Requests) HumanString(converter sdk.CoinsConverter) string {
	if len(rs) == 0 {
		return "[]"
	}

	var str string
	for _, r := range rs {
		str += r.HumanString(converter) + "\n"
	}

	return str
}

// Response defines a response
type Response struct {
	Provider                   sdk.AccAddress `json:"provider"`
	Consumer                   sdk.AccAddress `json:"consumer"`
	Result                     string         `json:"result"`
	Output                     string         `json:"output"`
	RequestContextID           cmn.HexBytes   `json:"request_context_id"`
	RequestContextBatchCounter uint64         `json:"request_context_batch_counter"`
}

// NewResponse creates a new Response instance
func NewResponse(
	provider,
	consumer sdk.AccAddress,
	result,
	output string,
	requestContextID cmn.HexBytes,
	batchCounter uint64,
) Response {
	return Response{
		Provider:                   provider,
		Consumer:                   consumer,
		Result:                     result,
		Output:                     output,
		RequestContextID:           requestContextID,
		RequestContextBatchCounter: batchCounter,
	}
}

// String implements Stringer
func (r Response) String() string {
	return fmt.Sprintf(`Response:
	Provider:                %s
	Consumer:                %s
	Result:                  %s
	Output:                  %s
	RequestContextID:        %s
	BatchCounter:            %d`,
		r.Provider, r.Consumer, r.Result, r.Output,
		r.RequestContextID.String(),
		r.RequestContextBatchCounter,
	)
}

// Responses represents a set of responses
type Responses []Response

// String implements Stringer
func (rs Responses) String() string {
	if len(rs) == 0 {
		return "[]"
	}

	var str string
	for _, r := range rs {
		str += r.String() + "\n"
	}

	return str
}

// Result defines a struct for the response result
type Result struct {
	Code    uint16 `json:"code"`
	Message string `json:"message"`
}

// ParseResult parses the given string to Result
func ParseResult(result string) (Result, sdk.Error) {
	var r Result

	if err := json.Unmarshal([]byte(result), &r); err != nil {
		return r, ErrInvalidResponseResult(DefaultCodespace, fmt.Sprintf("failed to unmarshal the result: %s", err))
	}

	return r, nil
}

// EarnedFees defines a struct for the fees earned by the provider
type EarnedFees struct {
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

// NewEarnedFees creates a new EarnedFees instance
func NewEarnedFees(address sdk.AccAddress, coins sdk.Coins) EarnedFees {
	return EarnedFees{
		Address: address,
		Coins:   coins,
	}
}

// String implements Stringer
func (e EarnedFees) String() string {
	return fmt.Sprintf(`EarnedFees:
	Address:                 %s
	Coins:                   %s`,
		e.Address, e.Coins.String(),
	)
}

// HumanString implements human stringer
func (e EarnedFees) HumanString(converter sdk.CoinsConverter) string {
	return fmt.Sprintf(`EarnedFees:
	Address:                 %s
	Coins:                   %s`,
		e.Address, converter.ToMainUnit(e.Coins),
	)
}

// RequestContextState defines the state for the request context
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

// RequestContextBatchState defines the current batch state for the request context
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

// ResponseCallback defines the response callback interface
type ResponseCallback func(ctx sdk.Context, requestContextID cmn.HexBytes, responses []string, err error)

const (
	requestIDLen = 42
)

// ConvertRequestID converts the given string to request ID
func ConvertRequestID(requestIDStr string) (cmn.HexBytes, error) {
	if len(requestIDStr) != 2*requestIDLen {
		return nil, errors.New("invalid request id")
	}

	requestID, err := hex.DecodeString(requestIDStr)
	if err != nil {
		return nil, errors.New("invalid request id")
	}

	return requestID, nil
}

// GenerateRequestContextID generates a unique ID for the request context from the specified params
func GenerateRequestContextID(blockHeight int64, intraCounter int16) cmn.HexBytes {
	bz := make([]byte, 10)

	binary.BigEndian.PutUint64(bz, uint64(blockHeight))
	binary.BigEndian.PutUint16(bz[8:], uint16(intraCounter))

	return sdk.SHA256(bz)
}

// GenerateRequestID generates a unique request ID from the given params
func GenerateRequestID(requestContextID cmn.HexBytes, requestContextBatchCounter uint64, batchRequestIndex int16) cmn.HexBytes {
	bz := make([]byte, 10)

	binary.BigEndian.PutUint64(bz, requestContextBatchCounter)
	binary.BigEndian.PutUint16(bz[8:], uint16(batchRequestIndex))

	return append(requestContextID, bz...)
}
