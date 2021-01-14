package types

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

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
	batchResponseCount uint32,
	batchResponseThreshold uint32,
	batchState RequestContextBatchState,
	state RequestContextState,
	responseThreshold uint32,
	moduleName string,
) RequestContext {
	pds := make([]string, len(providers))
	for i, provider := range providers {
		pds[i] = provider.String()
	}

	return RequestContext{
		ServiceName:            serviceName,
		Providers:              pds,
		Consumer:               consumer.String(),
		Input:                  input,
		ServiceFeeCap:          serviceFeeCap,
		Timeout:                timeout,
		SuperMode:              superMode,
		Repeated:               repeated,
		RepeatedFrequency:      repeatedFrequency,
		RepeatedTotal:          repeatedTotal,
		BatchCounter:           batchCounter,
		BatchRequestCount:      batchRequestCount,
		BatchResponseCount:     batchResponseCount,
		BatchResponseThreshold: batchResponseThreshold,
		BatchState:             batchState,
		State:                  state,
		ResponseThreshold:      responseThreshold,
		ModuleName:             moduleName,
	}
}

// Validate validates the request context
func (rc RequestContext) Validate() error {
	if err := ValidateServiceName(rc.ServiceName); err != nil {
		return err
	}

	pds := make([]sdk.AccAddress, len(rc.Providers))
	for i, provider := range rc.Providers {
		pd, err := sdk.AccAddressFromBech32(provider)
		if err != nil {
			return err
		}
		pds[i] = pd
	}

	if err := ValidateProviders(pds); err != nil {
		return err
	}

	if err := ValidateConsumer(rc.Consumer); err != nil {
		return err
	}

	if err := ValidateInput(rc.Input); err != nil {
		return err
	}

	return ValidateServiceFeeCap(rc.ServiceFeeCap)
}

// Empty returns true if empty
func (rc RequestContext) Empty() bool {
	// TODO: use rc.ID
	return len(rc.Consumer) == 0
}

// NewCompactRequest creates a new CompactRequest instance
func NewCompactRequest(
	requestContextID tmbytes.HexBytes,
	batchCounter uint64,
	provider sdk.AccAddress,
	serviceFee sdk.Coins,
	requestHeight int64,
	expirationHeight int64,
) CompactRequest {
	return CompactRequest{
		RequestContextId:           requestContextID.String(),
		RequestContextBatchCounter: batchCounter,
		Provider:                   provider.String(),
		ServiceFee:                 serviceFee,
		RequestHeight:              requestHeight,
		ExpirationHeight:           expirationHeight,
	}
}

// NewRequest creates a new Request instance
func NewRequest(
	id tmbytes.HexBytes,
	serviceName string,
	provider,
	consumer sdk.AccAddress,
	input string,
	serviceFee sdk.Coins,
	superMode bool,
	requestHeight int64,
	expirationHeight int64,
	requestContextID tmbytes.HexBytes,
	batchCounter uint64,
) Request {
	return Request{
		Id:                         id.String(),
		ServiceName:                serviceName,
		Provider:                   provider.String(),
		Consumer:                   consumer.String(),
		Input:                      input,
		ServiceFee:                 serviceFee,
		SuperMode:                  superMode,
		RequestHeight:              requestHeight,
		ExpirationHeight:           expirationHeight,
		RequestContextId:           requestContextID.String(),
		RequestContextBatchCounter: batchCounter,
	}
}

// Empty returns true if empty
func (r Request) Empty() bool {
	return len(r.Id) == 0
}

// NewResponse creates a new Response instance
func NewResponse(
	provider,
	consumer sdk.AccAddress,
	result,
	output string,
	requestContextID tmbytes.HexBytes,
	batchCounter uint64,
) Response {
	return Response{
		Provider:                   provider.String(),
		Consumer:                   consumer.String(),
		Result:                     result,
		Output:                     output,
		RequestContextId:           requestContextID.String(),
		RequestContextBatchCounter: batchCounter,
	}
}

// Empty returns true if empty
func (r Response) Empty() bool {
	return len(r.RequestContextId) == 0
}

// ResultCode defines the type for the result code
type ResultCode uint16

const (
	// ResultOK means that the request is well processed
	ResultOK ResultCode = 200

	// ResultBadRequest means that the request is incorrect
	ResultBadRequest ResultCode = 400

	// ResultServerError is the code resulting from the internal server error
	ResultServerError ResultCode = 500
)

// Result defines a struct for the response result
type Result struct {
	Code    ResultCode `json:"code"`
	Message string     `json:"message"`
}

// ParseResult parses the given string to Result
func ParseResult(result string) (Result, error) {
	var r Result

	if err := json.Unmarshal([]byte(result), &r); err != nil {
		return r, sdkerrors.Wrapf(ErrInvalidResponseResult, "failed to unmarshal the result: %s", err)
	}

	return r, nil
}

func RequestContextStateFromString(str string) (RequestContextState, error) {
	if state, ok := RequestContextState_value[strings.ToUpper(str)]; ok {
		return RequestContextState(state), nil
	}
	return RequestContextState(0xff), fmt.Errorf("'%s' is not a valid request context state", str)
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

// Marshal needed for protobuf compatibility
func (state RequestContextBatchState) Marshal() ([]byte, error) {
	return []byte{byte(state)}, nil
}

// Unmarshal needed for protobuf compatibility
func (state *RequestContextBatchState) Unmarshal(data []byte) error {
	*state = RequestContextBatchState(data[0])
	return nil
}

// ResponseCallback defines the response callback interface
type ResponseCallback func(ctx sdk.Context, requestContextID tmbytes.HexBytes, responses []string, err error)

// StateCallback defines the state callback interface
type StateCallback func(ctx sdk.Context, requestContextID tmbytes.HexBytes, cause string)

const (
	RequestIDLen = 116
	ContextIDLen = 80
)

// ConvertRequestID converts the given string to request ID
func ConvertRequestID(requestIDStr string) (tmbytes.HexBytes, error) {
	if len(requestIDStr) != RequestIDLen {
		return nil, errors.New("invalid request id")
	}

	requestID, err := hex.DecodeString(requestIDStr)
	if err != nil {
		return nil, errors.New("invalid request id")
	}

	return requestID, nil
}

// GenerateRequestContextID generates a unique ID for the request context from the specified params
func GenerateRequestContextID(txHash []byte, msgIndex int64) tmbytes.HexBytes {
	bz := make([]byte, 8)

	binary.BigEndian.PutUint64(bz, uint64(msgIndex))

	return append(txHash, bz...)
}

// SplitRequestContextID splits the given contextID to txHash and msgIndex
func SplitRequestContextID(contextID tmbytes.HexBytes) (tmbytes.HexBytes, int64, error) {
	if len(contextID) != ContextIDLen/2 {
		return nil, 0, errors.New("invalid request context ID")
	}

	txHash := contextID[0:32]
	msgIndex := int64(binary.BigEndian.Uint64(contextID[32:40]))

	return txHash, msgIndex, nil
}

// GenerateRequestID generates a unique request ID from the given params
func GenerateRequestID(requestContextID tmbytes.HexBytes, requestContextBatchCounter uint64, requestHeight int64, batchRequestIndex int16) tmbytes.HexBytes {
	contextID := make([]byte, len(requestContextID))
	copy(contextID, requestContextID)

	bz := make([]byte, 18)

	binary.BigEndian.PutUint64(bz, requestContextBatchCounter)
	binary.BigEndian.PutUint64(bz[8:], uint64(requestHeight))
	binary.BigEndian.PutUint16(bz[16:], uint16(batchRequestIndex))

	return append(contextID, bz...)
}

// SplitRequestID splits the given requestID to contextID, batchCounter, requestHeight, batchRequestIndex
func SplitRequestID(requestID tmbytes.HexBytes) (tmbytes.HexBytes, uint64, int64, int16, error) {
	if len(requestID) != RequestIDLen/2 {
		return nil, 0, 0, 0, errors.New("invalid request ID")
	}

	contextID := requestID[0:40]
	batchCounter := binary.BigEndian.Uint64(requestID[40:48])
	requestHeight := int64(binary.BigEndian.Uint64(requestID[48:56]))
	batchRequestIndex := int16(binary.BigEndian.Uint16(requestID[56:]))

	return contextID, batchCounter, requestHeight, batchRequestIndex, nil
}
