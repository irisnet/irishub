package types

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/irisnet/irishub/types"
)

// RequestContext defines a context which holds request-related data
type RequestContext struct {
	ServiceName        string                   `json:"service_name"`
	Providers          []sdk.AccAddress         `json:"providers"`
	Consumer           sdk.AccAddress           `json:"consumer"`
	Input              string                   `json:"input"`
	ServiceFeeCap      sdk.Coins                `json:"service_fee_cap"`
	Profiling          bool                     `json:"profiling"`
	Timeout            int64                    `json:"timeout"`
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
	profiling bool,
	timeout int64,
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
		Profiling:          profiling,
		Timeout:            timeout,
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

// HumanString implements human stringer
func (rc RequestContext) HumanString(converter sdk.CoinsConverter) string {
	providers := ""

	for i, p := range len(rc.Providers) {
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
	RequestContextID           []byte
	RequestContextBatchCounter uint64
	Provider                   sdk.AccAddress
	ServiceFee                 sdk.Coins
	RequestHeight              int64
}

// NewCompactRequest creates a new CompactRequest instance
func NewCompactRequest(
	requestContextID []byte,
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
	Profiling                  bool           `json:"profiling"`
	RequestHeight              int64          `json:"request_height"`
	ExpirationHeight           int64          `json:"expiration_height"`
	RequestContextID           []byte         `json:"request_context_id"`
	RequestContextBatchCounter uint64         `json:"request_context_batch_counter"`
}

// NewRequest creates a new Request instance
func NewRequest(
	serviceName string,
	provider,
	consumer sdk.AccAddress,
	input string,
	serviceFee sdk.Coins,
	profiling bool,
	requestHeight int64,
	expirationHeight int64,
	requestContextID []byte,
	batchCounter uint64,
) Request {
	return Request{
		ServiceName:                serviceName,
		Provider:                   provider,
		Consumer:                   consumer,
		Input:                      input,
		ServiceFee:                 serviceFee,
		Profiling:                  profiling,
		RequestHeight:              requestHeight,
		ExpirationHeight:           expirationHeight,
		RequestContextID:           requestContextID,
		RequestContextBatchCounter: batchCounter,
	}
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
		r.ServiceName, r.Provider, r.Consumer, rc.Input, converter.ToMainUnit(ServiceFee),
		r.SuperMode, r.RequestHeight, r.ExpirationHeight,
		hex.EncodeToString(r.RequestContextID), r.RequestContextBatchCounter,
	)
}

// Requests represents a set of requests
var Requests []Request

// String implements Stringer
func (rs Requests) String() string {
	if len(rs) == 0 {
		return "[]"
	}

	var str string
	for _, r := range rs {
		str += r.HumanString() + "\n"
	}

	return str
}

// Response defines a response
type Response struct {
	Provider                   sdk.AccAddress `json:"provider"`
	Consumer                   sdk.AccAddress `json:"consumer"`
	Output                     string         `json:"output"`
	Error                      string         `json:"error"`
	RequestContextID           []byte         `json:"request_context_id"`
	RequestContextBatchCounter uint64         `json:"request_context_batch_counter"`
}

// NewResponse creates a new Response instance
func NewResponse(
	provider,
	consumer sdk.AccAddress,
	output,
	err string,
	requestContextID []byte,
	batchCounter uint64,
) Response {
	return Response{
		Provider:                   provider,
		Consumer:                   consumer,
		Output:                     output,
		Error:                      err,
		RequestContextID:           requestContextID,
		RequestContextBatchCounter: batchCounter,
	}
}

// String implements Stringer
func (r Response) String() string {
	return fmt.Sprintf(`Response:
	Provider:                %s
	Consumer:                %s
	Output:                  %s
	Error:                   %s
	RequestContextID:        %s
	BatchCounter:            %d`,
		r.Provider, r.Consumer, r.Output, r.Error,
		hex.EncodeToString(r.RequestContextID),
		r.RequestContextBatchCounter,
	)
}

// Responses represents a set of responses
var Responses []Response

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
		RUNNING:   "running",
		PAUSED:    "paused",
		COMPLETED: "completed",
	}
	StringToRequestContextBatchStateMap = map[string]RequestContextBatchState{
		"running":   RUNNING,
		"paused":    PAUSED,
		"completed": COMPLETED,
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
type ResponseCallback func(ctx sdk.Context, requestContextID []byte, reponses []string)

// ConvertRequestID converts the given string to request ID
func ConvertRequestID(requestID string) ([]byte, error) {
	parts := strings.Split(requestID, "-")
	if len(parts) != 3 {
		return nil, errors.New("invalid request id")
	}

	requestContextID, err := hex.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}

	batchCounter, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, err
	}

	batchRequestIndex, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, err
	}

	return GenerateRequestID(requestContextID, uint64(batchCounter), int16(batchRequestIndex)), nil
}

// RequestIDToString returns the string representation of the given request ID
// request ID layout: 32+8+2 bytes
func RequestIDToString(requestID []byte) (string, error) {
	if len(requestID) != 42 {
		return "", errors.New("invalid request id")
	}

	requestContextID := requestID[0:32]
	batchCounter := requestID[32:40]
	batchRequestIndex := requestID[40:42]

	return fmt.Sprintf(
			"%s-%d-%d",
			hex.EncodeToString(requestContextID),
			binary.BigEndian.Uint64(batchCounter),
			binary.BigEndian.Uint16(batchRequestIndex),
		),
		nil
}

// GenerateRequestContextID generates a unique ID for the request context from the specified params
func GenerateRequestContextID(blockHeight int64, intraCounter int16) []byte {
	bz := make([]byte, 10)

	binary.BigEndian.PutUint64(bz, uint64(blockHeight))
	binary.BigEndian.PutUint16(bz[8:], uint16(intraCounter))

	return sdk.SHA256(bz)
}

// GenerateRequestID generates a unique request ID from the given params
func GenerateRequestID(requestContextID []byte, requestContextBatchCounter uint64, batchRequestIndex int16) []byte {
	bz := make([]byte, 10)

	binary.BigEndian.PutUint64(bz, requestContextBatchCounter)
	binary.BigEndian.PutUint16(bz[8:], uint16(batchRequestIndex))

	return append(requestContextID, bz...)
}
